package processor_test

import (
	"reflect"
	"testing"

	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/processor"
	"gitlab.com/gomidi/midi/v2"
)

func TestMIDINoteOffCreteaFromRegistry(t *testing.T) {
	registration, ok := processor.ProcessorRegistry["midi.note_off.create"]
	if !ok {
		t.Fatalf("midi.note_off.create processor not registered")
	}

	processorInstance, err := registration.New(config.ProcessorConfig{
		Type: "midi.note_off.create",
		Params: map[string]any{
			"channel":  "1",
			"note":     "60",
			"velocity": "100",
		},
	})

	if err != nil {
		t.Fatalf("failed to create midi.note_off.create processor: %s", err)
	}

	if processorInstance.Type() != "midi.note_off.create" {
		t.Fatalf("midi.note_off.create processor has wrong type: %s", processorInstance.Type())
	}
}

func TestGoodMIDINoteOffCretea(t *testing.T) {

	tests := []struct {
		name     string
		params   map[string]any
		payload  any
		expected any
	}{
		{
			name: "note_off message",
			params: map[string]any{
				"channel":  "1",
				"note":     "60",
				"velocity": "100",
			},
			payload:  "test",
			expected: midi.NoteOffVelocity(1, 60, 100),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			registration, ok := processor.ProcessorRegistry["midi.note_off.create"]
			if !ok {
				t.Fatalf("midi.note_off.create processor not registered")
			}

			processorInstance, err := registration.New(config.ProcessorConfig{
				Type:   "midi.note_off.create",
				Params: test.params,
			})

			if err != nil {
				t.Fatalf("midi.note_off.create failed to create processor: %s", err)
			}

			got, err := processorInstance.Process(t.Context(), common.GetWrappedPayload(t.Context(), test.payload))
			if err != nil {
				t.Fatalf("midi.note_off.create processing failed: %s", err)
			}

			gotMessage, ok := got.Payload.(midi.Message)
			if !ok {
				t.Fatalf("midi.note_off.create returned a %T payload: %+v", got, got)
			}

			if !reflect.DeepEqual(gotMessage, test.expected) {
				t.Fatalf("midi.note_off.create got %v, expected %v", gotMessage, test.expected)
			}
		})
	}
}

func TestBadMIDINoteOffCretea(t *testing.T) {

	tests := []struct {
		name        string
		params      map[string]any
		payload     any
		errorString string
	}{
		{
			name: "note_off message no channel",
			params: map[string]any{
				"type":     "note_off",
				"note":     "60",
				"velocity": "100",
			},
			payload:     "test",
			errorString: "midi.note_off.create channel error: not found",
		},
		{
			name: "note_off message no note",
			params: map[string]any{
				"type":     "note_off",
				"channel":  "1",
				"velocity": "100",
			},
			payload:     "test",
			errorString: "midi.note_off.create note error: not found",
		},
		{
			name: "note_off message no velocity",
			params: map[string]any{
				"type":    "note_off",
				"channel": "1",
				"note":    "60",
			},
			payload:     "test",
			errorString: "midi.note_off.create velocity error: not found",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			registration, ok := processor.ProcessorRegistry["midi.note_off.create"]
			if !ok {
				t.Fatalf("midi.note_off.create processor not registered")
			}

			processorInstance, err := registration.New(config.ProcessorConfig{
				Type:   "midi.note_off.create",
				Params: test.params,
			})

			if err != nil {
				if test.errorString != err.Error() {
					t.Fatalf("midi.note_off.create got error '%s', expected '%s'", err.Error(), test.errorString)
				}
				return
			}

			got, err := processorInstance.Process(t.Context(), common.GetWrappedPayload(t.Context(), test.payload))

			if err == nil {
				t.Fatalf("midi.note_off.create expected to fail but succeeded, got: %v", got)
			}

			if err.Error() != test.errorString {
				t.Fatalf("midi.note_off.create got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
