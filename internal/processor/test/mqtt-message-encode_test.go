package processor_test

import (
	"slices"
	"testing"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/processor"
)

func TestMQTTMessageEncodeFromRegistry(t *testing.T) {
	registration, ok := processor.ProcessorRegistry["mqtt.message.encode"]
	if !ok {
		t.Fatalf("mqtt.message.encode processor not registered")
	}

	processorInstance, err := registration.New(config.ProcessorConfig{
		Type: "mqtt.message.encode",
	})

	if err != nil {
		t.Fatalf("failed to create mqtt.message.encode processor: %s", err)
	}

	if processorInstance.Type() != "mqtt.message.encode" {
		t.Fatalf("mqtt.message.encode processor has wrong type: %s", processorInstance.Type())
	}
}

func TestGoodMQTTMessageEncode(t *testing.T) {
	stringEncoder := processor.MQTTMessageEncode{}
	tests := []struct {
		name     string
		payload  mqtt.Message
		expected []byte
	}{
		{
			name:     "basic string",
			payload:  processor.NewMQTTMessage("test/topic", 1, true, []byte("hello")),
			expected: []byte{0x68, 0x65, 0x6c, 0x6c, 0x6f},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := stringEncoder.Process(t.Context(), test.payload)

			gotBytes, ok := got.([]byte)
			if !ok {
				t.Fatalf("mqtt.message.encode returned a %T payload: %s", got, got)
			}
			if err != nil {
				t.Fatalf("mqtt.message.encode failed: %s", err)
			}
			if !slices.Equal(gotBytes, test.expected) {
				t.Fatalf("mqtt.message.encode got %s, expected %s", got, test.expected)
			}
		})
	}
}

func TestBadMQTTMessageEncode(t *testing.T) {
	stringEncoder := processor.MQTTMessageEncode{}
	tests := []struct {
		name        string
		payload     any
		errorString string
	}{
		{
			name:        "non-mqtt message input",
			payload:     []byte{0x68, 0x65, 0x6c, 0x6c, 0x6f},
			errorString: "mqtt.message.encode processor only accepts an mqtt.Message",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := stringEncoder.Process(t.Context(), test.payload)

			if err == nil {
				t.Fatalf("mqtt.message.encode expected to fail but got payload: %s", got)
			}
			if err.Error() != test.errorString {
				t.Fatalf("mqtt.message.encode got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
