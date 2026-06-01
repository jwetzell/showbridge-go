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
	"time"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/framer"
)

func init() {
	RegisterModule(ModuleRegistration{
		Type:  "net.tcp.server",
		Title: "TCP Server",
		ParamsSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"ip": {
					Title:       "IP",
					Description: "the IP address to bind the TCP server to",
					Type:        "string",
					Default:     json.RawMessage(`"0.0.0.0"`),
				},
				"port": {
					Title:       "Port",
					Description: "the port for the TCP server to listen on",
					Type:        "integer",
					Minimum:     jsonschema.Ptr[float64](1024),
					Maximum:     jsonschema.Ptr[float64](65535),
				},
				"framing": {
					Title:       "Framing Method",
					Description: "the method used to frame messages over the TCP connection",
					Type:        "string",
					Enum:        []any{"LF", "CR", "CRLF", "SLIP", "RAW"},
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

			inFramer := framer.GetFramer(framingMethodString)

			if inFramer == nil {
				return nil, fmt.Errorf("net.tcp.server unknown framing method: %s", framingMethodString)
			}

			outFramer := framer.GetFramer(framingMethodString)

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
			return &TCPServer{InFramer: inFramer, OutFramer: outFramer, framerType: framingMethodString, Addr: addr, config: moduleConfig, logger: CreateLogger(moduleConfig)}, nil
		},
	})
}

type tcpConnection struct {
	conn   *net.TCPConn
	framer framer.Framer
}

type TCPServer struct {
	config                config.ModuleConfig
	Addr                  *net.TCPAddr
	InFramer              framer.Framer
	OutFramer             framer.Framer
	framerType            string
	ctx                   context.Context
	inputHandler          common.InputHandler
	wg                    sync.WaitGroup
	connections           []tcpConnection
	connectionsMu         sync.RWMutex
	logger                *slog.Logger
	cancel                context.CancelFunc
	listener              *net.TCPListener
	listenerMu            sync.Mutex
	connectionShutdownCtx context.Context
	connectionShutdown    context.CancelFunc
}

func (ts *TCPServer) Id() string {
	return ts.config.Id
}

func (ts *TCPServer) Type() string {
	return ts.config.Type
}

func (ts *TCPServer) handleClient(client *net.TCPConn) {
	ts.connectionsMu.Lock()
	ts.connections = append(ts.connections, tcpConnection{conn: client, framer: framer.GetFramer(ts.framerType)})
	ts.connectionsMu.Unlock()
	ts.logger.Debug("connection accepted", "remoteAddr", client.RemoteAddr().String())
	defer func() {
		client.Close()
		ts.connectionsMu.Lock()
		for i := 0; i < len(ts.connections); i++ {
			if ts.connections[i].conn == client {
				ts.connections = slices.Delete(ts.connections, i, i+1)
				break
			}
		}
		ts.connectionsMu.Unlock()
		ts.logger.Debug("connection closed", "remoteAddr", client.RemoteAddr().String())
	}()

	buffer := make([]byte, 1024)
	for ts.ctx.Err() == nil && ts.connectionShutdownCtx.Err() == nil {
		client.SetDeadline(time.Now().Add(time.Millisecond * 200))
		byteCount, err := client.Read(buffer)
		if err != nil {
			if opErr, ok := err.(*net.OpError); ok {
				//NOTE(jwetzell) we hit deadline
				if opErr.Timeout() {
					continue
				}
			}
			break
		}
		if ts.InFramer != nil {
			if byteCount > 0 {
				messages := ts.InFramer.Decode(buffer[0:byteCount])
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

func (ts *TCPServer) Start(ctx context.Context, inputHandler common.InputHandler) error {
	ts.logger.Debug("running")
	ts.inputHandler = inputHandler
	moduleContext, cancel := context.WithCancel(ctx)
	ts.ctx = moduleContext
	ts.cancel = cancel

	shutdownCtx, shutdownCancel := context.WithCancel(context.Background())
	ts.connectionShutdownCtx = shutdownCtx
	ts.connectionShutdown = shutdownCancel

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
			if errors.Is(err, net.ErrClosed) {
				break AcceptLoop
			}
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
	<-ts.ctx.Done()
	ts.logger.Debug("done")
	return nil
}

func (ts *TCPServer) Output(ctx context.Context, payload any) error {
	payloadBytes, ok := common.GetAnyAsByteSlice(payload)

	if !ok {
		return errors.New("net.tcp.server is only able to output bytes")
	}
	ts.connectionsMu.Lock()
	defer ts.connectionsMu.Unlock()
	var errorString strings.Builder

	if ts.OutFramer == nil {
		return errors.New("no output framer configured")
	}

	outputBytes := ts.OutFramer.Encode(payloadBytes)

	for _, connection := range ts.connections {
		_, err := connection.conn.Write(outputBytes)
		if err != nil {
			fmt.Fprintf(&errorString, "%s\n", err.Error())
		}
	}

	if errorString.String() == "" {
		return nil
	}
	return fmt.Errorf("net.tcp.server error during output: %s", errorString.String())
}

func (ts *TCPServer) Stop() {
	if ts.cancel != nil {
		defer ts.cancel()
	}
	if ts.connectionShutdown != nil {
		ts.connectionShutdown()
	}

	ts.listenerMu.Lock()
	defer ts.listenerMu.Unlock()
	if ts.listener != nil {
		ts.listener.Close()
	}
	ts.logger.Debug("waiting for connections to close")
	ts.wg.Wait()
	ts.logger.Debug("all connections closed")
}
