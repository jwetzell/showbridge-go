package processor_test

import (
	"testing"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/processor"
)

func TestSipResponseDTMFCreateFromRegistry(t *testing.T) {
	registration, ok := processor.ProcessorRegistry["sip.response.dtmf.create"]
	if !ok {
		t.Fatalf("sip.response.dtmf.create processor not registered")
	}

	processorInstance, err := registration.New(config.ProcessorConfig{
		Type: "sip.response.dtmf.create",
		Params: map[string]any{
			"preWait":  0.0,
			"digits":   "good.wav",
			"postWait": 0.0,
		},
	})

	if err != nil {
		t.Fatalf("failed to filter sip.response.dtmf.create processor: %s", err)
	}

	if processorInstance.Type() != "sip.response.dtmf.create" {
		t.Fatalf("sip.response.dtmf.create processor has wrong type: %s", processorInstance.Type())
	}
}
