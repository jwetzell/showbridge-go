package showbridge

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"time"
)

type TCPClient struct {
	config ModuleConfig
	Host   string
	Port   uint16
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

			return TCPClient{Host: hostString, Port: uint16(portNum), config: config}, nil
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
						break READ
					}

					if byteCount > 0 {
						slog.Info(string(buffer[0:byteCount]))
					}
				}

			}
		}

	}

}
