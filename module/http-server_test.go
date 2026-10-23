package module_test

import (
	"testing"

	"github.com/jwetzell/showbridge-go/config"
	"github.com/jwetzell/showbridge-go/module"
)

// TODO(jwetzell): figure out how to test this
func TestGoodHTTPServer(t *testing.T) {

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

			registration, ok := module.GetModuleRegistration("http.server")
			if !ok {
				t.Fatalf("http.server module not registered")
			}

			_, err := registration.New(config.ModuleConfig{
				Id:     "test",
				Type:   "http.server",
				Params: test.params,
			})

			if err != nil {
				t.Fatalf("http.server failed to create module: %s", err)
			}
		})
	}
}

func TestBadHTTPServer(t *testing.T) {
	tests := []struct {
		name        string
		params      map[string]any
		errorString string
	}{
		{
			name:        "no port param",
			params:      map[string]any{},
			errorString: "http.server port error: not found",
		},
		{
			name:        "non-numeric port",
			params:      map[string]any{"port": "3000"},
			errorString: "http.server port error: not a number",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			registration, ok := module.GetModuleRegistration("http.server")
			if !ok {
				t.Fatalf("http.server module not registered")
			}

			moduleInstance, err := registration.New(config.ModuleConfig{
				Id:     "test",
				Type:   "http.server",
				Params: test.params,
			})

			if err != nil {
				if test.errorString != err.Error() {
					t.Fatalf("http.server got error '%s', expected '%s'", err.Error(), test.errorString)
				}
				return
			}

			err = moduleInstance.Start(t.Context(), nil)

			if err == nil {
				t.Fatalf("http.server expected to fail")
			}

			if err.Error() != test.errorString {
				t.Fatalf("http.server got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
