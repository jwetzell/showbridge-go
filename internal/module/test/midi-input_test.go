package module_test

import (
	"testing"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/module"
)

func TestMIDIInputFromRegistry(t *testing.T) {
	registration, ok := module.ModuleRegistry["midi.input"]
	if !ok {
		t.Fatalf("midi.input module not registered")
	}

	moduleInstance, err := registration.New(config.ModuleConfig{
		Id:   "test",
		Type: "midi.input",
		Params: map[string]any{
			"port": "test",
		},
	})

	if err != nil {
		t.Fatalf("failed to create midi.input module: %s", err)
	}

	if moduleInstance.Id() != "test" {
		t.Fatalf("midi.input module has wrong id: %s", moduleInstance.Id())
	}

	if moduleInstance.Type() != "midi.input" {
		t.Fatalf("midi.input module has wrong type: %s", moduleInstance.Type())
	}
}
