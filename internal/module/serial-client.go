//go:build cgo

package module

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/framer"
	"github.com/jwetzell/showbridge-go/internal/route"
	"go.bug.st/serial"
)

type SerialClient struct {
	config config.ModuleConfig
	ctx    context.Context
	router route.RouteIO
	Port   string
	Framer framer.Framer
	Mode   *serial.Mode
	port   serial.Port
	logger *slog.Logger
	cancel context.CancelFunc
}

func init() {
	RegisterModule(ModuleRegistration{
		Type: "serial.client",
		New: func(config config.ModuleConfig) (Module, error) {
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

func (sc *SerialClient) Id() string {
	return sc.config.Id
}

func (sc *SerialClient) Type() string {
	return sc.config.Type
}

func (sc *SerialClient) SetupPort() error {

	port, err := serial.Open(sc.Port, sc.Mode)
	if err != nil {
		return fmt.Errorf("serial.client can't open input port: %s", sc.Port)
	}

	sc.port = port

	return nil
}

func (sc *SerialClient) Start(ctx context.Context) error {
	sc.logger.Debug("running")
	router, ok := ctx.Value(route.RouterContextKey).(route.RouteIO)

	if !ok {
		return errors.New("serial.client unable to get router from context")
	}

	sc.router = router
	moduleContext, cancel := context.WithCancel(ctx)
	sc.ctx = moduleContext
	sc.cancel = cancel

	// TODO(jwetzell): shutdown with router.Context properly
	go func() {
		<-sc.ctx.Done()
		sc.logger.Debug("done")
		if sc.port != nil {
			sc.port.Close()
		}
	}()

	for {
		err := sc.SetupPort()
		if err != nil {
			if sc.ctx.Err() != nil {
				sc.logger.Debug("done")
				return nil
			}
			sc.logger.Error("port setup error", "error", err.Error())
			time.Sleep(time.Second * 2)
			continue
		}

		buffer := make([]byte, 1024)
		select {
		case <-sc.ctx.Done():
			sc.logger.Debug("done")
			return nil
		default:
		READ:
			for {
				select {
				case <-sc.ctx.Done():
					sc.logger.Debug("done")
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
								if sc.router != nil {
									sc.router.HandleInput(sc.ctx, sc.Id(), message)
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
}

func (sc *SerialClient) Output(ctx context.Context, payload any) error {

	payloadBytes, ok := payload.([]byte)

	if !ok {
		return errors.New("serial.client can only ouptut bytes")
	}

	_, err := sc.port.Write(sc.Framer.Encode(payloadBytes))
	return err
}

func (sc *SerialClient) Stop() {
	sc.cancel()
}
