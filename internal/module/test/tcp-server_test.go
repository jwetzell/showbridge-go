package module_test

import (
	"testing"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/module"
)

func TestTCPServerFromRegistry(t *testing.T) {
	registration, ok := module.ModuleRegistry["net.tcp.server"]
	if !ok {
		t.Fatalf("net.tcp.server module not registered")
	}

	moduleInstance, err := registration.New(config.ModuleConfig{
		Id:   "test",
		Type: "net.tcp.server",
		Params: map[string]any{
			"port":    8000.0,
			"framing": "LF",
		},
	})

	if err != nil {
		t.Fatalf("failed to create net.tcp.server module: %s", err)
	}

	if moduleInstance.Id() != "test" {
		t.Fatalf("net.tcp.server module has wrong id: %s", moduleInstance.Id())
	}

	if moduleInstance.Type() != "net.tcp.server" {
		t.Fatalf("net.tcp.server module has wrong type: %s", moduleInstance.Type())
	}
}

func TestBadTCPServer(t *testing.T) {
	tests := []struct {
		name        string
		params      map[string]any
		errorString string
	}{
		{
			name: "no port param",
			params: map[string]any{
				"framing": "LF",
			},
			errorString: "net.tcp.server requires a port parameter",
		},
		{
			name: "non-number port param",
			params: map[string]any{
				"port":    "8000",
				"framing": "LF",
			},
			errorString: "net.tcp.server port must be a number",
		},
		{
			name: "no framing param",
			params: map[string]any{
				"port": 8000.0,
			},
			errorString: "net.tcp.server requires a framing parameter",
		},
		{
			name: "non-string framing param",
			params: map[string]any{
				"port":    8000.0,
				"framing": 1,
			},
			errorString: "net.tcp.server framing method must be a string",
		},
		{
			name: "unkown framing method",
			params: map[string]any{
				"port":    8000.0,
				"framing": "asdfasdfasdfasdflkj",
			},
			errorString: "net.tcp.server unknown framing method: asdfasdfasdfasdflkj",
		},
		{
			name: "non-string ip param",
			params: map[string]any{
				"port":    8000.0,
				"framing": "LF",
				"ip":      123,
			},
			errorString: "net.tcp.server ip must be a string",
		},
		{
			name: "invalid addr",
			params: map[string]any{
				"ip":      "127.0.0.",
				"port":    8000.0,
				"framing": "LF",
			},
			errorString: "lookup 127.0.0.: no such host",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			registration, ok := module.ModuleRegistry["net.tcp.server"]
			if !ok {
				t.Fatalf("net.tcp.server module not registered")
			}

			moduleInstance, err := registration.New(config.ModuleConfig{
				Id:     "test",
				Type:   "net.tcp.server",
				Params: test.params,
			})

			if err != nil {
				if test.errorString != err.Error() {
					t.Fatalf("net.tcp.server got error '%s', expected '%s'", err.Error(), test.errorString)
				}
				return
			}

			err = moduleInstance.Start(t.Context())

			if err == nil {
				t.Fatalf("net.tcp.server expected to fail")
			}

			if err.Error() != test.errorString {
				t.Fatalf("net.tcp.server got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
