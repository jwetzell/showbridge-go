package processor_test

import (
	"testing"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/processor"
)

func TestHTTPRequestCreateFromRegistry(t *testing.T) {
	registration, ok := processor.ProcessorRegistry["http.request.do"]
	if !ok {
		t.Fatalf("http.request.do processor not registered")
	}

	processorInstance, err := registration.New(config.ProcessorConfig{
		Type: "http.request.do",
		Params: map[string]any{
			"method": "GET",
			"url":    "http://example.com",
		},
	})

	if err != nil {
		t.Fatalf("failed to create http.request.do processor: %s", err)
	}

	if processorInstance.Type() != "http.request.do" {
		t.Fatalf("http.request.do processor has wrong type: %s", processorInstance.Type())
	}
}
