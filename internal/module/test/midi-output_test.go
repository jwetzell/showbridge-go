package module_test

import (
	"testing"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/module"
)

func TestMIDIOutputFromRegistry(t *testing.T) {
	registration, ok := module.ModuleRegistry["midi.output"]
	if !ok {
		t.Fatalf("midi.output module not registered")
	}

	moduleInstance, err := registration.New(config.ModuleConfig{
		Type: "midi.output",
		Params: map[string]any{
			"port": "test",
		},
	})

	if err != nil {
		t.Fatalf("failed to create midi.output module: %s", err)
	}

	if moduleInstance.Type() != "midi.output" {
		t.Fatalf("midi.output module has wrong type: %s", moduleInstance.Type())
	}
}
