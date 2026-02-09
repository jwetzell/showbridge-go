package processor_test

import (
	"testing"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/processor"
)

func TestIntRandomFromRegistry(t *testing.T) {
	registration, ok := processor.ProcessorRegistry["int.random"]
	if !ok {
		t.Fatalf("int.random processor not registered")
	}

	processorInstance, err := registration.New(config.ProcessorConfig{
		Type: "int.random",
		Params: map[string]any{
			"min": 1.0,
			"max": 10.0,
		},
	})

	if err != nil {
		t.Fatalf("failed to create int.random processor: %s", err)
	}

	if processorInstance.Type() != "int.random" {
		t.Fatalf("int.random processor has wrong type: %s", processorInstance.Type())
	}
}

func TestIntRandomGoodConfig(t *testing.T) {
	registration, ok := processor.ProcessorRegistry["int.random"]
	if !ok {
		t.Fatalf("int.random processor not registered")
	}

	processorInstance, err := registration.New(config.ProcessorConfig{
		Type: "int.random",
		Params: map[string]any{
			"min": 1.0,
			"max": 10.0,
		},
	})

	if err != nil {
		t.Fatalf("int.random should have created processor but got error: %s", err)
	}

	payload := "12345"

	got, err := processorInstance.Process(t.Context(), payload)
	if err != nil {
		t.Fatalf("int.random processing failed: %s", err)
	}

	gotInt, ok := got.(int)
	if !ok {
		t.Fatalf("int.random returned a %T payload: %s", got, got)
	}

	if gotInt < 1 || gotInt > 10 {
		t.Fatalf("int.random got %d, expected between %d and %d", gotInt, 1, 10)
	}
}

func TestGoodIntRandom(t *testing.T) {
	tests := []struct {
		processor processor.Processor
		name      string
		payload   any
		min       int
		max       int
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
			intRandom := processor.IntRandom{
				Min: test.min,
				Max: test.max,
			}
			got, err := intRandom.Process(t.Context(), test.payload)
			gotInt, ok := got.(int)
			if !ok {
				t.Fatalf("int.random returned a %T payload: %s", got, got)
			}
			if err != nil {
				t.Fatalf("int.random failed: %s", err)
			}
			if gotInt < test.min || gotInt > test.max {
				t.Fatalf("int.random got %d, expected between %d and %d", gotInt, test.min, test.max)
			}
		})
	}
}

func TestBadIntRandom(t *testing.T) {
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
			errorString: "int.random requires a min parameter",
		},
		{
			name:        "no max param",
			payload:     "hello",
			params:      map[string]any{"min": 1.0},
			errorString: "int.random requires a max parameter",
		},
		{
			name:        "min param not a number",
			payload:     "hello",
			params:      map[string]any{"min": "1", "max": 10.0},
			errorString: "int.random min must be a number",
		},
		{
			name:        "max param not a number",
			payload:     "hello",
			params:      map[string]any{"min": 1.0, "max": "10"},
			errorString: "int.random max must be a number",
		},
		{
			name:        "max less than min",
			payload:     "hello",
			params:      map[string]any{"min": 1.0, "max": 0.0},
			errorString: "int.random max must be greater than min",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			registration, ok := processor.ProcessorRegistry["int.random"]
			if !ok {
				t.Fatalf("int.random processor not registered")
			}

			processorInstance, err := registration.New(config.ProcessorConfig{
				Type:   "int.random",
				Params: test.params,
			})

			if err != nil {
				if test.errorString != err.Error() {
					t.Fatalf("int.random got error '%s', expected '%s'", err.Error(), test.errorString)
				}
				return
			}

			got, err := processorInstance.Process(t.Context(), test.payload)

			if err == nil {
				t.Fatalf("int.random expected to fail but got payload: %s", got)
			}

			if err.Error() != test.errorString {
				t.Fatalf("int.random got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
