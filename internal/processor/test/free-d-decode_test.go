package processor_test

import (
	"reflect"
	"testing"

	freeD "github.com/jwetzell/free-d-go"
	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/processor"
)

func TestFreeDDecodeFromRegistry(t *testing.T) {
	registration, ok := processor.ProcessorRegistry["freed.decode"]
	if !ok {
		t.Fatalf("freed.decode processor not registered")
	}

	processorInstance, err := registration.New(config.ProcessorConfig{
		Type: "freed.decode",
	})

	if err != nil {
		t.Fatalf("failed to create freed.decode processor: %s", err)
	}

	if processorInstance.Type() != "freed.decode" {
		t.Fatalf("freed.decode processor has wrong type: %s", processorInstance.Type())
	}
}

func TestGoodFreeDDecode(t *testing.T) {
	packetEncoder := processor.FreeDDecode{}

	tests := []struct {
		name     string
		payload  []byte
		expected any
	}{
		{
			name: "basic freed",
			payload: []byte{0xd1, 0x01, 0x5a, 0x00, 0x00, 0x2d, 0x00, 0x00, 0xa6, 0x00, 0x00, 0x7f, 0xff, 0x40, 0x7f, 0xff, 0x80, 0x7f, 0xff,
				0xc0, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x00, 0x00, 50,
			},
			expected: freeD.FreeDPosition{
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
				t.Fatalf("freed.decode processing failed: %s", err)
			}

			//TODO(jwetzell): work out better way to compare the any/any
			if !reflect.DeepEqual(got.Payload, test.expected) {
				t.Fatalf("freed.decode got %+v (%T), expected %+v (%T)", got.Payload, got.Payload, test.expected, test.expected)
			}
		})
	}
}

func TestBadFreeDDecode(t *testing.T) {
	packetEncoder := processor.FreeDDecode{}
	tests := []struct {
		name        string
		payload     any
		errorString string
	}{
		{
			name:        "not a byte slice",
			payload:     "test",
			errorString: "freed.decode processor only accepts a []byte",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			got, err := packetEncoder.Process(t.Context(), common.GetWrappedPayload(t.Context(), test.payload))

			if err == nil {
				t.Fatalf("freed.decode expected to fail but succeeded, got: %v", got)

			}
			if err.Error() != test.errorString {
				t.Fatalf("freed.decode got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
