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

type MIDIMessageCreate struct {
	config      config.ProcessorConfig
	ProcessFunc func(ctx context.Context, payload any) (any, error)
}

func (mmd *MIDIMessageCreate) Process(ctx context.Context, payload any) (any, error) {
	return mmd.ProcessFunc(ctx, payload)
}

func (mmd *MIDIMessageCreate) Type() string {
	return mmd.config.Type
}

func newMidiNoteOnCreate(config config.ProcessorConfig) (Processor, error) {

	params := config.Params

	channel, ok := params["channel"]

	if !ok {
		return nil, fmt.Errorf("midi.message.create NoteOn requires a channel parameter")
	}

	channelString, ok := channel.(string)

	if !ok {
		return nil, fmt.Errorf("midi.message.create NoteOn channel must be a string")
	}

	channelTemplate, err := template.New("channel").Parse(channelString)

	if err != nil {
		return nil, err
	}

	note, ok := params["note"]

	if !ok {
		return nil, fmt.Errorf("midi.message.create NoteOn requires a note parameter")
	}

	noteString, ok := note.(string)

	if !ok {
		return nil, fmt.Errorf("midi.message.create NoteOn note must be a string")
	}

	noteTemplate, err := template.New("note").Parse(noteString)

	if err != nil {
		return nil, err
	}

	velocity, ok := params["velocity"]

	if !ok {
		return nil, fmt.Errorf("midi.message.create NoteOn requires a velocity parameter")
	}

	velocityString, ok := velocity.(string)

	if !ok {
		return nil, fmt.Errorf("midi.message.create NoteOn velocity must be a string")
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
		return nil, fmt.Errorf("midi.message.create NoteOn requires a channel parameter")
	}

	channelString, ok := channel.(string)

	if !ok {
		return nil, fmt.Errorf("midi.message.create NoteOn channel must be a string")
	}

	channelTemplate, err := template.New("channel").Parse(channelString)

	if err != nil {
		return nil, err
	}

	note, ok := params["note"]

	if !ok {
		return nil, fmt.Errorf("midi.message.create NoteOn requires a note parameter")
	}

	noteString, ok := note.(string)

	if !ok {
		return nil, fmt.Errorf("midi.message.create NoteOn note must be a string")
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

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type: "midi.message.create",
		New: func(config config.ProcessorConfig) (Processor, error) {
			params := config.Params

			msgType, ok := params["type"]

			if !ok {
				return nil, fmt.Errorf("midi.message.create requires a type parameter")
			}

			msgTypeString, ok := msgType.(string)

			if !ok {
				return nil, fmt.Errorf("midi.message.create type parameter must be a string")
			}

			switch msgTypeString {
			case "NoteOn":
				return newMidiNoteOnCreate(config)
			case "NoteOff":
				return newMidiNoteOffCreate(config)
			default:
				return nil, fmt.Errorf("midi.message.create does not support type %s", msgTypeString)
			}
		},
	})
}
