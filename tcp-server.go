package showbridge

import (
	"context"
	"fmt"
	"log/slog"
	"net"
)

type TCPServer struct {
	Port uint16
}

func init() {
	RegisterModule(ModuleRegistration{
		Type: "net.tcp.server",
		New: func(params map[string]any) (Module, error) {
			port, ok := params["port"]
			if !ok {
				return nil, fmt.Errorf("tcp server requires a port parameter")
			}

			portNum, ok := port.(float64)

			if !ok {
				return nil, fmt.Errorf("tcp server port must be uint16")
			}

			return TCPServer{Port: uint16(portNum)}, nil
		},
	})
}

func (ts TCPServer) HandleClient(ctx context.Context, client net.Conn) {
	slog.Info("handling connection", "remoteAddr", client.RemoteAddr().String())

	buffer := make([]byte, 1024)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			byteCount, err := client.Read(buffer)

			if err != nil {
				if err.Error() == "EOF" {
					slog.Debug("connection closed")
				}
				return
			}

			if byteCount > 0 {
				slog.Info(string(buffer[0:byteCount]))
			}
		}

	}
}

func (ts TCPServer) Run(ctx context.Context) error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", ts.Port))
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			client, err := listener.Accept()
			if err != nil {
				return err
			}
			go ts.HandleClient(ctx, client)
		}
	}

}
