package processor_test

import (
	"slices"
	"testing"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/processor"
	"gitlab.com/gomidi/midi/v2"
)

func TestMIDIMessageEncodeFromRegistry(t *testing.T) {
	registration, ok := processor.ProcessorRegistry["midi.message.encode"]
	if !ok {
		t.Fatalf("midi.message.encode processor not registered")
	}

	processorInstance, err := registration.New(config.ProcessorConfig{
		Type: "midi.message.encode",
	})

	if err != nil {
		t.Fatalf("failed to create midi.message.encode processor: %s", err)
	}

	if processorInstance.Type() != "midi.message.encode" {
		t.Fatalf("midi.message.encode processor has wrong type: %s", processorInstance.Type())
	}
}

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
			got, err := midiMessageEncoder.Process(t.Context(), test.payload)

			gotBytes, ok := got.([]byte)
			if !ok {
				t.Fatalf("midi.message.encode returned a %T payload: %s", got, got)
			}
			if err != nil {
				t.Fatalf("midi.message.encode failed: %s", err)
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
			got, err := midiMessageEncoder.Process(t.Context(), test.payload)

			if err == nil {
				t.Fatalf("midi.message.encode expected to fail but got payload: %s", got)
			}
			if err.Error() != test.errorString {
				t.Fatalf("midi.message.encode got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
