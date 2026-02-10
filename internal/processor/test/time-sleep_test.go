package processor_test

import (
	"testing"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/processor"
)

func TestTimeSleepFromRegistry(t *testing.T) {
	registration, ok := processor.ProcessorRegistry["time.sleep"]
	if !ok {
		t.Fatalf("time.sleep processor not registered")
	}

	processorInstance, err := registration.New(config.ProcessorConfig{
		Type: "time.sleep",
		Params: map[string]any{
			"duration": 1000.0,
		},
	})

	if err != nil {
		t.Fatalf("failed to create time.sleep processor: %s", err)
	}

	if processorInstance.Type() != "time.sleep" {
		t.Fatalf("time.sleep processor has wrong type: %s", processorInstance.Type())
	}
}

func TestBadTimeSleep(t *testing.T) {
	tests := []struct {
		name        string
		params      map[string]any
		payload     any
		errorString string
	}{}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			registration, ok := processor.ProcessorRegistry["time.sleep"]
			if !ok {
				t.Fatalf("time.sleep processor not registered")
			}

			processorInstance, err := registration.New(config.ProcessorConfig{
				Type:   "time.sleep",
				Params: test.params,
			})

			got, err := processorInstance.Process(t.Context(), test.payload)

			if err == nil {
				t.Fatalf("time.sleep expected to fail but succeeded, got: %v", got)

			}
			if err.Error() != test.errorString {
				t.Fatalf("time.sleep got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
