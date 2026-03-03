package processor_test

import (
	"reflect"
	"testing"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/processor"
	"gitlab.com/gomidi/midi/v2"
)

func TestMIDIMessageCreateFromRegistry(t *testing.T) {
	registration, ok := processor.ProcessorRegistry["midi.message.create"]
	if !ok {
		t.Fatalf("midi.message.create processor not registered")
	}

	processorInstance, err := registration.New(config.ProcessorConfig{
		Type: "midi.message.create",
		Params: map[string]any{
			"type":     "note_on",
			"channel":  "1",
			"note":     "60",
			"velocity": "100",
		},
	})

	if err != nil {
		t.Fatalf("failed to create midi.message.create processor: %s", err)
	}

	if processorInstance.Type() != "midi.message.create" {
		t.Fatalf("midi.message.create processor has wrong type: %s", processorInstance.Type())
	}
}

func TestGoodMIDIMessageCreate(t *testing.T) {

	tests := []struct {
		name     string
		params   map[string]any
		payload  any
		expected any
	}{
		{
			name: "note_on message",
			params: map[string]any{
				"type":     "note_on",
				"channel":  "1",
				"note":     "60",
				"velocity": "100",
			},
			payload:  "test",
			expected: midi.NoteOn(1, 60, 100),
		},
		{
			name: "note_off message",
			params: map[string]any{
				"type":     "note_off",
				"channel":  "1",
				"note":     "60",
				"velocity": "100",
			},
			payload:  "test",
			expected: midi.NoteOffVelocity(1, 60, 100),
		},
		{
			name: "control_change message",
			params: map[string]any{
				"type":    "control_change",
				"channel": "1",
				"control": "64",
				"value":   "127",
			},
			payload:  "test",
			expected: midi.ControlChange(1, 64, 127),
		},
		{
			name: "program_change message",
			params: map[string]any{
				"type":    "program_change",
				"channel": "1",
				"program": "10",
			},
			payload:  "test",
			expected: midi.ProgramChange(1, 10),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			registration, ok := processor.ProcessorRegistry["midi.message.create"]
			if !ok {
				t.Fatalf("midi.message.create processor not registered")
			}

			processorInstance, err := registration.New(config.ProcessorConfig{
				Type:   "midi.message.create",
				Params: test.params,
			})

			if err != nil {
				t.Fatalf("midi.message.create failed to create processor: %s", err)
			}

			got, err := processorInstance.Process(t.Context(), test.payload)
			if err != nil {
				t.Fatalf("midi.message.create processing failed: %s", err)
			}

			gotMessage, ok := got.(midi.Message)
			if !ok {
				t.Fatalf("midi.message.create returned a %T payload: %s", got, got)
			}

			if !reflect.DeepEqual(gotMessage, test.expected) {
				t.Fatalf("midi.message.create got %v, expected %v", gotMessage, test.expected)
			}
		})
	}
}

func TestBadMIDIMessageCreate(t *testing.T) {

	tests := []struct {
		name        string
		params      map[string]any
		payload     any
		errorString string
	}{
		{
			name:        "no type parameter",
			params:      map[string]any{},
			payload:     "test",
			errorString: "midi.message.create type error: not found",
		},
		{
			name: "non-string type parameter",
			params: map[string]any{
				"type": 1,
			},
			payload:     "test",
			errorString: "midi.message.create type error: not a string",
		},
		{
			name: "unknown type parameter",
			params: map[string]any{
				"type": "asdf",
			},
			payload:     "test",
			errorString: "midi.message.create does not support type asdf",
		},
		{
			name: "note_on message no channel",
			params: map[string]any{
				"type":     "note_on",
				"note":     "60",
				"velocity": "100",
			},
			payload:     "test",
			errorString: "midi.message.create channel error: not found",
		},
		{
			name: "note_on message no note",
			params: map[string]any{
				"type":     "note_on",
				"channel":  "1",
				"velocity": "100",
			},
			payload:     "test",
			errorString: "midi.message.create note error: not found",
		},
		{
			name: "note_on message no velocity",
			params: map[string]any{
				"type":    "note_on",
				"channel": "1",
				"note":    "60",
			},
			payload:     "test",
			errorString: "midi.message.create velocity error: not found",
		},
		{
			name: "note_off message no channel",
			params: map[string]any{
				"type":     "note_off",
				"note":     "60",
				"velocity": "100",
			},
			payload:     "test",
			errorString: "midi.message.create channel error: not found",
		},
		{
			name: "note_off message no note",
			params: map[string]any{
				"type":     "note_off",
				"channel":  "1",
				"velocity": "100",
			},
			payload:     "test",
			errorString: "midi.message.create note error: not found",
		},
		{
			name: "note_off message no velocity",
			params: map[string]any{
				"type":    "note_off",
				"channel": "1",
				"note":    "60",
			},
			payload:     "test",
			errorString: "midi.message.create velocity error: not found",
		},
		{
			name: "control_change no channel",
			params: map[string]any{
				"type":    "control_change",
				"control": "64",
				"value":   "127",
			},
			payload:     "test",
			errorString: "midi.message.create channel error: not found",
		},
		{
			name: "control_change no control",
			params: map[string]any{
				"type":    "control_change",
				"channel": "1",
				"value":   "127",
			},
			payload:     "test",
			errorString: "midi.message.create control error: not found",
		},
		{
			name: "control_change no value",
			params: map[string]any{
				"type":    "control_change",
				"channel": "1",
				"control": "64",
			},
			payload:     "test",
			errorString: "midi.message.create value error: not found",
		},
		{
			name: "program_change no channel",
			params: map[string]any{
				"type":    "program_change",
				"program": "64",
			},
			payload:     "test",
			errorString: "midi.message.create channel error: not found",
		},
		{
			name: "program_change no program",
			params: map[string]any{
				"type":    "program_change",
				"channel": "1",
			},
			payload:     "test",
			errorString: "midi.message.create program error: not found",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			registration, ok := processor.ProcessorRegistry["midi.message.create"]
			if !ok {
				t.Fatalf("midi.message.create processor not registered")
			}

			processorInstance, err := registration.New(config.ProcessorConfig{
				Type:   "midi.message.create",
				Params: test.params,
			})

			if err != nil {
				if test.errorString != err.Error() {
					t.Fatalf("string.create got error '%s', expected '%s'", err.Error(), test.errorString)
				}
				return
			}

			got, err := processorInstance.Process(t.Context(), test.payload)

			if err == nil {
				t.Fatalf("midi.message.create expected to fail but succeeded, got: %v", got)
			}

			if err.Error() != test.errorString {
				t.Fatalf("midi.message.create got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
