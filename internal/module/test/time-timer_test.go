package module_test

import (
	"testing"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/module"
)

func TestTimeTimerFromRegistry(t *testing.T) {
	registration, ok := module.ModuleRegistry["time.timer"]
	if !ok {
		t.Fatalf("time.timer module not registered")
	}

	moduleInstance, err := registration.New(config.ModuleConfig{
		Id:   "test",
		Type: "time.timer",
		Params: map[string]any{
			"duration": 1000.0,
		},
	})

	if err != nil {
		t.Fatalf("failed to create time.timer module: %s", err)
	}

	if moduleInstance.Id() != "test" {
		t.Fatalf("time.timer module has wrong id: %s", moduleInstance.Id())
	}

	if moduleInstance.Type() != "time.timer" {
		t.Fatalf("time.timer module has wrong type: %s", moduleInstance.Type())
	}
}
