package module_test

import (
	"testing"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/module"
)

func TestDbSqliteFromRegistry(t *testing.T) {
	registration, ok := module.ModuleRegistry["db.sqlite"]
	if !ok {
		t.Fatalf("db.sqlite module not registered")
	}

	moduleInstance, err := registration.New(config.ModuleConfig{
		Id:   "test",
		Type: "db.sqlite",
		Params: map[string]any{
			"dsn": ":memory:",
		},
	})

	if err != nil {
		t.Fatalf("failed to create db.sqlite module: %s", err)
	}

	if moduleInstance.Id() != "test" {
		t.Fatalf("db.sqlite module has wrong id: %s", moduleInstance.Id())
	}

	if moduleInstance.Type() != "db.sqlite" {
		t.Fatalf("db.sqlite module has wrong type: %s", moduleInstance.Type())
	}
}

func TestBadDbSqlite(t *testing.T) {
	tests := []struct {
		name        string
		params      map[string]any
		errorString string
	}{
		{
			name:        "no dsn param",
			params:      map[string]any{},
			errorString: "db.sqlite dsn error: not found",
		},
		{
			name:        "non-string dsn",
			params:      map[string]any{"dsn": 123},
			errorString: "db.sqlite dsn error: not a string",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			registration, ok := module.ModuleRegistry["db.sqlite"]
			if !ok {
				t.Fatalf("db.sqlite module not registered")
			}

			moduleInstance, err := registration.New(config.ModuleConfig{
				Id:     "test",
				Type:   "db.sqlite",
				Params: test.params,
			})

			if err != nil {
				if test.errorString != err.Error() {
					t.Fatalf("db.sqlite got error '%s', expected '%s'", err.Error(), test.errorString)
				}
				return
			}

			err = moduleInstance.Start(t.Context())

			if err == nil {
				t.Fatalf("db.sqlite expected to fail")
			}

			if err.Error() != test.errorString {
				t.Fatalf("db.sqlite got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
