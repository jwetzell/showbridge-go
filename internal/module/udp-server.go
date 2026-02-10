package module

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net"
	"time"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/route"
)

type UDPServer struct {
	Addr       *net.UDPAddr
	BufferSize int
	config     config.ModuleConfig
	ctx        context.Context
	router     route.RouteIO
	logger     *slog.Logger
	cancel     context.CancelFunc
}

func init() {
	RegisterModule(ModuleRegistration{
		Type: "net.udp.server",
		New: func(config config.ModuleConfig) (Module, error) {
			params := config.Params
			port, ok := params["port"]
			if !ok {
				return nil, errors.New("net.udp.server requires a port parameter")
			}

			portNum, ok := port.(float64)

			if !ok {
				return nil, errors.New("net.udp.server port must be a number")
			}

			ipString := "0.0.0.0"

			ip, ok := params["ip"]
			if ok {

				specificIpString, ok := ip.(string)

				if !ok {
					return nil, errors.New("net.udp.server ip must be a string")
				}
				ipString = specificIpString
			}

			addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", ipString, uint16(portNum)))
			if err != nil {
				log.Fatalf("error resolving UDP address: %v", err)
			}

			bufferSizeNum := 2048
			bufferSize, ok := params["bufferSize"]

			if ok {
				bufferSizeFloat, ok := bufferSize.(float64)

				if !ok {
					return nil, errors.New("net.udp.server bufferSize must be a number")
				}
				bufferSizeNum = int(bufferSizeFloat)
			}
			return &UDPServer{Addr: addr, BufferSize: bufferSizeNum, config: config, logger: CreateLogger(config)}, nil
		},
	})
}

func (us *UDPServer) Id() string {
	return us.config.Id
}

func (us *UDPServer) Type() string {
	return us.config.Type
}

func (us *UDPServer) Start(ctx context.Context) error {
	us.logger.Debug("running")
	router, ok := ctx.Value(route.RouterContextKey).(route.RouteIO)

	if !ok {
		return errors.New("net.udp.server unable to get router from context")
	}
	us.router = router
	moduleContext, cancel := context.WithCancel(ctx)
	us.ctx = moduleContext
	us.cancel = cancel

	listener, err := net.ListenUDP("udp", us.Addr)
	if err != nil {
		return err
	}

	defer listener.Close()

	buffer := make([]byte, us.BufferSize)
	for {
		select {
		case <-us.ctx.Done():
			// TODO(jwetzell): cleanup?
			us.logger.Debug("done")
			return nil
		default:
			listener.SetDeadline(time.Now().Add(time.Millisecond * 200))

			numBytes, _, err := listener.ReadFromUDP(buffer)
			if err != nil {
				//NOTE(jwetzell) we hit deadline
				if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
					continue
				}
				return err
			}
			message := buffer[:numBytes]
			if us.router != nil {
				us.router.HandleInput(us.ctx, us.Id(), message)
			} else {
				us.logger.Error("input received but no router is configured")
			}
		}
	}

}

func (us *UDPServer) Output(ctx context.Context, payload any) error {
	return errors.New("net.udp.server output is not implemented")
}

func (us *UDPServer) Stop() {
	us.cancel()
}
