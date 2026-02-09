package processor_test

import (
	"testing"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/processor"
)

func TestOSCMessageDecodeFromRegistry(t *testing.T) {
	registration, ok := processor.ProcessorRegistry["osc.message.decode"]
	if !ok {
		t.Fatalf("osc.message.decode processor not registered")
	}

	processorInstance, err := registration.New(config.ProcessorConfig{
		Type: "osc.message.decode",
	})

	if err != nil {
		t.Fatalf("failed to create osc.message.decode processor: %s", err)
	}

	if processorInstance.Type() != "osc.message.decode" {
		t.Fatalf("osc.message.decode processor has wrong type: %s", processorInstance.Type())
	}
}
