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
		Type:        "midi.control_change",
		Title:       "Create MIDI Control Change Message",
		Description: "Creates a MIDI Control Change message",
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
					Description: "channel number for the control change message",
					Type:        "string",
				},
				"control": {
					Title:       "Control",
					Description: "control number for the control change message",
					Type:        "string",
				},
				"value": {
					Title:       "Value",
					Description: "value for the control change message",
					Type:        "string",
				},
			},
			Required:             []string{"module", "channel", "control", "value"},
			AdditionalProperties: &jsonschema.Schema{Not: &jsonschema.Schema{}},
		},
		New: func(config config.ProcessorConfig) (routing.Processor, error) {

			processorParams := config.Params

			moduleIdString, err := processorParams.GetString("module")
			if err != nil {
				return nil, fmt.Errorf("midi.control_change module error: %w", err)
			}

			channelString, err := processorParams.GetString("channel")
			if err != nil {
				return nil, fmt.Errorf("midi.control_change channel error: %w", err)
			}

			channelTemplate, err := template.New("channel").Parse(channelString)

			if err != nil {
				return nil, err
			}

			controlString, err := processorParams.GetString("control")
			if err != nil {
				return nil, fmt.Errorf("midi.control_change control error: %w", err)
			}

			controlTemplate, err := template.New("control").Parse(controlString)

			if err != nil {
				return nil, err
			}

			valueString, err := processorParams.GetString("value")
			if err != nil {
				return nil, fmt.Errorf("midi.control_change value error: %w", err)
			}

			valueTemplate, err := template.New("value").Parse(valueString)

			if err != nil {
				return nil, err
			}
			return &MIDIControlChangeCreate{
				config:   config,
				ModuleId: moduleIdString,
				Channel:  channelTemplate,
				Control:  controlTemplate,
				Value:    valueTemplate,
			}, nil
		},
	})
}

type MIDIControlChangeCreate struct {
	config   config.ProcessorConfig
	ModuleId string
	module   routing.MIDIOutputModule
	mu       sync.Mutex
	Channel  *template.Template
	Control  *template.Template
	Value    *template.Template
}

func (mccc *MIDIControlChangeCreate) Process(next routing.ProcessFunc) routing.ProcessFunc {
	return func(ctx context.Context, wrappedPayload routing.WrappedPayload, moduleManager routing.ModuleManager, router routing.Router) (routing.WrappedPayload, error) {
		tmplData := newTemplateData(ctx, wrappedPayload, moduleManager)

		var channelBuffer bytes.Buffer
		err := mccc.Channel.Execute(&channelBuffer, tmplData)

		if err != nil {
			return wrappedPayload, err
		}

		channelValue, err := strconv.ParseUint(channelBuffer.String(), 10, 8)

		if err != nil {
			return wrappedPayload, err
		}

		var controlBuffer bytes.Buffer
		err = mccc.Control.Execute(&controlBuffer, tmplData)

		if err != nil {
			return wrappedPayload, err
		}

		controlValue, err := strconv.ParseUint(controlBuffer.String(), 10, 8)

		if err != nil {
			return wrappedPayload, err
		}

		var valueBuffer bytes.Buffer
		err = mccc.Value.Execute(&valueBuffer, tmplData)

		if err != nil {
			return wrappedPayload, err
		}

		valueValue, err := strconv.ParseUint(valueBuffer.String(), 10, 8)

		if err != nil {
			return wrappedPayload, err
		}

		payloadMessage := midi.ControlChange(uint8(channelValue), uint8(controlValue), uint8(valueValue))
		mccc.mu.Lock()
		defer mccc.mu.Unlock()

		if mccc.module == nil {
			if moduleManager == nil {
				return wrappedPayload, errors.New("midi.control_change module manager is nil")
			}

			module, ok := moduleManager.GetModule(mccc.ModuleId)
			if !ok {
				return wrappedPayload, fmt.Errorf("midi.control_change unable to find module with id: %s", mccc.ModuleId)
			}

			outputModule, ok := module.(routing.MIDIOutputModule)
			if !ok {
				return wrappedPayload, fmt.Errorf("midi.control_change module with id %s is not a MIDIOutputModule", mccc.ModuleId)
			}
			mccc.module = outputModule
		}

		err = mccc.module.Output(ctx, payloadMessage)

		if err != nil {
			return wrappedPayload, fmt.Errorf("midi.control_change failed to send output: %w", err)
		}
		return next(ctx, wrappedPayload, moduleManager, router)
	}
}
