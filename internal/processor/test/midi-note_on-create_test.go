package processor_test

import (
	"reflect"
	"testing"

	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/processor"
	"gitlab.com/gomidi/midi/v2"
)

func TestMIDINoteOnCreateFromRegistry(t *testing.T) {
	registration, ok := processor.ProcessorRegistry["midi.note_on.create"]
	if !ok {
		t.Fatalf("midi.note_on.create processor not registered")
	}

	processorInstance, err := registration.New(config.ProcessorConfig{
		Type: "midi.note_on.create",
		Params: map[string]any{
			"channel":  "1",
			"note":     "60",
			"velocity": "100",
		},
	})

	if err != nil {
		t.Fatalf("failed to create midi.note_on.create processor: %s", err)
	}

	if processorInstance.Type() != "midi.note_on.create" {
		t.Fatalf("midi.note_on.create processor has wrong type: %s", processorInstance.Type())
	}
}

func TestGoodMIDINoteOnCreate(t *testing.T) {

	tests := []struct {
		name     string
		params   map[string]any
		payload  any
		expected any
	}{
		{
			name: "note_on message",
			params: map[string]any{
				"channel":  "1",
				"note":     "60",
				"velocity": "100",
			},
			payload:  "test",
			expected: midi.NoteOn(1, 60, 100),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			registration, ok := processor.ProcessorRegistry["midi.note_on.create"]
			if !ok {
				t.Fatalf("midi.note_on.create processor not registered")
			}

			processorInstance, err := registration.New(config.ProcessorConfig{
				Type:   "midi.note_on.create",
				Params: test.params,
			})

			if err != nil {
				t.Fatalf("midi.note_on.create failed to create processor: %s", err)
			}

			got, err := processorInstance.Process(t.Context(), common.GetWrappedPayload(t.Context(), test.payload))
			if err != nil {
				t.Fatalf("midi.note_on.create processing failed: %s", err)
			}

			gotMessage, ok := got.Payload.(midi.Message)
			if !ok {
				t.Fatalf("midi.note_on.create returned a %T payload: %+v", got, got)
			}

			if !reflect.DeepEqual(gotMessage, test.expected) {
				t.Fatalf("midi.note_on.create got %v, expected %v", gotMessage, test.expected)
			}
		})
	}
}

func TestBadMIDINoteOnCreate(t *testing.T) {

	tests := []struct {
		name        string
		params      map[string]any
		payload     any
		errorString string
	}{
		{
			name: "note_on message no channel",
			params: map[string]any{
				"note":     "60",
				"velocity": "100",
			},
			payload:     "test",
			errorString: "midi.note_on.create channel error: not found",
		},
		{
			name: "note_on message no note",
			params: map[string]any{
				"channel":  "1",
				"velocity": "100",
			},
			payload:     "test",
			errorString: "midi.note_on.create note error: not found",
		},
		{
			name: "note_on message no velocity",
			params: map[string]any{
				"channel": "1",
				"note":    "60",
			},
			payload:     "test",
			errorString: "midi.note_on.create velocity error: not found",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			registration, ok := processor.ProcessorRegistry["midi.note_on.create"]
			if !ok {
				t.Fatalf("midi.note_on.create processor not registered")
			}

			processorInstance, err := registration.New(config.ProcessorConfig{
				Type:   "midi.note_on.create",
				Params: test.params,
			})

			if err != nil {
				if test.errorString != err.Error() {
					t.Fatalf("midi.note_on.create got error '%s', expected '%s'", err.Error(), test.errorString)
				}
				return
			}

			got, err := processorInstance.Process(t.Context(), common.GetWrappedPayload(t.Context(), test.payload))

			if err == nil {
				t.Fatalf("midi.note_on.create expected to fail but succeeded, got: %v", got)
			}

			if err.Error() != test.errorString {
				t.Fatalf("midi.note_on.create got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
