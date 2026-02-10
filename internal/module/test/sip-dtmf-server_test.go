package module_test

import (
	"testing"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/module"
)

func TestSIPDTMFServerFromRegistry(t *testing.T) {
	registration, ok := module.ModuleRegistry["sip.dtmf.server"]
	if !ok {
		t.Fatalf("sip.dtmf.server module not registered")
	}

	moduleInstance, err := registration.New(config.ModuleConfig{
		Type: "sip.dtmf.server",
		Params: map[string]any{
			"separator": "#",
		},
	})

	if err != nil {
		t.Fatalf("failed to create sip.dtmf.server module: %s", err)
	}

	if moduleInstance.Type() != "sip.dtmf.server" {
		t.Fatalf("sip.dtmf.server module has wrong type: %s", moduleInstance.Type())
	}
}
