package processor_test

import (
	"reflect"
	"testing"

	"github.com/jwetzell/osc-go"
	"github.com/jwetzell/showbridge-go/config"
	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/processor"
)

func TestOSCMessageCreateFromRegistry(t *testing.T) {
	registration, ok := processor.GetProcessorRegistration("osc.message.create")
	if !ok {
		t.Fatalf("osc.message.create processor not registered")
	}

	processorInstance, err := registration.New(config.ProcessorConfig{
		Id:   "test-id",
		Type: "osc.message.create",
		Params: map[string]any{
			"address": "/test",
		},
	})

	if err != nil {
		t.Fatalf("failed to create osc.message.create processor: %s", err)
	}

	if processorInstance.Id() != "test-id" {
		t.Fatalf("osc.message.create processor has wrong id: %s", processorInstance.Id())
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
			payload:  osc.Message{},
			expected: &osc.Message{Address: "/test"},
		},
		{
			name: "address with template and no args",
			params: map[string]any{
				"address": "/test/{{.Payload.Value}}",
			},
			payload:  map[string]any{"Value": "value"},
			expected: &osc.Message{Address: "/test/value"},
		},
		{
			name: "address with template and string arg",
			params: map[string]any{
				"address": "/test/{{.Payload.Value}}",
				"args":    []map[string]any{{"value": "arg1", "type": "s"}},
			},
			payload:  map[string]any{"Value": "value"},
			expected: &osc.Message{Address: "/test/value", Args: []osc.Arg{{Value: "arg1", Type: "s"}}},
		},
		{
			name: "address with template and mixed args",
			params: map[string]any{
				"address": "/test/{{.Payload.Value}}",
				"args": []map[string]any{
					{"value": "arg1", "type": "s"},
					{"value": "42", "type": "i"},
					{"value": "3.14", "type": "f"},
				},
			},
			payload: map[string]any{"Value": "value"},
			expected: &osc.Message{
				Address: "/test/value",
				Args: []osc.Arg{
					{Value: "arg1", Type: "s"},
					{Value: int32(42), Type: "i"},
					{Value: float32(3.14), Type: "f"},
				},
			},
		},
		{
			name: "address with template and int64 arg",
			params: map[string]any{
				"address": "/test/{{.Payload.Value}}",
				"args":    []map[string]any{{"value": "42", "type": "h"}},
			},
			payload:  map[string]any{"Value": "value"},
			expected: &osc.Message{Address: "/test/value", Args: []osc.Arg{{Value: int64(42), Type: "h"}}},
		},
		{
			name: "address with template and double arg",
			params: map[string]any{
				"address": "/test/{{.Payload.Value}}",
				"args":    []map[string]any{{"value": "42", "type": "d"}},
			},
			payload:  map[string]any{"Value": "value"},
			expected: &osc.Message{Address: "/test/value", Args: []osc.Arg{{Value: float64(42), Type: "d"}}},
		},
		{
			name: "address with template and true arg",
			params: map[string]any{
				"address": "/test/{{.Payload.Value}}",
				"args":    []map[string]any{{"value": "", "type": "T"}},
			},
			payload:  map[string]any{"Value": "value"},
			expected: &osc.Message{Address: "/test/value", Args: []osc.Arg{{Value: true, Type: "T"}}},
		},
		{
			name: "address with template and false arg",
			params: map[string]any{
				"address": "/test/{{.Payload.Value}}",
				"args":    []map[string]any{{"value": "", "type": "F"}},
			},
			payload:  map[string]any{"Value": "value"},
			expected: &osc.Message{Address: "/test/value", Args: []osc.Arg{{Value: false, Type: "F"}}},
		},
		{
			name: "address with template and nil arg",
			params: map[string]any{
				"address": "/test/{{.Payload.Value}}",
				"args":    []map[string]any{{"value": "", "type": "N"}},
			},
			payload:  map[string]any{"Value": "value"},
			expected: &osc.Message{Address: "/test/value", Args: []osc.Arg{{Value: nil, Type: "N"}}},
		},
		{
			name: "blob arg",
			params: map[string]any{
				"address": "/test",
				"args":    []map[string]any{{"value": "deadbeef", "type": "b"}},
			},
			payload:  "",
			expected: &osc.Message{Address: "/test", Args: []osc.Arg{{Value: []byte{0xde, 0xad, 0xbe, 0xef}, Type: "b"}}},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			registration, ok := processor.GetProcessorRegistration("osc.message.create")
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

			got, err := processorInstance.Process(t.Context(), common.WrappedPayload{Payload: test.payload})

			if err != nil {
				t.Fatalf("osc.message.create processing failed: %s", err)
			}

			if test.expected == nil {
				if got.Payload != nil {
					t.Fatalf("osc.message.create got %+v, expected nil", got)
				}
				return
			}

			gotMessage, ok := got.Payload.(*osc.Message)
			if !ok {
				t.Fatalf("osc.message.create returned a %T payload: %+v", got, got)
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
			},
			payload:     "test",
			errorString: "osc.message.create args error: not a slice",
		},
		{
			name: "args not an object array",
			params: map[string]any{
				"address": "/test",
				"args":    []any{"arg1"},
			},
			payload:     "test",
			errorString: "osc.message.create args error: not an object slice",
		},
		{
			name: "arg value not a string",
			params: map[string]any{
				"address": "/test",
				"args":    []map[string]any{{"value": 123, "type": "s"}},
			},
			payload:     "test",
			errorString: "osc.message.create arg value error: not a string",
		},
		{
			name: "bad arg template",
			params: map[string]any{
				"address": "/test",
				"args":    []map[string]any{{"value": "{{", "type": "s"}},
			},
			payload:     "test",
			errorString: "template: arg:1: unclosed action",
		},
		{
			name: "non-string type parameter",
			params: map[string]any{
				"address": "/test",
				"args":    []map[string]any{{"value": "arg1", "type": 123}},
			},
			payload:     "test",
			errorString: "osc.message.create arg type error: not a string",
		},
		{
			name: "invalid type in types parameter",
			params: map[string]any{
				"address": "/test",
				"args":    []map[string]any{{"value": "arg1", "type": "x"}},
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
			errorString: "template: address:1:8: executing \"address\" at <.missing>: can't evaluate field missing in type common.WrappedPayload",
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
			name: "arg template with missing field",
			params: map[string]any{
				"address": "/test",
				"args":    []map[string]any{{"value": "{{.missing}}", "type": "s"}},
			},
			payload:     "test",
			errorString: "template: arg:1:2: executing \"arg\" at <.missing>: can't evaluate field missing in type common.WrappedPayload",
		},
		{
			name: "wrong arg type for int arg",
			params: map[string]any{
				"address": "/test",
				"args":    []map[string]any{{"value": "{{.Payload}}", "type": "i"}},
			},
			payload:     "test",
			errorString: "strconv.ParseInt: parsing \"test\": invalid syntax",
		},
		{
			name: "wrong arg type for float arg",
			params: map[string]any{
				"address": "/test",
				"args":    []map[string]any{{"value": "{{.Payload}}", "type": "f"}},
			},
			payload:     "test",
			errorString: "strconv.ParseFloat: parsing \"test\": invalid syntax",
		},
		{
			name: "wrong arg type for blob arg",
			params: map[string]any{
				"address": "/test",
				"args":    []map[string]any{{"value": "{{.Payload}}", "type": "b"}},
			},
			payload:     "test",
			errorString: "encoding/hex: invalid byte: U+0074 't'",
		},
		{
			name: "wrong arg type for int64 arg",
			params: map[string]any{
				"address": "/test",
				"args":    []map[string]any{{"value": "{{.Payload}}", "type": "h"}},
			},
			payload:     "test",
			errorString: "strconv.ParseInt: parsing \"test\": invalid syntax",
		},
		{
			name: "wrong arg type for double arg",
			params: map[string]any{
				"address": "/test",
				"args":    []map[string]any{{"value": "{{.Payload}}", "type": "d"}},
			},
			payload:     "test",
			errorString: "strconv.ParseFloat: parsing \"test\": invalid syntax",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			registration, ok := processor.GetProcessorRegistration("osc.message.create")
			if !ok {
				t.Fatalf("osc.message.create processor not registered")
			}

			processorInstance, err := registration.New(config.ProcessorConfig{
				Type:   "osc.message.create",
				Params: test.params,
			})

			if err != nil {
				if test.errorString != err.Error() {
					t.Fatalf("osc.message.create got error '%s', expected '%s'", err.Error(), test.errorString)
				}
				return
			}

			got, err := processorInstance.Process(t.Context(), common.WrappedPayload{Payload: test.payload})

			if err == nil {
				t.Fatalf("osc.message.create expected to fail but succeeded, got: %v", got)
			}

			if err.Error() != test.errorString {
				t.Fatalf("osc.message.create got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}

func BenchmarkOSCMessageCreate(b *testing.B) {
	registration, ok := processor.GetProcessorRegistration("osc.message.create")
	if !ok {
		b.Fatalf("osc.message.create processor not registered")
	}

	processorInstance, err := registration.New(config.ProcessorConfig{
		Type: "osc.message.create",
		Params: map[string]any{
			"address": "/hello",
			"args":    []map[string]any{{"value": "{{.Payload}}", "type": "i"}},
		},
	})

	if err != nil {
		b.Fatalf("osc.message.create failed to create processor: %s", err)
	}

	count := 0
	for b.Loop() {
		_, err := processorInstance.Process(b.Context(), common.WrappedPayload{Payload: count})
		if err != nil {
			b.Fatalf("osc.message.create processing failed: %s", err)
		}
		count++
	}
}
