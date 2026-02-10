package module_test

import (
	"testing"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/module"
)

func TestUDPMulticastFromRegistry(t *testing.T) {
	registration, ok := module.ModuleRegistry["net.udp.multicast"]
	if !ok {
		t.Fatalf("udp.multicast module not registered")
	}

	moduleInstance, err := registration.New(config.ModuleConfig{
		Type: "net.udp.multicast",
		Params: map[string]any{
			"ip":   "236.10.10.10",
			"port": 56565.0,
		},
	})

	if err != nil {
		t.Fatalf("failed to create udp.multicast module: %s", err)
	}

	if moduleInstance.Type() != "net.udp.multicast" {
		t.Fatalf("udp.multicast module has wrong type: %s", moduleInstance.Type())
	}
}
