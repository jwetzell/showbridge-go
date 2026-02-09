package processor_test

import (
	"testing"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/processor"
)

func TestNATSMessageEncodeFromRegistry(t *testing.T) {
	registration, ok := processor.ProcessorRegistry["nats.message.encode"]
	if !ok {
		t.Fatalf("nats.message.encode processor not registered")
	}

	processorInstance, err := registration.New(config.ProcessorConfig{
		Type: "nats.message.encode",
	})

	if err != nil {
		t.Fatalf("failed to create nats.message.encode processor: %s", err)
	}

	if processorInstance.Type() != "nats.message.encode" {
		t.Fatalf("nats.message.encode processor has wrong type: %s", processorInstance.Type())
	}
}
