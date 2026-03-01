package processor_test

import (
	"testing"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/processor"
)

func TestUintRandomFromRegistry(t *testing.T) {
	registration, ok := processor.ProcessorRegistry["uint.random"]
	if !ok {
		t.Fatalf("uint.random processor not registered")
	}

	processorInstance, err := registration.New(config.ProcessorConfig{
		Type: "uint.random",
		Params: map[string]any{
			"min": 1.0,
			"max": 10.0,
		},
	})

	if err != nil {
		t.Fatalf("failed to create uint.random processor: %s", err)
	}

	if processorInstance.Type() != "uint.random" {
		t.Fatalf("uint.random processor has wrong type: %s", processorInstance.Type())
	}
}

func TestUintRandomGoodConfig(t *testing.T) {
	registration, ok := processor.ProcessorRegistry["uint.random"]
	if !ok {
		t.Fatalf("uint.random processor not registered")
	}

	processorInstance, err := registration.New(config.ProcessorConfig{
		Type: "uint.random",
		Params: map[string]any{
			"min": 1.0,
			"max": 10.0,
		},
	})

	if err != nil {
		t.Fatalf("uint.random should have created processor but got error: %s", err)
	}

	payload := "12345"

	got, err := processorInstance.Process(t.Context(), payload)
	if err != nil {
		t.Fatalf("uint.random processing failed: %s", err)
	}

	gotUint, ok := got.(uint)
	if !ok {
		t.Fatalf("uint.random returned a %T payload: %s", got, got)
	}

	if gotUint < 1 || gotUint > 10 {
		t.Fatalf("uint.random got %d, expected between %d and %d", gotUint, 1, 10)
	}
}

func TestGoodUintRandom(t *testing.T) {

	tests := []struct {
		name    string
		params  map[string]any
		payload any
	}{
		{
			name:    "1-10",
			params:  map[string]any{"min": 1.0, "max": 10.0},
			payload: "12345",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			registration, ok := processor.ProcessorRegistry["uint.random"]
			if !ok {
				t.Fatalf("uint.random processor not registered")
			}

			processorInstance, err := registration.New(config.ProcessorConfig{
				Type:   "uint.random",
				Params: test.params,
			})

			if err != nil {
				t.Fatalf("uint.random failed to create processor: %s", err)
			}

			got, err := processorInstance.Process(t.Context(), test.payload)

			gotUint, ok := got.(uint)
			if !ok {
				t.Fatalf("uint.random returned a %T payload: %s", got, got)
			}
			if err != nil {
				t.Fatalf("uint.random failed: %s", err)
			}
			minNum, ok := test.params["min"].(float64)
			if !ok {
				t.Fatalf("uint.random test min param is not a number")
			}
			maxNum, ok := test.params["max"].(float64)
			if !ok {
				t.Fatalf("uint.random test max param is not a number")
			}
			if gotUint < uint(minNum) || gotUint > uint(maxNum) {
				t.Fatalf("uint.random got %d, expected between %d and %d", gotUint, uint(minNum), uint(maxNum))
			}
		})
	}
}

func TestBadUintRandom(t *testing.T) {
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
			errorString: "uint.random min error: not found",
		},
		{
			name:        "no max param",
			payload:     "hello",
			params:      map[string]any{"min": 1.0},
			errorString: "uint.random max error: not found",
		},
		{
			name:        "min param not a number",
			payload:     "hello",
			params:      map[string]any{"min": "1", "max": 10.0},
			errorString: "uint.random min error: not a number",
		},
		{
			name:        "max param not a number",
			payload:     "hello",
			params:      map[string]any{"min": 1.0, "max": "10"},
			errorString: "uint.random max error: not a number",
		},
		{
			name:        "max less than min",
			payload:     "hello",
			params:      map[string]any{"min": 1.0, "max": 0.0},
			errorString: "uint.random max must be greater than min",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			registration, ok := processor.ProcessorRegistry["uint.random"]
			if !ok {
				t.Fatalf("uint.random processor not registered")
			}

			processorInstance, err := registration.New(config.ProcessorConfig{
				Type:   "uint.random",
				Params: test.params,
			})

			if err != nil {
				if test.errorString != err.Error() {
					t.Fatalf("uint.random got error '%s', expected '%s'", err.Error(), test.errorString)
				}
				return
			}

			got, err := processorInstance.Process(t.Context(), test.payload)

			if err == nil {
				t.Fatalf("uint.random expected to fail but got payload: %s", got)
			}

			if err.Error() != test.errorString {
				t.Fatalf("uint.random got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
