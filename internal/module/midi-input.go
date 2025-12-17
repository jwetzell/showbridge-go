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

type MIDIInput struct {
	config   config.ModuleConfig
	ctx      context.Context
	router   route.RouteIO
	Port     string
	SendFunc func(midi.Message) error
}

func init() {
	RegisterModule(ModuleRegistration{
		Type: "midi.input",
		New: func(ctx context.Context, config config.ModuleConfig, router route.RouteIO) (Module, error) {
			params := config.Params
			port, ok := params["port"]

			if !ok {
				return nil, fmt.Errorf("midi.input requires a port parameter")
			}

			portString, ok := port.(string)

			if !ok {
				return nil, fmt.Errorf("midi.input port must be a string")
			}

			return &MIDIInput{config: config, Port: portString, ctx: ctx, router: router}, nil
		},
	})
}

func (mc *MIDIInput) Id() string {
	return mc.config.Id
}

func (mc *MIDIInput) Type() string {
	return mc.config.Type
}

func (mc *MIDIInput) Run() error {
	defer midi.CloseDriver()

	in, err := midi.FindInPort(mc.Port)
	if err != nil {
		return fmt.Errorf("midi.input can't find input port: %s", mc.Port)
	}

	stop, err := midi.ListenTo(in, func(msg midi.Message, timestampms int32) {
		if mc.router != nil {
			mc.router.HandleInput(mc.Id(), msg)
		}
	}, midi.UseSysEx())

	if err != nil {
		return err
	}

	defer stop()

	<-mc.ctx.Done()
	slog.Debug("router context done in module", "id", mc.Id())
	return nil
}

func (mc *MIDIInput) Output(payload any) error {
	return fmt.Errorf("midi.input output is not implemented")
}
