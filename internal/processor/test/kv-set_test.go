package processor_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/processor"
	"github.com/jwetzell/showbridge-go/internal/test"
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

	got, err := processorInstance.Process(t.Context(), common.GetWrappedPayload(test.GetContextWithModules(
		t.Context(),
		map[string]common.Module{
			"test": &test.TestKVModule{},
		},
	), payload))
	if err != nil {
		t.Fatalf("kv.set processing failed: %s", err)
	}

	if got.Payload != expected {
		t.Fatalf("kv.set got %+v, expected %+v", got.Payload, expected)
	}
}

func TestGoodKvSet(t *testing.T) {

	testCases := []struct {
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
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			registration, ok := processor.ProcessorRegistry["kv.set"]
			if !ok {
				t.Fatalf("kv.set processor not registered")
			}

			processorInstance, err := registration.New(config.ProcessorConfig{
				Type:   "kv.set",
				Params: testCase.params,
			})

			if err != nil {
				t.Fatalf("kv.set failed to create processor: %s", err)
			}

			got, err := processorInstance.Process(t.Context(), common.GetWrappedPayload(test.GetContextWithModules(
				t.Context(),
				map[string]common.Module{
					"test": &test.TestKVModule{},
				},
			), testCase.payload))

			if err != nil {
				t.Fatalf("kv.set processing failed: %s", err)
			}

			if !reflect.DeepEqual(got.Payload, testCase.expected) {
				t.Fatalf("kv.set got payload: %+v, expected %+v", got.Payload, testCase.expected)
			}
		})
	}
}

func TestBadKvSet(t *testing.T) {
	testCases := []struct {
		name              string
		params            map[string]any
		payload           any
		wrappedPayloadCtx context.Context
		errorString       string
	}{
		{
			name:    "no module param",
			payload: test.TestStruct{Data: "hello"},
			params: map[string]any{
				"key":   "test",
				"value": "test",
			},
			wrappedPayloadCtx: test.GetContextWithModules(t.Context(), map[string]common.Module{
				"test": &test.TestKVModule{},
			}),
			errorString: "kv.set module error: not found",
		},
		{
			name:    "non string module",
			payload: test.TestStruct{Data: "hello"},
			params: map[string]any{
				"module": 1,
				"key":    "test",
				"value":  "test",
			},
			wrappedPayloadCtx: test.GetContextWithModules(t.Context(), map[string]common.Module{
				"test": &test.TestKVModule{},
			}),
			errorString: "kv.set module error: not a string",
		},
		{
			name:    "no key param",
			payload: test.TestStruct{Data: "hello"},
			params: map[string]any{
				"module": "test",
				"value":  "test",
			},
			wrappedPayloadCtx: test.GetContextWithModules(t.Context(), map[string]common.Module{
				"test": &test.TestKVModule{},
			}),
			errorString: "kv.set key error: not found",
		},
		{
			name:    "non string key",
			payload: test.TestStruct{Data: "hello"},
			params: map[string]any{
				"module": "test",
				"key":    1,
				"value":  "test",
			},
			wrappedPayloadCtx: test.GetContextWithModules(t.Context(), map[string]common.Module{
				"test": &test.TestKVModule{},
			}),
			errorString: "kv.set key error: not a string",
		},
		{
			name:    "no value param",
			payload: test.TestStruct{Data: "hello"},
			params: map[string]any{
				"module": "test",
				"key":    "test",
			},
			wrappedPayloadCtx: test.GetContextWithModules(t.Context(), map[string]common.Module{
				"test": &test.TestKVModule{},
			}),
			errorString: "kv.set value error: not found",
		},
		{
			name:    "non string value",
			payload: test.TestStruct{Data: "hello"},
			params: map[string]any{
				"module": "test",
				"key":    "test",
				"value":  1,
			},
			wrappedPayloadCtx: test.GetContextWithModules(t.Context(), map[string]common.Module{
				"test": &test.TestKVModule{},
			}),
			errorString: "kv.set value error: not a string",
		},
		{
			name:    "no modules in context",
			payload: test.TestStruct{Data: "hello"},
			params: map[string]any{
				"module": "test",
				"key":    "test",
				"value":  "hello",
			},
			wrappedPayloadCtx: t.Context(),
			errorString:       "kv.set wrapped payload has no modules",
		},
		{
			name:    "value template syntax error",
			payload: test.TestStruct{Data: "hello"},
			params: map[string]any{
				"module": "test",
				"key":    "test",
				"value":  "{{",
			},
			wrappedPayloadCtx: test.GetContextWithModules(t.Context(), map[string]common.Module{
				"test": &test.TestKVModule{},
			}),
			errorString: "template: template:1: unclosed action",
		},
		{
			name:    "value template execution error",
			payload: test.TestStruct{Data: "hello"},
			params: map[string]any{
				"module": "test",
				"key":    "test",
				"value":  "{{.Data}}",
			},
			wrappedPayloadCtx: test.GetContextWithModules(t.Context(), map[string]common.Module{
				"test": &test.TestKVModule{},
			}),
			errorString: "template: template:1:2: executing \"template\" at <.Data>: can't evaluate field Data in type common.WrappedPayload",
		},
		{
			name:    "module not found in context",
			payload: test.TestStruct{Data: "hello"},
			params: map[string]any{
				"module": "test",
				"key":    "test",
				"value":  "hello",
			},
			wrappedPayloadCtx: test.GetContextWithModules(t.Context(), map[string]common.Module{}),
			errorString:       "kv.set unable to find module with id: test",
		},
		{
			name:    "module not a kv module",
			payload: test.TestStruct{Data: "hello"},
			params: map[string]any{
				"module": "test",
				"key":    "test",
				"value":  "hello",
			},
			wrappedPayloadCtx: test.GetContextWithModules(t.Context(), map[string]common.Module{
				"test": test.NewTestDBModule("test"),
			}),
			errorString: "kv.set module with id test is not a KeyValueModule",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {

			registration, ok := processor.ProcessorRegistry["kv.set"]
			if !ok {
				t.Fatalf("kv.set processor not registered")
			}

			processorInstance, err := registration.New(config.ProcessorConfig{
				Type:   "kv.set",
				Params: testCase.params,
			})

			if err != nil {
				if testCase.errorString != err.Error() {
					t.Fatalf("kv.set got error '%s', expected '%s'", err.Error(), testCase.errorString)
				}
				return
			}

			got, err := processorInstance.Process(t.Context(), common.GetWrappedPayload(testCase.wrappedPayloadCtx, testCase.payload))

			if err == nil {
				t.Fatalf("kv.set expected to fail but got payload: %+v", got)
			}

			if err.Error() != testCase.errorString {
				t.Fatalf("kv.set got error '%s', expected '%s'", err.Error(), testCase.errorString)
			}
		})
	}
}
