package processor_test

import (
	"testing"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/processor"
)

func TestMIDIMessageDecodeFromRegistry(t *testing.T) {
	registration, ok := processor.ProcessorRegistry["midi.message.decode"]
	if !ok {
		t.Fatalf("midi.message.decode processor not registered")
	}

	processorInstance, err := registration.New(config.ProcessorConfig{
		Type: "midi.message.decode",
	})

	if err != nil {
		t.Fatalf("failed to create midi.message.decode processor: %s", err)
	}

	if processorInstance.Type() != "midi.message.decode" {
		t.Fatalf("midi.message.decode processor has wrong type: %s", processorInstance.Type())
	}
}
