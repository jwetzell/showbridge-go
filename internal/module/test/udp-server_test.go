package module_test

import (
	"testing"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/module"
)

func TestUDPServerFromRegistry(t *testing.T) {
	registration, ok := module.ModuleRegistry["net.udp.server"]
	if !ok {
		t.Fatalf("net.udp.server module not registered")
	}

	moduleInstance, err := registration.New(config.ModuleConfig{
		Type: "net.udp.server",
		Params: map[string]any{
			"port": 8000.0,
		},
	})

	if err != nil {
		t.Fatalf("failed to create udp.server module: %s", err)
	}

	if moduleInstance.Type() != "net.udp.server" {
		t.Fatalf("net.udp.server module has wrong type: %s", moduleInstance.Type())
	}
}
