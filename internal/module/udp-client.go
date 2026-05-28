package module

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"sync"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/config"
)

func init() {
	RegisterModule(ModuleRegistration{
		Type:  "net.udp.client",
		Title: "UDP Client",
		ParamsSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"host": {
					Title: "Host",
					Type:  "string",
				},
				"port": {
					Title:   "Port",
					Type:    "integer",
					Minimum: jsonschema.Ptr[float64](1),
					Maximum: jsonschema.Ptr[float64](65535),
				},
			},
			Required:             []string{"host", "port"},
			AdditionalProperties: &jsonschema.Schema{Not: &jsonschema.Schema{}},
		},
		New: func(config config.ModuleConfig) (common.Module, error) {
			params := config.Params
			hostString, err := params.GetString("host")
			if err != nil {
				return nil, fmt.Errorf("net.udp.client host error: %w", err)
			}

			portNum, err := params.GetInt("port")
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
	config       config.ModuleConfig
	Addr         *net.UDPAddr
	Port         uint16
	conn         *net.UDPConn
	ctx          context.Context
	inputHandler common.InputHandler
	logger       *slog.Logger
	cancel       context.CancelFunc
	connMu       sync.Mutex
}

func (uc *UDPClient) Id() string {
	return uc.config.Id
}

func (uc *UDPClient) Type() string {
	return uc.config.Type
}

func (uc *UDPClient) SetupConn() error {
	uc.connMu.Lock()
	defer uc.connMu.Unlock()
	client, err := net.DialUDP("udp", nil, uc.Addr)
	uc.conn = client
	return err
}

func (uc *UDPClient) Start(ctx context.Context, inputHandler common.InputHandler) error {
	uc.logger.Debug("running")
	uc.inputHandler = inputHandler
	moduleContext, cancel := context.WithCancel(ctx)
	uc.ctx = moduleContext
	uc.cancel = cancel

	err := uc.SetupConn()
	if err != nil {
		return err
	}

	<-uc.ctx.Done()
	uc.logger.Debug("done")
	return nil
}

func (uc *UDPClient) Output(ctx context.Context, payload any) error {
	uc.connMu.Lock()
	defer uc.connMu.Unlock()
	payloadBytes, ok := common.GetAnyAsByteSlice(payload)
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
	if uc.cancel != nil {
		defer uc.cancel()
	}
	uc.connMu.Lock()
	defer uc.connMu.Unlock()
	if uc.conn != nil {
		uc.conn.Close()
	}
}
