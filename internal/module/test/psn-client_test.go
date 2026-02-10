package module_test

import (
	"testing"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/module"
)

func TestPSNClientFromRegistry(t *testing.T) {
	registration, ok := module.ModuleRegistry["psn.client"]
	if !ok {
		t.Fatalf("psn.client module not registered")
	}

	moduleInstance, err := registration.New(config.ModuleConfig{
		Type: "psn.client",
	})

	if err != nil {
		t.Fatalf("failed to create psn.client module: %s", err)
	}

	if moduleInstance.Type() != "psn.client" {
		t.Fatalf("psn.client module has wrong type: %s", moduleInstance.Type())
	}
}
