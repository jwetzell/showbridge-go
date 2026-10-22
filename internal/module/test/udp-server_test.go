package module_test

import (
	"testing"

	"github.com/jwetzell/showbridge-go/config"
	"github.com/jwetzell/showbridge-go/internal/module"
)

// TODO(jwetzell): figure out how to test this
func TestGoodUDPServer(t *testing.T) {

	testCases := []struct {
		name   string
		params map[string]any
	}{
		{
			name: "minimal config",
			params: map[string]any{
				"port": 8000,
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {

			registration, ok := module.GetModuleRegistration("net.udp.server")
			if !ok {
				t.Fatalf("net.udp.server module not registered")
			}

			_, err := registration.New(config.ModuleConfig{
				Id:     "test",
				Type:   "net.udp.server",
				Params: test.params,
			})

			if err != nil {
				t.Fatalf("net.udp.server failed to create module: %s", err)
			}
		})
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
				"port": 8000,
				"ip":   123,
			},
			errorString: "net.udp.server ip error: not a string",
		},
		{
			name: "non-number bufferSize param",
			params: map[string]any{
				"port":       8000,
				"bufferSize": "1024",
			},
			errorString: "net.udp.server bufferSize error: not a number",
		},
		{
			name: "invalid addr",
			params: map[string]any{
				"ip":   "127.0.0.",
				"port": 8000,
			},
			errorString: "lookup 127.0.0.: no such host",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			registration, ok := module.GetModuleRegistration("net.udp.server")
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

			err = moduleInstance.Start(t.Context(), nil)

			if err == nil {
				t.Fatalf("net.udp.server expected to fail")
			}

			if err.Error() != test.errorString {
				t.Fatalf("net.udp.server got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
