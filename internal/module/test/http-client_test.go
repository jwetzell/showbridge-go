package module_test

import (
	"testing"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/module"
)

func TestHTTPClientFromRegistry(t *testing.T) {
	registration, ok := module.ModuleRegistry["http.client"]
	if !ok {
		t.Fatalf("http.client module not registered")
	}

	moduleInstance, err := registration.New(config.ModuleConfig{
		Id:   "test",
		Type: "http.client",
	})

	if err != nil {
		t.Fatalf("failed to create http.client module: %s", err)
	}

	if moduleInstance.Id() != "test" {
		t.Fatalf("http.client module has wrong id: %s", moduleInstance.Id())
	}

	if moduleInstance.Type() != "http.client" {
		t.Fatalf("http.client module has wrong type: %s", moduleInstance.Type())
	}
}
