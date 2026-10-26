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
	"github.com/jwetzell/showbridge-go/config"
	"github.com/jwetzell/showbridge-go/framer"
	"github.com/jwetzell/showbridge-go/params"
	"github.com/jwetzell/showbridge-go/routing"
)

func init() {
	RegisterModule(ModuleRegistration{
		Type:        "net.tcp.client",
		Title:       "TCP Client",
		Description: "Connects to the specified TCP server using the provided host and port",
		ParamsSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"host": {
					Title:       "Host",
					Description: "the hostname or IP address of the TCP server to connect to",
					Type:        "string",
				},
				"port": {
					Title:       "Port",
					Description: "the port of the TCP server to connect to",
					Type:        "integer",
					Minimum:     jsonschema.Ptr[float64](1),
					Maximum:     jsonschema.Ptr[float64](65535),
				},
				"framing": {
					Title:       "Framing Method",
					Description: "the method used to frame messages over the TCP connection",
					Type:        "string",
					Enum:        []any{"LF", "CR", "CRLF", "SLIP", "RAW"},
				},
			},
			Required:             []string{"host", "port", "framing"},
			AdditionalProperties: &jsonschema.Schema{Not: &jsonschema.Schema{}},
		},
		New: func(config config.ModuleConfig) (routing.Module, error) {
			moduleParams := config.Params
			hostString, err := moduleParams.GetString("host")
			if err != nil {
				return nil, fmt.Errorf("net.tcp.client host error: %w", err)
			}

			portNum, err := moduleParams.GetInt("port")
			if err != nil {
				return nil, fmt.Errorf("net.tcp.client port error: %w", err)
			}

			addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", hostString, uint16(portNum)))
			if err != nil {
				return nil, err
			}

			framingMethodString, err := moduleParams.GetString("framing")
			if err != nil {
				return nil, fmt.Errorf("net.tcp.client framing error: %w", err)
			}

			inFramer := framer.GetFramer(framingMethodString)

			if inFramer == nil {
				return nil, fmt.Errorf("net.tcp.client unknown framing method: %s", framingMethodString)
			}

			return &TCPClient{inFramer: inFramer, framerType: framingMethodString, Addr: addr, config: config, logger: CreateLogger(config)}, nil
		},
	})
}

type TCPClient struct {
	config          config.ModuleConfig
	framerType      string
	inFramer        framer.Framer
	outFramedWriter *framer.FramedWriter
	conn            *net.TCPConn
	ctx             context.Context
	router          routing.Router
	Addr            *net.TCPAddr
	logger          *slog.Logger
	cancel          context.CancelFunc
	connMu          sync.Mutex
}

func (tc *TCPClient) Start(ctx context.Context, router routing.Router) error {
	tc.logger.Debug("running")
	tc.router = router
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
			err := tc.conn.SetReadDeadline(time.Now().Add(time.Millisecond * 200))
			if err != nil {
				tc.logger.Error("tcp client set read deadline error", "error", err)
				continue
			}
			byteCount, err := tc.conn.Read(buffer)

			if err != nil {
				opErr, ok := err.(*net.OpError)
				if ok {
					// NOTE(jwetzell) we hit deadline
					if opErr.Timeout() {
						continue
					}
				}
				if errors.Is(err, net.ErrClosed) {
					break CONNECT_RETRY
				}
				break READ
			}

			if tc.inFramer != nil {
				if byteCount > 0 {
					messages := tc.inFramer.Decode(buffer[0:byteCount])
					for _, message := range messages {
						if tc.router != nil {
							contextWriter := routing.NewContextWriter(tc.outFramedWriter)
							contextWithWriter := context.WithValue(tc.ctx, routing.WriterContextKey, contextWriter)
							tc.router.HandleInput(contextWithWriter, tc.config.Id, message)
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
	tc.outFramedWriter = framer.NewFramedWrite(client, framer.GetFramer(tc.framerType))
	return err
}

func (tc *TCPClient) Output(ctx context.Context, payload any) error {
	tc.connMu.Lock()
	defer tc.connMu.Unlock()
	if tc.conn == nil {
		return errors.New("net.tcp.client client is not setup")
	}
	payloadBytes, ok := params.GetAnyAsByteSlice(payload)
	if !ok {
		return errors.New("net.tcp.client is only able to output bytes")
	}
	_, err := tc.outFramedWriter.Write(payloadBytes)
	return err
}

func (tc *TCPClient) Stop() {
	if tc.cancel != nil {
		defer tc.cancel()
	}
	tc.connMu.Lock()
	defer tc.connMu.Unlock()
	if tc.conn != nil {
		err := tc.conn.Close()
		if err != nil {
			tc.logger.Error("tcp client close error", "error", err)
		}
	}

	if tc.inFramer != nil {
		tc.inFramer.Clear()
	}
}
