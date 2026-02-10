package module_test

import (
	"testing"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/module"
)

func TestHTTPServerFromRegistry(t *testing.T) {
	registration, ok := module.ModuleRegistry["http.server"]
	if !ok {
		t.Fatalf("http.server module not registered")
	}

	moduleInstance, err := registration.New(config.ModuleConfig{
		Id:   "test",
		Type: "http.server",
		Params: map[string]any{
			"port": 3000.0,
		},
	})

	if err != nil {
		t.Fatalf("failed to create http.server module: %s", err)
	}

	if moduleInstance.Id() != "test" {
		t.Fatalf("http.server module has wrong id: %s", moduleInstance.Id())
	}

	if moduleInstance.Type() != "http.server" {
		t.Fatalf("http.server module has wrong type: %s", moduleInstance.Type())
	}
}

func TestBadHTTPServer(t *testing.T) {
	tests := []struct {
		name        string
		params      map[string]any
		errorString string
	}{
		{
			name:        "no port param",
			params:      map[string]any{},
			errorString: "http.server requires a port parameter",
		},
		{
			name:        "non-numeric port",
			params:      map[string]any{"port": "3000"},
			errorString: "http.server port must be a number",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			registration, ok := module.ModuleRegistry["http.server"]
			if !ok {
				t.Fatalf("http.server module not registered")
			}

			moduleInstance, err := registration.New(config.ModuleConfig{
				Id:     "test",
				Type:   "http.server",
				Params: test.params,
			})

			if err != nil {
				if test.errorString != err.Error() {
					t.Fatalf("http.server got error '%s', expected '%s'", err.Error(), test.errorString)
				}
				return
			}

			err = moduleInstance.Start(t.Context())

			if err == nil {
				t.Fatalf("http.server expected to fail")
			}

			if err.Error() != test.errorString {
				t.Fatalf("http.server got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
