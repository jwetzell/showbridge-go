package module_test

import (
	"testing"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/module"
)

func TestPSNClientFromRegistry(t *testing.T) {
	registration, ok := module.ModuleRegistry["psn.client"]
	if !ok {
		t.Fatalf("psn.client module not registered")
	}

	moduleInstance, err := registration.New(config.ModuleConfig{
		Id:   "test",
		Type: "psn.client",
	})

	if err != nil {
		t.Fatalf("failed to create psn.client module: %s", err)
	}

	if moduleInstance.Id() != "test" {
		t.Fatalf("psn.client module has wrong id: %s", moduleInstance.Id())
	}

	if moduleInstance.Type() != "psn.client" {
		t.Fatalf("psn.client module has wrong type: %s", moduleInstance.Type())
	}
}

func TestBadPSNClient(t *testing.T) {
	tests := []struct {
		name        string
		params      map[string]any
		errorString string
	}{}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			registration, ok := module.ModuleRegistry["psn.client"]
			if !ok {
				t.Fatalf("psn.client module not registered")
			}

			moduleInstance, err := registration.New(config.ModuleConfig{
				Id:     "test",
				Type:   "psn.client",
				Params: test.params,
			})

			if err != nil {
				if test.errorString != err.Error() {
					t.Fatalf("psn.client got error '%s', expected '%s'", err.Error(), test.errorString)
				}
				return
			}

			err = moduleInstance.Start(t.Context())

			if err == nil {
				t.Fatalf("psn.client expected to fail")
			}

			if err.Error() != test.errorString {
				t.Fatalf("psn.client got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
