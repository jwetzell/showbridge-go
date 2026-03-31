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

type MIDIProgramChangeCreate struct {
	config  config.ProcessorConfig
	Channel *template.Template
	Program *template.Template
}

func (mpcc *MIDIProgramChangeCreate) Process(ctx context.Context, wrappedPayload common.WrappedPayload) (common.WrappedPayload, error) {
	templateData := wrappedPayload

	var channelBuffer bytes.Buffer
	err := mpcc.Channel.Execute(&channelBuffer, templateData)

	if err != nil {
		wrappedPayload.End = true
		return wrappedPayload, err
	}

	channelValue, err := strconv.ParseUint(channelBuffer.String(), 10, 8)

	var programBuffer bytes.Buffer
	err = mpcc.Program.Execute(&programBuffer, templateData)

	if err != nil {
		wrappedPayload.End = true
		return wrappedPayload, err
	}

	programValue, err := strconv.ParseUint(programBuffer.String(), 10, 8)

	payloadMessage := midi.ProgramChange(uint8(channelValue), uint8(programValue))
	wrappedPayload.Payload = payloadMessage
	return wrappedPayload, nil
}

func (mpcc *MIDIProgramChangeCreate) Type() string {
	return mpcc.config.Type
}

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type:  "midi.program_change.create",
		Title: "Create MIDI Prgoram Change Message",
		ParamsSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"channel": {
					Title: "Channel",
					Type:  "string",
				},
				"program": {
					Title: "Program",
					Type:  "string",
				},
			},
			Required:             []string{"type", "channel", "program"},
			AdditionalProperties: &jsonschema.Schema{Not: &jsonschema.Schema{}},
		},
		New: func(config config.ProcessorConfig) (Processor, error) {
			params := config.Params

			channelString, err := params.GetString("channel")
			if err != nil {
				return nil, fmt.Errorf("midi.program_change.create channel error: %w", err)
			}

			channelTemplate, err := template.New("channel").Parse(channelString)

			if err != nil {
				return nil, err
			}

			programString, err := params.GetString("program")
			if err != nil {
				return nil, fmt.Errorf("midi.program_change.create program error: %w", err)
			}

			programTemplate, err := template.New("program").Parse(programString)

			if err != nil {
				return nil, err
			}
			return &MIDIProgramChangeCreate{
				config:  config,
				Channel: channelTemplate,
				Program: programTemplate,
			}, nil
		},
	})
}
