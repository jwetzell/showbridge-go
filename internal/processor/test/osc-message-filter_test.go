package processor_test

import (
	"testing"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/processor"
)

func TestOSCMessageFilterFromRegistry(t *testing.T) {
	registration, ok := processor.ProcessorRegistry["osc.message.filter"]
	if !ok {
		t.Fatalf("osc.message.filter processor not registered")
	}

	processorInstance, err := registration.New(config.ProcessorConfig{
		Type: "osc.message.filter",
		Params: map[string]any{
			"address": "/test*",
		},
	})

	if err != nil {
		t.Fatalf("failed to filter osc.message.filter processor: %s", err)
	}

	if processorInstance.Type() != "osc.message.filter" {
		t.Fatalf("osc.message.filter processor has wrong type: %s", processorInstance.Type())
	}
}
