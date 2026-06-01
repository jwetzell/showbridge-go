package module

import (
	"context"
	"encoding/json"
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
		Type:  "net.udp.server",
		Title: "UDP Server",
		ParamsSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"ip": {
					Title:       "IP",
					Description: "the IP address to bind the UDP server to",
					Type:        "string",
					Default:     json.RawMessage(`"0.0.0.0"`),
				},
				"port": {
					Title:       "Port",
					Description: "the port for the UDP server to listen on",
					Type:        "integer",
					Minimum:     jsonschema.Ptr[float64](1024),
					Maximum:     jsonschema.Ptr[float64](65535),
				},
				"bufferSize": {
					Title:       "Buffer Size",
					Description: "the size of the buffer for incoming UDP messages",
					Type:        "integer",
					Minimum:     jsonschema.Ptr[float64](1),
					Maximum:     jsonschema.Ptr[float64](65535),
					Default:     json.RawMessage("2048"),
				},
			},
			Required:             []string{"port"},
			AdditionalProperties: &jsonschema.Schema{Not: &jsonschema.Schema{}},
		},
		New: func(moduleConfig config.ModuleConfig) (common.Module, error) {
			params := moduleConfig.Params
			portNum, err := params.GetInt("port")
			if err != nil {
				return nil, fmt.Errorf("net.udp.server port error: %w", err)
			}

			ipString, err := params.GetString("ip")
			if err != nil {
				if errors.Is(err, config.ErrParamNotFound) {
					ipString = "0.0.0.0"
				} else {
					return nil, fmt.Errorf("net.udp.server ip error: %w", err)
				}
			}

			addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", ipString, uint16(portNum)))
			if err != nil {
				return nil, err
			}

			bufferSizeNum, err := params.GetInt("bufferSize")
			if err != nil {
				if errors.Is(err, config.ErrParamNotFound) {
					bufferSizeNum = 2048
				} else {
					return nil, fmt.Errorf("net.udp.server bufferSize error: %w", err)
				}
			}
			return &UDPServer{Addr: addr, BufferSize: bufferSizeNum, config: moduleConfig, logger: CreateLogger(moduleConfig)}, nil
		},
	})
}

type UDPServer struct {
	Addr         *net.UDPAddr
	BufferSize   int
	config       config.ModuleConfig
	ctx          context.Context
	inputHandler common.InputHandler
	logger       *slog.Logger
	cancel       context.CancelFunc
	listener     *net.UDPConn
	listenerMu   sync.Mutex
}

func (us *UDPServer) Id() string {
	return us.config.Id
}

func (us *UDPServer) Type() string {
	return us.config.Type
}

func (us *UDPServer) Start(ctx context.Context, inputHandler common.InputHandler) error {
	us.logger.Debug("running")
	us.inputHandler = inputHandler
	moduleContext, cancel := context.WithCancel(ctx)
	us.ctx = moduleContext
	us.cancel = cancel

	listener, err := net.ListenUDP("udp", us.Addr)
	if err != nil {
		return err
	}
	us.listenerMu.Lock()
	us.listener = listener
	us.listenerMu.Unlock()

	buffer := make([]byte, us.BufferSize)
	for us.ctx.Err() == nil {
		listener.SetDeadline(time.Now().Add(time.Millisecond * 200))

		numBytes, _, err := listener.ReadFromUDP(buffer)
		if err != nil {
			//NOTE(jwetzell) we hit deadline
			if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
				continue
			}
			break
		}
		message := buffer[:numBytes]
		if us.inputHandler != nil {
			us.inputHandler(us.ctx, us.Id(), message)
		} else {
			us.logger.Error("input received but no input handler is configured")
		}
	}
	<-us.ctx.Done()
	us.logger.Debug("done")
	return nil
}

func (us *UDPServer) Output(ctx context.Context, payload any) error {
	return errors.New("net.udp.server output is not implemented")
}

func (us *UDPServer) Stop() {
	if us.cancel != nil {
		defer us.cancel()
	}
	us.listenerMu.Lock()
	defer us.listenerMu.Unlock()
	if us.listener != nil {
		us.listener.Close()
	}
}
