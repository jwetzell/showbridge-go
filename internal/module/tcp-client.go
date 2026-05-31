package module

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"sync"
	"time"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/framer"
)

func init() {
	RegisterModule(ModuleRegistration{
		Type:  "net.tcp.client",
		Title: "TCP Client",
		ParamsSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"host": {
					Title: "Host",
					Type:  "string",
				},
				"port": {
					Title:   "Port",
					Type:    "integer",
					Minimum: jsonschema.Ptr[float64](1),
					Maximum: jsonschema.Ptr[float64](65535),
				},
				"framing": {
					Title: "Framing Method",
					Type:  "string",
					Enum:  []any{"LF", "CR", "CRLF", "SLIP", "RAW"},
				},
			},
			Required:             []string{"host", "port", "framing"},
			AdditionalProperties: &jsonschema.Schema{Not: &jsonschema.Schema{}},
		},
		New: func(config config.ModuleConfig) (common.Module, error) {
			params := config.Params
			hostString, err := params.GetString("host")
			if err != nil {
				return nil, fmt.Errorf("net.tcp.client host error: %w", err)
			}

			portNum, err := params.GetInt("port")
			if err != nil {
				return nil, fmt.Errorf("net.tcp.client port error: %w", err)
			}

			addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", hostString, uint16(portNum)))
			if err != nil {
				return nil, err
			}

			framingMethodString, err := params.GetString("framing")
			if err != nil {
				return nil, fmt.Errorf("net.tcp.client framing error: %w", err)
			}

			framer := framer.GetFramer(framingMethodString)

			if framer == nil {
				return nil, fmt.Errorf("net.tcp.client unknown framing method: %s", framingMethodString)
			}
			return &TCPClient{framer: framer, Addr: addr, config: config, logger: CreateLogger(config)}, nil
		},
	})
}

type TCPClient struct {
	config       config.ModuleConfig
	framer       framer.Framer
	conn         *net.TCPConn
	ctx          context.Context
	inputHandler common.InputHandler
	Addr         *net.TCPAddr
	logger       *slog.Logger
	cancel       context.CancelFunc
	connMu       sync.Mutex
}

func (tc *TCPClient) Id() string {
	return tc.config.Id
}

func (tc *TCPClient) Type() string {
	return tc.config.Type
}

func (tc *TCPClient) Start(ctx context.Context, inputHandler common.InputHandler) error {
	tc.logger.Debug("running")
	tc.inputHandler = inputHandler
	moduleContext, cancel := context.WithCancel(ctx)
	tc.ctx = moduleContext
	tc.cancel = cancel

CONNECT_RETRY:
	for tc.ctx.Err() == nil {
		err := tc.SetupConn()
		if err != nil {
			if tc.ctx.Err() != nil {
				break CONNECT_RETRY
			}
			tc.logger.Error("connection error", "error", err.Error())
			time.Sleep(time.Second * 2)
			continue
		}

		buffer := make([]byte, 1024)

	READ:
		for tc.ctx.Err() == nil {
			tc.conn.SetReadDeadline(time.Now().Add(time.Millisecond * 200))
			byteCount, err := tc.conn.Read(buffer)

			if err != nil {
				if opErr, ok := err.(*net.OpError); ok {
					//NOTE(jwetzell) we hit deadline
					if opErr.Timeout() {
						continue
					}
				}
				if errors.Is(err, net.ErrClosed) {
					break CONNECT_RETRY
				}
				break READ
			}

			if tc.framer != nil {
				if byteCount > 0 {
					messages := tc.framer.Decode(buffer[0:byteCount])
					for _, message := range messages {
						if tc.inputHandler != nil {
							tc.inputHandler(tc.ctx, tc.Id(), message)
						} else {
							tc.logger.Error("input received but no input handler is configured")
						}
					}
				}
			}
		}
	}
	<-tc.ctx.Done()
	tc.logger.Debug("done")
	return nil
}

func (tc *TCPClient) SetupConn() error {
	tc.connMu.Lock()
	defer tc.connMu.Unlock()
	client, err := net.DialTCP("tcp", nil, tc.Addr)
	tc.conn = client
	return err
}

func (tc *TCPClient) Output(ctx context.Context, payload any) error {
	tc.connMu.Lock()
	defer tc.connMu.Unlock()
	if tc.conn == nil {
		return errors.New("net.tcp.client client is not setup")
	}
	payloadBytes, ok := common.GetAnyAsByteSlice(payload)
	if !ok {
		return errors.New("net.tcp.client is only able to output bytes")
	}
	_, err := tc.conn.Write(tc.framer.Encode(payloadBytes))
	return err
}

func (tc *TCPClient) Stop() {
	if tc.cancel != nil {
		defer tc.cancel()
	}
	tc.connMu.Lock()
	defer tc.connMu.Unlock()
	if tc.conn != nil {
		tc.conn.Close()
	}
}
