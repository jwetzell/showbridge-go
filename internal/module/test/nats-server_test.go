package module_test

import (
	"testing"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/module"
)

func TestNATSServerFromRegistry(t *testing.T) {
	registration, ok := module.ModuleRegistry["nats.server"]
	if !ok {
		t.Fatalf("nats.server module not registered")
	}

	moduleInstance, err := registration.New(config.ModuleConfig{
		Id:   "test",
		Type: "nats.server",
		Params: map[string]any{
			"ip":   "127.0.0.1",
			"port": 4222,
		},
	})

	if err != nil {
		t.Fatalf("failed to create nats.server module: %s", err)
	}

	if moduleInstance.Id() != "test" {
		t.Fatalf("nats.server module has wrong id: %s", moduleInstance.Id())
	}

	if moduleInstance.Type() != "nats.server" {
		t.Fatalf("nats.server module has wrong type: %s", moduleInstance.Type())
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
			errorString: "nats.server ip must be a string",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			registration, ok := module.ModuleRegistry["nats.server"]
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

			err = moduleInstance.Start(t.Context())

			if err == nil {
				t.Fatalf("nats.server expected to fail")
			}

			if err.Error() != test.errorString {
				t.Fatalf("nats.server got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
