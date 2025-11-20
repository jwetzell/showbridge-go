package showbridge

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	"github.com/jwetzell/showbridge-go/internals/framing"
)

type TCPServer struct {
	config        ModuleConfig
	Port          uint16
	framingMethod string
	router        *Router
}

func init() {
	RegisterModule(ModuleRegistration{
		Type: "net.tcp.server",
		New: func(config ModuleConfig) (Module, error) {
			params := config.Params
			port, ok := params["port"]
			if !ok {
				return nil, fmt.Errorf("tcp server requires a port parameter")
			}

			portNum, ok := port.(float64)

			if !ok {
				return nil, fmt.Errorf("tcp server port must be uint16")
			}

			framingMethod, ok := params["framing"]
			if !ok {
				return nil, fmt.Errorf("tcp server requires a framing method")
			}

			framingMethodString, ok := framingMethod.(string)

			if !ok {
				return nil, fmt.Errorf("tcp framing method must be a string")
			}

			return &TCPServer{framingMethod: framingMethodString, Port: uint16(portNum), config: config}, nil
		},
	})
}

func (ts *TCPServer) Id() string {
	return ts.config.Id
}

func (ts *TCPServer) Type() string {
	return ts.config.Type
}

func (ts *TCPServer) RegisterRouter(router *Router) {
	slog.Debug("registering router", "id", ts.config.Id)
	ts.router = router
}

func (ts *TCPServer) HandleClient(ctx context.Context, client net.Conn) {
	slog.Debug("connection accepted", "id", ts.config.Id, "remoteAddr", client.RemoteAddr().String())

	var framer framing.Framer

	switch ts.framingMethod {
	case "LF":
		framer = framing.NewByteSeparatorFramer([]byte{'\n'})
	case "CR":
		framer = framing.NewByteSeparatorFramer([]byte{'\r'})
	case "CRLF":
		framer = framing.NewByteSeparatorFramer([]byte{'\r', '\n'})
	case "SLIP":
		framer = framing.NewSlipFramer()
	}

	buffer := make([]byte, 1024)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			byteCount, err := client.Read(buffer)

			if err != nil {
				if err.Error() == "EOF" {
					slog.Debug("connection closed", "id", ts.config.Id, "remoteAddr", client.RemoteAddr().String())
				}
				return
			}
			if framer != nil {
				if byteCount > 0 {
					messages := framer.Decode(buffer[0:byteCount])
					for _, message := range messages {
						if ts.router != nil {
							ts.router.HandleInput(ts.config.Id, message)
						} else {
							slog.Error("tcp-server has not router", "id", ts.config.Id)
						}
					}
				}
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

func (ts *TCPServer) Output(payload any) error {
	return fmt.Errorf("tcp-server output is not implemented")
}
