//go:build cgo

package processor

import (
	"bytes"
	"context"
	"errors"
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

	channel, ok := params["channel"]

	if !ok {
		return nil, errors.New("midi.message.create NoteOn requires a channel parameter")
	}

	channelString, ok := channel.(string)

	if !ok {
		return nil, errors.New("midi.message.create NoteOn channel must be a string")
	}

	channelTemplate, err := template.New("channel").Parse(channelString)

	if err != nil {
		return nil, err
	}

	note, ok := params["note"]

	if !ok {
		return nil, errors.New("midi.message.create NoteOn requires a note parameter")
	}

	noteString, ok := note.(string)

	if !ok {
		return nil, errors.New("midi.message.create NoteOn note must be a string")
	}

	noteTemplate, err := template.New("note").Parse(noteString)

	if err != nil {
		return nil, err
	}

	velocity, ok := params["velocity"]

	if !ok {
		return nil, errors.New("midi.message.create NoteOn requires a velocity parameter")
	}

	velocityString, ok := velocity.(string)

	if !ok {
		return nil, errors.New("midi.message.create NoteOn velocity must be a string")
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

	channel, ok := params["channel"]

	if !ok {
		return nil, errors.New("midi.message.create NoteOn requires a channel parameter")
	}

	channelString, ok := channel.(string)

	if !ok {
		return nil, errors.New("midi.message.create NoteOn channel must be a string")
	}

	channelTemplate, err := template.New("channel").Parse(channelString)

	if err != nil {
		return nil, err
	}

	note, ok := params["note"]

	if !ok {
		return nil, errors.New("midi.message.create NoteOn requires a note parameter")
	}

	noteString, ok := note.(string)

	if !ok {
		return nil, errors.New("midi.message.create NoteOn note must be a string")
	}

	noteTemplate, err := template.New("note").Parse(noteString)

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

		payloadMessage := midi.NoteOff(uint8(channelValue), uint8(noteValue))
		return payloadMessage, nil
	}}, nil
}

func newMidiControlChangeCreate(config config.ProcessorConfig) (Processor, error) {

	params := config.Params

	channel, ok := params["channel"]

	if !ok {
		return nil, errors.New("midi.message.create ControlChange requires a channel parameter")
	}

	channelString, ok := channel.(string)

	if !ok {
		return nil, errors.New("midi.message.create ControlChange channel must be a string")
	}

	channelTemplate, err := template.New("channel").Parse(channelString)

	if err != nil {
		return nil, err
	}

	control, ok := params["control"]

	if !ok {
		return nil, errors.New("midi.message.create ControlChange requires a control parameter")
	}

	controlString, ok := control.(string)

	if !ok {
		return nil, errors.New("midi.message.create ControlChange control must be a string")
	}

	controlTemplate, err := template.New("control").Parse(controlString)

	if err != nil {
		return nil, err
	}

	value, ok := params["value"]

	if !ok {
		return nil, errors.New("midi.message.create ControlChange requires a value parameter")
	}

	valueString, ok := value.(string)

	if !ok {
		return nil, errors.New("midi.message.create ControlChange value must be a string")
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

	channel, ok := params["channel"]

	if !ok {
		return nil, errors.New("midi.message.create ProgramChange requires a channel parameter")
	}

	channelString, ok := channel.(string)

	if !ok {
		return nil, errors.New("midi.message.create ProgramChange channel must be a string")
	}

	channelTemplate, err := template.New("channel").Parse(channelString)

	if err != nil {
		return nil, err
	}

	program, ok := params["program"]

	if !ok {
		return nil, errors.New("midi.message.create ProgramChange requires a program parameter")
	}

	programString, ok := program.(string)

	if !ok {
		return nil, errors.New("midi.message.create ProgramChange program must be a string")
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

			msgType, ok := params["type"]

			if !ok {
				return nil, errors.New("midi.message.create requires a type parameter")
			}

			msgTypeString, ok := msgType.(string)

			if !ok {
				return nil, errors.New("midi.message.create type parameter must be a string")
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
