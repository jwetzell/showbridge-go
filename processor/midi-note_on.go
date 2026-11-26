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
		Type:        "midi.note_on",
		Title:       "Create MIDI Note On Message",
		Description: "Creates a MIDI Note On message",
		ParamsSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"module": {
					Title:       "Module ID",
					Description: "ID of module to send output to",
					Type:        "string",
				},
				"channel": {
					Title:       "Channel",
					Description: "channel number for the note on message",
					Type:        "string",
				},
				"note": {
					Title:       "Note",
					Description: "note number for the note on message",
					Type:        "string",
				},
				"velocity": {
					Title:       "Velocity",
					Description: "velocity for the note on message",
					Type:        "string",
				},
			},
			Required:             []string{"module", "channel", "note", "velocity"},
			AdditionalProperties: &jsonschema.Schema{Not: &jsonschema.Schema{}},
		},
		New: func(config config.ProcessorConfig) (routing.Processor, error) {
			processorParams := config.Params

			moduleIdString, err := processorParams.GetString("module")
			if err != nil {
				return nil, fmt.Errorf("midi.note_on module error: %w", err)
			}

			channelString, err := processorParams.GetString("channel")
			if err != nil {
				return nil, fmt.Errorf("midi.note_on channel error: %w", err)
			}

			channelTemplate, err := template.New("channel").Parse(channelString)

			if err != nil {
				return nil, err
			}

			noteString, err := processorParams.GetString("note")
			if err != nil {
				return nil, fmt.Errorf("midi.note_on note error: %w", err)
			}

			noteTemplate, err := template.New("note").Parse(noteString)

			if err != nil {
				return nil, err
			}

			velocityString, err := processorParams.GetString("velocity")
			if err != nil {
				return nil, fmt.Errorf("midi.note_on velocity error: %w", err)
			}

			velocityTemplate, err := template.New("velocity").Parse(velocityString)

			if err != nil {
				return nil, err
			}
			return &MIDINoteOnCreate{
				config:   config,
				ModuleId: moduleIdString,
				Channel:  channelTemplate,
				Note:     noteTemplate,
				Velocity: velocityTemplate,
			}, nil
		},
	})
}

type MIDINoteOnCreate struct {
	config   config.ProcessorConfig
	ModuleId string
	module   routing.MIDIOutputModule
	mu       sync.Mutex
	Channel  *template.Template
	Note     *template.Template
	Velocity *template.Template
}

func (mnoc *MIDINoteOnCreate) Process(next routing.ProcessFunc) routing.ProcessFunc {
	return func(ctx context.Context, wrappedPayload routing.WrappedPayload, moduleManager routing.ModuleManager, router routing.Router) (routing.WrappedPayload, error) {
		tmplData := newTemplateData(ctx, wrappedPayload, moduleManager)

		var channelBuffer bytes.Buffer
		err := mnoc.Channel.Execute(&channelBuffer, tmplData)

		if err != nil {
			return wrappedPayload, err
		}

		channelValue, err := strconv.ParseUint(channelBuffer.String(), 10, 8)

		if err != nil {
			return wrappedPayload, err
		}

		var noteBuffer bytes.Buffer
		err = mnoc.Note.Execute(&noteBuffer, tmplData)

		if err != nil {
			return wrappedPayload, err
		}

		noteValue, err := strconv.ParseUint(noteBuffer.String(), 10, 8)

		if err != nil {
			return wrappedPayload, err
		}

		var velocityBuffer bytes.Buffer
		err = mnoc.Velocity.Execute(&velocityBuffer, tmplData)

		if err != nil {
			return wrappedPayload, err
		}

		velocityValue, err := strconv.ParseUint(velocityBuffer.String(), 10, 8)
		if err != nil {
			return wrappedPayload, err
		}

		payloadMessage := midi.NoteOn(uint8(channelValue), uint8(noteValue), uint8(velocityValue))
		mnoc.mu.Lock()
		defer mnoc.mu.Unlock()

		if mnoc.module == nil {
			if moduleManager == nil {
				return wrappedPayload, errors.New("midi.note_on module manager is nil")
			}

			module, ok := moduleManager.GetModule(mnoc.ModuleId)
			if !ok {
				return wrappedPayload, fmt.Errorf("midi.note_on unable to find module with id: %s", mnoc.ModuleId)
			}

			outputModule, ok := module.(routing.MIDIOutputModule)
			if !ok {
				return wrappedPayload, fmt.Errorf("midi.note_on module with id %s is not a MIDIOutputModule", mnoc.ModuleId)
			}
			mnoc.module = outputModule
		}

		err = mnoc.module.Output(ctx, payloadMessage)

		if err != nil {
			return wrappedPayload, fmt.Errorf("midi.note_on failed to send output: %w", err)
		}
		return next(ctx, wrappedPayload, moduleManager, router)
	}
}
