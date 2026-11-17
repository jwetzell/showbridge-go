//go:build cgo

package processor_test

import (
	"testing"

	"github.com/jwetzell/showbridge-go/config"
	"github.com/jwetzell/showbridge-go/processor"
	"github.com/jwetzell/showbridge-go/routing"
	"github.com/stretchr/testify/require"
	"gitlab.com/gomidi/midi/v2"
)

func TestGoodMIDIMessageDecode(t *testing.T) {
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
			registration, ok := processor.GetProcessorRegistration("midi.message.decode")
			require.True(t, ok, "processor registration not found")

			processorInstance, err := registration.New(config.ProcessorConfig{
				Type: "midi.message.decode",
			})

			got, err := EvaluateProcessFunc(processorInstance, t.Context(), routing.WrappedPayload{Payload: test.payload}, nil, nil)
			require.NoError(t, err, "midi.message.decode processing failed")

			gotMessage, ok := got.Payload.(midi.Message)
			require.True(t, ok, "midi.message.decode returned a non-midi.Message payload")

			require.Equal(t, test.expected, gotMessage, "midi.message.decode returned unexpected message")
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
			payload:     12345,
			errorString: "midi.message.decode processor only accepts a []byte",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			got, err := EvaluateProcessFunc(processorInstance, t.Context(), routing.WrappedPayload{Payload: test.payload}, nil, nil)

			if err == nil {
				t.Fatalf("midi.message.decode expected to fail but succeeded, got: %v", got)

			}
			if err.Error() != test.errorString {
				t.Fatalf("midi.message.decode got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
