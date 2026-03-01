package processor_test

import (
	"reflect"
	"testing"

	"github.com/jwetzell/artnet-go"
	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/processor"
)

func TestArtnetPacketFilterFromRegistry(t *testing.T) {
	registration, ok := processor.ProcessorRegistry["artnet.packet.filter"]
	if !ok {
		t.Fatalf("artnet.packet.filter processor not registered")
	}

	processorInstance, err := registration.New(config.ProcessorConfig{
		Type: "artnet.packet.filter",
		Params: map[string]any{
			"opCode": float64(artnet.OpTimeCode),
		},
	})

	if err != nil {
		t.Fatalf("failed to create artnet.packet.filter processor: %s", err)
	}

	if processorInstance.Type() != "artnet.packet.filter" {
		t.Fatalf("artnet.packet.filter processor has wrong type: %s", processorInstance.Type())
	}

	payload := &artnet.ArtTimeCode{
		ID:        []byte{'A', 'r', 't', '-', 'N', 'e', 't', 0x00},
		OpCode:    artnet.OpTimeCode,
		ProtVerHi: 0,
		ProtVerLo: 14,
		Filler1:   0,
		StreamId:  0,
		Frames:    11,
		Seconds:   17,
		Minutes:   3,
		Hours:     0,
		Type:      0,
	}
	expected := &artnet.ArtTimeCode{
		ID:        []byte{'A', 'r', 't', '-', 'N', 'e', 't', 0x00},
		OpCode:    artnet.OpTimeCode,
		ProtVerHi: 0,
		ProtVerLo: 14,
		Filler1:   0,
		StreamId:  0,
		Frames:    11,
		Seconds:   17,
		Minutes:   3,
		Hours:     0,
		Type:      0,
	}

	got, err := processorInstance.Process(t.Context(), payload)
	if err != nil {
		t.Fatalf("artnet.packet.filter processing failed: %s", err)
	}

	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("artnet.packet.filter got %+v, expected %+v", got, expected)
	}
}

func TestGoodArtnetPacketFilter(t *testing.T) {

	tests := []struct {
		name     string
		params   map[string]any
		payload  any
		expected artnet.ArtNetPacket
	}{
		{
			name: "tiemcode packet with matching opCode",
			params: map[string]any{
				"opCode": float64(artnet.OpTimeCode),
			},
			payload: &artnet.ArtTimeCode{
				ID:        []byte{'A', 'r', 't', '-', 'N', 'e', 't', 0x00},
				OpCode:    artnet.OpTimeCode,
				ProtVerHi: 0,
				ProtVerLo: 14,
				Filler1:   0,
				StreamId:  0,
				Frames:    11,
				Seconds:   17,
				Minutes:   3,
				Hours:     0,
				Type:      0,
			},
			expected: &artnet.ArtTimeCode{
				ID:        []byte{'A', 'r', 't', '-', 'N', 'e', 't', 0x00},
				OpCode:    artnet.OpTimeCode,
				ProtVerHi: 0,
				ProtVerLo: 14,
				Filler1:   0,
				StreamId:  0,
				Frames:    11,
				Seconds:   17,
				Minutes:   3,
				Hours:     0,
				Type:      0,
			},
		},
		{
			name: "timecode packet with mismatching opCode",
			params: map[string]any{
				"opCode": float64(artnet.OpDmx),
			},
			payload: &artnet.ArtTimeCode{
				ID:        []byte{'A', 'r', 't', '-', 'N', 'e', 't', 0x00},
				OpCode:    artnet.OpTimeCode,
				ProtVerHi: 0,
				ProtVerLo: 14,
				Filler1:   0,
				StreamId:  0,
				Frames:    11,
				Seconds:   17,
				Minutes:   3,
				Hours:     0,
				Type:      0,
			},
			expected: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			registration, ok := processor.ProcessorRegistry["artnet.packet.filter"]
			if !ok {
				t.Fatalf("artnet.packet.filter processor not registered")
			}

			processorInstance, err := registration.New(config.ProcessorConfig{
				Type:   "artnet.packet.filter",
				Params: test.params,
			})

			if err != nil {
				t.Fatalf("artnet.packet.filter failed to create processor: %s", err)
			}

			got, err := processorInstance.Process(t.Context(), test.payload)

			if err != nil {
				t.Fatalf("artnet.packet.filter failed: %s", err)
			}

			if test.expected == nil {
				if got != nil {
					t.Fatalf("artnet.packet.filter got %+v, expected nil", got)
				}
				return
			}

			gotPacket, ok := got.(artnet.ArtNetPacket)
			if !ok {
				t.Fatalf("artnet.packet.filter returned a %T payload: %s", got, got)
			}

			if !reflect.DeepEqual(gotPacket, test.expected) {
				t.Fatalf("artnet.packet.filter got %+v, expected %+v", gotPacket, test.expected)
			}
		})
	}
}

func TestBadArtnetPacketFilter(t *testing.T) {
	tests := []struct {
		name        string
		params      map[string]any
		payload     any
		errorString string
	}{
		{
			name:        "non-artnet input",
			payload:     []byte{0x01},
			params:      map[string]any{"opCode": float64(artnet.OpTimeCode)},
			errorString: "artnet.packet.filter processor only accepts an ArtNetPacket",
		},
		{
			name: "no opCode param",
			payload: &artnet.ArtTimeCode{
				ID:        []byte{'A', 'r', 't', '-', 'N', 'e', 't', 0x00},
				OpCode:    artnet.OpTimeCode,
				ProtVerHi: 0,
				ProtVerLo: 14,
				Filler1:   0,
				StreamId:  0,
				Frames:    11,
				Seconds:   17,
				Minutes:   3,
				Hours:     0,
				Type:      0,
			},
			params:      map[string]any{},
			errorString: "artnet.packet.filter opCode error: not found",
		},
		{
			name: "opCode not a number",
			payload: &artnet.ArtTimeCode{
				ID:        []byte{'A', 'r', 't', '-', 'N', 'e', 't', 0x00},
				OpCode:    artnet.OpTimeCode,
				ProtVerHi: 0,
				ProtVerLo: 14,
				Filler1:   0,
				StreamId:  0,
				Frames:    11,
				Seconds:   17,
				Minutes:   3,
				Hours:     0,
				Type:      0,
			},
			params:      map[string]any{"opCode": "100"},
			errorString: "artnet.packet.filter opCode error: not a number",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			registration, ok := processor.ProcessorRegistry["artnet.packet.filter"]
			if !ok {
				t.Fatalf("artnet.packet.filter processor not registered")
			}

			processorInstance, err := registration.New(config.ProcessorConfig{
				Type:   "artnet.packet.filter",
				Params: test.params,
			})

			if err != nil {
				if test.errorString != err.Error() {
					t.Fatalf("artnet.packet.filter got error '%s', expected '%s'", err.Error(), test.errorString)
				}
				return
			}

			got, err := processorInstance.Process(t.Context(), test.payload)

			if err == nil {
				t.Fatalf("artnet.packet.filter expected to fail but got payload: %s", got)
			}

			if err.Error() != test.errorString {
				t.Fatalf("artnet.packet.filter got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
