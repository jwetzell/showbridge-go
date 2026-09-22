package module_test

import (
	"testing"

	"github.com/jwetzell/showbridge-go/config"
	"github.com/jwetzell/showbridge-go/internal/module"
)

func TestTimeCronFromRegistry(t *testing.T) {
	registration, ok := module.GetModuleRegistration("time.cron")
	if !ok {
		t.Fatalf("time.cron module not registered")
	}

	moduleInstance, err := registration.New(config.ModuleConfig{
		Id:   "test",
		Type: "time.cron",
		Params: map[string]any{
			"cron": "* * * * *",
		},
	})

	if err != nil {
		t.Fatalf("failed to create time.cron module: %s", err)
	}

	if moduleInstance.Id() != "test" {
		t.Fatalf("time.cron module has wrong id: %s", moduleInstance.Id())
	}

	if moduleInstance.Type() != "time.cron" {
		t.Fatalf("time.cron module has wrong type: %s", moduleInstance.Type())
	}
}

func TestGoodTimeCron(t *testing.T) {

	testCases := []struct {
		name   string
		params map[string]any
	}{
		{
			name: "minimal config",
			params: map[string]any{
				"cron": "* * * * *",
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {

			registration, ok := module.GetModuleRegistration("time.cron")
			if !ok {
				t.Fatalf("time.cron module not registered")
			}

			_, err := registration.New(config.ModuleConfig{
				Id:     "test",
				Type:   "time.cron",
				Params: test.params,
			})

			if err != nil {
				t.Fatalf("time.cron failed to create module: %s", err)
			}
			// TODO(jwetzell) figure out how to test the actual cron execution
		})
	}
}

func TestBadTimeCron(t *testing.T) {
	tests := []struct {
		name        string
		params      map[string]any
		errorString string
	}{
		{
			name:        "no cron param",
			params:      map[string]any{},
			errorString: "time.cron cron error: not found",
		},
		{
			name: "non-string cron param",
			params: map[string]any{
				"cron": 8000,
			},
			errorString: "time.cron cron error: not a string",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			registration, ok := module.GetModuleRegistration("time.cron")
			if !ok {
				t.Fatalf("time.cron module not registered")
			}

			moduleInstance, err := registration.New(config.ModuleConfig{
				Id:     "test",
				Type:   "time.cron",
				Params: test.params,
			})

			if err != nil {
				if test.errorString != err.Error() {
					t.Fatalf("time.cron got error '%s', expected '%s'", err.Error(), test.errorString)
				}
				return
			}

			err = moduleInstance.Start(t.Context(), nil)

			if err == nil {
				t.Fatalf("time.cron expected to fail")
			}

			if err.Error() != test.errorString {
				t.Fatalf("time.cron got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
