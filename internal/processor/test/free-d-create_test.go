package processor_test

import (
	"testing"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/processor"
)

func TestFreeDCreateFromRegistry(t *testing.T) {
	registration, ok := processor.ProcessorRegistry["freed.create"]
	if !ok {
		t.Fatalf("freed.create processor not registered")
	}

	processorInstance, err := registration.New(config.ProcessorConfig{
		Type: "freed.create",
		Params: map[string]any{
			"id":    "0",
			"pan":   "0",
			"tilt":  "0",
			"roll":  "0",
			"posX":  "0",
			"posY":  "0",
			"posZ":  "0",
			"zoom":  "0",
			"focus": "0",
		},
	})

	if err != nil {
		t.Fatalf("failed to create freed.create processor: %s", err)
	}

	if processorInstance.Type() != "freed.create" {
		t.Fatalf("freed.create processor has wrong type: %s", processorInstance.Type())
	}
}
