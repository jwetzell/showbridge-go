package processor_test

import (
	"reflect"
	"testing"

	"github.com/jwetzell/osc-go"
	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/processor"
)

func TestOSCMessageFilterFromRegistry(t *testing.T) {
	registration, ok := processor.ProcessorRegistry["osc.message.filter"]
	if !ok {
		t.Fatalf("osc.message.filter processor not registered")
	}

	processorInstance, err := registration.New(config.ProcessorConfig{
		Type: "osc.message.filter",
		Params: map[string]any{
			"address": "/test*",
		},
	})

	if err != nil {
		t.Fatalf("failed to filter osc.message.filter processor: %s", err)
	}

	if processorInstance.Type() != "osc.message.filter" {
		t.Fatalf("osc.message.filter processor has wrong type: %s", processorInstance.Type())
	}
}

func TestGoodOSCMessageFilter(t *testing.T) {
	tests := []struct {
		name     string
		payload  osc.OSCMessage
		params   map[string]any
		expected any
	}{
		{
			name: "basic address match",
			params: map[string]any{
				"address": "/test",
			},
			payload:  osc.OSCMessage{Address: "/test"},
			expected: osc.OSCMessage{Address: "/test"},
		},
		{
			name: "basic address no match",
			params: map[string]any{
				"address": "/test",
			},
			payload:  osc.OSCMessage{Address: "/testing"},
			expected: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			registration, ok := processor.ProcessorRegistry["osc.message.filter"]
			if !ok {
				t.Fatalf("osc.message.filter processor not registered")
			}

			processorInstance, err := registration.New(config.ProcessorConfig{
				Type:   "osc.message.filter",
				Params: test.params,
			})

			if err != nil {
				t.Fatalf("osc.message.filter failed to create processor: %s", err)
			}

			got, err := processorInstance.Process(t.Context(), test.payload)

			if err != nil {
				t.Fatalf("osc.message.filter process failed: %s", err)
			}

			if test.expected == nil {
				if got != nil {
					t.Fatalf("osc.message.filter got %+v, expected nil", got)
				}
				return
			}

			gotMessage, ok := got.(osc.OSCMessage)
			if !ok {
				t.Fatalf("osc.message.filter returned a %T payload: %s", got, got)
			}

			if !reflect.DeepEqual(gotMessage, test.expected) {
				t.Fatalf("osc.message.filter got %+v, expected %+v", gotMessage, test.expected)
			}
		})
	}
}

func TestBadOSCMessageFilter(t *testing.T) {
	tests := []struct {
		name        string
		params      map[string]any
		payload     any
		errorString string
	}{
		{
			name:        "no address parameter",
			params:      map[string]any{},
			payload:     osc.OSCMessage{Address: "/test"},
			errorString: "osc.message.filter address error: not found",
		},
		{
			name: "non-string address parameter",
			params: map[string]any{
				"address": 123,
			},
			payload:     osc.OSCMessage{Address: "/test"},
			errorString: "osc.message.filter address error: not a string",
		},
		{
			name: "bad address pattern",
			params: map[string]any{
				"address": "[",
			},
			payload:     osc.OSCMessage{Address: "/test"},
			errorString: "error parsing regexp: missing closing ]: `[$`",
		},
		{
			name: "non-osc input",
			params: map[string]any{
				"address": "/test",
			},
			payload:     []byte("hello"),
			errorString: "osc.message.filter can only operate on OSCMessage payloads",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			registration, ok := processor.ProcessorRegistry["osc.message.filter"]
			if !ok {
				t.Fatalf("osc.message.filter processor not registered")
			}

			processorInstance, err := registration.New(config.ProcessorConfig{
				Type:   "osc.message.filter",
				Params: test.params,
			})

			if err != nil {
				if test.errorString != err.Error() {
					t.Fatalf("string.create got error '%s', expected '%s'", err.Error(), test.errorString)
				}
				return
			}

			got, err := processorInstance.Process(t.Context(), test.payload)

			if err == nil {
				t.Fatalf("osc.message.filter expected to fail but succeeded, got: %v", got)
			}

			if err.Error() != test.errorString {
				t.Fatalf("osc.message.filter got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
