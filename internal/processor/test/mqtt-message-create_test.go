package processor_test

import (
	"testing"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/processor"
)

func TestMQTTMessageCreateFromRegistry(t *testing.T) {
	registration, ok := processor.ProcessorRegistry["mqtt.message.create"]
	if !ok {
		t.Fatalf("mqtt.message.create processor not registered")
	}

	processorInstance, err := registration.New(config.ProcessorConfig{
		Type: "mqtt.message.create",
		Params: map[string]any{
			"topic":    "test/topic",
			"payload":  "Hello, World!",
			"qos":      1.0,
			"retained": true,
		},
	})

	if err != nil {
		t.Fatalf("failed to create mqtt.message.create processor: %s", err)
	}

	if processorInstance.Type() != "mqtt.message.create" {
		t.Fatalf("mqtt.message.create processor has wrong type: %s", processorInstance.Type())
	}
}
