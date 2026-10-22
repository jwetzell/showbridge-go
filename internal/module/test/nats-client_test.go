package module_test

import (
	"testing"

	"github.com/jwetzell/showbridge-go/config"
	"github.com/jwetzell/showbridge-go/internal/module"
)

func TestBadNATSClient(t *testing.T) {
	tests := []struct {
		name        string
		params      map[string]any
		errorString string
	}{
		{
			name: "no url param",
			params: map[string]any{
				"subject": "test/subject",
			},
			errorString: "nats.client url error: not found",
		},
		{
			name: "non-string url",
			params: map[string]any{
				"url":     123,
				"subject": "test/subject",
			},
			errorString: "nats.client url error: not a string",
		},
		{
			name: "no subject param",
			params: map[string]any{
				"url": "nats://127.0.0.1:4222",
			},
			errorString: "nats.client subject error: not found",
		},
		{
			name: "non-string subject",
			params: map[string]any{
				"url":     "nats://127.0.0.1:4222",
				"subject": 123,
			},
			errorString: "nats.client subject error: not a string",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			registration, ok := module.GetModuleRegistration("nats.client")
			if !ok {
				t.Fatalf("nats.client module not registered")
			}

			moduleInstance, err := registration.New(config.ModuleConfig{
				Id:     "test",
				Type:   "nats.client",
				Params: test.params,
			})

			if err != nil {
				if test.errorString != err.Error() {
					t.Fatalf("nats.client got error '%s', expected '%s'", err.Error(), test.errorString)
				}
				return
			}

			err = moduleInstance.Start(t.Context(), nil)

			if err == nil {
				t.Fatalf("nats.client expected to fail")
			}

			if err.Error() != test.errorString {
				t.Fatalf("nats.client got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
