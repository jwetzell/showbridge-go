//go:build cgo

package module

import (
	"context"
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
}

func init() {
	RegisterModule(ModuleRegistration{
		//TODO(jwetzell): find a better namespace than "misc"
		Type: "serial.client",
		New: func(ctx context.Context, config config.ModuleConfig, router route.RouteIO) (Module, error) {
			params := config.Params
			port, ok := params["port"]

			if !ok {
				return nil, fmt.Errorf("serial.client requires a port parameter")
			}

			portString, ok := port.(string)

			if !ok {
				return nil, fmt.Errorf("serial.client port must be a string")
			}

			framingMethod := "RAW"

			framingMethodRaw, ok := params["framing"]

			if ok {
				framingMethodString, ok := framingMethodRaw.(string)

				if !ok {
					return nil, fmt.Errorf("serial.client framing method must be a string")
				}
				framingMethod = framingMethodString
			}

			framer, err := framer.GetFramer(framingMethod)

			if err != nil {
				return nil, err
			}

			buadRate, ok := params["baudRate"]
			if !ok {
				return nil, fmt.Errorf("serial.client requires a baudRate parameter")
			}

			baudRateNum, ok := buadRate.(float64)
			if !ok {
				return nil, fmt.Errorf("serial.client baudRate must be a number")
			}

			mode := serial.Mode{
				BaudRate: int(baudRateNum),
			}

			return &SerialClient{config: config, Port: portString, Framer: framer, Mode: &mode, ctx: ctx, router: router}, nil
		},
	})
}

func (mc *SerialClient) Id() string {
	return mc.config.Id
}

func (mc *SerialClient) Type() string {
	return mc.config.Type
}

func (mc *SerialClient) SetupPort() error {

	port, err := serial.Open(mc.Port, mc.Mode)
	if err != nil {
		return fmt.Errorf("serial.client can't open input port: %s", mc.Port)
	}

	mc.port = port

	return nil
}

func (mc *SerialClient) Run() error {

	// TODO(jwetzell): shutdown with router.Context properly
	go func() {
		<-mc.ctx.Done()
		slog.Debug("router context done in module", "id", mc.Id())
		if mc.port != nil {
			mc.port.Close()
		}
	}()

	for {
		err := mc.SetupPort()
		if err != nil {
			if mc.ctx.Err() != nil {
				slog.Debug("router context done in module", "id", mc.Id())
				return nil
			}
			slog.Error("serial.client", "id", mc.Id(), "error", err.Error())
			time.Sleep(time.Second * 2)
			continue
		}

		buffer := make([]byte, 1024)
		select {
		case <-mc.ctx.Done():
			slog.Debug("router context done in module", "id", mc.Id())
			return nil
		default:
		READ:
			for {
				select {
				case <-mc.ctx.Done():
					slog.Debug("router context done in module", "id", mc.Id())
					return nil
				default:
					byteCount, err := mc.port.Read(buffer)

					if err != nil {
						mc.Framer.Clear()
						break READ
					}

					if mc.Framer != nil {
						if byteCount > 0 {
							messages := mc.Framer.Decode(buffer[0:byteCount])
							for _, message := range messages {
								if mc.router != nil {
									mc.router.HandleInput(mc.Id(), message)
								} else {
									slog.Error("serial.client has no router", "id", mc.Id())
								}
							}
						}
					}
				}
			}
		}
	}
}

func (mc *SerialClient) Output(payload any) error {

	payloadBytes, ok := payload.([]byte)

	if !ok {
		return fmt.Errorf("serial.client can only ouptut bytes")
	}

	_, err := mc.port.Write(mc.Framer.Encode(payloadBytes))
	return err
}
