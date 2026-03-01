package processor_test

import (
	"reflect"
	"testing"

	"github.com/jwetzell/osc-go"
	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/processor"
)

func TestOSCMessageCreateFromRegistry(t *testing.T) {
	registration, ok := processor.ProcessorRegistry["osc.message.create"]
	if !ok {
		t.Fatalf("osc.message.create processor not registered")
	}

	processorInstance, err := registration.New(config.ProcessorConfig{
		Type: "osc.message.create",
		Params: map[string]any{
			"address": "/test",
		},
	})

	if err != nil {
		t.Fatalf("failed to create osc.message.create processor: %s", err)
	}

	if processorInstance.Type() != "osc.message.create" {
		t.Fatalf("osc.message.create processor has wrong type: %s", processorInstance.Type())
	}
}

func TestGoodOSCMessageCreate(t *testing.T) {
	tests := []struct {
		name     string
		payload  any
		params   map[string]any
		expected any
	}{
		{
			name: "basic address and no args",
			params: map[string]any{
				"address": "/test",
			},
			payload:  osc.OSCMessage{},
			expected: osc.OSCMessage{Address: "/test"},
		},
		{
			name: "address with template and no args",
			params: map[string]any{
				"address": "/test/{{.Value}}",
			},
			payload:  map[string]any{"Value": "value"},
			expected: osc.OSCMessage{Address: "/test/value"},
		},
		{
			name: "address with template and string arg",
			params: map[string]any{
				"address": "/test/{{.Value}}",
				"args":    []interface{}{"arg1"},
				"types":   "s",
			},
			payload:  map[string]any{"Value": "value"},
			expected: osc.OSCMessage{Address: "/test/value", Args: []osc.OSCArg{{Value: "arg1", Type: "s"}}},
		},
		{
			name: "address with template and mixed args",
			params: map[string]any{
				"address": "/test/{{.Value}}",
				"args":    []interface{}{"arg1", "42", "3.14"},
				"types":   "sif",
			},
			payload: map[string]any{"Value": "value"},
			expected: osc.OSCMessage{
				Address: "/test/value",
				Args: []osc.OSCArg{
					{Value: "arg1", Type: "s"},
					{Value: int32(42), Type: "i"},
					{Value: float32(3.14), Type: "f"},
				},
			},
		},
		{
			name: "address with template and int64 arg",
			params: map[string]any{
				"address": "/test/{{.Value}}",
				"args":    []interface{}{"42"},
				"types":   "h",
			},
			payload:  map[string]any{"Value": "value"},
			expected: osc.OSCMessage{Address: "/test/value", Args: []osc.OSCArg{{Value: int64(42), Type: "h"}}},
		},
		{
			name: "address with template and double arg",
			params: map[string]any{
				"address": "/test/{{.Value}}",
				"args":    []interface{}{"42"},
				"types":   "d",
			},
			payload:  map[string]any{"Value": "value"},
			expected: osc.OSCMessage{Address: "/test/value", Args: []osc.OSCArg{{Value: float64(42), Type: "d"}}},
		},
		{
			name: "address with template and true arg",
			params: map[string]any{
				"address": "/test/{{.Value}}",
				"args":    []interface{}{""},
				"types":   "T",
			},
			payload:  map[string]any{"Value": "value"},
			expected: osc.OSCMessage{Address: "/test/value", Args: []osc.OSCArg{{Value: true, Type: "T"}}},
		},
		{
			name: "address with template and false arg",
			params: map[string]any{
				"address": "/test/{{.Value}}",
				"args":    []interface{}{""},
				"types":   "F",
			},
			payload:  map[string]any{"Value": "value"},
			expected: osc.OSCMessage{Address: "/test/value", Args: []osc.OSCArg{{Value: false, Type: "F"}}},
		},
		{
			name: "address with template and nil arg",
			params: map[string]any{
				"address": "/test/{{.Value}}",
				"args":    []interface{}{""},
				"types":   "N",
			},
			payload:  map[string]any{"Value": "value"},
			expected: osc.OSCMessage{Address: "/test/value", Args: []osc.OSCArg{{Value: nil, Type: "N"}}},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			registration, ok := processor.ProcessorRegistry["osc.message.create"]
			if !ok {
				t.Fatalf("osc.message.create processor not registered")
			}

			processorInstance, err := registration.New(config.ProcessorConfig{
				Type:   "osc.message.create",
				Params: test.params,
			})

			if err != nil {
				t.Fatalf("osc.message.create failed to create processor: %s", err)
			}

			got, err := processorInstance.Process(t.Context(), test.payload)

			if err != nil {
				t.Fatalf("osc.message.create process failed: %s", err)
			}

			if test.expected == nil {
				if got != nil {
					t.Fatalf("osc.message.create got %+v, expected nil", got)
				}
				return
			}

			gotMessage, ok := got.(osc.OSCMessage)
			if !ok {
				t.Fatalf("osc.message.create returned a %T payload: %s", got, got)
			}

			if !reflect.DeepEqual(gotMessage, test.expected) {
				t.Fatalf("osc.message.create got %+v, expected %+v", gotMessage, test.expected)
			}
		})
	}
}

func TestBadOSCMessageCreate(t *testing.T) {
	tests := []struct {
		name        string
		params      map[string]any
		payload     any
		errorString string
	}{
		{
			name:        "no address parameter",
			params:      map[string]any{},
			payload:     "test",
			errorString: "osc.message.create address error: not found",
		},
		{
			name: "non-string address parameter",
			params: map[string]any{
				"address": 123,
			},
			payload:     "test",
			errorString: "osc.message.create address error: not a string",
		},
		{
			name: "bad address template",
			params: map[string]any{
				"address": "{{",
			},
			payload:     "test",
			errorString: "template: address:1: unclosed action",
		},
		{
			name: "non-array args parameter",
			params: map[string]any{
				"address": "/test",
				"args":    "not an array",
				"types":   "s",
			},
			payload:     "test",
			errorString: "osc.message.create address must be an array found string",
		},
		{
			name: "args without types parameter",
			params: map[string]any{
				"address": "/test",
				"args":    []interface{}{"arg1"},
			},
			payload:     "test",
			errorString: "osc.message.create types error: not found",
		},
		{
			name: "args and types length mismatch",
			params: map[string]any{
				"address": "/test",
				"args":    []interface{}{"arg1", "arg2"},
				"types":   "s",
			},
			payload:     "test",
			errorString: "osc.message.create args and types must be the same length",
		},
		{
			name: "non-string arg",
			params: map[string]any{
				"address": "/test",
				"args":    []interface{}{"arg1", 123},
				"types":   "ss",
			},
			payload:     "test",
			errorString: "osc.message.create arg error: not a string",
		},
		{
			name: "bad arg template",
			params: map[string]any{
				"address": "/test",
				"args":    []interface{}{"{{"},
				"types":   "s",
			},
			payload:     "test",
			errorString: "template: arg:1: unclosed action",
		},
		{
			name: "non-string types parameter",
			params: map[string]any{
				"address": "/test",
				"args":    []interface{}{"arg1"},
				"types":   123,
			},
			payload:     "test",
			errorString: "osc.message.create types error: not a string",
		},
		{
			name: "invalid type in types parameter",
			params: map[string]any{
				"address": "/test",
				"args":    []interface{}{"arg1"},
				"types":   "x",
			},
			payload:     "test",
			errorString: "osc.message.create unhandled osc type: x",
		},
		{
			name: "empty address template",
			params: map[string]any{
				"address": "",
			},
			payload:     "test",
			errorString: "osc.message.create address must not be empty",
		},
		{
			name: "address template with missing value",
			params: map[string]any{
				"address": "/test/{{.missing}}",
			},
			payload:     "test",
			errorString: "template: address:1:8: executing \"address\" at <.missing>: can't evaluate field missing in type string",
		},
		{
			name: "address doesn't start with slash",
			params: map[string]any{
				"address": "test",
			},
			payload:     "test",
			errorString: "osc.message.create address must start with '/'",
		},
		{
			name: "address template with missing field",
			params: map[string]any{
				"address": "/test",
				"args":    []interface{}{"{{.missing}}"},
				"types":   "s",
			},
			payload:     "test",
			errorString: "template: arg:1:2: executing \"arg\" at <.missing>: can't evaluate field missing in type string",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			registration, ok := processor.ProcessorRegistry["osc.message.create"]
			if !ok {
				t.Fatalf("osc.message.create processor not registered")
			}

			processorInstance, err := registration.New(config.ProcessorConfig{
				Type:   "osc.message.create",
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
				t.Fatalf("osc.message.create expected to fail but succeeded, got: %v", got)
			}

			if err.Error() != test.errorString {
				t.Fatalf("osc.message.create got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
