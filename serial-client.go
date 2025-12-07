//go:build cgo

package showbridge

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/framing"
	"go.bug.st/serial"
)

type SerialClient struct {
	config config.ModuleConfig
	router *Router
	Port   string
	Framer framing.Framer
	Mode   *serial.Mode
	port   serial.Port
}

func init() {
	RegisterModule(ModuleRegistration{
		//TODO(jwetzell): find a better namespace than "misc"
		Type: "misc.serial.client",
		New: func(config config.ModuleConfig, router *Router) (Module, error) {
			params := config.Params
			port, ok := params["port"]

			if !ok {
				return nil, fmt.Errorf("misc.serial.client requires a port parameter")
			}

			portString, ok := port.(string)

			if !ok {
				return nil, fmt.Errorf("misc.serial.client port must be a string")
			}

			framingMethod, ok := params["framing"]
			if !ok {
				return nil, fmt.Errorf("misc.serial.client requires a framing method")
			}

			framingMethodString, ok := framingMethod.(string)

			if !ok {
				return nil, fmt.Errorf("misc.serial.client framing method must be a string")
			}

			framer, err := framing.GetFramer(framingMethodString)

			if err != nil {
				return nil, err
			}

			buadRate, ok := params["baudRate"]
			if !ok {
				return nil, fmt.Errorf("misc.serial.client requires a baudRate parameter")
			}

			baudRateNum, ok := buadRate.(float64)
			if !ok {
				return nil, fmt.Errorf("misc.serial.client baudRate must be a number")
			}

			mode := serial.Mode{
				BaudRate: int(baudRateNum),
			}

			return &SerialClient{config: config, Port: portString, Framer: framer, Mode: &mode, router: router}, nil
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
		return fmt.Errorf("misc.serial.client can't open input port: %s", mc.Port)
	}

	mc.port = port

	return nil
}

func (mc *SerialClient) Run() error {

	// TODO(jwetzell): shutdown with router.Context properly
	go func() {
		<-mc.router.Context.Done()
		slog.Debug("router context done in module", "id", mc.config.Id)
		if mc.port != nil {
			mc.port.Close()
		}
	}()

	for {
		err := mc.SetupPort()
		if err != nil {
			if mc.router.Context.Err() != nil {
				slog.Debug("router context done in module", "id", mc.config.Id)
				return nil
			}
			slog.Error("misc.serial.client", "id", mc.config.Id, "error", err.Error())
			time.Sleep(time.Second * 2)
			continue
		}

		buffer := make([]byte, 1024)
		select {
		case <-mc.router.Context.Done():
			slog.Debug("router context done in module", "id", mc.config.Id)
			return nil
		default:
		READ:
			for {
				select {
				case <-mc.router.Context.Done():
					slog.Debug("router context done in module", "id", mc.config.Id)
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
									mc.router.HandleInput(mc.config.Id, message)
								} else {
									slog.Error("misc.serial.client has no router", "id", mc.config.Id)
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
		return fmt.Errorf("misc.serial.client can only ouptut bytes")
	}

	_, err := mc.port.Write(mc.Framer.Encode(payloadBytes))
	return err
}
