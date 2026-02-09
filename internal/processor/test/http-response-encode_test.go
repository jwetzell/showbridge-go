package processor_test

import (
	"testing"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/processor"
)

func TestHTTPResponseEncodeFromRegistry(t *testing.T) {
	registration, ok := processor.ProcessorRegistry["http.response.encode"]
	if !ok {
		t.Fatalf("http.response.encode processor not registered")
	}

	processorInstance, err := registration.New(config.ProcessorConfig{
		Type: "http.response.encode",
	})

	if err != nil {
		t.Fatalf("failed to create http.response.encode processor: %s", err)
	}

	if processorInstance.Type() != "http.response.encode" {
		t.Fatalf("http.response.encode processor has wrong type: %s", processorInstance.Type())
	}
}
