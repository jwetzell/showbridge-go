package module_test

import (
	"testing"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/module"
)

func TestSerialClientFromRegistry(t *testing.T) {
	registration, ok := module.ModuleRegistry["serial.client"]
	if !ok {
		t.Fatalf("serial.client module not registered")
	}

	moduleInstance, err := registration.New(config.ModuleConfig{
		Id:   "test",
		Type: "serial.client",
		Params: map[string]any{
			"port":     "/dev/ttyUSB0",
			"framing":  "LF",
			"baudRate": 9600.0,
		},
	})

	if err != nil {
		t.Fatalf("failed to create serial.client module: %s", err)
	}

	if moduleInstance.Id() != "test" {
		t.Fatalf("serial.client module has wrong id: %s", moduleInstance.Id())
	}

	if moduleInstance.Type() != "serial.client" {
		t.Fatalf("serial.client module has wrong type: %s", moduleInstance.Type())
	}
}

func TestBadSerialClient(t *testing.T) {
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
			errorString: "serial.client requires a port parameter",
		},
		{
			name: "non-string port param",
			params: map[string]any{
				"port":    8000,
				"framing": "LF",
			},
			errorString: "serial.client port must be a string",
		},
		{
			name: "no framing param",
			params: map[string]any{
				"port": "/dev/ttyUSB0",
			},
			errorString: "serial.client requires a framing parameter",
		},
		{
			name: "non-string framing param",
			params: map[string]any{
				"port":    "/dev/ttyUSB0",
				"framing": 1,
			},
			errorString: "serial.client framing method must be a string",
		},
		{
			name: "unkown framing method",
			params: map[string]any{
				"port":    "/dev/ttyUSB0",
				"framing": "asdfasdfasdfasdflkj",
			},
			errorString: "serial.client unknown framing method: asdfasdfasdfasdflkj",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			registration, ok := module.ModuleRegistry["serial.client"]
			if !ok {
				t.Fatalf("serial.client module not registered")
			}

			moduleInstance, err := registration.New(config.ModuleConfig{
				Id:     "test",
				Type:   "serial.client",
				Params: test.params,
			})

			if err != nil {
				if test.errorString != err.Error() {
					t.Fatalf("serial.client got error '%s', expected '%s'", err.Error(), test.errorString)
				}
				return
			}

			err = moduleInstance.Start(t.Context())

			if err == nil {
				t.Fatalf("serial.client expected to fail")
			}

			if err.Error() != test.errorString {
				t.Fatalf("serial.client got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
