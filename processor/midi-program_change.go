//go:build cgo

package processor

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"text/template"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/jwetzell/showbridge-go/config"
	"github.com/jwetzell/showbridge-go/routing"
	"gitlab.com/gomidi/midi/v2"
)

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type:        "midi.program_change",
		Title:       "Output MIDI Program Change Message",
		Description: "Outputs a MIDI Program Change message",
		ParamsSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"module": {
					Title:       "Module ID",
					Description: "Output module for the program change message",
					Type:        "string",
				},
				"channel": {
					Title:       "Channel",
					Description: "channel number for the program change message",
					Type:        "string",
				},
				"program": {
					Title:       "Program",
					Description: "program number for the program change message",
					Type:        "string",
				},
			},
			Required:             []string{"module", "channel", "program"},
			AdditionalProperties: &jsonschema.Schema{Not: &jsonschema.Schema{}},
		},
		New: func(config config.ProcessorConfig) (routing.Processor, error) {
			processorParams := config.Params

			moduleIdString, err := processorParams.GetString("module")
			if err != nil {
				return nil, fmt.Errorf("midi.program_change output error: %w", err)
			}

			channelString, err := processorParams.GetString("channel")
			if err != nil {
				return nil, fmt.Errorf("midi.program_change channel error: %w", err)
			}

			channelTemplate, err := template.New("channel").Parse(channelString)

			if err != nil {
				return nil, err
			}

			programString, err := processorParams.GetString("program")
			if err != nil {
				return nil, fmt.Errorf("midi.program_change program error: %w", err)
			}

			programTemplate, err := template.New("program").Parse(programString)

			if err != nil {
				return nil, err
			}
			return &MIDIProgramChangeCreate{
				config:   config,
				ModuleId: moduleIdString,
				Channel:  channelTemplate,
				Program:  programTemplate,
			}, nil
		},
	})
}

type MIDIProgramChangeCreate struct {
	config   config.ProcessorConfig
	ModuleId string
	Channel  *template.Template
	Program  *template.Template
	module   routing.MIDIOutputModule
	mu       sync.Mutex
}

func (mpcc *MIDIProgramChangeCreate) Process(next routing.ProcessFunc) routing.ProcessFunc {
	return func(ctx context.Context, wrappedPayload routing.WrappedPayload, moduleManager routing.ModuleManager, router routing.Router) (routing.WrappedPayload, error) {
		templateData := wrappedPayload

		var channelBuffer bytes.Buffer
		err := mpcc.Channel.Execute(&channelBuffer, templateData)

		if err != nil {
			return wrappedPayload, err
		}

		channelValue, err := strconv.ParseUint(channelBuffer.String(), 10, 8)

		if err != nil {
			return wrappedPayload, err
		}

		var programBuffer bytes.Buffer
		err = mpcc.Program.Execute(&programBuffer, templateData)

		if err != nil {
			return wrappedPayload, err
		}

		programValue, err := strconv.ParseUint(programBuffer.String(), 10, 8)

		if err != nil {
			return wrappedPayload, err
		}

		payloadMessage := midi.ProgramChange(uint8(channelValue), uint8(programValue))

		mpcc.mu.Lock()
		defer mpcc.mu.Unlock()

		if mpcc.module == nil {
			if moduleManager == nil {
				return wrappedPayload, errors.New("midi.program_change module manager is nil")
			}

			module, ok := moduleManager.GetModule(mpcc.ModuleId)
			if !ok {
				return wrappedPayload, fmt.Errorf("midi.program_change unable to find module with id: %s", mpcc.ModuleId)
			}

			outputModule, ok := module.(routing.MIDIOutputModule)
			if !ok {
				return wrappedPayload, fmt.Errorf("midi.program_change module with id %s is not a MIDIOutputModule", mpcc.ModuleId)
			}
			mpcc.module = outputModule
		}

		err = mpcc.module.Output(ctx, payloadMessage)

		if err != nil {
			return wrappedPayload, fmt.Errorf("midi.program_change failed to send output: %w", err)
		}
		return next(ctx, wrappedPayload, moduleManager, router)
	}
}
