package processor_test

import (
	"testing"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/processor"
)

func TestSipResponseAudioCreateFromRegistry(t *testing.T) {
	registration, ok := processor.ProcessorRegistry["sip.response.audio.create"]
	if !ok {
		t.Fatalf("sip.response.audio.create processor not registered")
	}

	processorInstance, err := registration.New(config.ProcessorConfig{
		Type: "sip.response.audio.create",
		Params: map[string]any{
			"preWait":   0.0,
			"audioFile": "good.wav",
			"postWait":  0.0,
		},
	})

	if err != nil {
		t.Fatalf("failed to filter sip.response.audio.create processor: %s", err)
	}

	if processorInstance.Type() != "sip.response.audio.create" {
		t.Fatalf("sip.response.audio.create processor has wrong type: %s", processorInstance.Type())
	}
}
