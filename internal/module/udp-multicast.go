//go:build !js

package module

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"time"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/config"
)

type UDPMulticast struct {
	config config.ModuleConfig
	conn   *net.UDPConn
	ctx    context.Context
	router common.RouteIO
	Addr   *net.UDPAddr
	logger *slog.Logger
	cancel context.CancelFunc
}

func init() {
	RegisterModule(ModuleRegistration{
		Type:  "net.udp.multicast",
		Title: "UDP Multicast",
		ParamsSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"ip": {
					Title: "IP",
					Type:  "string",
				},
				"port": {
					Title:   "Port",
					Type:    "integer",
					Minimum: jsonschema.Ptr[float64](1024),
					Maximum: jsonschema.Ptr[float64](65535),
				},
			},
			Required:             []string{"ip", "port"},
			AdditionalProperties: nil,
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

func (um *UDPMulticast) Id() string {
	return um.config.Id
}

func (um *UDPMulticast) Type() string {
	return um.config.Type
}

func (um *UDPMulticast) Start(ctx context.Context) error {
	um.logger.Debug("running")
	router, ok := ctx.Value(common.RouterContextKey).(common.RouteIO)

	if !ok {
		return errors.New("net.udp.multicast unable to get router from context")
	}
	um.router = router
	moduleContext, cancel := context.WithCancel(ctx)
	um.ctx = moduleContext
	um.cancel = cancel

	client, err := net.ListenMulticastUDP("udp", nil, um.Addr)
	if err != nil {
		return err
	}
	defer client.Close()

	um.conn = client

	buffer := make([]byte, 2048)
	for {
		select {
		case <-um.ctx.Done():
			// TODO(jwetzell): cleanup?
			um.logger.Debug("done")
			return nil
		default:
			um.conn.SetDeadline(time.Now().Add(time.Millisecond * 200))

			numBytes, _, err := um.conn.ReadFromUDP(buffer)
			if err != nil {
				//NOTE(jwetzell) we hit deadline
				if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
					continue
				}
				return err
			}

			if numBytes > 0 {
				message := buffer[:numBytes]

				if um.router != nil {
					um.router.HandleInput(um.ctx, um.Id(), message)
				} else {
					um.logger.Error("input received but no router is configured")
				}
			}
		}
	}
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
	um.cancel()
}
