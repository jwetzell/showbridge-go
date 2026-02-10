package module_test

import (
	"testing"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/module"
)

func TestUDPClientFromRegistry(t *testing.T) {
	registration, ok := module.ModuleRegistry["net.udp.client"]
	if !ok {
		t.Fatalf("udp.client module not registered")
	}

	moduleInstance, err := registration.New(config.ModuleConfig{
		Id:   "test",
		Type: "net.udp.client",
		Params: map[string]any{
			"host":    "localhost",
			"port":    8000.0,
			"framing": "LF",
		},
	})

	if err != nil {
		t.Fatalf("failed to create net.udp.client module: %s", err)
	}

	if moduleInstance.Id() != "test" {
		t.Fatalf("net.udp.client module has wrong id: %s", moduleInstance.Id())
	}

	if moduleInstance.Type() != "net.udp.client" {
		t.Fatalf("net.udp.client module has wrong type: %s", moduleInstance.Type())
	}
}
