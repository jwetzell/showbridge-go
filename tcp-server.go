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
	Ip            string
	Port          uint16
	framingMethod string
	router        *Router
	quit          chan interface{}
	wg            sync.WaitGroup
	connections   []net.Conn
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

			return &TCPServer{framingMethod: framingMethodString, Port: uint16(portNum), Ip: ipString, config: config, quit: make(chan interface{})}, nil
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
	ts.connectionsMu.Lock()
	ts.connections = append(ts.connections, client)
	ts.connectionsMu.Unlock()
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
				if opErr, ok := err.(*net.OpError); ok {
					if opErr.Timeout() {
						continue ClientRead
					}
					fmt.Println(opErr.Err)
					if errors.Is(opErr, syscall.ECONNRESET) {
						ts.connectionsMu.Lock()
						for i := 0; i < len(ts.connections); i++ {
							if ts.connections[i] == client {
								ts.connections = slices.Delete(ts.connections, i, i+1)
								break
							}
						}
						slog.Debug("connection closed", "id", ts.config.Id, "remoteAddr", client.RemoteAddr().String())
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
					slog.Debug("connection closed", "id", ts.config.Id, "remoteAddr", client.RemoteAddr().String())
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
	// TODO(jwetzell): switch to net.ListenTCP and move addr resolution to init
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", ts.Ip, ts.Port))
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
