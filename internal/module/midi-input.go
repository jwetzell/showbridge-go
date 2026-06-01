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

func init() {
	RegisterModule(ModuleRegistration{
		Type:  "midi.input",
		Title: "MIDI Input",
		ParamsSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"port": {
					Title:       "Port",
					Description: "the name of the MIDI port to listen to",
					Type:        "string",
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

type MIDIInput struct {
	config       config.ModuleConfig
	ctx          context.Context
	inputHandler common.InputHandler
	Port         string
	logger       *slog.Logger
	cancel       context.CancelFunc
	stop         func()
}

func (mi *MIDIInput) Id() string {
	return mi.config.Id
}

func (mi *MIDIInput) Type() string {
	return mi.config.Type
}

func (mi *MIDIInput) Start(ctx context.Context, inputHandler common.InputHandler) error {
	mi.logger.Debug("running")
	mi.inputHandler = inputHandler
	moduleContext, cancel := context.WithCancel(ctx)
	mi.ctx = moduleContext
	mi.cancel = cancel

	in, err := midi.FindInPort(mi.Port)
	if err != nil {
		return fmt.Errorf("midi.input can't find input port: %s", mi.Port)
	}

	stop, err := midi.ListenTo(in, func(msg midi.Message, timestampms int32) {
		if mi.inputHandler != nil {
			mi.inputHandler(mi.ctx, mi.Id(), msg)
		}
	}, midi.UseSysEx())

	if err != nil {
		return err
	}
	mi.stop = stop

	<-mi.ctx.Done()
	mi.logger.Debug("done")
	return nil
}

func (mi *MIDIInput) Stop() {
	if mi.cancel != nil {
		defer mi.cancel()
	}
	if mi.stop != nil {
		mi.stop()
	}
	midi.CloseDriver()
}
