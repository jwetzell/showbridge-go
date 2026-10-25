package processor_test

import (
	"reflect"
	"testing"

	"github.com/jwetzell/artnet-go"
	"github.com/jwetzell/showbridge-go/processor"
	"github.com/jwetzell/showbridge-go/routing"
)

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
				Sequence: 237,
				Physical: 0,
				SubUni:   1,
				Net:      0,
				Data:     make([]uint8, 512),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			got, err := packetEncoder.Process(t.Context(), routing.WrappedPayload{Payload: test.payload}, nil, nil)

			if err != nil {
				t.Fatalf("artnet.packet.encode processing failed: %s", err)
			}

			// TODO(jwetzell): work out better way to compare the any/any
			if !reflect.DeepEqual(got.Payload, test.expected) {
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

			got, err := packetEncoder.Process(t.Context(), routing.WrappedPayload{Payload: test.payload}, nil, nil)

			if err == nil {
				t.Fatalf("artnet.packet.encode expected to fail but succeeded, got: %v", got)

			}
			if err.Error() != test.errorString {
				t.Fatalf("artnet.packet.encode got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
