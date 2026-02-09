package processor_test

import (
	"testing"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/processor"
)

func TestMQTTMessageEncodeFromRegistry(t *testing.T) {
	registration, ok := processor.ProcessorRegistry["mqtt.message.encode"]
	if !ok {
		t.Fatalf("mqtt.message.encode processor not registered")
	}

	processorInstance, err := registration.New(config.ProcessorConfig{
		Type: "mqtt.message.encode",
	})

	if err != nil {
		t.Fatalf("failed to create mqtt.message.encode processor: %s", err)
	}

	if processorInstance.Type() != "mqtt.message.encode" {
		t.Fatalf("mqtt.message.encode processor has wrong type: %s", processorInstance.Type())
	}
}
