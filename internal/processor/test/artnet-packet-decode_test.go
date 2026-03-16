package processor_test

import (
	"reflect"
	"testing"

	"github.com/jwetzell/artnet-go"
	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/processor"
)

func TestArtnetPacketDecodeFromRegistry(t *testing.T) {
	registration, ok := processor.ProcessorRegistry["artnet.packet.decode"]
	if !ok {
		t.Fatalf("artnet.packet.decode processor not registered")
	}

	processorInstance, err := registration.New(config.ProcessorConfig{
		Type: "artnet.packet.decode",
	})

	if err != nil {
		t.Fatalf("failed to create artnet.packet.decode processor: %s", err)
	}

	if processorInstance.Type() != "artnet.packet.decode" {
		t.Fatalf("artnet.packet.decode processor has wrong type: %s", processorInstance.Type())
	}
}

func TestGoodArtnetPacketDecode(t *testing.T) {
	packetDecoder := processor.ArtNetPacketDecode{}

	tests := []struct {
		name     string
		payload  any
		expected artnet.ArtNetPacket
	}{
		{
			name:    "number",
			payload: []byte{65, 114, 116, 45, 78, 101, 116, 0, 0, 80, 0, 14, 237, 0, 1, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			expected: &artnet.ArtDmx{
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

			got, err := packetDecoder.Process(t.Context(), common.GetWrappedPayload(t.Context(), test.payload))

			if err != nil {
				t.Fatalf("artnet.packet.decode processing failed: %s", err)
			}

			//TODO(jwetzell): work out better way to compare the any/any
			if !reflect.DeepEqual(got.Payload, test.expected) {
				t.Fatalf("artnet.packet.decode got %+v (%T), expected %+v (%T)", got.Payload, got.Payload, test.expected, test.expected)
			}
		})
	}
}

func TestBadArtnetPacketDecode(t *testing.T) {
	packetDecoder := processor.ArtNetPacketDecode{}
	tests := []struct {
		name        string
		payload     any
		errorString string
	}{
		{
			name:        "not a byte slice",
			payload:     "not a byte slice",
			errorString: "artnet.packet.decode processor only accepts a []byte",
		},
		{
			name:        "not enough bytes",
			payload:     []byte{1, 2, 3},
			errorString: "ArtNet packet must be at least 12 bytes",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			got, err := packetDecoder.Process(t.Context(), common.GetWrappedPayload(t.Context(), test.payload))

			if err == nil {
				t.Fatalf("artnet.packet.decode expected to fail but succeeded, got: %v", got)

			}
			if err.Error() != test.errorString {
				t.Fatalf("artnet.packet.decode got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
