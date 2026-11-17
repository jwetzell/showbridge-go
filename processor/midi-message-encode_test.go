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

func TestGoodMIDIMessageEncode(t *testing.T) {
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
			registration, ok := processor.GetProcessorRegistration("midi.message.encode")
			require.True(t, ok, "processor registration not found")

			processorInstance, err := registration.New(config.ProcessorConfig{
				Type: "midi.message.encode",
			})
			got, err := EvaluateProcessFunc(processorInstance, t.Context(), routing.WrappedPayload{Payload: test.payload}, nil, nil)
			require.NoError(t, err, "midi.message.encode processing failed")

			gotBytes, ok := got.Payload.([]byte)
			require.True(t, ok, "midi.message.encode returned a non-[]byte payload")

			require.Equal(t, test.expected, gotBytes, "midi.message.encode returned unexpected bytes")
		})
	}
}

func TestBadMIDIMessageEncode(t *testing.T) {
	midiMessageEncoder := &processor.MIDIMessageEncode{}
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
			got, err := EvaluateProcessFunc(midiMessageEncoder, t.Context(), routing.WrappedPayload{Payload: test.payload}, nil, nil)

			if err == nil {
				t.Fatalf("midi.message.encode expected to fail but got payload: %+v", got)
			}
			if err.Error() != test.errorString {
				t.Fatalf("midi.message.encode got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
