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
	Host   string
	Port   uint16
	framer framing.Framer
	conn   net.Conn
	router *Router
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
				return nil, fmt.Errorf("net.tcp.client host must be uint16")
			}

			port, ok := params["port"]
			if !ok {
				return nil, fmt.Errorf("net.tcp.client requires a port parameter")
			}

			portNum, ok := port.(float64)

			if !ok {
				return nil, fmt.Errorf("net.tcp.client port must be uint16")
			}

			framingMethod, ok := params["framing"]
			if !ok {
				return nil, fmt.Errorf("net.tcp.client requires a framing method")
			}

			framingMethodString, ok := framingMethod.(string)

			if !ok {
				return nil, fmt.Errorf("tcp framing method must be a string")
			}

			var framer framing.Framer

			switch framingMethodString {
			case "CR":
				framer = framing.NewByteSeparatorFramer([]byte{'\r'})
			case "LF":
				framer = framing.NewByteSeparatorFramer([]byte{'\n'})
			case "CRLF":
				framer = framing.NewByteSeparatorFramer([]byte{'\r', '\n'})
			case "SLIP":
				framer = framing.NewSlipFramer()
			default:
				return nil, fmt.Errorf("unknown framing method: %s", framingMethodString)
			}

			return &TCPClient{framer: framer, Host: hostString, Port: uint16(portNum), config: config}, nil
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
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", tc.Host, tc.Port))
	if err != nil {
		return err
	}

	// TODO(jwetzell): shutdown with router.Context properly
	go func() {
		<-tc.router.Context.Done()
		slog.Debug("router context done in module", "id", tc.config.Id)
		if tc.conn != nil {
			tc.conn.Close()
		}
	}()

	for {
		client, err := net.DialTCP("tcp", nil, addr)
		if err != nil {
			if tc.router.Context.Err() != nil {
				slog.Debug("router context done in module", "id", tc.config.Id)
				return nil
			}
			slog.Error(err.Error())
			time.Sleep(time.Second * 2)
			continue
		}

		tc.conn = client

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
					byteCount, err := client.Read(buffer)

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

func (tc *TCPClient) Output(payload any) error {
	if tc.conn != nil {
		payloadBytes, ok := payload.([]byte)
		if !ok {
			return fmt.Errorf("net.tcp.client is only able to output bytes")
		}
		_, err := tc.conn.Write(tc.framer.Encode(payloadBytes))
		return err
	}
	return nil
}
