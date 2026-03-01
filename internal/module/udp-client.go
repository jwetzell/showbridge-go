package module

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/route"
)

type UDPClient struct {
	config config.ModuleConfig
	Addr   *net.UDPAddr
	Port   uint16
	conn   *net.UDPConn
	ctx    context.Context
	router route.RouteIO
	logger *slog.Logger
	cancel context.CancelFunc
}

func init() {
	RegisterModule(ModuleRegistration{
		Type: "net.udp.client",
		New: func(config config.ModuleConfig) (Module, error) {
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

func (uc *UDPClient) Id() string {
	return uc.config.Id
}

func (uc *UDPClient) Type() string {
	return uc.config.Type
}

func (uc *UDPClient) SetupConn() error {
	client, err := net.DialUDP("udp", nil, uc.Addr)
	uc.conn = client
	return err
}

func (uc *UDPClient) Start(ctx context.Context) error {
	uc.logger.Debug("running")
	router, ok := ctx.Value(route.RouterContextKey).(route.RouteIO)

	if !ok {
		return errors.New("net.udp.client unable to get router from context")
	}
	uc.router = router
	moduleContext, cancel := context.WithCancel(ctx)
	uc.ctx = moduleContext
	uc.cancel = cancel

	err := uc.SetupConn()
	if err != nil {
		return err
	}

	<-uc.ctx.Done()
	uc.logger.Debug("done")
	if uc.conn != nil {
		uc.conn.Close()
	}
	return nil
}

func (uc *UDPClient) Output(ctx context.Context, payload any) error {

	payloadBytes, ok := payload.([]byte)
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
	uc.cancel()
}
