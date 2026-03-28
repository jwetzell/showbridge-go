//go:build cgo

package processor

import (
	"bytes"
	"context"
	"fmt"
	"strconv"
	"text/template"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/config"
	"gitlab.com/gomidi/midi/v2"
)

type MIDIControlChangeCreate struct {
	config  config.ProcessorConfig
	Channel *template.Template
	Control *template.Template
	Value   *template.Template
}

func (mccc *MIDIControlChangeCreate) Process(ctx context.Context, wrappedPayload common.WrappedPayload) (common.WrappedPayload, error) {
	templateData := wrappedPayload

	var channelBuffer bytes.Buffer
	err := mccc.Channel.Execute(&channelBuffer, templateData)

	if err != nil {
		wrappedPayload.End = true
		return wrappedPayload, err
	}

	channelValue, err := strconv.ParseUint(channelBuffer.String(), 10, 8)

	var controlBuffer bytes.Buffer
	err = mccc.Control.Execute(&controlBuffer, templateData)

	if err != nil {
		wrappedPayload.End = true
		return wrappedPayload, err
	}

	controlValue, err := strconv.ParseUint(controlBuffer.String(), 10, 8)

	var valueBuffer bytes.Buffer
	err = mccc.Value.Execute(&valueBuffer, templateData)

	if err != nil {
		wrappedPayload.End = true
		return wrappedPayload, err
	}

	valueValue, err := strconv.ParseUint(valueBuffer.String(), 10, 8)

	payloadMessage := midi.ControlChange(uint8(channelValue), uint8(controlValue), uint8(valueValue))
	wrappedPayload.Payload = payloadMessage
	return wrappedPayload, nil
}

func (mccc *MIDIControlChangeCreate) Type() string {
	return mccc.config.Type
}

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type:  "midi.control_change.create",
		Title: "Create MIDI Control Change Message",
		ParamsSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"channel": {
					Title: "Channel",
					Type:  "string",
				},
				"control": {
					Title: "Control",
					Type:  "string",
				},
				"value": {
					Title: "Value",
					Type:  "string",
				},
			},
			Required:             []string{"channel", "control", "value"},
			AdditionalProperties: nil,
		},
		New: func(config config.ProcessorConfig) (Processor, error) {

			params := config.Params

			channelString, err := params.GetString("channel")
			if err != nil {
				return nil, fmt.Errorf("midi.control_change.create channel error: %w", err)
			}

			channelTemplate, err := template.New("channel").Parse(channelString)

			if err != nil {
				return nil, err
			}

			controlString, err := params.GetString("control")
			if err != nil {
				return nil, fmt.Errorf("midi.control_change.create control error: %w", err)
			}

			controlTemplate, err := template.New("control").Parse(controlString)

			if err != nil {
				return nil, err
			}

			valueString, err := params.GetString("value")
			if err != nil {
				return nil, fmt.Errorf("midi.control_change.create value error: %w", err)
			}

			valueTemplate, err := template.New("value").Parse(valueString)

			if err != nil {
				return nil, err
			}
			return &MIDIControlChangeCreate{
				config:  config,
				Channel: channelTemplate,
				Control: controlTemplate,
				Value:   valueTemplate,
			}, nil
		},
	})
}
