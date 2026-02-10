package module_test

import (
	"testing"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/module"
)

func TestHTTPServerFromRegistry(t *testing.T) {
	registration, ok := module.ModuleRegistry["http.server"]
	if !ok {
		t.Fatalf("http.server module not registered")
	}

	moduleInstance, err := registration.New(config.ModuleConfig{
		Id:   "test",
		Type: "http.server",
		Params: map[string]any{
			"port": 3000.0,
		},
	})

	if err != nil {
		t.Fatalf("failed to create http.server module: %s", err)
	}

	if moduleInstance.Id() != "test" {
		t.Fatalf("http.server module has wrong id: %s", moduleInstance.Id())
	}

	if moduleInstance.Type() != "http.server" {
		t.Fatalf("http.server module has wrong type: %s", moduleInstance.Type())
	}
}
