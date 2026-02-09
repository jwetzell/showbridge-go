package processor_test

import (
	"testing"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/processor"
)

func TestMIDIMessageFilterFromRegistry(t *testing.T) {
	registration, ok := processor.ProcessorRegistry["midi.message.filter"]
	if !ok {
		t.Fatalf("midi.message.filter processor not registered")
	}

	processorInstance, err := registration.New(config.ProcessorConfig{
		Type: "midi.message.filter",
		Params: map[string]any{
			"type": "NoteOn",
		},
	})

	if err != nil {
		t.Fatalf("failed to create midi.message.filter processor: %s", err)
	}

	if processorInstance.Type() != "midi.message.filter" {
		t.Fatalf("midi.message.filter processor has wrong type: %s", processorInstance.Type())
	}
}
