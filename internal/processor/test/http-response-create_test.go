package processor_test

import (
	"testing"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/processor"
)

func TestHTTPResponseCreateFromRegistry(t *testing.T) {
	registration, ok := processor.ProcessorRegistry["http.response.create"]
	if !ok {
		t.Fatalf("http.response.create processor not registered")
	}

	processorInstance, err := registration.New(config.ProcessorConfig{
		Type: "http.response.create",
		Params: map[string]any{
			"status":       200.0,
			"bodyTemplate": "Hello, World!",
		},
	})

	if err != nil {
		t.Fatalf("failed to create http.response.create processor: %s", err)
	}

	if processorInstance.Type() != "http.response.create" {
		t.Fatalf("http.response.create processor has wrong type: %s", processorInstance.Type())
	}
}
