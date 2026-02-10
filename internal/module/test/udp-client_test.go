package module_test

import (
	"testing"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/module"
)

func TestUDPClientFromRegistry(t *testing.T) {
	registration, ok := module.ModuleRegistry["net.udp.client"]
	if !ok {
		t.Fatalf("udp.client module not registered")
	}

	moduleInstance, err := registration.New(config.ModuleConfig{
		Id:   "test",
		Type: "net.udp.client",
		Params: map[string]any{
			"host":    "localhost",
			"port":    8000.0,
			"framing": "LF",
		},
	})

	if err != nil {
		t.Fatalf("failed to create net.udp.client module: %s", err)
	}

	if moduleInstance.Id() != "test" {
		t.Fatalf("net.udp.client module has wrong id: %s", moduleInstance.Id())
	}

	if moduleInstance.Type() != "net.udp.client" {
		t.Fatalf("net.udp.client module has wrong type: %s", moduleInstance.Type())
	}
}

func TestBadUDPClient(t *testing.T) {
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
			errorString: "net.udp.client requires a port parameter",
		},
		{
			name: "non-number port param",
			params: map[string]any{
				"host": "localhost",
				"port": "8000",
			},
			errorString: "net.udp.client port must be a number",
		},
		{
			name: "no host param",
			params: map[string]any{
				"port": 8000.0,
			},
			errorString: "net.udp.client requires a host parameter",
		},
		{
			name: "non-string host param",
			params: map[string]any{
				"host": 123,
				"port": 8000.0,
			},
			errorString: "net.udp.client host must be a string",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			registration, ok := module.ModuleRegistry["net.udp.client"]
			if !ok {
				t.Fatalf("net.udp.client module not registered")
			}

			moduleInstance, err := registration.New(config.ModuleConfig{
				Id:     "test",
				Type:   "net.udp.client",
				Params: test.params,
			})

			if err != nil {
				if test.errorString != err.Error() {
					t.Fatalf("net.udp.client got error '%s', expected '%s'", err.Error(), test.errorString)
				}
				return
			}

			err = moduleInstance.Start(t.Context())

			if err == nil {
				t.Fatalf("net.udp.client expected to fail")
			}

			if err.Error() != test.errorString {
				t.Fatalf("net.udp.client got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
