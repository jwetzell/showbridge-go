package module_test

import (
	"testing"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/module"
)

func TestSIPCallServerFromRegistry(t *testing.T) {
	registration, ok := module.ModuleRegistry["sip.call.server"]
	if !ok {
		t.Fatalf("sip.call.server module not registered")
	}

	moduleInstance, err := registration.New(config.ModuleConfig{
		Id:   "test",
		Type: "sip.call.server",
	})

	if err != nil {
		t.Fatalf("failed to create sip.call.server module: %s", err)
	}

	if moduleInstance.Id() != "test" {
		t.Fatalf("sip.call.server module has wrong id: %s", moduleInstance.Id())
	}

	if moduleInstance.Type() != "sip.call.server" {
		t.Fatalf("sip.call.server module has wrong type: %s", moduleInstance.Type())
	}
}
