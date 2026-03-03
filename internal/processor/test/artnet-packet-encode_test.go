package processor_test

import (
	"reflect"
	"testing"

	"github.com/jwetzell/artnet-go"
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

func TestGoodArtnetPacketEncode(t *testing.T) {
	packetEncoder := processor.ArtNetPacketEncode{}

	tests := []struct {
		name     string
		expected []byte
		payload  any
	}{
		{
			name:     "number",
			expected: []byte{65, 114, 116, 45, 78, 101, 116, 0, 0, 80, 0, 14, 237, 0, 1, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			payload: &artnet.ArtDmx{
				ID:        []byte{'A', 'r', 't', '-', 'N', 'e', 't', 0x00},
				OpCode:    artnet.OpDmx,
				ProtVerHi: 0,
				ProtVerLo: 14,
				Sequence:  237,
				Physical:  0,
				SubUni:    1,
				Net:       0,
				Length:    512,
				Data:      make([]uint8, 512),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			got, err := packetEncoder.Process(t.Context(), test.payload)

			if err != nil {
				t.Fatalf("artnet.packet.encode processing failed: %s", err)
			}

			//TODO(jwetzell): work out better way to compare the any/any
			if !reflect.DeepEqual(got, test.expected) {
				t.Fatalf("artnet.packet.encode got %+v (%T), expected %+v (%T)", got, got, test.expected, test.expected)
			}
		})
	}
}

func TestBadArtnetPacketEncode(t *testing.T) {
	packetEncoder := processor.ArtNetPacketEncode{}
	tests := []struct {
		name        string
		payload     any
		errorString string
	}{
		{
			name:        "not an ArtNet packet",
			payload:     "test",
			errorString: "artnet.packet.encode processor only accepts an ArtNetPacket",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			got, err := packetEncoder.Process(t.Context(), test.payload)

			if err == nil {
				t.Fatalf("artnet.packet.encode expected to fail but succeeded, got: %v", got)

			}
			if err.Error() != test.errorString {
				t.Fatalf("artnet.packet.encode got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
