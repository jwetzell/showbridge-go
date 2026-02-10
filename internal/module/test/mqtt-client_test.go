package module_test

import (
	"testing"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/module"
)

func TestMQTTClientFromRegistry(t *testing.T) {
	registration, ok := module.ModuleRegistry["mqtt.client"]
	if !ok {
		t.Fatalf("mqtt.client module not registered")
	}

	moduleInstance, err := registration.New(config.ModuleConfig{
		Type: "mqtt.client",
		Params: map[string]any{
			"broker":   "mqtt://localhost:1883",
			"topic":    "test/topic",
			"clientId": "test",
		},
	})

	if err != nil {
		t.Fatalf("failed to create mqtt.client module: %s", err)
	}

	if moduleInstance.Type() != "mqtt.client" {
		t.Fatalf("mqtt.client module has wrong type: %s", moduleInstance.Type())
	}
}
