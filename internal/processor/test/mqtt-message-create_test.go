package processor_test

import (
	"reflect"
	"testing"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/processor"
)

func TestMQTTMessageCreateFromRegistry(t *testing.T) {
	registration, ok := processor.ProcessorRegistry["mqtt.message.create"]
	if !ok {
		t.Fatalf("mqtt.message.create processor not registered")
	}

	processorInstance, err := registration.New(config.ProcessorConfig{
		Type: "mqtt.message.create",
		Params: map[string]any{
			"topic":    "test/topic",
			"payload":  "Hello, World!",
			"qos":      1,
			"retained": true,
		},
	})

	if err != nil {
		t.Fatalf("failed to create mqtt.message.create processor: %s", err)
	}

	if processorInstance.Type() != "mqtt.message.create" {
		t.Fatalf("mqtt.message.create processor has wrong type: %s", processorInstance.Type())
	}
}

func TestGoodMQTTMessageCreate(t *testing.T) {
	tests := []struct {
		name     string
		payload  any
		params   map[string]any
		expected any
	}{
		{
			name: "basic topic and string payload",
			params: map[string]any{
				"topic":    "test/topic",
				"payload":  "Hello, World!",
				"qos":      1,
				"retained": true,
			},
			payload:  "test",
			expected: processor.NewMQTTMessage("test/topic", 1, true, []byte("Hello, World!")),
		},
		{
			name: "basic topic and []byte payload",
			params: map[string]any{
				"topic":    "test/topic",
				"payload":  []byte{72, 101, 108, 108, 111},
				"qos":      1,
				"retained": true,
			},
			payload:  "test",
			expected: processor.NewMQTTMessage("test/topic", 1, true, []byte("Hello")),
		},
		{
			name: "basic topic and []int payload",
			params: map[string]any{
				"topic":    "test/topic",
				"payload":  []int{72, 101, 108, 108, 111},
				"qos":      1,
				"retained": true,
			},
			payload:  "test",
			expected: processor.NewMQTTMessage("test/topic", 1, true, []byte("Hello")),
		},
		{
			name: "basic topic and []uint payload",
			params: map[string]any{
				"topic":    "test/topic",
				"payload":  []uint{72, 101, 108, 108, 111},
				"qos":      1,
				"retained": true,
			},
			payload:  "test",
			expected: processor.NewMQTTMessage("test/topic", 1, true, []byte("Hello")),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			registration, ok := processor.ProcessorRegistry["mqtt.message.create"]
			if !ok {
				t.Fatalf("mqtt.message.create processor not registered")
			}

			processorInstance, err := registration.New(config.ProcessorConfig{
				Type:   "mqtt.message.create",
				Params: test.params,
			})

			if err != nil {
				t.Fatalf("mqtt.message.create failed to create processor: %s", err)
			}

			got, err := processorInstance.Process(t.Context(), test.payload)

			if err != nil {
				t.Fatalf("mqtt.message.create processing failed: %s", err)
			}

			if test.expected == nil {
				if got != nil {
					t.Fatalf("mqtt.message.create got %+v, expected nil", got)
				}
				return
			}

			gotMessage, ok := got.(mqtt.Message)
			if !ok {
				t.Fatalf("mqtt.message.create returned a %T payload: %s", got, got)
			}

			if !reflect.DeepEqual(gotMessage, test.expected) {
				t.Fatalf("mqtt.message.create got %+v, expected %+v", gotMessage, test.expected)
			}
		})
	}
}

func TestBadMQTTMessageCreate(t *testing.T) {
	tests := []struct {
		name        string
		params      map[string]any
		payload     any
		errorString string
	}{
		{
			name:        "no topic parameter",
			params:      map[string]any{},
			payload:     "test",
			errorString: "mqtt.message.create topic error: not found",
		},
		{
			name: "non-string topic parameter",
			params: map[string]any{
				"topic": 1,
			},
			payload:     "test",
			errorString: "mqtt.message.create topic error: not a string",
		},
		{
			name: "no qos parameter",
			params: map[string]any{
				"topic": "test/topic",
			},
			payload:     "test",
			errorString: "mqtt.message.create qos error: not found",
		},
		{
			name: "non-number qos parameter",
			params: map[string]any{
				"topic": "test/topic",
				"qos":   "1",
			},
			payload:     "test",
			errorString: "mqtt.message.create qos error: not a number",
		},
		{
			name: "no retained parameter",
			params: map[string]any{
				"topic": "test/topic",
				"qos":   1,
			},
			payload:     "test",
			errorString: "mqtt.message.create retained error: not found",
		},
		{
			name: "non-bool retained parameter",
			params: map[string]any{
				"topic":    "test/topic",
				"qos":      1,
				"retained": "1",
			},
			payload:     "test",
			errorString: "mqtt.message.create retained error: not a boolean",
		},
		{
			name: "no payload parameter",
			params: map[string]any{
				"topic":    "test/topic",
				"qos":      1,
				"retained": true,
			},
			payload:     "test",
			errorString: "mqtt.message.create payload error: not found",
		},
		{
			name: "non-string payload parameter",
			params: map[string]any{
				"topic":    "test/topic",
				"qos":      1,
				"retained": true,
				"payload":  123,
			},
			payload:     1,
			errorString: "mqtt.message.create payload error: not a slice",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			registration, ok := processor.ProcessorRegistry["mqtt.message.create"]
			if !ok {
				t.Fatalf("mqtt.message.create processor not registered")
			}

			processorInstance, err := registration.New(config.ProcessorConfig{
				Type:   "mqtt.message.create",
				Params: test.params,
			})

			if err != nil {
				if test.errorString != err.Error() {
					t.Fatalf("mqtt.message.create got error '%s', expected '%s'", err.Error(), test.errorString)
				}
				return
			}

			got, err := processorInstance.Process(t.Context(), test.payload)

			if err == nil {
				t.Fatalf("mqtt.message.create expected to fail but succeeded, got: %v", got)
			}

			if err.Error() != test.errorString {
				t.Fatalf("mqtt.message.create got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
