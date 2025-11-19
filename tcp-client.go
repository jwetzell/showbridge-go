package showbridge

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"time"

	"github.com/jwetzell/showbridge-go/internals/framing"
)

type TCPClient struct {
	config ModuleConfig
	Host   string
	Port   uint16
	framer framing.Framer
}

func init() {
	RegisterModule(ModuleRegistration{
		Type: "net.tcp.client",
		New: func(config ModuleConfig) (Module, error) {
			params := config.Params
			host, ok := params["host"]

			if !ok {
				return nil, fmt.Errorf("tcp client requires a host parameter")
			}

			hostString, ok := host.(string)

			if !ok {
				return nil, fmt.Errorf("tcp client host must be uint16")
			}

			port, ok := params["port"]
			if !ok {
				return nil, fmt.Errorf("tcp client requires a port parameter")
			}

			portNum, ok := port.(float64)

			if !ok {
				return nil, fmt.Errorf("tcp client port must be uint16")
			}

			framingMethod, ok := params["framing"]
			if !ok {
				return nil, fmt.Errorf("tcp client requires a framing method")
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
			default:
				return nil, fmt.Errorf("unknown framing method: %s", framingMethodString)
			}

			return TCPClient{framer: framer, Host: hostString, Port: uint16(portNum), config: config}, nil
		},
	})
}

func (tc TCPClient) Id() string {
	return tc.config.Id
}

func (tc TCPClient) Type() string {
	return tc.config.Type
}

func (tc TCPClient) Run(ctx context.Context) error {
	for {
		client, err := net.Dial("tcp", fmt.Sprintf(":%d", tc.Port))
		if err != nil {
			slog.Error(err.Error())
			time.Sleep(time.Second * 2)
			continue
		}

		buffer := make([]byte, 1024)
		select {
		case <-ctx.Done():
			return nil
		default:
		READ:
			for {
				select {
				case <-ctx.Done():
					return nil
				default:
					byteCount, err := client.Read(buffer)

					if err != nil {
						slog.Debug("connection closed")
						tc.framer.Clear()
						break READ
					}

					if tc.framer != nil {
						if byteCount > 0 {
							messages := tc.framer.Frame(buffer[0:byteCount])
							for _, message := range messages {
								slog.Debug("tcp-client message", "bytes", message)
							}
						}
					}
				}

			}
		}

	}

}
