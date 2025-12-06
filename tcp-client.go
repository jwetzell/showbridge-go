package showbridge

import (
	"fmt"
	"log/slog"
	"net"
	"time"

	"github.com/jwetzell/showbridge-go/internal/framing"
)

type TCPClient struct {
	config ModuleConfig
	framer framing.Framer
	conn   *net.TCPConn
	router *Router
	Addr   *net.TCPAddr
}

func init() {
	RegisterModule(ModuleRegistration{
		Type: "net.tcp.client",
		New: func(config ModuleConfig) (Module, error) {
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

			framingMethod, ok := params["framing"]
			if !ok {
				return nil, fmt.Errorf("net.tcp.client requires a framing method")
			}

			framingMethodString, ok := framingMethod.(string)

			if !ok {
				return nil, fmt.Errorf("net.tcp.client framing method must be a string")
			}

			framer, err := framing.GetFramer(framingMethodString)

			if err != nil {
				return nil, err
			}

			return &TCPClient{framer: framer, Addr: addr, config: config}, nil
		},
	})
}

func (tc *TCPClient) Id() string {
	return tc.config.Id
}

func (tc *TCPClient) Type() string {
	return tc.config.Type
}

func (tc *TCPClient) RegisterRouter(router *Router) {
	tc.router = router
}

func (tc *TCPClient) Run() error {

	// TODO(jwetzell): shutdown with router.Context properly
	go func() {
		<-tc.router.Context.Done()
		slog.Debug("router context done in module", "id", tc.config.Id)
		if tc.conn != nil {
			tc.conn.Close()
		}
	}()

	for {
		err := tc.SetupConn()
		if err != nil {
			if tc.router.Context.Err() != nil {
				slog.Debug("router context done in module", "id", tc.config.Id)
				return nil
			}
			slog.Error("net.tcp.client", "id", tc.config.Id, "error", err.Error())
			time.Sleep(time.Second * 2)
			continue
		}

		buffer := make([]byte, 1024)
		select {
		case <-tc.router.Context.Done():
			slog.Debug("router context done in module", "id", tc.config.Id)
			return nil
		default:
		READ:
			for {
				select {
				case <-tc.router.Context.Done():
					slog.Debug("router context done in module", "id", tc.config.Id)
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
									tc.router.HandleInput(tc.config.Id, message)
								} else {
									slog.Error("net.tcp.client has no router", "id", tc.config.Id)
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
