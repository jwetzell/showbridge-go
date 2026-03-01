//go:build cgo

package processor

import (
	"bytes"
	"context"
	"fmt"
	"strconv"
	"text/template"

	"github.com/jwetzell/showbridge-go/internal/config"
	"gitlab.com/gomidi/midi/v2"
)

// TODO(jwetzell): support using numbers in config file treated as hardcoded values
type MIDIMessageCreate struct {
	config      config.ProcessorConfig
	ProcessFunc func(ctx context.Context, payload any) (any, error)
}

func (mmc *MIDIMessageCreate) Process(ctx context.Context, payload any) (any, error) {
	return mmc.ProcessFunc(ctx, payload)
}

func (mmc *MIDIMessageCreate) Type() string {
	return mmc.config.Type
}

func newMidiNoteOnCreate(config config.ProcessorConfig) (Processor, error) {

	params := config.Params

	channelString, err := params.GetString("channel")
	if err != nil {
		return nil, fmt.Errorf("midi.message.create channel error: %w", err)
	}

	channelTemplate, err := template.New("channel").Parse(channelString)

	if err != nil {
		return nil, err
	}

	noteString, err := params.GetString("note")
	if err != nil {
		return nil, fmt.Errorf("midi.message.create note error: %w", err)
	}

	noteTemplate, err := template.New("note").Parse(noteString)

	if err != nil {
		return nil, err
	}

	velocityString, err := params.GetString("velocity")
	if err != nil {
		return nil, fmt.Errorf("midi.message.create velocity error: %w", err)
	}

	velocityTemplate, err := template.New("velocity").Parse(velocityString)

	if err != nil {
		return nil, err
	}

	return &MIDIMessageCreate{config: config, ProcessFunc: func(ctx context.Context, payload any) (any, error) {

		var channelBuffer bytes.Buffer
		err := channelTemplate.Execute(&channelBuffer, payload)

		if err != nil {
			return nil, err
		}

		channelValue, err := strconv.ParseUint(channelBuffer.String(), 10, 8)

		var noteBuffer bytes.Buffer
		err = noteTemplate.Execute(&noteBuffer, payload)

		if err != nil {
			return nil, err
		}

		noteValue, err := strconv.ParseUint(noteBuffer.String(), 10, 8)

		var velocityBuffer bytes.Buffer
		err = velocityTemplate.Execute(&velocityBuffer, payload)

		if err != nil {
			return nil, err
		}

		velocityValue, err := strconv.ParseUint(velocityBuffer.String(), 10, 8)
		payloadMessage := midi.NoteOn(uint8(channelValue), uint8(noteValue), uint8(velocityValue))
		return payloadMessage, nil
	}}, nil
}

func newMidiNoteOffCreate(config config.ProcessorConfig) (Processor, error) {

	params := config.Params

	channelString, err := params.GetString("channel")
	if err != nil {
		return nil, fmt.Errorf("midi.message.create channel error: %w", err)
	}

	channelTemplate, err := template.New("channel").Parse(channelString)

	if err != nil {
		return nil, err
	}

	noteString, err := params.GetString("note")
	if err != nil {
		return nil, fmt.Errorf("midi.message.create note error: %w", err)
	}

	noteTemplate, err := template.New("note").Parse(noteString)

	if err != nil {
		return nil, err
	}

	velocityString, err := params.GetString("velocity")
	if err != nil {
		return nil, fmt.Errorf("midi.message.create velocity error: %w", err)
	}

	velocityTemplate, err := template.New("velocity").Parse(velocityString)

	if err != nil {
		return nil, err
	}

	return &MIDIMessageCreate{config: config, ProcessFunc: func(ctx context.Context, payload any) (any, error) {

		var channelBuffer bytes.Buffer
		err := channelTemplate.Execute(&channelBuffer, payload)

		if err != nil {
			return nil, err
		}

		channelValue, err := strconv.ParseUint(channelBuffer.String(), 10, 8)

		var noteBuffer bytes.Buffer
		err = noteTemplate.Execute(&noteBuffer, payload)

		if err != nil {
			return nil, err
		}

		noteValue, err := strconv.ParseUint(noteBuffer.String(), 10, 8)

		var velocityBuffer bytes.Buffer
		err = velocityTemplate.Execute(&velocityBuffer, payload)

		if err != nil {
			return nil, err
		}

		velocityValue, err := strconv.ParseUint(velocityBuffer.String(), 10, 8)

		payloadMessage := midi.NoteOffVelocity(uint8(channelValue), uint8(noteValue), uint8(velocityValue))
		return payloadMessage, nil
	}}, nil
}

func newMidiControlChangeCreate(config config.ProcessorConfig) (Processor, error) {

	params := config.Params

	channelString, err := params.GetString("channel")
	if err != nil {
		return nil, fmt.Errorf("midi.message.create channel error: %w", err)
	}

	channelTemplate, err := template.New("channel").Parse(channelString)

	if err != nil {
		return nil, err
	}

	controlString, err := params.GetString("control")
	if err != nil {
		return nil, fmt.Errorf("midi.message.create control error: %w", err)
	}

	controlTemplate, err := template.New("control").Parse(controlString)

	if err != nil {
		return nil, err
	}

	valueString, err := params.GetString("value")
	if err != nil {
		return nil, fmt.Errorf("midi.message.create value error: %w", err)
	}

	valueTemplate, err := template.New("value").Parse(valueString)

	if err != nil {
		return nil, err
	}

	return &MIDIMessageCreate{config: config, ProcessFunc: func(ctx context.Context, payload any) (any, error) {

		var channelBuffer bytes.Buffer
		err := channelTemplate.Execute(&channelBuffer, payload)

		if err != nil {
			return nil, err
		}

		channelValue, err := strconv.ParseUint(channelBuffer.String(), 10, 8)

		var controlBuffer bytes.Buffer
		err = controlTemplate.Execute(&controlBuffer, payload)

		if err != nil {
			return nil, err
		}

		controlValue, err := strconv.ParseUint(controlBuffer.String(), 10, 8)

		var valueBuffer bytes.Buffer
		err = valueTemplate.Execute(&valueBuffer, payload)

		if err != nil {
			return nil, err
		}

		valueValue, err := strconv.ParseUint(valueBuffer.String(), 10, 8)

		payloadMessage := midi.ControlChange(uint8(channelValue), uint8(controlValue), uint8(valueValue))
		return payloadMessage, nil
	}}, nil
}

func newMidiProgramChangeCreate(config config.ProcessorConfig) (Processor, error) {

	params := config.Params

	channelString, err := params.GetString("channel")
	if err != nil {
		return nil, fmt.Errorf("midi.message.create channel error: %w", err)
	}

	channelTemplate, err := template.New("channel").Parse(channelString)

	if err != nil {
		return nil, err
	}

	programString, err := params.GetString("program")
	if err != nil {
		return nil, fmt.Errorf("midi.message.create program error: %w", err)
	}

	programTemplate, err := template.New("program").Parse(programString)

	if err != nil {
		return nil, err
	}

	return &MIDIMessageCreate{config: config, ProcessFunc: func(ctx context.Context, payload any) (any, error) {

		var channelBuffer bytes.Buffer
		err := channelTemplate.Execute(&channelBuffer, payload)

		if err != nil {
			return nil, err
		}

		channelValue, err := strconv.ParseUint(channelBuffer.String(), 10, 8)

		var programBuffer bytes.Buffer
		err = programTemplate.Execute(&programBuffer, payload)

		if err != nil {
			return nil, err
		}

		programValue, err := strconv.ParseUint(programBuffer.String(), 10, 8)

		payloadMessage := midi.ProgramChange(uint8(channelValue), uint8(programValue))
		return payloadMessage, nil
	}}, nil
}

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type: "midi.message.create",
		New: func(config config.ProcessorConfig) (Processor, error) {
			params := config.Params

			msgTypeString, err := params.GetString("type")
			if err != nil {
				return nil, fmt.Errorf("midi.message.create type error: %w", err)
			}

			switch msgTypeString {
			case "NoteOn", "noteon", "note_on":
				return newMidiNoteOnCreate(config)
			case "NoteOff", "noteoff", "note_off":
				return newMidiNoteOffCreate(config)
			case "ControlChange", "controlchange", "control_change":
				return newMidiControlChangeCreate(config)
			case "ProgramChange", "programchange", "program_change":
				return newMidiProgramChangeCreate(config)
			default:
				return nil, fmt.Errorf("midi.message.create does not support type %s", msgTypeString)
			}
		},
	})
}
