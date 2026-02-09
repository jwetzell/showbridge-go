package processor_test

import (
	"testing"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/processor"
)

func TestMIDIMessageCreateFromRegistry(t *testing.T) {
	registration, ok := processor.ProcessorRegistry["midi.message.create"]
	if !ok {
		t.Fatalf("midi.message.create processor not registered")
	}

	processorInstance, err := registration.New(config.ProcessorConfig{
		Type: "midi.message.create",
		Params: map[string]any{
			"type":     "note_on",
			"channel":  "1",
			"note":     "60",
			"velocity": "100",
		},
	})

	if err != nil {
		t.Fatalf("failed to create midi.message.create processor: %s", err)
	}

	if processorInstance.Type() != "midi.message.create" {
		t.Fatalf("midi.message.create processor has wrong type: %s", processorInstance.Type())
	}
}
