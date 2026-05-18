//go:build cgo

package module

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/config"
	"gitlab.com/gomidi/midi/v2"
	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv"
)

type MIDIInput struct {
	config config.ModuleConfig
	ctx    context.Context
	router common.RouteIO
	Port   string
	logger *slog.Logger
	cancel context.CancelFunc
	stop   func()
}

func init() {
	RegisterModule(ModuleRegistration{
		Type:  "midi.input",
		Title: "MIDI Input",
		ParamsSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"port": {
					Title: "Port",
					Type:  "string",
				},
			},
			Required:             []string{"port"},
			AdditionalProperties: &jsonschema.Schema{Not: &jsonschema.Schema{}},
		},
		New: func(config config.ModuleConfig) (common.Module, error) {
			params := config.Params
			portString, err := params.GetString("port")
			if err != nil {
				return nil, fmt.Errorf("midi.input port error: %w", err)
			}

			return &MIDIInput{config: config, Port: portString, logger: CreateLogger(config)}, nil
		},
	})
}

func (mi *MIDIInput) Id() string {
	return mi.config.Id
}

func (mi *MIDIInput) Type() string {
	return mi.config.Type
}

func (mi *MIDIInput) Start(ctx context.Context, router common.RouteIO) error {
	mi.logger.Debug("running")
	mi.router = router
	moduleContext, cancel := context.WithCancel(ctx)
	mi.ctx = moduleContext
	mi.cancel = cancel

	in, err := midi.FindInPort(mi.Port)
	if err != nil {
		return fmt.Errorf("midi.input can't find input port: %s", mi.Port)
	}

	stop, err := midi.ListenTo(in, func(msg midi.Message, timestampms int32) {
		if mi.router != nil {
			mi.router.HandleInput(mi.ctx, mi.Id(), msg)
		}
	}, midi.UseSysEx())

	if err != nil {
		return err
	}
	mi.stop = stop

	<-mi.ctx.Done()
	return nil
}

func (mi *MIDIInput) Stop() {
	if mi.cancel != nil {
		mi.cancel()
	}
	if mi.stop != nil {
		mi.stop()
		mi.stop = nil
	}
	midi.CloseDriver()
	mi.logger.Debug("done")
}
