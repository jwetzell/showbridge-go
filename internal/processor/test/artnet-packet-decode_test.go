package processor_test

import (
	"testing"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/processor"
)

func TestArtnetPacketCreateFromRegistry(t *testing.T) {
	registration, ok := processor.ProcessorRegistry["artnet.packet.decode"]
	if !ok {
		t.Fatalf("artnet.packet.decode processor not registered")
	}

	processorInstance, err := registration.New(config.ProcessorConfig{
		Type: "artnet.packet.decode",
	})

	if err != nil {
		t.Fatalf("failed to decode artnet.packet.decode processor: %s", err)
	}

	if processorInstance.Type() != "artnet.packet.decode" {
		t.Fatalf("artnet.packet.decode processor has wrong type: %s", processorInstance.Type())
	}
}
