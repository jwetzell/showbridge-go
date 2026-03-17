package processor_test

import (
	"reflect"
	"testing"

	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/processor"
)

func TestNATSMessageCreateFromRegistry(t *testing.T) {
	registration, ok := processor.ProcessorRegistry["nats.message.create"]
	if !ok {
		t.Fatalf("nats.message.create processor not registered")
	}

	processorInstance, err := registration.New(config.ProcessorConfig{
		Type: "nats.message.create",
		Params: map[string]any{
			"subject": "test",
			"payload": "Hello, World!",
		},
	})

	if err != nil {
		t.Fatalf("failed to create nats.message.create processor: %s", err)
	}

	if processorInstance.Type() != "nats.message.create" {
		t.Fatalf("nats.message.create processor has wrong type: %s", processorInstance.Type())
	}
}

func TestGoodNATSMessageCreate(t *testing.T) {

	tests := []struct {
		name     string
		expected processor.NATSMessage
		params   map[string]any
		payload  any
	}{
		{
			name: "simple payload",
			params: map[string]any{
				"subject": "test",
				"payload": "Hello, World!",
			},
			payload: nil,
			expected: processor.NATSMessage{
				Subject: "test",
				Payload: []byte("Hello, World!"),
			},
		},
		{
			name: "payload with template",
			params: map[string]any{
				"subject": "test",
				"payload": "Hello, {{.Payload.Name}}!",
			},
			payload: map[string]any{
				"Name": "Alice",
			},
			expected: processor.NATSMessage{
				Subject: "test",
				Payload: []byte("Hello, Alice!"),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			registration, ok := processor.ProcessorRegistry["nats.message.create"]
			if !ok {
				t.Fatalf("nats.message.create processor not registered")
			}

			processorInstance, err := registration.New(config.ProcessorConfig{
				Type:   "nats.message.create",
				Params: test.params,
			})

			if err != nil {
				t.Fatalf("nats.message.create failed to create processor: %s", err)
			}

			got, err := processorInstance.Process(t.Context(), common.GetWrappedPayload(t.Context(), test.payload))
			if err != nil {
				t.Fatalf("nats.message.create processing failed: %s", err)
			}

			if !reflect.DeepEqual(got.Payload, test.expected) {
				t.Fatalf("nats.message.create got %+v (%T), expected %+v (%T)", got.Payload, got.Payload, test.expected, test.expected)
			}
		})
	}
}

func TestBadNATSMessageCreate(t *testing.T) {
	tests := []struct {
		name        string
		params      map[string]any
		payload     any
		errorString string
	}{
		{
			name: "missing subject param",
			params: map[string]any{
				"payload": "Hello, World!",
			},
			payload:     nil,
			errorString: "nats.message.create subject error: not found",
		},
		{
			name: "subject param not a string",
			params: map[string]any{
				"subject": 123,
				"payload": "Hello, World!",
			},
			payload:     nil,
			errorString: "nats.message.create subject error: not a string",
		},
		{
			name: "missing payload param",
			params: map[string]any{
				"subject": "test",
			},
			payload:     nil,
			errorString: "nats.message.create payload error: not found",
		},
		{
			name: "payload param not a string",
			params: map[string]any{
				"subject": "test",
				"payload": 123,
			},
			payload:     nil,
			errorString: "nats.message.create payload error: not a string",
		},
		{
			name: "payload template error",
			params: map[string]any{
				"subject": "test",
				"payload": "Hello, {{.Payload.Name}}!",
			},
			payload:     nil,
			errorString: "template: payload:1:17: executing \"payload\" at <.Payload.Name>: nil pointer evaluating interface {}.Name",
		},
		{
			name: "subject template error",
			params: map[string]any{
				"subject": "test.{{.Payload.Name}}",
				"payload": "Hello, World!",
			},
			payload:     nil,
			errorString: "template: subject:1:15: executing \"subject\" at <.Payload.Name>: nil pointer evaluating interface {}.Name",
		},
		{
			name: "subject template syntax error",
			params: map[string]any{
				"subject": "{{.Payload.Name",
				"payload": "Hello, World!",
			},
			payload:     nil,
			errorString: "template: subject:1: unclosed action",
		},
		{
			name: "payload template syntax error",
			params: map[string]any{
				"subject": "test",
				"payload": "Hello, {{.Payload.Name",
			},
			payload:     nil,
			errorString: "template: payload:1: unclosed action",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			registration, ok := processor.ProcessorRegistry["nats.message.create"]
			if !ok {
				t.Fatalf("nats.message.create processor not registered")
			}

			processorInstance, err := registration.New(config.ProcessorConfig{
				Type:   "nats.message.create",
				Params: test.params,
			})

			if err != nil {
				if test.errorString != err.Error() {
					t.Fatalf("nats.message.create got error '%s', expected '%s'", err.Error(), test.errorString)
				}
				return
			}

			got, err := processorInstance.Process(t.Context(), common.GetWrappedPayload(t.Context(), test.payload))

			if err == nil {
				t.Fatalf("nats.message.create expected to fail but succeeded, got: %v", got)
			}

			if err.Error() != test.errorString {
				t.Fatalf("nats.message.create got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
