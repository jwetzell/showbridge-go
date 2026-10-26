package module

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"sync"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/jwetzell/showbridge-go/config"
	"github.com/jwetzell/showbridge-go/params"
	"github.com/jwetzell/showbridge-go/routing"
)

func init() {
	RegisterModule(ModuleRegistration{
		Type:        "net.udp.client",
		Title:       "UDP Client",
		Description: "Sends UDP packets to the specified server",
		ParamsSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"host": {
					Title:       "Host",
					Description: "the hostname or IP address of the UDP server to connect to",
					Type:        "string",
				},
				"port": {
					Title:       "Port",
					Description: "the port of the UDP server to connect to",
					Type:        "integer",
					Minimum:     jsonschema.Ptr[float64](1),
					Maximum:     jsonschema.Ptr[float64](65535),
				},
			},
			Required:             []string{"host", "port"},
			AdditionalProperties: &jsonschema.Schema{Not: &jsonschema.Schema{}},
		},
		New: func(config config.ModuleConfig) (routing.Module, error) {
			moduleParams := config.Params
			hostString, err := moduleParams.GetString("host")
			if err != nil {
				return nil, fmt.Errorf("net.udp.client host error: %w", err)
			}

			portNum, err := moduleParams.GetInt("port")
			if err != nil {
				return nil, fmt.Errorf("net.udp.client port error: %w", err)
			}

			addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", hostString, uint16(portNum)))
			if err != nil {
				return nil, err
			}
			return &UDPClient{Addr: addr, config: config, logger: CreateLogger(config)}, nil
		},
	})
}

type UDPClient struct {
	config config.ModuleConfig
	Addr   *net.UDPAddr
	Port   uint16
	conn   *net.UDPConn
	ctx    context.Context
	router routing.Router
	logger *slog.Logger
	cancel context.CancelFunc
	mu     sync.Mutex
}

func (uc *UDPClient) Start(ctx context.Context, router routing.Router) error {
	uc.logger.Debug("running")
	uc.mu.Lock()
	uc.router = router
	moduleContext, cancel := context.WithCancel(ctx)
	uc.ctx = moduleContext
	uc.cancel = cancel

	client, err := net.DialUDP("udp", nil, uc.Addr)
	if err != nil {
		uc.mu.Unlock()
		return err
	}
	uc.conn = client
	uc.mu.Unlock()

	<-uc.ctx.Done()
	uc.logger.Debug("done")
	return nil
}

func (uc *UDPClient) Output(ctx context.Context, payload any) error {
	uc.mu.Lock()
	defer uc.mu.Unlock()
	payloadBytes, ok := params.GetAnyAsByteSlice(payload)
	if !ok {
		return errors.New("net.udp.client is only able to output bytes")
	}
	if uc.conn != nil {
		_, err := uc.conn.Write(payloadBytes)

		if err != nil {
			return err
		}
	} else {
		return errors.New("net.udp.client client is not setup")
	}
	return nil
}

func (uc *UDPClient) Stop() {
	uc.mu.Lock()
	defer uc.mu.Unlock()
	if uc.cancel != nil {
		defer uc.cancel()
	}
	if uc.conn != nil {
		uc.conn.Close()
	}
}
