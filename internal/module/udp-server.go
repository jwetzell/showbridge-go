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
		New: func(moduleConfig config.ModuleConfig) (Module, error) {
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
