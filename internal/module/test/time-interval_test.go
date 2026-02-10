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
