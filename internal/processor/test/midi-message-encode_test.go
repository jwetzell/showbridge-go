//go:build cgo

package processor_test

import (
	"slices"
	"testing"

	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/processor"
	"gitlab.com/gomidi/midi/v2"
)

func TestGoodMIDIMessageEncode(t *testing.T) {
	midiMessageEncoder := processor.MIDIMessageEncode{}
	tests := []struct {
		name     string
		payload  any
		expected []byte
	}{
		{
			name:     "note on message",
			payload:  midi.NoteOn(1, 60, 127),
			expected: []byte{0x91, 0x3c, 0x7f},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := midiMessageEncoder.Process(t.Context(), common.WrappedPayload{Payload: test.payload})
			if err != nil {
				t.Fatalf("midi.message.encode processing failed: %s", err)
			}

			gotBytes, ok := got.Payload.([]byte)
			if !ok {
				t.Fatalf("midi.message.encode returned a %T payload: %+v", got, got)
			}

			if !slices.Equal(gotBytes, test.expected) {
				t.Fatalf("midi.message.encode got %+v, expected %+v", got, test.expected)
			}
		})
	}
}

func TestBadMIDIMessageEncode(t *testing.T) {
	midiMessageEncoder := processor.MIDIMessageEncode{}
	tests := []struct {
		name        string
		payload     any
		errorString string
	}{
		{
			name:        "non-midi message input",
			payload:     []byte{0x68, 0x65, 0x6c, 0x6c, 0x6f},
			errorString: "midi.message.encode processor only accepts a midi.Message",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := midiMessageEncoder.Process(t.Context(), common.WrappedPayload{Payload: test.payload})

			if err == nil {
				t.Fatalf("midi.message.encode expected to fail but got payload: %+v", got)
			}
			if err.Error() != test.errorString {
				t.Fatalf("midi.message.encode got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
