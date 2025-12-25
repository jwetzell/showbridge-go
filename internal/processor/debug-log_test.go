package processor_test

import (
	"testing"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/processor"
)

func TestDebugLogFromRegistry(t *testing.T) {
	registration, ok := processor.ProcessorRegistry["debug.log"]
	if !ok {
		t.Fatalf("debug.log processor not registered")
	}

	processorInstance, err := registration.New(config.ProcessorConfig{
		Type: "debug.log",
	})

	if err != nil {
		t.Fatalf("failed to create debug.log processor: %s", err)
	}

	if processorInstance.Type() != "debug.log" {
		t.Fatalf("debug.log processor has wrong type: %s", processorInstance.Type())
	}

	payload := "test"
	expected := "test"

	got, err := processorInstance.Process(t.Context(), payload)
	if err != nil {
		t.Fatalf("debug.log processing failed: %s", err)
	}

	if got != expected {
		t.Fatalf("debug.log got %+v, expected %+v", got, expected)
	}
}
