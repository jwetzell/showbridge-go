package processor_test

import (
	"testing"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/processor"
)

func TestMIDIMessageEncodeFromRegistry(t *testing.T) {
	registration, ok := processor.ProcessorRegistry["midi.message.encode"]
	if !ok {
		t.Fatalf("midi.message.encode processor not registered")
	}

	processorInstance, err := registration.New(config.ProcessorConfig{
		Type: "midi.message.encode",
	})

	if err != nil {
		t.Fatalf("failed to create midi.message.encode processor: %s", err)
	}

	if processorInstance.Type() != "midi.message.encode" {
		t.Fatalf("midi.message.encode processor has wrong type: %s", processorInstance.Type())
	}
}
