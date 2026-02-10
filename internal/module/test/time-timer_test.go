package module_test

import (
	"testing"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/module"
)

func TestTimeTimerFromRegistry(t *testing.T) {
	registration, ok := module.ModuleRegistry["time.timer"]
	if !ok {
		t.Fatalf("time.timer module not registered")
	}

	moduleInstance, err := registration.New(config.ModuleConfig{
		Id:   "test",
		Type: "time.timer",
		Params: map[string]any{
			"duration": 1000.0,
		},
	})

	if err != nil {
		t.Fatalf("failed to create time.timer module: %s", err)
	}

	if moduleInstance.Id() != "test" {
		t.Fatalf("time.timer module has wrong id: %s", moduleInstance.Id())
	}

	if moduleInstance.Type() != "time.timer" {
		t.Fatalf("time.timer module has wrong type: %s", moduleInstance.Type())
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
			errorString: "time.timer requires a duration parameter",
		},
		{
			name: "non-number duration param",
			params: map[string]any{
				"duration": "8000",
			},
			errorString: "time.timer duration must be a number",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			registration, ok := module.ModuleRegistry["time.timer"]
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

			err = moduleInstance.Start(t.Context())

			if err == nil {
				t.Fatalf("time.timer expected to fail")
			}

			if err.Error() != test.errorString {
				t.Fatalf("time.timer got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
