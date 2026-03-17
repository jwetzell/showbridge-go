package processor_test

import (
	"reflect"
	"testing"

	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/processor"
	"gitlab.com/gomidi/midi/v2"
)

func TestMIDIMessageUnpackFromRegistry(t *testing.T) {
	registration, ok := processor.ProcessorRegistry["midi.message.unpack"]
	if !ok {
		t.Fatalf("midi.message.unpack processor not registered")
	}

	processorInstance, err := registration.New(config.ProcessorConfig{
		Type: "midi.message.unpack",
	})

	if err != nil {
		t.Fatalf("failed to create midi.message.unpack processor: %s", err)
	}

	if processorInstance.Type() != "midi.message.unpack" {
		t.Fatalf("midi.message.unpack processor has wrong type: %s", processorInstance.Type())
	}
}

func TestGoodMIDIMessageUnpack(t *testing.T) {
	processorInstance := &processor.MIDIMessageUnpack{}
	tests := []struct {
		name     string
		expected any
		params   map[string]any
		payload  any
	}{
		{
			name:    "note on",
			payload: midi.NoteOn(1, 60, 127),
			expected: processor.MIDINoteOn{
				Channel:  1,
				Note:     60,
				Velocity: 127,
			},
		},
		{
			name:    "note off",
			payload: midi.NoteOffVelocity(1, 60, 127),
			expected: processor.MIDINoteOff{
				Channel:  1,
				Note:     60,
				Velocity: 127,
			},
		},
		{
			name:    "control change",
			payload: midi.ControlChange(1, 64, 127),
			expected: processor.MIDIControlChange{
				Channel: 1,
				Control: 64,
				Value:   127,
			},
		},
		{
			name:    "program change",
			payload: midi.ProgramChange(1, 10),
			expected: processor.MIDIProgramChange{
				Channel: 1,
				Program: 10,
			},
		},
		{
			name:    "pitch bend",
			payload: midi.Message{0xE1, 0x00, 0x40},
			expected: processor.MIDIPitchBend{
				Channel:  1,
				Relative: 0,
				Absolute: 8192,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := processorInstance.Process(t.Context(), common.GetWrappedPayload(t.Context(), test.payload))
			if err != nil {
				t.Fatalf("midi.message.unpack processing failed: %s", err)
			}

			if !reflect.DeepEqual(got.Payload, test.expected) {
				t.Fatalf("midi.message.unpack got %+v (%T), expected %+v (%T)", got.Payload, got.Payload, test.expected, test.expected)
			}
		})
	}
}

func TestBadMIDIMessageUnpack(t *testing.T) {
	tests := []struct {
		name        string
		params      map[string]any
		payload     any
		errorString string
	}{
		{
			name:        "payload not a MIDI message",
			payload:     "not a MIDI message",
			errorString: "midi.message.unpack processor only accepts a midi.Message",
			params:      nil,
		},
		{
			name:        "unsupported MIDI message type",
			payload:     midi.Message{0x00, 0x00, 0x00},
			errorString: "midi.message.unpack message type not supported UnknownType",
			params:      nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			registration, ok := processor.ProcessorRegistry["midi.message.unpack"]
			if !ok {
				t.Fatalf("midi.message.unpack processor not registered")
			}

			processorInstance, err := registration.New(config.ProcessorConfig{
				Type:   "midi.message.unpack",
				Params: test.params,
			})

			if err != nil {
				if test.errorString != err.Error() {
					t.Fatalf("midi.message.unpack got error '%s', expected '%s'", err.Error(), test.errorString)
				}
				return
			}

			got, err := processorInstance.Process(t.Context(), common.GetWrappedPayload(t.Context(), test.payload))

			if err == nil {
				t.Fatalf("midi.message.unpack expected to fail but succeeded, got: %v", got)
			}

			if err.Error() != test.errorString {
				t.Fatalf("midi.message.unpack got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
