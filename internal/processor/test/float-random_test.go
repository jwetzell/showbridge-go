package processor_test

import (
	"testing"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/processor"
)

func TestFloatRandomFromRegistry(t *testing.T) {
	registration, ok := processor.ProcessorRegistry["float.random"]
	if !ok {
		t.Fatalf("float.random processor not registered")
	}

	processorInstance, err := registration.New(config.ProcessorConfig{
		Type: "float.random",
		Params: map[string]any{
			"min": 1.0,
			"max": 10.0,
		},
	})

	if err != nil {
		t.Fatalf("failed to create float.random processor: %s", err)
	}

	if processorInstance.Type() != "float.random" {
		t.Fatalf("float.random processor has wrong type: %s", processorInstance.Type())
	}
}

func TestGoodFloatRandom(t *testing.T) {
	tests := []struct {
		name    string
		payload any
		params  map[string]any
	}{
		{
			name: "1-10",
			params: map[string]any{
				"min": 1.0,
				"max": 10.0,
			},
			payload: "12345",
		},
		{
			name: "1-1",
			params: map[string]any{
				"min": 1.0,
				"max": 1.0,
			},
			payload: "12345",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			registration, ok := processor.ProcessorRegistry["float.random"]
			if !ok {
				t.Fatalf("float.random processor not registered")
			}

			processorInstance, err := registration.New(config.ProcessorConfig{
				Type:   "float.random",
				Params: test.params,
			})

			if err != nil {
				t.Fatalf("float.random failed to create processor: %s", err)
			}

			got, err := processorInstance.Process(t.Context(), test.payload)
			if err != nil {
				t.Fatalf("float.random processing failed: %s", err)
			}

			bitSize, ok := test.params["bitSize"].(int)
			if !ok {
				bitSize = 32
			}

			var gotFloat float64
			if bitSize == 32 {
				gotFloat32, ok := got.(float32)
				if !ok {
					t.Fatalf("float.random returned a %T payload: %s", got, got)
				}
				gotFloat = float64(gotFloat32)
			}
			if bitSize == 64 {
				gotFloat64, ok := got.(float64)
				if !ok {
					t.Fatalf("float.random returned a %T payload: %s", got, got)
				}
				gotFloat = gotFloat64
			}

			minNum, ok := test.params["min"].(float64)
			if !ok {
				t.Fatalf("float.random test min param is not a number: %s", test.params["min"])
			}
			maxNum, ok := test.params["max"].(float64)
			if !ok {
				t.Fatalf("float.random test max param is not a number: %s", test.params["max"])
			}

			if gotFloat < minNum || gotFloat > maxNum {
				t.Fatalf("float.random got %f, expected between %f and %f", gotFloat, minNum, maxNum)
			}
		})
	}
}

func TestBadFloatRandom(t *testing.T) {
	tests := []struct {
		name        string
		params      map[string]any
		payload     any
		errorString string
	}{
		{
			name:        "no min param",
			payload:     "hello",
			params:      map[string]any{"max": 10.0},
			errorString: "float.random min error: not found",
		},
		{
			name:        "no max param",
			payload:     "hello",
			params:      map[string]any{"min": 1.0},
			errorString: "float.random max error: not found",
		},
		{
			name:        "min param not a number",
			payload:     "hello",
			params:      map[string]any{"min": "1", "max": 10.0},
			errorString: "float.random min error: not a number",
		},
		{
			name:        "max param not a number",
			payload:     "hello",
			params:      map[string]any{"min": 1.0, "max": "10"},
			errorString: "float.random max error: not a number",
		},
		{
			name:        "max less than min",
			payload:     "hello",
			params:      map[string]any{"min": 1.0, "max": 0.0},
			errorString: "float.random max must be greater than min",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			registration, ok := processor.ProcessorRegistry["float.random"]
			if !ok {
				t.Fatalf("float.random processor not registered")
			}

			processorInstance, err := registration.New(config.ProcessorConfig{
				Type:   "float.random",
				Params: test.params,
			})

			if err != nil {
				if test.errorString != err.Error() {
					t.Fatalf("float.random got error '%s', expected '%s'", err.Error(), test.errorString)
				}
				return
			}

			got, err := processorInstance.Process(t.Context(), test.payload)

			if err == nil {
				t.Fatalf("float.random expected to fail but got payload: %s", got)
			}

			if err.Error() != test.errorString {
				t.Fatalf("float.random got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
