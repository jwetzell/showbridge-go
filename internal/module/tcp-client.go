package module

import (
	"context"
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
}

func init() {
	RegisterModule(ModuleRegistration{
		Type: "net.tcp.client",
		New: func(ctx context.Context, config config.ModuleConfig, router route.RouteIO) (Module, error) {
			params := config.Params
			host, ok := params["host"]

			if !ok {
				return nil, fmt.Errorf("net.tcp.client requires a host parameter")
			}

			hostString, ok := host.(string)

			if !ok {
				return nil, fmt.Errorf("net.tcp.client host must be string")
			}

			port, ok := params["port"]
			if !ok {
				return nil, fmt.Errorf("net.tcp.client requires a port parameter")
			}

			portNum, ok := port.(float64)

			if !ok {
				return nil, fmt.Errorf("net.tcp.client port must be a number")
			}

			addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", hostString, uint16(portNum)))
			if err != nil {
				return nil, err
			}

			framingMethod := "RAW"

			framingMethodRaw, ok := params["framing"]

			if ok {
				framingMethodString, ok := framingMethodRaw.(string)

				if !ok {
					return nil, fmt.Errorf("misc.serial.client framing method must be a string")
				}
				framingMethod = framingMethodString
			}

			framer, err := framer.GetFramer(framingMethod)

			if err != nil {
				return nil, err
			}

			return &TCPClient{framer: framer, Addr: addr, config: config, ctx: ctx, router: router}, nil
		},
	})
}

func (tc *TCPClient) Id() string {
	return tc.config.Id
}

func (tc *TCPClient) Type() string {
	return tc.config.Type
}

func (tc *TCPClient) Run() error {

	// TODO(jwetzell): shutdown with router.Context properly
	go func() {
		<-tc.ctx.Done()
		slog.Debug("router context done in module", "id", tc.Id())
		if tc.conn != nil {
			tc.conn.Close()
		}
	}()

	for {
		err := tc.SetupConn()
		if err != nil {
			if tc.ctx.Err() != nil {
				slog.Debug("router context done in module", "id", tc.Id())
				return nil
			}
			slog.Error("net.tcp.client", "id", tc.Id(), "error", err.Error())
			time.Sleep(time.Second * 2)
			continue
		}

		buffer := make([]byte, 1024)
		select {
		case <-tc.ctx.Done():
			slog.Debug("router context done in module", "id", tc.Id())
			return nil
		default:
		READ:
			for {
				select {
				case <-tc.ctx.Done():
					slog.Debug("router context done in module", "id", tc.Id())
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
									tc.router.HandleInput(tc.Id(), message)
								} else {
									slog.Error("net.tcp.client has no router", "id", tc.Id())
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

func (tc *TCPClient) Output(payload any) error {
	// NOTE(jwetzell): not sure how this would occur but
	if tc.conn == nil {
		err := tc.SetupConn()
		if err != nil {
			return err
		}
	}
	payloadBytes, ok := payload.([]byte)
	if !ok {
		return fmt.Errorf("net.tcp.client is only able to output bytes")
	}
	_, err := tc.conn.Write(tc.framer.Encode(payloadBytes))
	return err
}
