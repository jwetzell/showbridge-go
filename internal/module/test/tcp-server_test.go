package module_test

import (
	"testing"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/module"
)

func TestTCPServerFromRegistry(t *testing.T) {
	registration, ok := module.ModuleRegistry["net.tcp.server"]
	if !ok {
		t.Fatalf("net.tcp.server module not registered")
	}

	moduleInstance, err := registration.New(config.ModuleConfig{
		Type: "net.tcp.server",
		Params: map[string]any{
			"port":    8000.0,
			"framing": "LF",
		},
	})

	if err != nil {
		t.Fatalf("failed to create net.tcp.server module: %s", err)
	}

	if moduleInstance.Type() != "net.tcp.server" {
		t.Fatalf("net.tcp.server module has wrong type: %s", moduleInstance.Type())
	}
}
