//go:build cgo

package processor

import (
	"context"
	"errors"
	"fmt"

	"github.com/jwetzell/showbridge-go/internal/common"
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

func (mmu *MIDIMessageUnpack) Process(ctx context.Context, wrappedPayload common.WrappedPayload) (common.WrappedPayload, error) {
	payload := wrappedPayload.Payload
	payloadMidi, ok := common.GetAnyAs[midi.Message](payload)

	if !ok {
		wrappedPayload.End = true
		return wrappedPayload, errors.New("midi.message.unpack processor only accepts a midi.Message")
	}

	switch payloadMidi.Type() {
	case midi.NoteOnMsg:
		noteOnMsg := MIDINoteOn{}
		payloadMidi.GetNoteOn(&noteOnMsg.Channel, &noteOnMsg.Note, &noteOnMsg.Velocity)
		wrappedPayload.Payload = noteOnMsg
		return wrappedPayload, nil
	case midi.NoteOffMsg:
		noteOffMsg := MIDINoteOff{}
		payloadMidi.GetNoteOff(&noteOffMsg.Channel, &noteOffMsg.Note, &noteOffMsg.Velocity)
		wrappedPayload.Payload = noteOffMsg
		return wrappedPayload, nil
	case midi.ControlChangeMsg:
		controlChangeMsg := MIDIControlChange{}
		payloadMidi.GetControlChange(&controlChangeMsg.Channel, &controlChangeMsg.Control, &controlChangeMsg.Value)
		wrappedPayload.Payload = controlChangeMsg
		return wrappedPayload, nil
	case midi.ProgramChangeMsg:
		programChangeMsg := MIDIProgramChange{}
		payloadMidi.GetProgramChange(&programChangeMsg.Channel, &programChangeMsg.Program)
		wrappedPayload.Payload = programChangeMsg
		return wrappedPayload, nil
	case midi.PitchBendMsg:
		pitchBendMsg := MIDIPitchBend{}
		payloadMidi.GetPitchBend(&pitchBendMsg.Channel, &pitchBendMsg.Relative, &pitchBendMsg.Absolute)
		wrappedPayload.Payload = pitchBendMsg
		return wrappedPayload, nil
	default:
		wrappedPayload.End = true
		return wrappedPayload, fmt.Errorf("midi.message.unpack message type not supported %v", payloadMidi.Type())
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
