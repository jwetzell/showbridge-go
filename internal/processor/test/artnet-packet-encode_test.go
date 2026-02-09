package processor_test

import (
	"testing"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/processor"
)

func TestArtnetPacketEncodeFromRegistry(t *testing.T) {
	registration, ok := processor.ProcessorRegistry["artnet.packet.encode"]
	if !ok {
		t.Fatalf("artnet.packet.encode processor not registered")
	}

	processorInstance, err := registration.New(config.ProcessorConfig{
		Type: "artnet.packet.encode",
	})

	if err != nil {
		t.Fatalf("failed to create artnet.packet.encode processor: %s", err)
	}

	if processorInstance.Type() != "artnet.packet.encode" {
		t.Fatalf("artnet.packet.encode processor has wrong type: %s", processorInstance.Type())
	}
}
