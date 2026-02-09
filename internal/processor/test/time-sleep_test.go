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
