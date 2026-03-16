package processor_test

import (
	"slices"
	"testing"

	"github.com/jwetzell/osc-go"
	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/processor"
)

func TestOSCMessageEncodeFromRegistry(t *testing.T) {
	registration, ok := processor.ProcessorRegistry["osc.message.encode"]
	if !ok {
		t.Fatalf("osc.message.encode processor not registered")
	}

	processorInstance, err := registration.New(config.ProcessorConfig{
		Type: "osc.message.encode",
	})

	if err != nil {
		t.Fatalf("failed to create osc.message.encode processor: %s", err)
	}

	if processorInstance.Type() != "osc.message.encode" {
		t.Fatalf("osc.message.encode processor has wrong type: %s", processorInstance.Type())
	}
}

func TestGoodOSCMessageEncode(t *testing.T) {
	processorInstance := processor.OSCMessageEncode{}
	tests := []struct {
		name     string
		payload  any
		expected []byte
	}{
		{
			name: "basic OSC message",
			payload: &osc.OSCMessage{
				Address: "/test",
			},
			expected: []byte{47, 116, 101, 115, 116, 0, 0, 0, 44, 0, 0, 0},
		},
		{
			name: "basic OSC message with argument",
			payload: &osc.OSCMessage{
				Address: "/test",
				Args: []osc.OSCArg{
					{
						Type:  "i",
						Value: int32(42),
					},
				},
			},
			expected: []byte{47, 116, 101, 115, 116, 0, 0, 0, 44, 105, 0, 0, 0, 0, 0, 42},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := processorInstance.Process(t.Context(), common.GetWrappedPayload(t.Context(), test.payload))
			if err != nil {
				t.Fatalf("osc.message.encode processing failed: %s", err)
			}

			gotBytes, ok := got.Payload.([]byte)
			if !ok {
				t.Fatalf("osc.message.encode returned a %T payload: %+v", got, got)
			}

			if !slices.Equal(gotBytes, test.expected) {
				t.Fatalf("osc.message.encode got %+v, expected %+v", got, test.expected)
			}
		})
	}
}

func TestBadOSCMessageEncode(t *testing.T) {
	processorInstance := processor.OSCMessageEncode{}
	tests := []struct {
		name        string
		payload     any
		errorString string
	}{
		{
			name:        "non-osc message input",
			payload:     "test",
			errorString: "osc.message.encode processor only accepts an *OSCMessage",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := processorInstance.Process(t.Context(), common.GetWrappedPayload(t.Context(), test.payload))

			if err == nil {
				t.Fatalf("osc.message.encode expected to fail but got payload: %+v", got)
			}
			if err.Error() != test.errorString {
				t.Fatalf("osc.message.encode got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
