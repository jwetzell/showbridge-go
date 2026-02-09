package processor_test

import (
	"testing"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/processor"
)

func TestHTTPRequestEncodeFromRegistry(t *testing.T) {
	registration, ok := processor.ProcessorRegistry["http.request.encode"]
	if !ok {
		t.Fatalf("http.request.encode processor not registered")
	}

	processorInstance, err := registration.New(config.ProcessorConfig{
		Type: "http.request.encode",
	})

	if err != nil {
		t.Fatalf("failed to create http.request.encode processor: %s", err)
	}

	if processorInstance.Type() != "http.request.encode" {
		t.Fatalf("http.request.encode processor has wrong type: %s", processorInstance.Type())
	}
}
