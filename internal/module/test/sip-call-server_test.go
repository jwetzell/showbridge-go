package module_test

import (
	"testing"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/module"
)

func TestSIPCallServerFromRegistry(t *testing.T) {
	registration, ok := module.ModuleRegistry["sip.call.server"]
	if !ok {
		t.Fatalf("sip.call.server module not registered")
	}

	moduleInstance, err := registration.New(config.ModuleConfig{
		Id:   "test",
		Type: "sip.call.server",
	})

	if err != nil {
		t.Fatalf("failed to create sip.call.server module: %s", err)
	}

	if moduleInstance.Id() != "test" {
		t.Fatalf("sip.call.server module has wrong id: %s", moduleInstance.Id())
	}

	if moduleInstance.Type() != "sip.call.server" {
		t.Fatalf("sip.call.server module has wrong type: %s", moduleInstance.Type())
	}
}

func TestBadSIPCallServer(t *testing.T) {
	tests := []struct {
		name        string
		params      map[string]any
		errorString string
	}{
		{
			name: "non-number port param",
			params: map[string]any{
				"port": "8000",
			},
			errorString: "sip.call.server port must be a number",
		},
		{
			name: "non-string ip param",
			params: map[string]any{
				"ip": 123,
			},
			errorString: "sip.call.server ip must be a string",
		},
		{
			name: "non-string transport param",
			params: map[string]any{
				"transport": 123,
			},
			errorString: "sip.call.server transport must be a string",
		},
		{
			name: "non-string userAgent param",
			params: map[string]any{
				"userAgent": 123,
			},
			errorString: "sip.call.server userAgent must be a string",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			registration, ok := module.ModuleRegistry["sip.call.server"]
			if !ok {
				t.Fatalf("sip.call.server module not registered")
			}

			moduleInstance, err := registration.New(config.ModuleConfig{
				Id:     "test",
				Type:   "sip.call.server",
				Params: test.params,
			})

			if err != nil {
				if test.errorString != err.Error() {
					t.Fatalf("sip.call.server got error '%s', expected '%s'", err.Error(), test.errorString)
				}
				return
			}

			err = moduleInstance.Start(t.Context())

			if err == nil {
				t.Fatalf("sip.call.server expected to fail")
			}

			if err.Error() != test.errorString {
				t.Fatalf("sip.call.server got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
