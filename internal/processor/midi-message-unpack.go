//go:build cgo

package processor

import (
	"context"
	"errors"
	"fmt"

	"github.com/jwetzell/showbridge-go/internal/config"
	"gitlab.com/gomidi/midi/v2"
)

type MIDIMessageUnpack struct {
	config config.ProcessorConfig
}

type MIDINoteOn struct {
	Channel  uint8
	Note     uint8
	Velocity uint8
}

type MIDINoteOff struct {
	Channel  uint8
	Note     uint8
	Velocity uint8
}

type MIDIControlChange struct {
	Channel uint8
	Control uint8
	Value   uint8
}

type MIDIProgramChange struct {
	Channel uint8
	Program uint8
}

type MIDIPitchBend struct {
	Channel  uint8
	Relative int16
	Absolute uint16
}

func (mmu *MIDIMessageUnpack) Process(ctx context.Context, payload any) (any, error) {
	payloadMidi, ok := payload.(midi.Message)

	if !ok {
		return nil, errors.New("midi.message.unpack processor only accepts a midi.Message")
	}

	switch payloadMidi.Type() {
	case midi.NoteOnMsg:
		noteOnMsg := MIDINoteOn{}
		payloadMidi.GetNoteOn(&noteOnMsg.Channel, &noteOnMsg.Note, &noteOnMsg.Velocity)
		return noteOnMsg, nil
	case midi.NoteOffMsg:
		noteOffMsg := MIDINoteOff{}
		payloadMidi.GetNoteOff(&noteOffMsg.Channel, &noteOffMsg.Note, &noteOffMsg.Velocity)
		return noteOffMsg, nil
	case midi.ControlChangeMsg:
		controlChangeMsg := MIDIControlChange{}
		payloadMidi.GetControlChange(&controlChangeMsg.Channel, &controlChangeMsg.Control, &controlChangeMsg.Value)
		return controlChangeMsg, nil
	case midi.ProgramChangeMsg:
		programChangeMsg := MIDIProgramChange{}
		payloadMidi.GetProgramChange(&programChangeMsg.Channel, &programChangeMsg.Program)
		return programChangeMsg, nil
	case midi.PitchBendMsg:
		pitchBendMsg := MIDIPitchBend{}
		payloadMidi.GetPitchBend(&pitchBendMsg.Channel, &pitchBendMsg.Relative, &pitchBendMsg.Absolute)
		return pitchBendMsg, nil
	default:
		return nil, fmt.Errorf("midi.message.unpack message type not supported %v", payloadMidi.Type())
	}
}

func (mmu *MIDIMessageUnpack) Type() string {
	return mmu.config.Type
}

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type: "midi.message.unpack",
		New: func(config config.ProcessorConfig) (Processor, error) {
			return &MIDIMessageUnpack{config: config}, nil
		},
	})
}
