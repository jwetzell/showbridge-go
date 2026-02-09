package processor_test

import (
	"testing"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/processor"
)

func TestHTTPRequestCreateFromRegistry(t *testing.T) {
	registration, ok := processor.ProcessorRegistry["http.request.create"]
	if !ok {
		t.Fatalf("http.request.create processor not registered")
	}

	processorInstance, err := registration.New(config.ProcessorConfig{
		Type: "http.request.create",
		Params: map[string]any{
			"method": "GET",
			"url":    "http://example.com",
		},
	})

	if err != nil {
		t.Fatalf("failed to create http.request.create processor: %s", err)
	}

	if processorInstance.Type() != "http.request.create" {
		t.Fatalf("http.request.create processor has wrong type: %s", processorInstance.Type())
	}
}
