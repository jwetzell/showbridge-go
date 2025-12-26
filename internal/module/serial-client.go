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
}

func init() {
	RegisterModule(ModuleRegistration{
		Type: "serial.client",
		New: func(ctx context.Context, config config.ModuleConfig, router route.RouteIO) (Module, error) {
			params := config.Params
			port, ok := params["port"]

			if !ok {
				return nil, errors.New("serial.client requires a port parameter")
			}

			portString, ok := port.(string)

			if !ok {
				return nil, errors.New("serial.client port must be a string")
			}

			framingMethod := "RAW"

			framingMethodRaw, ok := params["framing"]

			if ok {
				framingMethodString, ok := framingMethodRaw.(string)

				if !ok {
					return nil, errors.New("serial.client framing method must be a string")
				}
				framingMethod = framingMethodString
			}

			framer := framer.GetFramer(framingMethod)

			if framer == nil {
				return nil, fmt.Errorf("serial.client unknown framing method: %s", framingMethod)
			}

			buadRate, ok := params["baudRate"]
			if !ok {
				return nil, errors.New("serial.client requires a baudRate parameter")
			}

			baudRateNum, ok := buadRate.(float64)
			if !ok {
				return nil, errors.New("serial.client baudRate must be a number")
			}

			mode := serial.Mode{
				BaudRate: int(baudRateNum),
			}

			return &SerialClient{config: config, Port: portString, Framer: framer, Mode: &mode, ctx: ctx, router: router, logger: CreateLogger(config)}, nil
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

func (sc *SerialClient) Run() error {

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
									sc.router.HandleInput(sc.Id(), message)
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

func (sc *SerialClient) Output(payload any) error {

	payloadBytes, ok := payload.([]byte)

	if !ok {
		return errors.New("serial.client can only ouptut bytes")
	}

	_, err := sc.port.Write(sc.Framer.Encode(payloadBytes))
	return err
}
