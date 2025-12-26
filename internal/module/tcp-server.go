package module

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"slices"
	"sync"
	"syscall"
	"time"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/framer"
	"github.com/jwetzell/showbridge-go/internal/route"
)

type TCPServer struct {
	config        config.ModuleConfig
	Addr          *net.TCPAddr
	Framer        framer.Framer
	ctx           context.Context
	router        route.RouteIO
	quit          chan interface{}
	wg            sync.WaitGroup
	connections   []*net.TCPConn
	connectionsMu sync.RWMutex
	logger        *slog.Logger
}

func init() {
	RegisterModule(ModuleRegistration{
		Type: "net.tcp.server",
		New: func(ctx context.Context, config config.ModuleConfig, router route.RouteIO) (Module, error) {
			params := config.Params
			port, ok := params["port"]
			if !ok {
				return nil, errors.New("net.tcp.server requires a port parameter")
			}

			portNum, ok := port.(float64)

			if !ok {
				return nil, errors.New("net.tcp.server port must be a number")
			}

			framingMethod := "RAW"

			framingMethodRaw, ok := params["framing"]

			if ok {
				framingMethodString, ok := framingMethodRaw.(string)

				if !ok {
					return nil, errors.New("net.tcp.server framing method must be a string")
				}
				framingMethod = framingMethodString
			}

			framer := framer.GetFramer(framingMethod)

			if framer == nil {
				return nil, fmt.Errorf("net.tcp.server unknown framing method: %s", framingMethod)
			}

			ipString := "0.0.0.0"

			ip, ok := params["ip"]
			if ok {

				specificIpString, ok := ip.(string)

				if !ok {
					return nil, errors.New("net.tcp.server ip must be a string")
				}
				ipString = specificIpString
			}

			addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", ipString, uint16(portNum)))
			if err != nil {
				return nil, err
			}

			return &TCPServer{Framer: framer, Addr: addr, config: config, quit: make(chan interface{}), ctx: ctx, router: router, logger: CreateLogger(config)}, nil
		},
	})
}

func (ts *TCPServer) Id() string {
	return ts.config.Id
}

func (ts *TCPServer) Type() string {
	return ts.config.Type
}

func (ts *TCPServer) handleClient(client *net.TCPConn) {
	ts.connectionsMu.Lock()
	ts.connections = append(ts.connections, client)
	ts.connectionsMu.Unlock()
	ts.logger.Debug("net.tcp.server connection accepted", "remoteAddr", client.RemoteAddr().String())
	defer client.Close()

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
						ts.logger.Debug("net.tcp.server connection reset", "remoteAddr", client.RemoteAddr().String())
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
					ts.logger.Debug("net.tcp.server stream ended", "remoteAddr", client.RemoteAddr().String())
					ts.connectionsMu.Unlock()
				}
				return
			}
			if ts.Framer != nil {
				if byteCount > 0 {
					messages := ts.Framer.Decode(buffer[0:byteCount])
					for _, message := range messages {
						if ts.router != nil {
							ts.router.HandleInput(ts.Id(), message)
						} else {
							ts.logger.Error("net.tcp.server has no router")
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
		<-ts.ctx.Done()
		close(ts.quit)
		listener.Close()
		ts.logger.Debug("router context done in module")
	}()

AcceptLoop:
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			select {
			case <-ts.quit:
				break AcceptLoop
			default:
				ts.logger.Debug("net.tcp.server problem with listener", "error", err)
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
		return errors.New("net.tcp.server is only able to output bytes")
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
	return fmt.Errorf("net.tcp.server error during output: %s", errorString)
}
