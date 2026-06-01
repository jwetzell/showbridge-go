//go:build cgo

package module

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/framer"
	"go.bug.st/serial"
)

func init() {
	RegisterModule(ModuleRegistration{
		Type:  "serial.client",
		Title: "Serial Client",
		ParamsSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"port": {
					Title:       "Port",
					Description: "the name of the serial port to connect to",
					Type:        "string",
				},
				"baudRate": {
					Title:       "Baud Rate",
					Description: "the baud rate to use when connecting to the serial port",
					Type:        "integer",
				},
				"framing": {
					Title:       "Framing Method",
					Description: "the method to use for framing messages on the serial port",
					Type:        "string",
					Enum:        []any{"LF", "CR", "CRLF", "SLIP", "RAW"},
				},
			},
			Required:             []string{"port", "baudRate", "framing"},
			AdditionalProperties: &jsonschema.Schema{Not: &jsonschema.Schema{}},
		},
		New: func(config config.ModuleConfig) (common.Module, error) {
			params := config.Params
			portString, err := params.GetString("port")
			if err != nil {
				return nil, fmt.Errorf("serial.client port error: %w", err)
			}

			framingMethodString, err := params.GetString("framing")
			if err != nil {
				return nil, fmt.Errorf("serial.client framing error: %w", err)
			}

			framer := framer.GetFramer(framingMethodString)

			if framer == nil {
				return nil, fmt.Errorf("serial.client unknown framing method: %s", framingMethodString)
			}

			baudRateInt, err := params.GetInt("baudRate")
			if err != nil {
				return nil, fmt.Errorf("serial.client baudRate error: %w", err)
			}

			mode := serial.Mode{
				BaudRate: baudRateInt,
			}

			return &SerialClient{config: config, Port: portString, Framer: framer, Mode: &mode, logger: CreateLogger(config)}, nil
		},
	})
}

type SerialClient struct {
	config       config.ModuleConfig
	ctx          context.Context
	inputHandler common.InputHandler
	Port         string
	Framer       framer.Framer
	Mode         *serial.Mode
	port         serial.Port
	logger       *slog.Logger
	cancel       context.CancelFunc
	portMu       sync.Mutex
}

func (sc *SerialClient) Id() string {
	return sc.config.Id
}

func (sc *SerialClient) Type() string {
	return sc.config.Type
}

func (sc *SerialClient) SetupPort() error {
	sc.portMu.Lock()
	defer sc.portMu.Unlock()
	port, err := serial.Open(sc.Port, sc.Mode)
	if err != nil {
		return err
	}

	sc.port = port

	return nil
}

func (sc *SerialClient) Start(ctx context.Context, inputHandler common.InputHandler) error {
	sc.logger.Debug("running")
	sc.inputHandler = inputHandler
	moduleContext, cancel := context.WithCancel(ctx)
	sc.ctx = moduleContext
	sc.cancel = cancel

	for sc.ctx.Err() == nil {
		err := sc.SetupPort()
		if err != nil {
			if sc.ctx.Err() != nil {
				return nil
			}
			sc.logger.Error("port setup error", "port", sc.Port, "error", err.Error())
			time.Sleep(time.Second * 2)
			continue
		}

		buffer := make([]byte, 1024)
		select {
		case <-sc.ctx.Done():
			return nil
		default:
		READ:
			for sc.ctx.Err() == nil {
				select {
				case <-sc.ctx.Done():
					return nil
				default:
					byteCount, err := sc.port.Read(buffer)

					if err != nil {
						sc.Framer.Clear()
						break READ
					}

					if sc.Framer != nil {
						if byteCount > 0 {
							messages := sc.Framer.Decode(buffer[0:byteCount])
							for _, message := range messages {
								if sc.inputHandler != nil {
									sc.inputHandler(sc.ctx, sc.Id(), message)
								} else {
									sc.logger.Error("input received but no router is configured")
								}
							}
						}
					}
				}
			}
		}
	}
	<-sc.ctx.Done()
	sc.logger.Debug("done")
	return nil
}

func (sc *SerialClient) Output(ctx context.Context, payload any) error {

	payloadBytes, ok := common.GetAnyAsByteSlice(payload)

	if !ok {
		return errors.New("serial.client can only output bytes")
	}

	_, err := sc.port.Write(sc.Framer.Encode(payloadBytes))
	return err
}

func (sc *SerialClient) Stop() {
	if sc.cancel != nil {
		defer sc.cancel()
	}
	sc.portMu.Lock()
	defer sc.portMu.Unlock()
	if sc.port != nil {
		sc.port.Close()
	}
}
