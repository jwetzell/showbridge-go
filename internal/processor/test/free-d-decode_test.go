package processor_test

import (
	"testing"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/processor"
)

func TestFreeDDecodeFromRegistry(t *testing.T) {
	registration, ok := processor.ProcessorRegistry["freed.decode"]
	if !ok {
		t.Fatalf("freed.decode processor not registered")
	}

	processorInstance, err := registration.New(config.ProcessorConfig{
		Type: "freed.decode",
	})

	if err != nil {
		t.Fatalf("failed to create freed.decode processor: %s", err)
	}

	if processorInstance.Type() != "freed.decode" {
		t.Fatalf("freed.decode processor has wrong type: %s", processorInstance.Type())
	}
}
