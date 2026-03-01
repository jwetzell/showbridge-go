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
		Id:   "test",
		Type: "midi.output",
		Params: map[string]any{
			"port": "test",
		},
	})

	if err != nil {
		t.Fatalf("failed to create midi.output module: %s", err)
	}

	if moduleInstance.Id() != "test" {
		t.Fatalf("midi.output module has wrong id: %s", moduleInstance.Id())
	}

	if moduleInstance.Type() != "midi.output" {
		t.Fatalf("midi.output module has wrong type: %s", moduleInstance.Type())
	}
}

func TestBadMIDIOutput(t *testing.T) {
	tests := []struct {
		name        string
		params      map[string]any
		errorString string
	}{
		{
			name:        "no port param",
			params:      map[string]any{},
			errorString: "midi.output port error: not found",
		},
		{
			name:        "non-string port",
			params:      map[string]any{"port": 123},
			errorString: "midi.output port error: not a string",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			registration, ok := module.ModuleRegistry["midi.output"]
			if !ok {
				t.Fatalf("midi.output module not registered")
			}

			moduleInstance, err := registration.New(config.ModuleConfig{
				Id:     "test",
				Type:   "midi.output",
				Params: test.params,
			})

			if err != nil {
				if test.errorString != err.Error() {
					t.Fatalf("midi.output got error '%s', expected '%s'", err.Error(), test.errorString)
				}
				return
			}

			err = moduleInstance.Start(t.Context())

			if err == nil {
				t.Fatalf("midi.output expected to fail")
			}

			if err.Error() != test.errorString {
				t.Fatalf("midi.output got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
