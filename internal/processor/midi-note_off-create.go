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

type MIDINoteOffCreate struct {
	config   config.ProcessorConfig
	Channel  *template.Template
	Note     *template.Template
	Velocity *template.Template
}

func (mnoc *MIDINoteOffCreate) Process(ctx context.Context, wrappedPayload common.WrappedPayload) (common.WrappedPayload, error) {
	templateData := wrappedPayload

	var channelBuffer bytes.Buffer
	err := mnoc.Channel.Execute(&channelBuffer, templateData)

	if err != nil {
		wrappedPayload.End = true
		return wrappedPayload, err
	}

	channelValue, err := strconv.ParseUint(channelBuffer.String(), 10, 8)

	var noteBuffer bytes.Buffer
	err = mnoc.Note.Execute(&noteBuffer, templateData)

	if err != nil {
		wrappedPayload.End = true
		return wrappedPayload, err
	}

	noteValue, err := strconv.ParseUint(noteBuffer.String(), 10, 8)

	var velocityBuffer bytes.Buffer
	err = mnoc.Velocity.Execute(&velocityBuffer, templateData)

	if err != nil {
		wrappedPayload.End = true
		return wrappedPayload, err
	}

	velocityValue, err := strconv.ParseUint(velocityBuffer.String(), 10, 8)
	payloadMessage := midi.NoteOffVelocity(uint8(channelValue), uint8(noteValue), uint8(velocityValue))
	wrappedPayload.Payload = payloadMessage
	return wrappedPayload, nil
}

func (mnoc *MIDINoteOffCreate) Type() string {
	return mnoc.config.Type
}

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type:  "midi.note_off.create",
		Title: "Create MIDI Note Off Message",
		ParamsSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"channel": {
					Title: "Channel",
					Type:  "string",
				},
				"note": {
					Title: "Note",
					Type:  "string",
				},
				"velocity": {
					Title: "Velocity",
					Type:  "string",
				},
			},
			Required:             []string{"channel", "note", "velocity"},
			AdditionalProperties: &jsonschema.Schema{Not: &jsonschema.Schema{}},
		},
		New: func(config config.ProcessorConfig) (Processor, error) {
			params := config.Params

			channelString, err := params.GetString("channel")
			if err != nil {
				return nil, fmt.Errorf("midi.note_off.create channel error: %w", err)
			}

			channelTemplate, err := template.New("channel").Parse(channelString)

			if err != nil {
				return nil, err
			}

			noteString, err := params.GetString("note")
			if err != nil {
				return nil, fmt.Errorf("midi.note_off.create note error: %w", err)
			}

			noteTemplate, err := template.New("note").Parse(noteString)

			if err != nil {
				return nil, err
			}

			velocityString, err := params.GetString("velocity")
			if err != nil {
				return nil, fmt.Errorf("midi.note_off.create velocity error: %w", err)
			}

			velocityTemplate, err := template.New("velocity").Parse(velocityString)

			if err != nil {
				return nil, err
			}
			return &MIDINoteOffCreate{
				config:   config,
				Channel:  channelTemplate,
				Note:     noteTemplate,
				Velocity: velocityTemplate,
			}, nil
		},
	})
}
