package module_test

import (
	"testing"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/module"
)

func TestTCPClientFromRegistry(t *testing.T) {
	registration, ok := module.ModuleRegistry["net.tcp.client"]
	if !ok {
		t.Fatalf("net.tcp.client module not registered")
	}

	moduleInstance, err := registration.New(config.ModuleConfig{
		Id:   "test",
		Type: "net.tcp.client",
		Params: map[string]any{
			"host":    "localhost",
			"port":    8000.0,
			"framing": "LF",
		},
	})

	if err != nil {
		t.Fatalf("failed to create net.tcp.client module: %s", err)
	}

	if moduleInstance.Id() != "test" {
		t.Fatalf("net.tcp.client module has wrong id: %s", moduleInstance.Id())
	}

	if moduleInstance.Type() != "net.tcp.client" {
		t.Fatalf("net.tcp.client module has wrong type: %s", moduleInstance.Type())
	}
}

func TestBadTCPClient(t *testing.T) {
	tests := []struct {
		name        string
		params      map[string]any
		errorString string
	}{
		{
			name: "no port param",
			params: map[string]any{
				"host": "localhost",
			},
			errorString: "net.tcp.client port error: not found",
		},
		{
			name: "non-number port param",
			params: map[string]any{
				"host": "localhost",
				"port": "8000",
			},
			errorString: "net.tcp.client port error: not a number",
		},
		{
			name: "no host param",
			params: map[string]any{
				"port": 8000.0,
			},
			errorString: "net.tcp.client host error: not found",
		},
		{
			name: "non-string host param",
			params: map[string]any{
				"host": 123,
				"port": 8000.0,
			},
			errorString: "net.tcp.client host error: not a string",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			registration, ok := module.ModuleRegistry["net.tcp.client"]
			if !ok {
				t.Fatalf("net.tcp.client module not registered")
			}

			moduleInstance, err := registration.New(config.ModuleConfig{
				Id:     "test",
				Type:   "net.tcp.client",
				Params: test.params,
			})

			if err != nil {
				if test.errorString != err.Error() {
					t.Fatalf("net.tcp.client got error '%s', expected '%s'", err.Error(), test.errorString)
				}
				return
			}

			err = moduleInstance.Start(t.Context())

			if err == nil {
				t.Fatalf("net.tcp.client expected to fail")
			}

			if err.Error() != test.errorString {
				t.Fatalf("net.tcp.client got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
