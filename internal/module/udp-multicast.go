package module

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"sync"
	"time"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/config"
)

func init() {
	RegisterModule(ModuleRegistration{
		Type:  "net.udp.multicast",
		Title: "UDP Multicast",
		ParamsSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"ip": {
					Title:       "IP",
					Description: "the multicast address to listen on",
					Type:        "string",
				},
				"port": {
					Title:       "Port",
					Description: "the port to listen on",
					Type:        "integer",
					Minimum:     jsonschema.Ptr[float64](1024),
					Maximum:     jsonschema.Ptr[float64](65535),
				},
			},
			Required:             []string{"ip", "port"},
			AdditionalProperties: &jsonschema.Schema{Not: &jsonschema.Schema{}},
		},
		New: func(moduleConfig config.ModuleConfig) (common.Module, error) {
			params := moduleConfig.Params
			ipString, err := params.GetString("ip")
			if err != nil {
				return nil, fmt.Errorf("net.udp.multicast ip error: %w", err)
			}

			portNum, err := params.GetInt("port")
			if err != nil {
				return nil, fmt.Errorf("net.udp.multicast port error: %w", err)
			}

			addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", ipString, uint16(portNum)))
			if err != nil {
				return nil, err
			}
			return &UDPMulticast{config: moduleConfig, Addr: addr, logger: CreateLogger(moduleConfig)}, nil
		},
	})
}

type UDPMulticast struct {
	config       config.ModuleConfig
	conn         *net.UDPConn
	ctx          context.Context
	inputHandler common.InputHandler
	Addr         *net.UDPAddr
	logger       *slog.Logger
	cancel       context.CancelFunc
	connMu       sync.Mutex
}

func (um *UDPMulticast) Id() string {
	return um.config.Id
}

func (um *UDPMulticast) Type() string {
	return um.config.Type
}

func (um *UDPMulticast) Start(ctx context.Context, inputHandler common.InputHandler) error {
	um.logger.Debug("running")
	um.inputHandler = inputHandler
	moduleContext, cancel := context.WithCancel(ctx)
	um.ctx = moduleContext
	um.cancel = cancel

	client, err := net.ListenMulticastUDP("udp", nil, um.Addr)
	if err != nil {
		return err
	}
	defer client.Close()

	um.connMu.Lock()
	um.conn = client
	um.connMu.Unlock()

	buffer := make([]byte, 2048)
	for um.ctx.Err() == nil {
		um.conn.SetDeadline(time.Now().Add(time.Millisecond * 200))

		numBytes, _, err := um.conn.ReadFromUDP(buffer)
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				break
			}
			//NOTE(jwetzell) we hit deadline
			if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
				continue
			}
			return err
		}

		if numBytes > 0 {
			message := buffer[:numBytes]

			if um.inputHandler != nil {
				um.inputHandler(um.ctx, um.Id(), message)
			} else {
				um.logger.Error("input received but no input handler is configured")
			}
		}
	}
	<-um.ctx.Done()
	um.logger.Debug("done")
	return nil
}

func (um *UDPMulticast) Output(ctx context.Context, payload any) error {

	payloadBytes, ok := common.GetAnyAsByteSlice(payload)
	if !ok {
		return errors.New("net.udp.multicast can only output bytes")
	}

	if um.conn == nil {
		return errors.New("net.udp.multicast connection is not setup")
	}

	_, err := um.conn.Write(payloadBytes)
	return err
}

func (um *UDPMulticast) Stop() {
	if um.cancel != nil {
		defer um.cancel()
	}
	um.connMu.Lock()
	defer um.connMu.Unlock()
	if um.conn != nil {
		um.conn.Close()
	}
}
