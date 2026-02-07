package module

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"time"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/framer"
	"github.com/jwetzell/showbridge-go/internal/route"
)

type TCPClient struct {
	config config.ModuleConfig
	framer framer.Framer
	conn   *net.TCPConn
	ctx    context.Context
	router route.RouteIO
	Addr   *net.TCPAddr
	logger *slog.Logger
	cancel context.CancelFunc
}

func init() {
	RegisterModule(ModuleRegistration{
		Type: "net.tcp.client",
		New: func(config config.ModuleConfig) (Module, error) {
			params := config.Params
			host, ok := params["host"]

			if !ok {
				return nil, errors.New("net.tcp.client requires a host parameter")
			}

			hostString, ok := host.(string)

			if !ok {
				return nil, errors.New("net.tcp.client host must be string")
			}

			port, ok := params["port"]
			if !ok {
				return nil, errors.New("net.tcp.client requires a port parameter")
			}

			portNum, ok := port.(float64)

			if !ok {
				return nil, errors.New("net.tcp.client port must be a number")
			}

			addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", hostString, uint16(portNum)))
			if err != nil {
				return nil, err
			}

			framingMethod, ok := params["framing"]

			if !ok {
				return nil, errors.New("net.tcp.client requires a framing parameter")
			}

			framingMethodString, ok := framingMethod.(string)

			if !ok {
				return nil, errors.New("net.tcp.client framing method must be a string")
			}

			framer := framer.GetFramer(framingMethodString)

			if framer == nil {
				return nil, fmt.Errorf("net.tcp.client unknown framing method: %s", framingMethod)
			}
			return &TCPClient{framer: framer, Addr: addr, config: config, logger: CreateLogger(config)}, nil
		},
	})
}

func (tc *TCPClient) Id() string {
	return tc.config.Id
}

func (tc *TCPClient) Type() string {
	return tc.config.Type
}

func (tc *TCPClient) Run(ctx context.Context) error {

	router, ok := ctx.Value(route.RouterContextKey).(route.RouteIO)

	if !ok {
		return errors.New("net.tcp.client unable to get router from context")
	}
	tc.router = router
	moduleContext, cancel := context.WithCancel(ctx)
	tc.ctx = moduleContext
	tc.cancel = cancel

	// TODO(jwetzell): shutdown with router.Context properly
	go func() {
		<-tc.ctx.Done()
		tc.logger.Debug("done")
		if tc.conn != nil {
			tc.conn.Close()
		}
	}()

	for {
		err := tc.SetupConn()
		if err != nil {
			if tc.ctx.Err() != nil {
				tc.logger.Debug("done")
				return nil
			}
			tc.logger.Error("connection error", "error", err.Error())
			time.Sleep(time.Second * 2)
			continue
		}

		buffer := make([]byte, 1024)
		select {
		case <-tc.ctx.Done():
			tc.logger.Debug("done")
			return nil
		default:
		READ:
			for {
				select {
				case <-tc.ctx.Done():
					tc.logger.Debug("done")
					return nil
				default:
					byteCount, err := tc.conn.Read(buffer)

					if err != nil {
						tc.framer.Clear()
						break READ
					}

					if tc.framer != nil {
						if byteCount > 0 {
							messages := tc.framer.Decode(buffer[0:byteCount])
							for _, message := range messages {
								if tc.router != nil {
									tc.router.HandleInput(tc.ctx, tc.Id(), message)
								} else {
									tc.logger.Error("input received but no router is configured")
								}
							}
						}
					}
				}
			}
		}
	}
}

func (tc *TCPClient) SetupConn() error {
	client, err := net.DialTCP("tcp", nil, tc.Addr)
	tc.conn = client
	return err
}

func (tc *TCPClient) Output(ctx context.Context, payload any) error {
	// NOTE(jwetzell): not sure how this would occur but
	if tc.conn == nil {
		err := tc.SetupConn()
		if err != nil {
			return err
		}
	}
	payloadBytes, ok := payload.([]byte)
	if !ok {
		return errors.New("net.tcp.client is only able to output bytes")
	}
	_, err := tc.conn.Write(tc.framer.Encode(payloadBytes))
	return err
}

func (tc *TCPClient) Stop() {
	tc.cancel()
}
