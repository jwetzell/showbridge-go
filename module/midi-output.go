//go:build cgo && !android

package module

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/jwetzell/showbridge-go/config"
	"github.com/jwetzell/showbridge-go/routing"
	"gitlab.com/gomidi/midi/v2"
	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv"
)

func init() {
	RegisterModule(ModuleRegistration{
		Type:        "midi.output",
		Title:       "MIDI Output",
		Description: "Sends MIDI messages to the specified port",
		ParamsSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"port": {
					Title:       "Port",
					Description: "the name of the MIDI port to send messages to",
					Type:        "string",
				},
			},
			Required:             []string{"port"},
			AdditionalProperties: &jsonschema.Schema{Not: &jsonschema.Schema{}},
		},
		New: func(config config.ModuleConfig) (routing.Module, error) {
			moduleParams := config.Params

			portString, err := moduleParams.GetString("port")
			if err != nil {
				return nil, fmt.Errorf("midi.output port error: %w", err)
			}

			return &MIDIOutput{config: config, Port: portString, logger: CreateLogger(config)}, nil
		},
	})
}

type MIDIOutput struct {
	config     config.ModuleConfig
	ctx        context.Context
	router     routing.Router
	Port       string
	sendFunc   func(midi.Message) error
	logger     *slog.Logger
	cancel     context.CancelFunc
	sendFuncMu sync.Mutex
}

func (mo *MIDIOutput) Start(ctx context.Context, router routing.Router) error {
	mo.logger.Debug("running")
	mo.router = router
	moduleContext, cancel := context.WithCancel(ctx)
	mo.ctx = moduleContext
	mo.cancel = cancel

	out, err := midi.FindOutPort(mo.Port)

	if err != nil {
		return fmt.Errorf("midi.output can't find output port: %s", mo.Port)
	}

	send, err := midi.SendTo(out)
	if err != nil {
		return err
	}

	mo.sendFuncMu.Lock()
	mo.sendFunc = send
	mo.sendFuncMu.Unlock()

	<-mo.ctx.Done()
	mo.logger.Debug("done")
	return nil
}

func (mo *MIDIOutput) Output(ctx context.Context, payload any) error {
	mo.sendFuncMu.Lock()
	defer mo.sendFuncMu.Unlock()
	if mo.sendFunc == nil {
		return errors.New("midi.output output is not setup")
	}

	payloadMessage, ok := payload.(midi.Message)

	if !ok {
		return errors.New("midi.output can only output midi.Message")
	}

	return mo.sendFunc(payloadMessage)
}

func (mo *MIDIOutput) Stop() {
	if mo.cancel != nil {
		defer mo.cancel()
	}
	midi.CloseDriver()
}
