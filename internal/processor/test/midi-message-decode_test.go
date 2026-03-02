package processor_test

import (
	"reflect"
	"testing"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/processor"
	"gitlab.com/gomidi/midi/v2"
)

func TestMIDIMessageDecodeFromRegistry(t *testing.T) {
	registration, ok := processor.ProcessorRegistry["midi.message.decode"]
	if !ok {
		t.Fatalf("midi.message.decode processor not registered")
	}

	processorInstance, err := registration.New(config.ProcessorConfig{
		Type: "midi.message.decode",
	})

	if err != nil {
		t.Fatalf("failed to create midi.message.decode processor: %s", err)
	}

	if processorInstance.Type() != "midi.message.decode" {
		t.Fatalf("midi.message.decode processor has wrong type: %s", processorInstance.Type())
	}
}

func TestGoodMIDIMessageDecode(t *testing.T) {
	processorInstance := &processor.MIDIMessageDecode{}
	tests := []struct {
		name     string
		payload  any
		expected any
	}{
		{
			name:     "note on message",
			payload:  []byte{0x90, 0x40, 0x7F},
			expected: midi.NoteOn(0, 64, 127),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := processorInstance.Process(t.Context(), test.payload)
			if err != nil {
				t.Fatalf("midi.message.decode failed: %s", err)
			}

			gotMessage, ok := got.(midi.Message)
			if !ok {
				t.Fatalf("midi.message.decode returned a %T payload: %s", got, got)
			}

			if !reflect.DeepEqual(gotMessage, test.expected) {
				t.Fatalf("midi.message.decode got %+v, expected %+v", gotMessage, test.expected)
			}
		})
	}
}

func TestBadMIDIMessageDecode(t *testing.T) {
	processorInstance := &processor.MIDIMessageDecode{}
	tests := []struct {
		name        string
		payload     any
		errorString string
	}{
		{
			name:        "non-byte slice payload",
			payload:     "12345",
			errorString: "midi.message.decode processor only accepts a []byte",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			got, err := processorInstance.Process(t.Context(), test.payload)

			if err == nil {
				t.Fatalf("midi.message.decode expected to fail but succeeded, got: %v", got)

			}
			if err.Error() != test.errorString {
				t.Fatalf("midi.message.decode got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
