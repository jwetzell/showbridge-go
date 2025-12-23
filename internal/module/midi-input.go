//go:build cgo

package module

import (
	"context"
	"errors"
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
	logger   *slog.Logger
}

func init() {
	RegisterModule(ModuleRegistration{
		Type: "midi.input",
		New: func(ctx context.Context, config config.ModuleConfig, router route.RouteIO) (Module, error) {
			params := config.Params
			port, ok := params["port"]

			if !ok {
				return nil, errors.New("midi.input requires a port parameter")
			}

			portString, ok := port.(string)

			if !ok {
				return nil, errors.New("midi.input port must be a string")
			}

			return &MIDIInput{config: config, Port: portString, ctx: ctx, router: router, logger: slog.Default().With("component", "module", "id", config.Id)}, nil
		},
	})
}

func (mi *MIDIInput) Id() string {
	return mi.config.Id
}

func (mi *MIDIInput) Type() string {
	return mi.config.Type
}

func (mi *MIDIInput) Run() error {
	defer midi.CloseDriver()

	in, err := midi.FindInPort(mi.Port)
	if err != nil {
		return fmt.Errorf("midi.input can't find input port: %s", mi.Port)
	}

	stop, err := midi.ListenTo(in, func(msg midi.Message, timestampms int32) {
		if mi.router != nil {
			mi.router.HandleInput(mi.Id(), msg)
		}
	}, midi.UseSysEx())

	if err != nil {
		return err
	}

	defer stop()

	<-mi.ctx.Done()
	mi.logger.Debug("router context done in module")
	return nil
}

func (mi *MIDIInput) Output(payload any) error {
	return errors.New("midi.input output is not implemented")
}
