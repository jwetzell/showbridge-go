package processor_test

import (
	"reflect"
	"testing"

	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/processor"
	"gitlab.com/gomidi/midi/v2"
)

func TestMIDIProgramChangeCreateFromRegistry(t *testing.T) {
	registration, ok := processor.ProcessorRegistry["midi.program_change.create"]
	if !ok {
		t.Fatalf("midi.program_change.create processor not registered")
	}

	processorInstance, err := registration.New(config.ProcessorConfig{
		Type: "midi.program_change.create",
		Params: map[string]any{
			"channel": "1",
			"program": "60",
		},
	})

	if err != nil {
		t.Fatalf("failed to create midi.program_change.create processor: %s", err)
	}

	if processorInstance.Type() != "midi.program_change.create" {
		t.Fatalf("midi.program_change.create processor has wrong type: %s", processorInstance.Type())
	}
}

func TestGoodMIDIProgramChangeCreate(t *testing.T) {

	tests := []struct {
		name     string
		params   map[string]any
		payload  any
		expected any
	}{
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
			registration, ok := processor.ProcessorRegistry["midi.program_change.create"]
			if !ok {
				t.Fatalf("midi.program_change.create processor not registered")
			}

			processorInstance, err := registration.New(config.ProcessorConfig{
				Type:   "midi.program_change.create",
				Params: test.params,
			})

			if err != nil {
				t.Fatalf("midi.program_change.create failed to create processor: %s", err)
			}

			got, err := processorInstance.Process(t.Context(), common.GetWrappedPayload(t.Context(), test.payload))
			if err != nil {
				t.Fatalf("midi.program_change.create processing failed: %s", err)
			}

			gotMessage, ok := got.Payload.(midi.Message)
			if !ok {
				t.Fatalf("midi.program_change.create returned a %T payload: %+v", got, got)
			}

			if !reflect.DeepEqual(gotMessage, test.expected) {
				t.Fatalf("midi.program_change.create got %v, expected %v", gotMessage, test.expected)
			}
		})
	}
}

func TestBadMIDIProgramChangeCreate(t *testing.T) {

	tests := []struct {
		name        string
		params      map[string]any
		payload     any
		errorString string
	}{
		{
			name: "program_change no channel",
			params: map[string]any{
				"type":    "program_change",
				"program": "64",
			},
			payload:     "test",
			errorString: "midi.program_change.create channel error: not found",
		},
		{
			name: "program_change no program",
			params: map[string]any{
				"type":    "program_change",
				"channel": "1",
			},
			payload:     "test",
			errorString: "midi.program_change.create program error: not found",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			registration, ok := processor.ProcessorRegistry["midi.program_change.create"]
			if !ok {
				t.Fatalf("midi.program_change.create processor not registered")
			}

			processorInstance, err := registration.New(config.ProcessorConfig{
				Type:   "midi.program_change.create",
				Params: test.params,
			})

			if err != nil {
				if test.errorString != err.Error() {
					t.Fatalf("midi.program_change.create got error '%s', expected '%s'", err.Error(), test.errorString)
				}
				return
			}

			got, err := processorInstance.Process(t.Context(), common.GetWrappedPayload(t.Context(), test.payload))

			if err == nil {
				t.Fatalf("midi.program_change.create expected to fail but succeeded, got: %v", got)
			}

			if err.Error() != test.errorString {
				t.Fatalf("midi.program_change.create got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
