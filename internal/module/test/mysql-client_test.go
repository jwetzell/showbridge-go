package module_test

import (
	"testing"
	"time"

	"github.com/jwetzell/showbridge-go/config"
	"github.com/jwetzell/showbridge-go/internal/module"
)

func TestMySQLClientFromRegistry(t *testing.T) {
	registration, ok := module.GetModuleRegistration("mysql.client")
	if !ok {
		t.Fatalf("mysql.client module not registered")
	}

	moduleInstance, err := registration.New(config.ModuleConfig{
		Id:   "test",
		Type: "mysql.client",
		Params: map[string]any{
			"dsn": "mysql:mysql@tcp(127.0.0.1:3306)/test",
		},
	})

	if err != nil {
		t.Fatalf("failed to create mysql.client module: %s", err)
	}

	if moduleInstance.Id() != "test" {
		t.Fatalf("mysql.client module has wrong id: %s", moduleInstance.Id())
	}

	if moduleInstance.Type() != "mysql.client" {
		t.Fatalf("mysql.client module has wrong type: %s", moduleInstance.Type())
	}
}

func TestGoodMySQLClient(t *testing.T) {

	testCases := []struct {
		name   string
		params map[string]any
	}{}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {

			registration, ok := module.GetModuleRegistration("mysql.client")
			if !ok {
				t.Fatalf("mysql.client module not registered")
			}

			moduleInstance, err := registration.New(config.ModuleConfig{
				Id:     "test",
				Type:   "mysql.client",
				Params: test.params,
			})

			if err != nil {
				t.Fatalf("mysql.client failed to create module: %s", err)
			}
			// TODO(jwetzell) this is kind of hacky
			go func() {
				time.Sleep(1 * time.Second)
				moduleInstance.Stop()
			}()
			err = moduleInstance.Start(t.Context(), nil)

			if err != nil {
				t.Fatalf("mysql.client failed to start: %s", err)
			}
		})
	}
}

func TestBadMySQLClient(t *testing.T) {
	tests := []struct {
		name        string
		params      map[string]any
		errorString string
	}{
		{
			name:        "no dsn param",
			params:      map[string]any{},
			errorString: "mysql.client dsn error: not found",
		},
		{
			name:        "non-string dsn",
			params:      map[string]any{"dsn": 123},
			errorString: "mysql.client dsn error: not a string",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			registration, ok := module.GetModuleRegistration("mysql.client")
			if !ok {
				t.Fatalf("mysql.client module not registered")
			}

			moduleInstance, err := registration.New(config.ModuleConfig{
				Id:     "test",
				Type:   "mysql.client",
				Params: test.params,
			})

			if err != nil {
				if test.errorString != err.Error() {
					t.Fatalf("mysql.client got error '%s', expected '%s'", err.Error(), test.errorString)
				}
				return
			}

			err = moduleInstance.Start(t.Context(), nil)

			if err == nil {
				t.Fatalf("mysql.client expected to fail")
			}

			if err.Error() != test.errorString {
				t.Fatalf("mysql.client got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
