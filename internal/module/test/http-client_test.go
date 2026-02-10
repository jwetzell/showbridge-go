package module_test

import (
	"testing"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/module"
)

func TestHTTPClientFromRegistry(t *testing.T) {
	registration, ok := module.ModuleRegistry["http.client"]
	if !ok {
		t.Fatalf("http.client module not registered")
	}

	moduleInstance, err := registration.New(config.ModuleConfig{
		Id:   "test",
		Type: "http.client",
	})

	if err != nil {
		t.Fatalf("failed to create http.client module: %s", err)
	}

	if moduleInstance.Id() != "test" {
		t.Fatalf("http.client module has wrong id: %s", moduleInstance.Id())
	}

	if moduleInstance.Type() != "http.client" {
		t.Fatalf("http.client module has wrong type: %s", moduleInstance.Type())
	}
}

func TestBadHTTPClient(t *testing.T) {
	tests := []struct {
		name        string
		params      map[string]any
		errorString string
	}{}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			registration, ok := module.ModuleRegistry["http.client"]
			if !ok {
				t.Fatalf("http.client module not registered")
			}

			moduleInstance, err := registration.New(config.ModuleConfig{
				Id:     "test",
				Type:   "http.client",
				Params: test.params,
			})

			if err != nil {
				if test.errorString != err.Error() {
					t.Fatalf("http.client got error '%s', expected '%s'", err.Error(), test.errorString)
				}
				return
			}

			err = moduleInstance.Start(t.Context())

			if err == nil {
				t.Fatalf("http.client expected to fail")
			}

			if err.Error() != test.errorString {
				t.Fatalf("http.client got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
