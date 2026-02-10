package module_test

import (
	"testing"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/module"
)

func TestTimeIntervalFromRegistry(t *testing.T) {
	registration, ok := module.ModuleRegistry["time.interval"]
	if !ok {
		t.Fatalf("time.interval module not registered")
	}

	moduleInstance, err := registration.New(config.ModuleConfig{
		Id:   "test",
		Type: "time.interval",
		Params: map[string]any{
			"duration": 1000.0,
		},
	})

	if err != nil {
		t.Fatalf("failed to create time.interval module: %s", err)
	}

	if moduleInstance.Id() != "test" {
		t.Fatalf("time.interval module has wrong id: %s", moduleInstance.Id())
	}

	if moduleInstance.Type() != "time.interval" {
		t.Fatalf("time.interval module has wrong type: %s", moduleInstance.Type())
	}
}

func TestBadTimeInterval(t *testing.T) {
	tests := []struct {
		name        string
		params      map[string]any
		errorString string
	}{
		{
			name:        "no duration param",
			params:      map[string]any{},
			errorString: "time.interval requires a duration parameter",
		},
		{
			name: "non-number duration param",
			params: map[string]any{
				"duration": "8000",
			},
			errorString: "time.interval duration must be a number",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			registration, ok := module.ModuleRegistry["time.interval"]
			if !ok {
				t.Fatalf("time.interval module not registered")
			}

			moduleInstance, err := registration.New(config.ModuleConfig{
				Id:     "test",
				Type:   "time.interval",
				Params: test.params,
			})

			if err != nil {
				if test.errorString != err.Error() {
					t.Fatalf("time.interval got error '%s', expected '%s'", err.Error(), test.errorString)
				}
				return
			}

			err = moduleInstance.Start(t.Context())

			if err == nil {
				t.Fatalf("time.interval expected to fail")
			}

			if err.Error() != test.errorString {
				t.Fatalf("time.interval got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
