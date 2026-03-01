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
	cancel        context.CancelFunc
}

func init() {
	RegisterModule(ModuleRegistration{
		Type: "net.tcp.server",
		New: func(moduleConfig config.ModuleConfig) (Module, error) {
			params := moduleConfig.Params
			portNum, err := params.GetInt("port")
			if err != nil {
				return nil, fmt.Errorf("net.tcp.server port error: %w", err)
			}

			framingMethodString, err := params.GetString("framing")
			if err != nil {
				return nil, fmt.Errorf("net.tcp.server framing error: %w", err)
			}

			framer := framer.GetFramer(framingMethodString)

			if framer == nil {
				return nil, fmt.Errorf("net.tcp.server unknown framing method: %s", framingMethodString)
			}

			ipString, err := params.GetString("ip")
			if err != nil {
				if errors.Is(err, config.ErrParamNotFound) {
					ipString = "0.0.0.0"
				} else {
					return nil, fmt.Errorf("net.tcp.server ip error: %w", err)
				}
			}

			addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", ipString, uint16(portNum)))
			if err != nil {
				return nil, err
			}
			return &TCPServer{Framer: framer, Addr: addr, config: moduleConfig, quit: make(chan interface{}), logger: CreateLogger(moduleConfig)}, nil
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
			client.Close()
			ts.connectionsMu.Lock()
			for i := 0; i < len(ts.connections); i++ {
				if ts.connections[i] == client {
					ts.connections = slices.Delete(ts.connections, i, i+1)
					break
				}
			}
			ts.connectionsMu.Unlock()
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
						ts.logger.Debug("connection reset", "remoteAddr", client.RemoteAddr().String())
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
					ts.logger.Debug("stream ended", "remoteAddr", client.RemoteAddr().String())
					ts.connectionsMu.Unlock()
				}
				return
			}
			if ts.Framer != nil {
				if byteCount > 0 {
					messages := ts.Framer.Decode(buffer[0:byteCount])
					for _, message := range messages {
						if ts.router != nil {
							ts.router.HandleInput(ts.ctx, ts.Id(), message)
						} else {
							ts.logger.Error("input received but no router is configured")
						}
					}
				}
			}
		}
	}
}

func (ts *TCPServer) Start(ctx context.Context) error {
	ts.logger.Debug("running")
	router, ok := ctx.Value(route.RouterContextKey).(route.RouteIO)

	if !ok {
		return errors.New("net.tcp.server unable to get router from context")
	}
	ts.router = router
	moduleContext, cancel := context.WithCancel(ctx)
	ts.ctx = moduleContext
	ts.cancel = cancel

	listener, err := net.ListenTCP("tcp", ts.Addr)
	if err != nil {
		return err
	}
	ts.wg.Add(1)

	go func() {
		<-ts.ctx.Done()
		close(ts.quit)
		listener.Close()
		ts.logger.Debug("done")
	}()

AcceptLoop:
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			select {
			case <-ts.quit:
				break AcceptLoop
			default:
				ts.logger.Debug("problem with listener", "error", err)
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

func (ts *TCPServer) Output(ctx context.Context, payload any) error {
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

func (ts *TCPServer) Stop() {
	ts.cancel()
	ts.wg.Wait()
}
