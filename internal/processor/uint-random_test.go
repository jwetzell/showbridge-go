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
		processor processor.Processor
		name      string
		payload   any
		min       uint
		max       uint
	}{
		{
			name:    "1-10",
			payload: "12345",
			min:     1,
			max:     10,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			uintRandom := processor.UintRandom{
				Min: test.min,
				Max: test.max,
			}
			got, err := uintRandom.Process(t.Context(), test.payload)

			gotUint, ok := got.(uint)
			if !ok {
				t.Fatalf("uint.random returned a %T payload: %s", got, got)
			}
			if err != nil {
				t.Fatalf("uint.random failed: %s", err)
			}
			if gotUint < test.min || gotUint > test.max {
				t.Fatalf("uint.random got %d, expected between %d and %d", gotUint, test.min, test.max)
			}
		})
	}
}
