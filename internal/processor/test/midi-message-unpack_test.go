package processor_test

import (
	"testing"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/processor"
)

func TestMIDIMessageUnpackFromRegistry(t *testing.T) {
	registration, ok := processor.ProcessorRegistry["midi.message.unpack"]
	if !ok {
		t.Fatalf("midi.message.unpack processor not registered")
	}

	processorInstance, err := registration.New(config.ProcessorConfig{
		Type: "midi.message.unpack",
	})

	if err != nil {
		t.Fatalf("failed to create midi.message.unpack processor: %s", err)
	}

	if processorInstance.Type() != "midi.message.unpack" {
		t.Fatalf("midi.message.unpack processor has wrong type: %s", processorInstance.Type())
	}
}
