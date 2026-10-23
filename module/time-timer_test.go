package module_test

import (
	"testing"

	"github.com/jwetzell/showbridge-go/config"
	"github.com/jwetzell/showbridge-go/module"
)

// TODO(jwetzell): figure out how to test this
func TestGoodTimeTimer(t *testing.T) {

	testCases := []struct {
		name   string
		params map[string]any
	}{
		{
			name: "minimal config",
			params: map[string]any{
				"duration": 2000,
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {

			registration, ok := module.GetModuleRegistration("time.timer")
			if !ok {
				t.Fatalf("time.timer module not registered")
			}

			_, err := registration.New(config.ModuleConfig{
				Id:     "test",
				Type:   "time.timer",
				Params: test.params,
			})

			if err != nil {
				t.Fatalf("time.timer failed to create module: %s", err)
			}
		})
	}
}

func TestBadTimeTimer(t *testing.T) {
	tests := []struct {
		name        string
		params      map[string]any
		errorString string
	}{
		{
			name:        "no duration param",
			params:      map[string]any{},
			errorString: "time.timer duration error: not found",
		},
		{
			name: "non-number duration param",
			params: map[string]any{
				"duration": "8000",
			},
			errorString: "time.timer duration error: not a number",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			registration, ok := module.GetModuleRegistration("time.timer")
			if !ok {
				t.Fatalf("time.timer module not registered")
			}

			moduleInstance, err := registration.New(config.ModuleConfig{
				Id:     "test",
				Type:   "time.timer",
				Params: test.params,
			})

			if err != nil {
				if test.errorString != err.Error() {
					t.Fatalf("time.timer got error '%s', expected '%s'", err.Error(), test.errorString)
				}
				return
			}

			err = moduleInstance.Start(t.Context(), nil)

			if err == nil {
				t.Fatalf("time.timer expected to fail")
			}

			if err.Error() != test.errorString {
				t.Fatalf("time.timer got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
