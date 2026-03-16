package processor_test

import (
	"testing"

	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/processor"
)

func TestIntScaleFromRegistry(t *testing.T) {
	registration, ok := processor.ProcessorRegistry["int.scale"]
	if !ok {
		t.Fatalf("int.scale processor not registered")
	}

	processorInstance, err := registration.New(config.ProcessorConfig{
		Type: "int.scale",
		Params: map[string]any{
			"inMin":  0,
			"inMax":  10,
			"outMin": 0,
			"outMax": 127,
		},
	})

	if err != nil {
		t.Fatalf("failed to create int.scale processor: %s", err)
	}

	if processorInstance.Type() != "int.scale" {
		t.Fatalf("int.scale processor has wrong type: %s", processorInstance.Type())
	}
}

func TestGoodIntScale(t *testing.T) {
	tests := []struct {
		name     string
		payload  any
		params   map[string]any
		expected int
	}{
		{
			name: "0-10 -> 0-127",
			params: map[string]any{
				"inMin":  0,
				"inMax":  10,
				"outMin": 0,
				"outMax": 127,
			},
			payload:  5,
			expected: 63,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			registration, ok := processor.ProcessorRegistry["int.scale"]
			if !ok {
				t.Fatalf("int.scale processor not registered")
			}

			processorInstance, err := registration.New(config.ProcessorConfig{
				Type:   "int.scale",
				Params: test.params,
			})

			if err != nil {
				t.Fatalf("int.scale failed to create processor: %s", err)
			}

			got, err := processorInstance.Process(t.Context(), common.GetWrappedPayload(t.Context(), test.payload))
			if err != nil {
				t.Fatalf("int.scale processing failed: %s", err)
			}

			gotInt, ok := got.Payload.(int)
			if !ok {
				t.Fatalf("int.scale returned a %T payload: %+v", got, got)
			}

			if gotInt != test.expected {
				t.Fatalf("int.scale got %d, expected %d", gotInt, test.expected)
			}
		})
	}
}

func TestBadIntScale(t *testing.T) {
	tests := []struct {
		name        string
		params      map[string]any
		payload     any
		errorString string
	}{
		{
			name:        "no inMin param",
			payload:     "hello",
			params:      map[string]any{"inMax": 10, "outMin": 0, "outMax": 127},
			errorString: "int.scale inMin error: not found",
		},
		{
			name:        "no inMax param",
			payload:     "hello",
			params:      map[string]any{"inMin": 0, "outMin": 0, "outMax": 127},
			errorString: "int.scale inMax error: not found",
		},
		{
			name:        "no outMin param",
			payload:     "hello",
			params:      map[string]any{"inMin": 0, "inMax": 10, "outMax": 127},
			errorString: "int.scale outMin error: not found",
		},
		{
			name:        "no outMax param",
			payload:     "hello",
			params:      map[string]any{"inMin": 0, "inMax": 10, "outMin": 0},
			errorString: "int.scale outMax error: not found",
		},
		{
			name:        "inMin param not a number",
			payload:     "hello",
			params:      map[string]any{"inMin": "0", "max": 10, "outMin": 0, "outMax": 127},
			errorString: "int.scale inMin error: not a number",
		},
		{
			name:        "inMax param not a number",
			payload:     "hello",
			params:      map[string]any{"inMin": 0, "inMax": "10", "outMin": 0, "outMax": 127},
			errorString: "int.scale inMax error: not a number",
		},
		{
			name:        "inMax less than inMin",
			payload:     "hello",
			params:      map[string]any{"inMin": 10, "inMax": 0, "outMin": 0, "outMax": 127},
			errorString: "int.scale inMax must be greater than inMin",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			registration, ok := processor.ProcessorRegistry["int.scale"]
			if !ok {
				t.Fatalf("int.scale processor not registered")
			}

			processorInstance, err := registration.New(config.ProcessorConfig{
				Type:   "int.scale",
				Params: test.params,
			})

			if err != nil {
				if test.errorString != err.Error() {
					t.Fatalf("int.scale got error '%s', expected '%s'", err.Error(), test.errorString)
				}
				return
			}

			got, err := processorInstance.Process(t.Context(), common.GetWrappedPayload(t.Context(), test.payload))

			if err == nil {
				t.Fatalf("int.scale expected to fail but got payload: %+v", got)
			}

			if err.Error() != test.errorString {
				t.Fatalf("int.scale got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
