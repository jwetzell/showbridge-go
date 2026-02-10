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
		Id:   "test",
		Type: "net.udp.multicast",
		Params: map[string]any{
			"ip":   "236.10.10.10",
			"port": 56565.0,
		},
	})

	if err != nil {
		t.Fatalf("failed to create net.udp.multicast module: %s", err)
	}

	if moduleInstance.Id() != "test" {
		t.Fatalf("net.udp.multicast module has wrong id: %s", moduleInstance.Id())
	}

	if moduleInstance.Type() != "net.udp.multicast" {
		t.Fatalf("net.udp.multicast module has wrong type: %s", moduleInstance.Type())
	}
}
