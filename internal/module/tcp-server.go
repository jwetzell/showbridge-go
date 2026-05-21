package module

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"slices"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/framer"
)

type TCPServer struct {
	config        config.ModuleConfig
	Addr          *net.TCPAddr
	Framer        framer.Framer
	ctx           context.Context
	inputHandler  common.InputHandler
	wg            sync.WaitGroup
	connections   []*net.TCPConn
	connectionsMu sync.RWMutex
	logger        *slog.Logger
	cancel        context.CancelFunc
	listener      *net.TCPListener
	listenerMu    sync.Mutex
}

func init() {
	RegisterModule(ModuleRegistration{
		Type:  "net.tcp.server",
		Title: "TCP Server",
		ParamsSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"ip": {
					Title:   "IP",
					Type:    "string",
					Default: json.RawMessage(`"0.0.0.0"`),
				},
				"port": {
					Title:   "Port",
					Type:    "integer",
					Minimum: jsonschema.Ptr[float64](1024),
					Maximum: jsonschema.Ptr[float64](65535),
				},
				"framing": {
					Title: "Framing Method",
					Type:  "string",
					Enum:  []any{"LF", "CR", "CRLF", "SLIP", "RAW"},
				},
			},
			Required:             []string{"port", "framing"},
			AdditionalProperties: &jsonschema.Schema{Not: &jsonschema.Schema{}},
		},
		New: func(moduleConfig config.ModuleConfig) (common.Module, error) {
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
			return &TCPServer{Framer: framer, Addr: addr, config: moduleConfig, logger: CreateLogger(moduleConfig)}, nil
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
	ts.logger.Debug("connection accepted", "remoteAddr", client.RemoteAddr().String())
	defer client.Close()

	buffer := make([]byte, 1024)
ClientRead:
	for ts.ctx.Err() == nil {
		select {
		case <-ts.ctx.Done():
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
					ts.connectionsMu.Unlock()
				}
				return
			}
			if ts.Framer != nil {
				if byteCount > 0 {
					messages := ts.Framer.Decode(buffer[0:byteCount])
					for _, message := range messages {
						if ts.inputHandler != nil {
							ts.inputHandler(ts.ctx, ts.Id(), message)
						} else {
							ts.logger.Error("input received but no input handler is configured")
						}
					}
				}
			}
		}
	}
}

func (ts *TCPServer) Start(ctx context.Context, inputHandler common.InputHandler) error {
	ts.logger.Debug("running")
	ts.inputHandler = inputHandler
	moduleContext, cancel := context.WithCancel(ctx)
	ts.ctx = moduleContext
	ts.cancel = cancel

	listener, err := net.ListenTCP("tcp", ts.Addr)
	if err != nil {
		return err
	}
	ts.listenerMu.Lock()
	ts.listener = listener
	ts.listenerMu.Unlock()
	ts.wg.Add(1)

AcceptLoop:
	for ts.ctx.Err() == nil {
		conn, err := listener.AcceptTCP()
		if err != nil {
			select {
			case <-ts.ctx.Done():
				break AcceptLoop
			default:
				ts.logger.Debug("problem with listener", "error", err)
			}
		} else {
			ts.wg.Go(func() {
				ts.handleClient(conn)
			})
		}
	}
	ts.wg.Done()
	return nil
}

func (ts *TCPServer) Output(ctx context.Context, payload any) error {
	payloadBytes, ok := common.GetAnyAsByteSlice(payload)

	if !ok {
		return errors.New("net.tcp.server is only able to output bytes")
	}
	ts.connectionsMu.Lock()
	var errorString strings.Builder

	for _, connection := range ts.connections {
		_, err := connection.Write(payloadBytes)
		if err != nil {
			fmt.Fprintf(&errorString, "%s\n", err.Error())
		}
	}
	ts.connectionsMu.Unlock()

	if errorString.String() == "" {
		return nil
	}
	return fmt.Errorf("net.tcp.server error during output: %s", errorString.String())
}

func (ts *TCPServer) Stop() {
	if ts.cancel != nil {
		ts.cancel()
	}
	ts.listenerMu.Lock()
	defer ts.listenerMu.Unlock()
	if ts.listener != nil {
		ts.listener.Close()
		ts.listener = nil
	}
	ts.wg.Wait()
	ts.logger.Debug("done")
}
