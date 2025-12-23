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
}

func init() {
	RegisterModule(ModuleRegistration{
		Type: "net.udp.client",
		New: func(ctx context.Context, config config.ModuleConfig, router route.RouteIO) (Module, error) {
			params := config.Params
			host, ok := params["host"]

			if !ok {
				return nil, errors.New("net.udp.client requires a host parameter")
			}

			hostString, ok := host.(string)

			if !ok {
				return nil, errors.New("net.udp.client host must be a string")
			}

			port, ok := params["port"]
			if !ok {
				return nil, errors.New("net.udp.client requires a port parameter")
			}

			portNum, ok := port.(float64)

			if !ok {
				return nil, errors.New("net.udp.client port must be a number")
			}

			addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", hostString, uint16(portNum)))
			if err != nil {
				return nil, err
			}

			return &UDPClient{Addr: addr, config: config, ctx: ctx, router: router, logger: slog.Default().With("component", "module", "id", config.Id)}, nil
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

func (uc *UDPClient) Run() error {

	err := uc.SetupConn()
	if err != nil {
		return err
	}

	<-uc.ctx.Done()
	uc.logger.Debug("router context done in module")
	if uc.conn != nil {
		uc.conn.Close()
	}
	return nil
}

func (uc *UDPClient) Output(payload any) error {

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
