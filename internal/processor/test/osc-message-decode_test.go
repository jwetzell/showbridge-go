package processor_test

import (
	"reflect"
	"testing"

	"github.com/jwetzell/osc-go"
	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/processor"
)

func TestOSCMessageDecodeFromRegistry(t *testing.T) {
	registration, ok := processor.ProcessorRegistry["osc.message.decode"]
	if !ok {
		t.Fatalf("osc.message.decode processor not registered")
	}

	processorInstance, err := registration.New(config.ProcessorConfig{
		Type: "osc.message.decode",
	})

	if err != nil {
		t.Fatalf("failed to create osc.message.decode processor: %s", err)
	}

	if processorInstance.Type() != "osc.message.decode" {
		t.Fatalf("osc.message.decode processor has wrong type: %s", processorInstance.Type())
	}
}

func TestGoodOSCMessageDecode(t *testing.T) {
	processorInstance := processor.OSCMessageDecode{}
	tests := []struct {
		name     string
		payload  []byte
		expected *osc.OSCMessage
	}{
		{
			name:    "basic OSC message",
			payload: []byte{47, 116, 101, 115, 116, 0, 0, 0, 44, 0, 0, 0},
			expected: &osc.OSCMessage{
				Address: "/test",
				Args:    []osc.OSCArg{},
			},
		},
		{
			name:    "basic OSC message with argument",
			payload: []byte{47, 116, 101, 115, 116, 0, 0, 0, 44, 105, 0, 0, 0, 0, 0, 42},
			expected: &osc.OSCMessage{
				Address: "/test",
				Args: []osc.OSCArg{
					{
						Type:  "i",
						Value: int32(42),
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := processorInstance.Process(t.Context(), common.GetWrappedPayload(t.Context(), test.payload))

			if err != nil {
				t.Fatalf("osc.message.decode processing failed: %s", err)
			}

			gotMessage, ok := got.Payload.(*osc.OSCMessage)
			if !ok {
				t.Fatalf("osc.message.decode returned a %T payload: %+v", got, got)
			}

			if !reflect.DeepEqual(gotMessage, test.expected) {
				t.Fatalf("osc.message.decode got %+v, expected %+v", gotMessage, test.expected)
			}
		})
	}
}

func TestBadOSCMessageDecode(t *testing.T) {
	processorInstance := processor.OSCMessageDecode{}
	tests := []struct {
		name        string
		payload     any
		errorString string
	}{
		{
			name:        "non []byte payload",
			payload:     "test",
			errorString: "osc.message.decode processor only accepts a []byte payload",
		},
		{
			name:        "empty []byte payload",
			payload:     []byte{},
			errorString: "osc.message.decode processor can't work on empty []byte",
		},
		{
			name:        "wrong start byte in payload",
			payload:     []byte{48, 116, 101, 115, 116, 0, 0, 0, 44, 105, 0, 0, 0, 0, 0, 42},
			errorString: "osc.message.decode processor needs an OSC looking []byte",
		},
		{
			name:        "invalid OSC payload",
			payload:     []byte{47, 116, 101, 115, 116, 0},
			errorString: "osc.message.decode processor failed to decode OSC message: string data is not properly padded",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := processorInstance.Process(t.Context(), common.GetWrappedPayload(t.Context(), test.payload))

			if err == nil {
				t.Fatalf("osc.message.decode expected to fail but got payload: %+v", got)
			}
			if err.Error() != test.errorString {
				t.Fatalf("osc.message.decode got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
