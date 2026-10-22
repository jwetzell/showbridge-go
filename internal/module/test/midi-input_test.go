//go:build cgo

package module_test

import (
	"testing"

	"github.com/jwetzell/showbridge-go/config"
	"github.com/jwetzell/showbridge-go/internal/module"
)

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

			registration, ok := module.GetModuleRegistration("midi.input")
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

			err = moduleInstance.Start(t.Context(), nil)

			if err == nil {
				t.Fatalf("midi.input expected to fail")
			}

			if err.Error() != test.errorString {
				t.Fatalf("midi.input got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
