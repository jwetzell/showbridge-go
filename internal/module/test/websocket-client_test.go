package module_test

import (
	"testing"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/module"
)

func TestWebSocketClientFromRegistry(t *testing.T) {
	registration, ok := module.GetModuleRegistration("websocket.client")
	if !ok {
		t.Fatalf("websocket.client module not registered")
	}

	moduleInstance, err := registration.New(config.ModuleConfig{
		Id:   "test",
		Type: "websocket.client",
		Params: map[string]any{
			"url": "ws://localhost",
		},
	})

	if err != nil {
		t.Fatalf("failed to create websocket.client module: %s", err)
	}

	if moduleInstance.Id() != "test" {
		t.Fatalf("websocket.client module has wrong id: %s", moduleInstance.Id())
	}

	if moduleInstance.Type() != "websocket.client" {
		t.Fatalf("websocket.client module has wrong type: %s", moduleInstance.Type())
	}
}

func TestBadWebSocketClient(t *testing.T) {
	tests := []struct {
		name        string
		params      map[string]any
		errorString string
	}{
		{
			name:        "no url param",
			params:      map[string]any{},
			errorString: "websocket.client url error: not found",
		},
		{
			name: "non-string url param",
			params: map[string]any{
				"url": 123,
			},
			errorString: "websocket.client url error: not a string",
		},
		{
			name: "invalid url param",
			params: map[string]any{
				"url": "invalid-url",
			},
			errorString: "websocket.client url error: scheme must be ws or wss",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			registration, ok := module.GetModuleRegistration("websocket.client")
			if !ok {
				t.Fatalf("websocket.client module not registered")
			}

			moduleInstance, err := registration.New(config.ModuleConfig{
				Id:     "test",
				Type:   "websocket.client",
				Params: test.params,
			})

			if err != nil {
				if test.errorString != err.Error() {
					t.Fatalf("websocket.client got error '%s', expected '%s'", err.Error(), test.errorString)
				}
				return
			}

			err = moduleInstance.Start(t.Context(), nil)

			if err == nil {
				t.Fatalf("websocket.client expected to fail")
			}

			if err.Error() != test.errorString {
				t.Fatalf("websocket.client got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
