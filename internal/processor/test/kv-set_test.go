package processor_test

import (
	"reflect"
	"testing"

	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/processor"
)

func TestKvSetFromRegistry(t *testing.T) {
	registration, ok := processor.ProcessorRegistry["kv.set"]
	if !ok {
		t.Fatalf("kv.set processor not registered")
	}

	processorInstance, err := registration.New(config.ProcessorConfig{
		Type: "kv.set",
		Params: map[string]any{
			"module": "test",
			"key":    "test",
			"value":  "hello",
		},
	})
	if err != nil {
		t.Fatalf("failed to create kv.set processor: %s", err)
	}

	if processorInstance.Type() != "kv.set" {
		t.Fatalf("kv.set processor has wrong type: %s", processorInstance.Type())
	}

	payload := ""
	expected := ""

	got, err := processorInstance.Process(t.Context(), common.GetWrappedPayload(GetTestContext(t.Context()), payload))
	if err != nil {
		t.Fatalf("kv.set processing failed: %s", err)
	}

	if got.Payload != expected {
		t.Fatalf("kv.set got %+v, expected %+v", got.Payload, expected)
	}
}

func TestGoodKvSet(t *testing.T) {

	tests := []struct {
		name     string
		params   map[string]any
		payload  any
		expected any
	}{
		{
			name: "basic key/value",
			params: map[string]any{
				"module": "test",
				"key":    "test",
				"value":  "hello",
			},
			payload:  "",
			expected: "",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			registration, ok := processor.ProcessorRegistry["kv.set"]
			if !ok {
				t.Fatalf("kv.set processor not registered")
			}

			processorInstance, err := registration.New(config.ProcessorConfig{
				Type:   "kv.set",
				Params: test.params,
			})

			if err != nil {
				t.Fatalf("kv.set failed to create processor: %s", err)
			}

			got, err := processorInstance.Process(t.Context(), common.GetWrappedPayload(GetTestContext(t.Context()), test.payload))

			if err != nil {
				t.Fatalf("kv.set processing failed: %s", err)
			}

			if !reflect.DeepEqual(got.Payload, test.expected) {
				t.Fatalf("kv.set got payload: %+v, expected %+v", got.Payload, test.expected)
			}
		})
	}
}

func TestBadKvSet(t *testing.T) {
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
				"key":   "test",
				"value": "test",
			},
			errorString: "kv.set module error: not found",
		},
		{
			name:    "non string module",
			payload: TestStruct{Data: "hello"},
			params: map[string]any{
				"module": 1,
				"key":    "test",
				"value":  "test",
			},
			errorString: "kv.set module error: not a string",
		},
		{
			name:    "no key param",
			payload: TestStruct{Data: "hello"},
			params: map[string]any{
				"module": "test",
				"value":  "test",
			},
			errorString: "kv.set key error: not found",
		},
		{
			name:    "non string key",
			payload: TestStruct{Data: "hello"},
			params: map[string]any{
				"module": "test",
				"key":    1,
				"value":  "test",
			},
			errorString: "kv.set key error: not a string",
		},
		{
			name:    "no value param",
			payload: TestStruct{Data: "hello"},
			params: map[string]any{
				"module": "test",
				"key":    "test",
			},
			errorString: "kv.set value error: not found",
		},
		{
			name:    "non string value",
			payload: TestStruct{Data: "hello"},
			params: map[string]any{
				"module": "test",
				"key":    "test",
				"value":  1,
			},
			errorString: "kv.set value error: not a string",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			registration, ok := processor.ProcessorRegistry["kv.set"]
			if !ok {
				t.Fatalf("kv.set processor not registered")
			}

			processorInstance, err := registration.New(config.ProcessorConfig{
				Type:   "kv.set",
				Params: test.params,
			})

			if err != nil {
				if test.errorString != err.Error() {
					t.Fatalf("kv.set got error '%s', expected '%s'", err.Error(), test.errorString)
				}
				return
			}

			got, err := processorInstance.Process(t.Context(), common.GetWrappedPayload(GetTestContext(t.Context()), test.payload))

			if err == nil {
				t.Fatalf("kv.set expected to fail but got payload: %+v", got)
			}

			if err.Error() != test.errorString {
				t.Fatalf("kv.set got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
