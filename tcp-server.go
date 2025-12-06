package showbridge

import (
	"errors"
	"fmt"
	"log/slog"
	"net"
	"slices"
	"sync"
	"syscall"
	"time"

	"github.com/jwetzell/showbridge-go/internal/framing"
)

type TCPServer struct {
	config        ModuleConfig
	Addr          *net.TCPAddr
	FramingMethod string
	router        *Router
	quit          chan interface{}
	wg            sync.WaitGroup
	connections   []*net.TCPConn
	connectionsMu sync.RWMutex
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
				return nil, fmt.Errorf("net.tcp.server port must be a number")
			}

			framingMethod, ok := params["framing"]
			if !ok {
				return nil, fmt.Errorf("net.tcp.server requires a framing method")
			}

			framingMethodString, ok := framingMethod.(string)

			if !ok {
				return nil, fmt.Errorf("net.tcp.server framing method must be a string")
			}

			ipString := "0.0.0.0"

			ip, ok := params["ip"]
			if ok {

				specificIpString, ok := ip.(string)

				if !ok {
					return nil, fmt.Errorf("net.tcp.server ip must be a string")
				}
				ipString = specificIpString
			}

			addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", ipString, uint16(portNum)))
			if err != nil {
				return nil, err
			}

			return &TCPServer{FramingMethod: framingMethodString, Addr: addr, config: config, quit: make(chan interface{})}, nil
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

func (ts *TCPServer) handleClient(client *net.TCPConn) {
	ts.connectionsMu.Lock()
	ts.connections = append(ts.connections, client)
	ts.connectionsMu.Unlock()
	slog.Debug("net.tcp.server connection accepted", "id", ts.config.Id, "remoteAddr", client.RemoteAddr().String())
	defer client.Close()
	var framer framing.Framer

	switch ts.FramingMethod {
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
				if opErr, ok := err.(*net.OpError); ok {
					//NOTE(jwetzell) we hit deadline
					if opErr.Timeout() {
						continue ClientRead
					}
					if errors.Is(opErr, syscall.ECONNRESET) {
						ts.connectionsMu.Lock()
						for i := 0; i < len(ts.connections); i++ {
							if ts.connections[i] == client {
								ts.connections = slices.Delete(ts.connections, i, i+1)
								break
							}
						}
						slog.Debug("net.tcp.server connection reset", "id", ts.config.Id, "remoteAddr", client.RemoteAddr().String())
						ts.connectionsMu.Unlock()
					}
				}

				if err.Error() == "EOF" {
					ts.connectionsMu.Lock()
					for i := 0; i < len(ts.connections); i++ {
						if ts.connections[i] == client {
							ts.connections = slices.Delete(ts.connections, i, i+1)
							break
						}
					}
					slog.Debug("net.tcp.server stream ended", "id", ts.config.Id, "remoteAddr", client.RemoteAddr().String())
					ts.connectionsMu.Unlock()
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
							slog.Error("net.tcp.server has no router", "id", ts.config.Id)
						}
					}
				}
			}
		}
	}
}

func (ts *TCPServer) Run() error {
	listener, err := net.ListenTCP("tcp", ts.Addr)
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
		conn, err := listener.AcceptTCP()
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
	payloadBytes, ok := payload.([]byte)

	if !ok {
		return fmt.Errorf("net.tcp.server is only able to output bytes")
	}
	ts.connectionsMu.Lock()
	errorString := ""

	for _, connection := range ts.connections {
		_, err := connection.Write(payloadBytes)
		if err != nil {
			errorString += fmt.Sprintf("%s\n", err.Error())
		}
	}
	ts.connectionsMu.Unlock()

	if errorString == "" {
		return nil
	}
	return fmt.Errorf("%s", errorString)
}
