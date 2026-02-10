package module_test

import (
	"testing"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/module"
)

func TestMQTTClientFromRegistry(t *testing.T) {
	registration, ok := module.ModuleRegistry["mqtt.client"]
	if !ok {
		t.Fatalf("mqtt.client module not registered")
	}

	moduleInstance, err := registration.New(config.ModuleConfig{
		Id:   "test",
		Type: "mqtt.client",
		Params: map[string]any{
			"broker":   "mqtt://localhost:1883",
			"topic":    "test/topic",
			"clientId": "test",
		},
	})

	if err != nil {
		t.Fatalf("failed to create mqtt.client module: %s", err)
	}

	if moduleInstance.Id() != "test" {
		t.Fatalf("mqtt.client module has wrong id: %s", moduleInstance.Id())
	}

	if moduleInstance.Type() != "mqtt.client" {
		t.Fatalf("mqtt.client module has wrong type: %s", moduleInstance.Type())
	}
}

func TestBadMQTTClient(t *testing.T) {
	tests := []struct {
		name        string
		params      map[string]any
		errorString string
	}{
		{
			name: "no broker param",
			params: map[string]any{
				"topic":    "test/topic",
				"clientId": "test",
			},
			errorString: "mqtt.client requires a broker parameter",
		},
		{
			name: "non-string broker",
			params: map[string]any{
				"broker":   123,
				"topic":    "test/topic",
				"clientId": "test",
			},
			errorString: "mqtt.client broker must be a string",
		},
		{
			name: "no topic param",
			params: map[string]any{
				"broker":   "mqtt://localhost:1883",
				"clientId": "test",
			},
			errorString: "mqtt.client requires a topic parameter",
		},
		{
			name: "non-string topic",
			params: map[string]any{
				"broker":   "mqtt://localhost:1883",
				"topic":    123,
				"clientId": "test",
			},
			errorString: "mqtt.client topic must be a string",
		},
		{
			name: "no clientId param",
			params: map[string]any{
				"broker": "mqtt://localhost:1883",
				"topic":  "test/topic",
			},
			errorString: "mqtt.client requires a clientId parameter",
		},
		{
			name: "non-string clientId",
			params: map[string]any{
				"broker":   "mqtt://localhost:1883",
				"topic":    "test/topic",
				"clientId": 123,
			},
			errorString: "mqtt.client clientId must be a string",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			registration, ok := module.ModuleRegistry["mqtt.client"]
			if !ok {
				t.Fatalf("mqtt.client module not registered")
			}

			moduleInstance, err := registration.New(config.ModuleConfig{
				Id:     "test",
				Type:   "mqtt.client",
				Params: test.params,
			})

			if err != nil {
				if test.errorString != err.Error() {
					t.Fatalf("mqtt.client got error '%s', expected '%s'", err.Error(), test.errorString)
				}
				return
			}

			err = moduleInstance.Start(t.Context())

			if err == nil {
				t.Fatalf("mqtt.client expected to fail")
			}

			if err.Error() != test.errorString {
				t.Fatalf("mqtt.client got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
