package module_test

import (
	"testing"

	"github.com/jwetzell/showbridge-go/config"
	"github.com/jwetzell/showbridge-go/module"
)

// TODO(jwetzell): figure out how to test this
func TestGoodPostgresClient(t *testing.T) {

	testCases := []struct {
		name   string
		params map[string]any
	}{}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {

			registration, ok := module.GetModuleRegistration("postgres.client")
			if !ok {
				t.Fatalf("postgres.client module not registered")
			}

			_, err := registration.New(config.ModuleConfig{
				Id:     "test",
				Type:   "postgres.client",
				Params: test.params,
			})

			if err != nil {
				t.Fatalf("postgres.client failed to create module: %s", err)
			}
		})
	}
}

func TestBadPostgresClient(t *testing.T) {
	tests := []struct {
		name        string
		params      map[string]any
		errorString string
	}{
		{
			name:        "no url param",
			params:      map[string]any{},
			errorString: "postgres.client url error: not found",
		},
		{
			name:        "non-string url",
			params:      map[string]any{"url": 123},
			errorString: "postgres.client url error: not a string",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			registration, ok := module.GetModuleRegistration("postgres.client")
			if !ok {
				t.Fatalf("postgres.client module not registered")
			}

			moduleInstance, err := registration.New(config.ModuleConfig{
				Id:     "test",
				Type:   "postgres.client",
				Params: test.params,
			})

			if err != nil {
				if test.errorString != err.Error() {
					t.Fatalf("postgres.client got error '%s', expected '%s'", err.Error(), test.errorString)
				}
				return
			}

			err = moduleInstance.Start(t.Context(), nil)

			if err == nil {
				t.Fatalf("postgres.client expected to fail")
			}

			if err.Error() != test.errorString {
				t.Fatalf("postgres.client got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
