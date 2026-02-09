package processor_test

import (
	"testing"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/processor"
)

func TestOSCMessageCreateFromRegistry(t *testing.T) {
	registration, ok := processor.ProcessorRegistry["osc.message.create"]
	if !ok {
		t.Fatalf("osc.message.create processor not registered")
	}

	processorInstance, err := registration.New(config.ProcessorConfig{
		Type: "osc.message.create",
		Params: map[string]any{
			"address": "/test",
		},
	})

	if err != nil {
		t.Fatalf("failed to create osc.message.create processor: %s", err)
	}

	if processorInstance.Type() != "osc.message.create" {
		t.Fatalf("osc.message.create processor has wrong type: %s", processorInstance.Type())
	}
}
