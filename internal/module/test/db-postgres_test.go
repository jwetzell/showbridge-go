package module_test

import (
	"testing"
	"time"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/module"
)

func TestDbPostgresFromRegistry(t *testing.T) {
	registration, ok := module.GetModuleRegistration("db.postgres")
	if !ok {
		t.Fatalf("db.postgres module not registered")
	}

	moduleInstance, err := registration.New(config.ModuleConfig{
		Id:   "test",
		Type: "db.postgres",
		Params: map[string]any{
			"url": "postgres://localhost:5432",
		},
	})

	if err != nil {
		t.Fatalf("failed to create db.postgres module: %s", err)
	}

	if moduleInstance.Id() != "test" {
		t.Fatalf("db.postgres module has wrong id: %s", moduleInstance.Id())
	}

	if moduleInstance.Type() != "db.postgres" {
		t.Fatalf("db.postgres module has wrong type: %s", moduleInstance.Type())
	}
}

func TestGoodDbPostgres(t *testing.T) {

	testCases := []struct {
		name   string
		params map[string]any
	}{}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {

			registration, ok := module.GetModuleRegistration("db.postgres")
			if !ok {
				t.Fatalf("db.postgres module not registered")
			}

			moduleInstance, err := registration.New(config.ModuleConfig{
				Id:     "test",
				Type:   "db.postgres",
				Params: test.params,
			})

			if err != nil {
				t.Fatalf("db.postgres failed to create module: %s", err)
			}
			// TODO(jwetzell) this is kind of hacky
			go func() {
				time.Sleep(1 * time.Second)
				moduleInstance.Stop()
			}()
			err = moduleInstance.Start(t.Context(), nil)

			if err != nil {
				t.Fatalf("db.postgres failed to start: %s", err)
			}
		})
	}
}

func TestBadDbPostgres(t *testing.T) {
	tests := []struct {
		name        string
		params      map[string]any
		errorString string
	}{
		{
			name:        "no url param",
			params:      map[string]any{},
			errorString: "db.postgres url error: not found",
		},
		{
			name:        "non-string url",
			params:      map[string]any{"url": 123},
			errorString: "db.postgres url error: not a string",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			registration, ok := module.GetModuleRegistration("db.postgres")
			if !ok {
				t.Fatalf("db.postgres module not registered")
			}

			moduleInstance, err := registration.New(config.ModuleConfig{
				Id:     "test",
				Type:   "db.postgres",
				Params: test.params,
			})

			if err != nil {
				if test.errorString != err.Error() {
					t.Fatalf("db.postgres got error '%s', expected '%s'", err.Error(), test.errorString)
				}
				return
			}

			err = moduleInstance.Start(t.Context(), nil)

			if err == nil {
				t.Fatalf("db.postgres expected to fail")
			}

			if err.Error() != test.errorString {
				t.Fatalf("db.postgres got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
