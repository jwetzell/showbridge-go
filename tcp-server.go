package showbridge

import (
	"fmt"
	"log/slog"
	"net"
	"sync"
	"time"

	"github.com/jwetzell/showbridge-go/internal/framing"
)

type TCPServer struct {
	config        ModuleConfig
	Port          uint16
	framingMethod string
	router        *Router
	quit          chan interface{}
	wg            sync.WaitGroup
}

func init() {
	RegisterModule(ModuleRegistration{
		Type: "net.tcp.server",
		New: func(config ModuleConfig) (Module, error) {
			params := config.Params
			port, ok := params["port"]
			if !ok {
				return nil, fmt.Errorf("net.tcp.server requires a port parameter")
			}

			portNum, ok := port.(float64)

			if !ok {
				return nil, fmt.Errorf("net.tcp.server port must be uint16")
			}

			framingMethod, ok := params["framing"]
			if !ok {
				return nil, fmt.Errorf("net.tcp.server requires a framing method")
			}

			framingMethodString, ok := framingMethod.(string)

			if !ok {
				return nil, fmt.Errorf("tcp framing method must be a string")
			}

			return &TCPServer{framingMethod: framingMethodString, Port: uint16(portNum), config: config, quit: make(chan interface{})}, nil
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
	ts.router = router
}

func (ts *TCPServer) handleClient(client net.Conn) {
	slog.Debug("connection accepted", "id", ts.config.Id, "remoteAddr", client.RemoteAddr().String())
	defer client.Close()
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
ClientRead:
	for {
		select {
		case <-ts.quit:
			return
		default:
			client.SetDeadline(time.Now().Add(time.Millisecond * 200))
			byteCount, err := client.Read(buffer)

			if err != nil {
				//NOTE(jwetzell) we hit deadline
				if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
					continue ClientRead
				}
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
							slog.Error("tcp-server has no router", "id", ts.config.Id)
						}
					}
				}
			}
		}
	}
}

func (ts *TCPServer) Run() error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", ts.Port))
	if err != nil {
		return err
	}
	ts.wg.Add(1)

	go func() {
		<-ts.router.Context.Done()
		close(ts.quit)
		listener.Close()
		slog.Debug("router context done in module", "id", ts.config.Id)
	}()

AcceptLoop:
	for {
		conn, err := listener.Accept()
		if err != nil {
			select {
			case <-ts.quit:
				break AcceptLoop
			default:
				slog.Debug("net.tcp.server problem with listener", "error", err)
			}
		} else {
			ts.wg.Add(1)
			go func() {
				ts.handleClient(conn)
				ts.wg.Done()
			}()
		}
	}
	ts.wg.Done()
	ts.wg.Wait()
	return nil
}

func (ts *TCPServer) Output(payload any) error {
	return fmt.Errorf("net.tcp.server output is not implemented")
}
