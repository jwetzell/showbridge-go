package processor_test

import (
	"reflect"
	"testing"

	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/processor"
)

func TestKvGetFromRegistry(t *testing.T) {
	registration, ok := processor.ProcessorRegistry["kv.get"]
	if !ok {
		t.Fatalf("kv.get processor not registered")
	}

	processorInstance, err := registration.New(config.ProcessorConfig{
		Type: "kv.get",
		Params: map[string]any{
			"module": "test",
			"key":    "test",
		},
	})
	if err != nil {
		t.Fatalf("failed to create kv.get processor: %s", err)
	}

	if processorInstance.Type() != "kv.get" {
		t.Fatalf("kv.get processor has wrong type: %s", processorInstance.Type())
	}

	payload := "hello"
	expected := "test"

	got, err := processorInstance.Process(t.Context(), common.GetWrappedPayload(GetTestContext(t.Context()), payload))
	if err != nil {
		t.Fatalf("kv.get processing failed: %s", err)
	}

	if got.Payload != expected {
		t.Fatalf("kv.get got %+v, expected %+v", got.Payload, expected)
	}
}

func TestGoodKvGet(t *testing.T) {

	tests := []struct {
		name     string
		params   map[string]any
		payload  any
		expected any
	}{
		{
			name: "basic value",
			params: map[string]any{
				"module": "test",
				"key":    "test",
				"value":  "hello",
			},
			payload:  "hello",
			expected: "test",
		},
		{
			name: "template value",
			params: map[string]any{
				"module": "test",
				"key":    "test",
				"value":  "{{.Payload}}",
			},
			payload:  "hello",
			expected: "test",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			registration, ok := processor.ProcessorRegistry["kv.get"]
			if !ok {
				t.Fatalf("kv.get processor not registered")
			}

			processorInstance, err := registration.New(config.ProcessorConfig{
				Type:   "kv.get",
				Params: test.params,
			})

			if err != nil {
				t.Fatalf("kv.get failed to create processor: %s", err)
			}

			got, err := processorInstance.Process(t.Context(), common.GetWrappedPayload(GetTestContext(t.Context()), test.payload))

			if err != nil {
				t.Fatalf("kv.get processing failed: %s", err)
			}

			if !reflect.DeepEqual(got.Payload, test.expected) {
				t.Fatalf("kv.get got payload: %+v, expected %+v", got.Payload, test.expected)
			}
		})
	}
}

func TestBadKvGet(t *testing.T) {
	tests := []struct {
		name        string
		params      map[string]any
		payload     any
		errorString string
	}{
		{
			name:    "no module param",
			payload: TestStruct{Data: "hello"},
			params: map[string]any{
				"key": "test",
			},
			errorString: "kv.get module error: not found",
		},
		{
			name:    "non string module",
			payload: TestStruct{Data: "hello"},
			params: map[string]any{
				"module": 1,
				"key":    "test",
			},
			errorString: "kv.get module error: not a string",
		},
		{
			name:    "no key param",
			payload: TestStruct{Data: "hello"},
			params: map[string]any{
				"module": "test",
			},
			errorString: "kv.get key error: not found",
		},
		{
			name:    "non string key",
			payload: TestStruct{Data: "hello"},
			params: map[string]any{
				"module": "test",
				"key":    1,
			},
			errorString: "kv.get key error: not a string",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			registration, ok := processor.ProcessorRegistry["kv.get"]
			if !ok {
				t.Fatalf("kv.get processor not registered")
			}

			processorInstance, err := registration.New(config.ProcessorConfig{
				Type:   "kv.get",
				Params: test.params,
			})

			if err != nil {
				if test.errorString != err.Error() {
					t.Fatalf("kv.get got error '%s', expected '%s'", err.Error(), test.errorString)
				}
				return
			}

			got, err := processorInstance.Process(t.Context(), common.GetWrappedPayload(GetTestContext(t.Context()), test.payload))

			if err == nil {
				t.Fatalf("kv.get expected to fail but got payload: %+v", got)
			}

			if err.Error() != test.errorString {
				t.Fatalf("kv.get got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
