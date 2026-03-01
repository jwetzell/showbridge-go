package processor_test

import (
	"reflect"
	"testing"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/processor"
	"gitlab.com/gomidi/midi/v2"
)

func TestMIDIMessageFilterFromRegistry(t *testing.T) {
	registration, ok := processor.ProcessorRegistry["midi.message.filter"]
	if !ok {
		t.Fatalf("midi.message.filter processor not registered")
	}

	processorInstance, err := registration.New(config.ProcessorConfig{
		Type: "midi.message.filter",
		Params: map[string]any{
			"type": "NoteOn",
		},
	})

	if err != nil {
		t.Fatalf("failed to create midi.message.filter processor: %s", err)
	}

	if processorInstance.Type() != "midi.message.filter" {
		t.Fatalf("midi.message.filter processor has wrong type: %s", processorInstance.Type())
	}
}

func TestGoodMIDIMessageFilter(t *testing.T) {
	tests := []struct {
		name     string
		params   map[string]any
		payload  midi.Message
		expected midi.Message
	}{
		{
			name:     "matches pattern",
			payload:  midi.NoteOn(1, 60, 127),
			params:   map[string]any{"type": "NoteOn"},
			expected: midi.NoteOn(1, 60, 127),
		},
		{
			name:     "does not match pattern",
			payload:  midi.NoteOn(1, 60, 127),
			params:   map[string]any{"type": "NoteOff"},
			expected: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			registration, ok := processor.ProcessorRegistry["midi.message.filter"]
			if !ok {
				t.Fatalf("midi.message.filter processor not registered")
			}

			processorInstance, err := registration.New(config.ProcessorConfig{
				Type:   "midi.message.filter",
				Params: test.params,
			})

			if err != nil {
				t.Fatalf("midi.message.filter failed to create processor: %s", err)
			}

			got, err := processorInstance.Process(t.Context(), test.payload)

			if err != nil {
				t.Fatalf("midi.message.filter failed: %s", err)
			}

			if test.expected == nil {
				if got != nil {
					t.Fatalf("midi.message.filter got %+v, expected nil", got)
				}
				return
			}

			gotMIDIMessage, ok := got.(midi.Message)
			if !ok {
				t.Fatalf("midi.message.filter returned a %T payload: %s", got, got)
			}

			if !reflect.DeepEqual(gotMIDIMessage, test.expected) {
				t.Fatalf("midi.message.filter got %+v, expected %+v", gotMIDIMessage, test.expected)
			}
		})
	}
}

func TestBadMIDIMessageFilter(t *testing.T) {
	tests := []struct {
		name        string
		params      map[string]any
		payload     any
		errorString string
	}{
		{
			name:        "no type param",
			params:      map[string]any{},
			payload:     midi.NoteOn(1, 60, 127),
			errorString: "midi.message.filter type error: not found",
		},
		{
			name: "non-string type param",
			params: map[string]any{
				"type": 123,
			},
			payload:     "hello",
			errorString: "midi.message.filter type error: not a string",
		},
		{
			name: "non-midi message input",
			params: map[string]any{
				"type": "NoteOn",
			},
			payload:     []byte{0x68, 0x65, 0x6c, 0x6c, 0x6f},
			errorString: "midi.message.filter processor only accepts a midi.Message",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			registration, ok := processor.ProcessorRegistry["midi.message.filter"]
			if !ok {
				t.Fatalf("midi.message.filter processor not registered")
			}

			processorInstance, err := registration.New(config.ProcessorConfig{
				Type:   "midi.message.filter",
				Params: test.params,
			})

			if err != nil {
				if test.errorString != err.Error() {
					t.Fatalf("midi.message.filter got error '%s', expected '%s'", err.Error(), test.errorString)
				}
				return
			}

			got, err := processorInstance.Process(t.Context(), test.payload)

			if err == nil {
				t.Fatalf("midi.message.filter expected to fail but got payload: %s", got)
			}

			if err.Error() != test.errorString {
				t.Fatalf("midi.message.filter got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
