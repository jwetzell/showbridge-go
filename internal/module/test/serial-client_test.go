//go:build cgo

package module_test

import (
	"testing"

	"github.com/jwetzell/showbridge-go/config"
	"github.com/jwetzell/showbridge-go/internal/module"
)

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
			errorString: "serial.client port error: not found",
		},
		{
			name: "non-string port param",
			params: map[string]any{
				"port":    8000,
				"framing": "LF",
			},
			errorString: "serial.client port error: not a string",
		},
		{
			name: "no framing param",
			params: map[string]any{
				"port": "/dev/ttyUSB0",
			},
			errorString: "serial.client framing error: not found",
		},
		{
			name: "non-string framing param",
			params: map[string]any{
				"port":    "/dev/ttyUSB0",
				"framing": 1,
			},
			errorString: "serial.client framing error: not a string",
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

			registration, ok := module.GetModuleRegistration("serial.client")
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

			err = moduleInstance.Start(t.Context(), nil)

			if err == nil {
				t.Fatalf("serial.client expected to fail")
			}

			if err.Error() != test.errorString {
				t.Fatalf("serial.client got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
