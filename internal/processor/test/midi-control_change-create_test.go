package processor_test

import (
	"reflect"
	"testing"

	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/processor"
	"gitlab.com/gomidi/midi/v2"
)

func TestMIDIControlChangeCreateFromRegistry(t *testing.T) {
	registration, ok := processor.ProcessorRegistry["midi.control_change.create"]
	if !ok {
		t.Fatalf("midi.control_change.create processor not registered")
	}

	processorInstance, err := registration.New(config.ProcessorConfig{
		Type: "midi.control_change.create",
		Params: map[string]any{
			"channel": "1",
			"control": "60",
			"value":   "100",
		},
	})

	if err != nil {
		t.Fatalf("failed to create midi.control_change.create processor: %s", err)
	}

	if processorInstance.Type() != "midi.control_change.create" {
		t.Fatalf("midi.control_change.create processor has wrong type: %s", processorInstance.Type())
	}
}

func TestGoodMIDIControlChangeCreate(t *testing.T) {

	tests := []struct {
		name     string
		params   map[string]any
		payload  any
		expected any
	}{
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
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			registration, ok := processor.ProcessorRegistry["midi.control_change.create"]
			if !ok {
				t.Fatalf("midi.control_change.create processor not registered")
			}

			processorInstance, err := registration.New(config.ProcessorConfig{
				Type:   "midi.control_change.create",
				Params: test.params,
			})

			if err != nil {
				t.Fatalf("midi.control_change.create failed to create processor: %s", err)
			}

			got, err := processorInstance.Process(t.Context(), common.GetWrappedPayload(t.Context(), test.payload))
			if err != nil {
				t.Fatalf("midi.control_change.create processing failed: %s", err)
			}

			gotMessage, ok := got.Payload.(midi.Message)
			if !ok {
				t.Fatalf("midi.control_change.create returned a %T payload: %+v", got, got)
			}

			if !reflect.DeepEqual(gotMessage, test.expected) {
				t.Fatalf("midi.control_change.create got %v, expected %v", gotMessage, test.expected)
			}
		})
	}
}

func TestBadMIDIControlChangeCreate(t *testing.T) {

	tests := []struct {
		name        string
		params      map[string]any
		payload     any
		errorString string
	}{
		{
			name: "control_change no channel",
			params: map[string]any{
				"type":    "control_change",
				"control": "64",
				"value":   "127",
			},
			payload:     "test",
			errorString: "midi.control_change.create channel error: not found",
		},
		{
			name: "control_change no control",
			params: map[string]any{
				"type":    "control_change",
				"channel": "1",
				"value":   "127",
			},
			payload:     "test",
			errorString: "midi.control_change.create control error: not found",
		},
		{
			name: "control_change no value",
			params: map[string]any{
				"type":    "control_change",
				"channel": "1",
				"control": "64",
			},
			payload:     "test",
			errorString: "midi.control_change.create value error: not found",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			registration, ok := processor.ProcessorRegistry["midi.control_change.create"]
			if !ok {
				t.Fatalf("midi.control_change.create processor not registered")
			}

			processorInstance, err := registration.New(config.ProcessorConfig{
				Type:   "midi.control_change.create",
				Params: test.params,
			})

			if err != nil {
				if test.errorString != err.Error() {
					t.Fatalf("midi.control_change.create got error '%s', expected '%s'", err.Error(), test.errorString)
				}
				return
			}

			got, err := processorInstance.Process(t.Context(), common.GetWrappedPayload(t.Context(), test.payload))

			if err == nil {
				t.Fatalf("midi.control_change.create expected to fail but succeeded, got: %v", got)
			}

			if err.Error() != test.errorString {
				t.Fatalf("midi.control_change.create got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
