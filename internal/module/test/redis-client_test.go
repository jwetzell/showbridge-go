package module_test

import (
	"testing"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/module"
)

func TestRedisClientFromRegistry(t *testing.T) {
	registration, ok := module.ModuleRegistry["redis.client"]
	if !ok {
		t.Fatalf("redis.client module not registered")
	}

	moduleInstance, err := registration.New(config.ModuleConfig{
		Id:   "test",
		Type: "redis.client",
		Params: map[string]any{
			"host": "localhost",
			"port": 6379,
		},
	})

	if err != nil {
		t.Fatalf("failed to create redis.client module: %s", err)
	}

	if moduleInstance.Id() != "test" {
		t.Fatalf("redis.client module has wrong id: %s", moduleInstance.Id())
	}

	if moduleInstance.Type() != "redis.client" {
		t.Fatalf("redis.client module has wrong type: %s", moduleInstance.Type())
	}
}

func TestBadRedisClient(t *testing.T) {
	tests := []struct {
		name        string
		params      map[string]any
		errorString string
	}{
		{
			name: "no host param",
			params: map[string]any{
				"port": 6379,
			},
			errorString: "redis.client host error: not found",
		},
		{
			name: "non-string host",
			params: map[string]any{
				"host": 123,
				"port": 6379,
			},
			errorString: "redis.client host error: not a string",
		},
		{
			name: "no port param",
			params: map[string]any{
				"host": "localhost",
			},
			errorString: "redis.client port error: not found",
		},
		{
			name: "non-number port",
			params: map[string]any{
				"host": "localhost",
				"port": "6379",
			},
			errorString: "redis.client port error: not a number",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			registration, ok := module.ModuleRegistry["redis.client"]
			if !ok {
				t.Fatalf("redis.client module not registered")
			}

			moduleInstance, err := registration.New(config.ModuleConfig{
				Id:     "test",
				Type:   "redis.client",
				Params: test.params,
			})

			if err != nil {
				if test.errorString != err.Error() {
					t.Fatalf("redis.client got error '%s', expected '%s'", err.Error(), test.errorString)
				}
				return
			}

			err = moduleInstance.Start(t.Context())

			if err == nil {
				t.Fatalf("redis.client expected to fail")
			}

			if err.Error() != test.errorString {
				t.Fatalf("redis.client got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
