package processor_test

import (
	"reflect"
	"testing"

	freeD "github.com/jwetzell/free-d-go"
	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/processor"
)

func TestFreeDEncodeFromRegistry(t *testing.T) {
	registration, ok := processor.ProcessorRegistry["freed.encode"]
	if !ok {
		t.Fatalf("freed.encode processor not registered")
	}

	processorInstance, err := registration.New(config.ProcessorConfig{
		Type: "freed.encode",
	})

	if err != nil {
		t.Fatalf("failed to create freed.encode processor: %s", err)
	}

	if processorInstance.Type() != "freed.encode" {
		t.Fatalf("freed.encode processor has wrong type: %s", processorInstance.Type())
	}
}

func TestGoodFreeDEncode(t *testing.T) {
	packetEncoder := processor.FreeDEncode{}

	tests := []struct {
		name     string
		expected []byte
		payload  any
	}{
		{
			name: "basic freed",
			expected: []byte{0xd1, 0x01, 0x5a, 0x00, 0x00, 0x2d, 0x00, 0x00, 0xa6, 0x00, 0x00, 0x7f, 0xff, 0x40, 0x7f, 0xff, 0x80, 0x7f, 0xff,
				0xc0, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x00, 0x00, 50,
			},
			payload: freeD.FreeDPosition{
				ID:    1,
				Pan:   180,
				Tilt:  90,
				Roll:  -180,
				PosX:  131069,
				PosY:  131070,
				PosZ:  131071,
				Zoom:  66051,
				Focus: 263430,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			got, err := packetEncoder.Process(t.Context(), common.GetWrappedPayload(t.Context(), test.payload))

			if err != nil {
				t.Fatalf("freed.encode processing failed: %s", err)
			}

			//TODO(jwetzell): work out better way to compare the any/any
			if !reflect.DeepEqual(got.Payload, test.expected) {
				t.Fatalf("freed.encode got %+v (%T), expected %+v (%T)", got.Payload, got.Payload, test.expected, test.expected)
			}
		})
	}
}

func TestBadFreeDEncode(t *testing.T) {
	packetEncoder := processor.FreeDEncode{}
	tests := []struct {
		name        string
		payload     any
		errorString string
	}{
		{
			name:        "not a FreeD packet",
			payload:     "test",
			errorString: "freed.encode processor only accepts a FreeDPosition",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			got, err := packetEncoder.Process(t.Context(), common.GetWrappedPayload(t.Context(), test.payload))

			if err == nil {
				t.Fatalf("freed.encode expected to fail but succeeded, got: %v", got)

			}
			if err.Error() != test.errorString {
				t.Fatalf("freed.encode got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
