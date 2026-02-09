package processor_test

import (
	"testing"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/processor"
)

func TestFreeDEncodeFromRegistry(t *testing.T) {
	registration, ok := processor.ProcessorRegistry["freed.encode"]
	if !ok {
		t.Fatalf("freed.encode processor not registered")
	}

	processorInstance, err := registration.New(config.ProcessorConfig{
		Type: "freed.encode",
	})

	if err != nil {
		t.Fatalf("failed to create freed.encode processor: %s", err)
	}

	if processorInstance.Type() != "freed.encode" {
		t.Fatalf("freed.encode processor has wrong type: %s", processorInstance.Type())
	}
}
