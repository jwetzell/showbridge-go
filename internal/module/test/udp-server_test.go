package module_test

import (
	"testing"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/module"
)

func TestUDPServerFromRegistry(t *testing.T) {
	registration, ok := module.ModuleRegistry["net.udp.server"]
	if !ok {
		t.Fatalf("net.udp.server module not registered")
	}

	moduleInstance, err := registration.New(config.ModuleConfig{
		Id:   "test",
		Type: "net.udp.server",
		Params: map[string]any{
			"port": 8000.0,
		},
	})

	if err != nil {
		t.Fatalf("failed to create udp.server module: %s", err)
	}

	if moduleInstance.Id() != "test" {
		t.Fatalf("udp.server module has wrong id: %s", moduleInstance.Id())
	}

	if moduleInstance.Type() != "net.udp.server" {
		t.Fatalf("net.udp.server module has wrong type: %s", moduleInstance.Type())
	}
}

func TestBadUDPServer(t *testing.T) {
	tests := []struct {
		name        string
		params      map[string]any
		errorString string
	}{
		{
			name:        "no port param",
			params:      map[string]any{},
			errorString: "net.udp.server port error: not found",
		},
		{
			name: "non-number port param",
			params: map[string]any{
				"port": "8000",
			},
			errorString: "net.udp.server port error: not a number",
		},
		{
			name: "non-string ip param",
			params: map[string]any{
				"port": 8000.0,
				"ip":   123,
			},
			errorString: "net.udp.server ip error: not a string",
		},
		{
			name: "non-number bufferSize param",
			params: map[string]any{
				"port":       8000.0,
				"bufferSize": "1024",
			},
			errorString: "net.udp.server bufferSize error: not a number",
		},
		{
			name: "invalid addr",
			params: map[string]any{
				"ip":   "127.0.0.",
				"port": 8000.0,
			},
			errorString: "lookup 127.0.0.: no such host",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			registration, ok := module.ModuleRegistry["net.udp.server"]
			if !ok {
				t.Fatalf("net.udp.server module not registered")
			}

			moduleInstance, err := registration.New(config.ModuleConfig{
				Id:     "test",
				Type:   "net.udp.server",
				Params: test.params,
			})

			if err != nil {
				if test.errorString != err.Error() {
					t.Fatalf("net.udp.server got error '%s', expected '%s'", err.Error(), test.errorString)
				}
				return
			}

			err = moduleInstance.Start(t.Context())

			if err == nil {
				t.Fatalf("net.udp.server expected to fail")
			}

			if err.Error() != test.errorString {
				t.Fatalf("net.udp.server got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
