package module_test

import (
	"testing"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/module"
)

func TestSerialClientFromRegistry(t *testing.T) {
	registration, ok := module.ModuleRegistry["serial.client"]
	if !ok {
		t.Fatalf("serial.client module not registered")
	}

	moduleInstance, err := registration.New(config.ModuleConfig{
		Id:   "test",
		Type: "serial.client",
		Params: map[string]any{
			"port":     "/dev/ttyUSB0",
			"framing":  "LF",
			"baudRate": 9600.0,
		},
	})

	if err != nil {
		t.Fatalf("failed to create serial.client module: %s", err)
	}

	if moduleInstance.Id() != "test" {
		t.Fatalf("serial.client module has wrong id: %s", moduleInstance.Id())
	}

	if moduleInstance.Type() != "serial.client" {
		t.Fatalf("serial.client module has wrong type: %s", moduleInstance.Type())
	}
}
