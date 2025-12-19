//go:build cgo

package module

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/route"
	"gitlab.com/gomidi/midi/v2"
	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv"
)

type MIDIOutput struct {
	config   config.ModuleConfig
	ctx      context.Context
	router   route.RouteIO
	Port     string
	SendFunc func(midi.Message) error
}

func init() {
	RegisterModule(ModuleRegistration{
		Type: "midi.output",
		New: func(ctx context.Context, config config.ModuleConfig, router route.RouteIO) (Module, error) {
			params := config.Params

			port, ok := params["port"]

			if !ok {
				return nil, fmt.Errorf("midi.output requires a port parameter")
			}

			portString, ok := port.(string)

			if !ok {
				return nil, fmt.Errorf("midi.output port must be a string")
			}

			return &MIDIOutput{config: config, Port: portString, ctx: ctx, router: router}, nil
		},
	})
}

func (mc *MIDIOutput) Id() string {
	return mc.config.Id
}

func (mc *MIDIOutput) Type() string {
	return mc.config.Type
}

func (mc *MIDIOutput) Run() error {
	defer midi.CloseDriver()

	out, err := midi.FindOutPort(mc.Port)

	if err != nil {
		return fmt.Errorf("midi.output can't find output port: %s", mc.Port)
	}

	send, err := midi.SendTo(out)
	if err != nil {
		return err
	}

	mc.SendFunc = send

	<-mc.ctx.Done()
	slog.Debug("router context done in module", "id", mc.Id())
	return nil
}

func (mc *MIDIOutput) Output(payload any) error {
	if mc.SendFunc == nil {
		return fmt.Errorf("midi.output output is not setup")
	}

	payloadMessage, ok := payload.(midi.Message)

	if !ok {
		return fmt.Errorf("midi.output can only ouptut midi.Message")
	}

	return mc.SendFunc(payloadMessage)
}
