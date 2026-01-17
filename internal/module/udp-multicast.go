package module

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"time"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/route"
)

type UDPMulticast struct {
	config config.ModuleConfig
	conn   *net.UDPConn
	ctx    context.Context
	router route.RouteIO
	Addr   *net.UDPAddr
	logger *slog.Logger
}

func init() {
	RegisterModule(ModuleRegistration{
		Type: "net.udp.multicast",
		New: func(config config.ModuleConfig) (Module, error) {
			params := config.Params
			ip, ok := params["ip"]

			if !ok {
				return nil, errors.New("net.udp.multicast requires an ip parameter")
			}

			ipString, ok := ip.(string)

			if !ok {
				return nil, errors.New("net.udp.multicast ip must be a string")
			}

			port, ok := params["port"]
			if !ok {
				return nil, errors.New("net.udp.multicast requires a port parameter")
			}

			portNum, ok := port.(float64)

			if !ok {
				return nil, errors.New("net.udp.multicast port must be a number")
			}

			addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", ipString, uint16(portNum)))
			if err != nil {
				return nil, err
			}
			return &UDPMulticast{config: config, Addr: addr, logger: CreateLogger(config)}, nil
		},
	})
}

func (um *UDPMulticast) Id() string {
	return um.config.Id
}

func (um *UDPMulticast) Type() string {
	return um.config.Type
}

func (um *UDPMulticast) Run(ctx context.Context) error {

	router, ok := ctx.Value(route.RouterContextKey).(route.RouteIO)

	if !ok {
		return errors.New("net.udp.multicast unable to get router from context")
	}
	um.router = router
	um.ctx = ctx

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

	payloadBytes, ok := payload.([]byte)
	if !ok {
		return errors.New("net.udp.multicast can only output bytes")
	}

	if um.conn == nil {
		return errors.New("net.udp.multicast connection is not setup")
	}

	_, err := um.conn.Write(payloadBytes)
	return err
}
