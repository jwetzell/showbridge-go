package module_test

import (
	"testing"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/module"
)

func TestNATSClientFromRegistry(t *testing.T) {
	registration, ok := module.ModuleRegistry["nats.client"]
	if !ok {
		t.Fatalf("nats.client module not registered")
	}

	moduleInstance, err := registration.New(config.ModuleConfig{
		Id:   "test",
		Type: "nats.client",
		Params: map[string]any{
			"url":     "nats://127.0.0.1:4222",
			"subject": "test",
		},
	})

	if err != nil {
		t.Fatalf("failed to create nats.client module: %s", err)
	}

	if moduleInstance.Id() != "test" {
		t.Fatalf("nats.client module has wrong id: %s", moduleInstance.Id())
	}

	if moduleInstance.Type() != "nats.client" {
		t.Fatalf("nats.client module has wrong type: %s", moduleInstance.Type())
	}
}
