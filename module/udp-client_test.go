package module_test

import (
	"testing"

	"github.com/jwetzell/showbridge-go/config"
	"github.com/jwetzell/showbridge-go/module"
)

// TODO(jwetzell): figure out how to test this
func TestGoodUDPClient(t *testing.T) {

	testCases := []struct {
		name   string
		params map[string]any
	}{
		{
			name: "minimal config",
			params: map[string]any{
				"host": "localhost",
				"port": 8000,
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {

			registration, ok := module.GetModuleRegistration("net.udp.client")
			if !ok {
				t.Fatalf("net.udp.client module not registered")
			}

			_, err := registration.New(config.ModuleConfig{
				Id:     "test",
				Type:   "net.udp.client",
				Params: test.params,
			})

			if err != nil {
				t.Fatalf("net.udp.client failed to create module: %s", err)
			}
		})
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
			errorString: "net.udp.client port error: not found",
		},
		{
			name: "non-number port param",
			params: map[string]any{
				"host": "localhost",
				"port": "8000",
			},
			errorString: "net.udp.client port error: not a number",
		},
		{
			name: "no host param",
			params: map[string]any{
				"port": 8000,
			},
			errorString: "net.udp.client host error: not found",
		},
		{
			name: "non-string host param",
			params: map[string]any{
				"host": 123,
				"port": 8000,
			},
			errorString: "net.udp.client host error: not a string",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			registration, ok := module.GetModuleRegistration("net.udp.client")
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

			err = moduleInstance.Start(t.Context(), nil)

			if err == nil {
				t.Fatalf("net.udp.client expected to fail")
			}

			if err.Error() != test.errorString {
				t.Fatalf("net.udp.client got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
