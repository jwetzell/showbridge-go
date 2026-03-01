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

func TestBadMIDIInput(t *testing.T) {
	tests := []struct {
		name        string
		params      map[string]any
		errorString string
	}{
		{
			name:        "no port param",
			params:      map[string]any{},
			errorString: "midi.input port error: not found",
		},
		{
			name:        "non-string port",
			params:      map[string]any{"port": 123},
			errorString: "midi.input port error: not a string",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			registration, ok := module.ModuleRegistry["midi.input"]
			if !ok {
				t.Fatalf("midi.input module not registered")
			}

			moduleInstance, err := registration.New(config.ModuleConfig{
				Id:     "test",
				Type:   "midi.input",
				Params: test.params,
			})

			if err != nil {
				if test.errorString != err.Error() {
					t.Fatalf("midi.input got error '%s', expected '%s'", err.Error(), test.errorString)
				}
				return
			}

			err = moduleInstance.Start(t.Context())

			if err == nil {
				t.Fatalf("midi.input expected to fail")
			}

			if err.Error() != test.errorString {
				t.Fatalf("midi.input got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
