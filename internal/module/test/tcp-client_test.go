package module_test

import (
	"testing"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/module"
)

func TestTCPClientFromRegistry(t *testing.T) {
	registration, ok := module.ModuleRegistry["net.tcp.client"]
	if !ok {
		t.Fatalf("net.tcp.client module not registered")
	}

	moduleInstance, err := registration.New(config.ModuleConfig{
		Type: "net.tcp.client",
		Params: map[string]any{
			"host":    "localhost",
			"port":    8000.0,
			"framing": "LF",
		},
	})

	if err != nil {
		t.Fatalf("failed to create net.tcp.client module: %s", err)
	}

	if moduleInstance.Type() != "net.tcp.client" {
		t.Fatalf("net.tcp.client module has wrong type: %s", moduleInstance.Type())
	}
}
