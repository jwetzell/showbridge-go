package processor_test

import (
	"testing"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/processor"
)

func TestHTTPRequestFilterFromRegistry(t *testing.T) {
	registration, ok := processor.ProcessorRegistry["http.request.filter"]
	if !ok {
		t.Fatalf("http.request.filter processor not registered")
	}

	processorInstance, err := registration.New(config.ProcessorConfig{
		Type: "http.request.filter",
		Params: map[string]any{
			"method": "GET",
			"path":   "/test",
		},
	})

	if err != nil {
		t.Fatalf("failed to create http.request.filter processor: %s", err)
	}

	if processorInstance.Type() != "http.request.filter" {
		t.Fatalf("http.request.filter processor has wrong type: %s", processorInstance.Type())
	}
}
