package module_test

import (
	"testing"
	"time"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/module"
)

func TestDbMySQLFromRegistry(t *testing.T) {
	registration, ok := module.GetModuleRegistration("db.mysql")
	if !ok {
		t.Fatalf("db.mysql module not registered")
	}

	moduleInstance, err := registration.New(config.ModuleConfig{
		Id:   "test",
		Type: "db.mysql",
		Params: map[string]any{
			"dsn": "mysql:mysql@tcp(127.0.0.1:3306)/test",
		},
	})

	if err != nil {
		t.Fatalf("failed to create db.mysql module: %s", err)
	}

	if moduleInstance.Id() != "test" {
		t.Fatalf("db.mysql module has wrong id: %s", moduleInstance.Id())
	}

	if moduleInstance.Type() != "db.mysql" {
		t.Fatalf("db.mysql module has wrong type: %s", moduleInstance.Type())
	}
}

func TestGoodDbMySQL(t *testing.T) {

	testCases := []struct {
		name   string
		params map[string]any
	}{}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {

			registration, ok := module.GetModuleRegistration("db.mysql")
			if !ok {
				t.Fatalf("db.mysql module not registered")
			}

			moduleInstance, err := registration.New(config.ModuleConfig{
				Id:     "test",
				Type:   "db.mysql",
				Params: test.params,
			})

			if err != nil {
				t.Fatalf("db.mysql failed to create module: %s", err)
			}
			// TODO(jwetzell) this is kind of hacky
			go func() {
				time.Sleep(1 * time.Second)
				moduleInstance.Stop()
			}()
			err = moduleInstance.Start(t.Context(), nil)

			if err != nil {
				t.Fatalf("db.mysql failed to start: %s", err)
			}
		})
	}
}

func TestBadDbMySQL(t *testing.T) {
	tests := []struct {
		name        string
		params      map[string]any
		errorString string
	}{
		{
			name:        "no dsn param",
			params:      map[string]any{},
			errorString: "db.mysql dsn error: not found",
		},
		{
			name:        "non-string dsn",
			params:      map[string]any{"dsn": 123},
			errorString: "db.mysql dsn error: not a string",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			registration, ok := module.GetModuleRegistration("db.mysql")
			if !ok {
				t.Fatalf("db.mysql module not registered")
			}

			moduleInstance, err := registration.New(config.ModuleConfig{
				Id:     "test",
				Type:   "db.mysql",
				Params: test.params,
			})

			if err != nil {
				if test.errorString != err.Error() {
					t.Fatalf("db.mysql got error '%s', expected '%s'", err.Error(), test.errorString)
				}
				return
			}

			err = moduleInstance.Start(t.Context(), nil)

			if err == nil {
				t.Fatalf("db.mysql expected to fail")
			}

			if err.Error() != test.errorString {
				t.Fatalf("db.mysql got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
