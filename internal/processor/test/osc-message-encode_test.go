package processor_test

import (
	"testing"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/processor"
)

func TestOSCMessageEncodeFromRegistry(t *testing.T) {
	registration, ok := processor.ProcessorRegistry["osc.message.encode"]
	if !ok {
		t.Fatalf("osc.message.encode processor not registered")
	}

	processorInstance, err := registration.New(config.ProcessorConfig{
		Type: "osc.message.encode",
	})

	if err != nil {
		t.Fatalf("failed to create osc.message.encode processor: %s", err)
	}

	if processorInstance.Type() != "osc.message.encode" {
		t.Fatalf("osc.message.encode processor has wrong type: %s", processorInstance.Type())
	}
}
