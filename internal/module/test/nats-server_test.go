package module_test

import (
	"testing"

	"github.com/jwetzell/showbridge-go/config"
	"github.com/jwetzell/showbridge-go/internal/module"
)

func TestGoodNATSServer(t *testing.T) {

	testCases := []struct {
		name   string
		params map[string]any
	}{
		{
			name: "minimal config",
			params: map[string]any{
				"port": 4222,
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {

			registration, ok := module.GetModuleRegistration("nats.server")
			if !ok {
				t.Fatalf("nats.server module not registered")
			}

			_, err := registration.New(config.ModuleConfig{
				Id:     "test",
				Type:   "nats.server",
				Params: test.params,
			})

			if err != nil {
				t.Fatalf("nats.server failed to create module: %s", err)
			}
		})
	}
}

func TestBadNATSServer(t *testing.T) {
	tests := []struct {
		name        string
		params      map[string]any
		errorString string
	}{
		{
			name: "non-string ip",
			params: map[string]any{
				"ip": 123,
			},
			errorString: "nats.server ip error: not a string",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			registration, ok := module.GetModuleRegistration("nats.server")
			if !ok {
				t.Fatalf("nats.server module not registered")
			}

			moduleInstance, err := registration.New(config.ModuleConfig{
				Id:     "test",
				Type:   "nats.server",
				Params: test.params,
			})

			if err != nil {
				if test.errorString != err.Error() {
					t.Fatalf("nats.server got error '%s', expected '%s'", err.Error(), test.errorString)
				}
				return
			}

			err = moduleInstance.Start(t.Context(), nil)

			if err == nil {
				t.Fatalf("nats.server expected to fail")
			}

			if err.Error() != test.errorString {
				t.Fatalf("nats.server got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
