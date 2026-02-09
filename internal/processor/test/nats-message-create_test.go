package processor_test

import (
	"testing"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/processor"
)

func TestNATSMessageCreateFromRegistry(t *testing.T) {
	registration, ok := processor.ProcessorRegistry["nats.message.create"]
	if !ok {
		t.Fatalf("nats.message.create processor not registered")
	}

	processorInstance, err := registration.New(config.ProcessorConfig{
		Type: "nats.message.create",
		Params: map[string]any{
			"subject": "test",
			"payload": "Hello, World!",
		},
	})

	if err != nil {
		t.Fatalf("failed to create nats.message.create processor: %s", err)
	}

	if processorInstance.Type() != "nats.message.create" {
		t.Fatalf("nats.message.create processor has wrong type: %s", processorInstance.Type())
	}
}
