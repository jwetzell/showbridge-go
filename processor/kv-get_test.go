package processor_test

import (
	"reflect"
	"testing"

	"github.com/jwetzell/showbridge-go/common"
	"github.com/jwetzell/showbridge-go/config"
	"github.com/jwetzell/showbridge-go/internal/test"
	"github.com/jwetzell/showbridge-go/processor"
)

func TestGoodKvGet(t *testing.T) {

	testCases := []struct {
		name         string
		presetValues map[string]any
		params       map[string]any
		payload      any
		expected     any
	}{
		{
			name: "basic key",
			presetValues: map[string]any{
				"test": "hello",
			},
			params: map[string]any{
				"module": "test",
				"key":    "test",
			},
			payload:  "hello",
			expected: "hello",
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			registration, ok := processor.GetProcessorRegistration("kv.get")
			if !ok {
				t.Fatalf("kv.get processor not registered")
			}

			processorInstance, err := registration.New(config.ProcessorConfig{
				Type:   "kv.get",
				Params: testCase.params,
			})

			if err != nil {
				t.Fatalf("kv.get failed to create processor: %s", err)
			}

			got, err := processorInstance.Process(t.Context(), common.WrappedPayload{
				Payload: testCase.payload,
			}, testModuleManager(map[string]common.Module{
				"test": test.NewTestKVModule("test", testCase.presetValues),
			}), nil)

			if err != nil {
				t.Fatalf("kv.get processing failed: %s", err)
			}

			if !reflect.DeepEqual(got.Payload, testCase.expected) {
				t.Fatalf("kv.get got payload: %+v, expected %+v", got.Payload, testCase.expected)
			}
		})
	}
}

func TestBadKvGet(t *testing.T) {
	tests := []struct {
		name          string
		params        map[string]any
		payload       any
		moduleManager common.ModuleManager
		errorString   string
	}{
		{
			name:    "no module param",
			payload: test.TestStruct{Data: "hello"},
			params: map[string]any{
				"key": "test",
			},
			moduleManager: testModuleManager(map[string]common.Module{
				"test": test.NewTestKVModule("test", nil),
			}),
			errorString: "kv.get module error: not found",
		},
		{
			name:    "non string module",
			payload: test.TestStruct{Data: "hello"},
			params: map[string]any{
				"module": 1,
				"key":    "test",
			},
			moduleManager: testModuleManager(map[string]common.Module{
				"test": test.NewTestKVModule("test", nil),
			}),
			errorString: "kv.get module error: not a string",
		},
		{
			name:    "no key param",
			payload: test.TestStruct{Data: "hello"},
			params: map[string]any{
				"module": "test",
			},
			moduleManager: testModuleManager(map[string]common.Module{
				"test": test.NewTestKVModule("test", nil),
			}),
			errorString: "kv.get key error: not found",
		},
		{
			name:    "non string key",
			payload: test.TestStruct{Data: "hello"},
			params: map[string]any{
				"module": "test",
				"key":    1,
			},
			moduleManager: testModuleManager(map[string]common.Module{
				"test": test.NewTestKVModule("test", nil),
			}),
			errorString: "kv.get key error: not a string",
		},
		{
			name:    "no modules",
			payload: test.TestStruct{Data: "hello"},
			params: map[string]any{
				"module": "test",
				"key":    "test",
			},
			moduleManager: nil,
			errorString:   "kv.get module manager is nil",
		},
		{
			name:    "module not found",
			payload: test.TestStruct{Data: "hello"},
			params: map[string]any{
				"module": "test",
				"key":    "test",
			},
			moduleManager: testModuleManager(map[string]common.Module{}),
			errorString:   "kv.get unable to find module with id: test",
		},
		{
			name:    "module not a kv module",
			payload: test.TestStruct{Data: "hello"},
			params: map[string]any{
				"module": "test",
				"key":    "test",
			},
			moduleManager: testModuleManager(map[string]common.Module{
				"test": test.NewTestDBModule("test"),
			}),
			errorString: "kv.get module with id test is not a KeyValueModule",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			registration, ok := processor.GetProcessorRegistration("kv.get")
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

			got, err := processorInstance.Process(t.Context(), common.WrappedPayload{Payload: test.payload}, test.moduleManager, nil)

			if err == nil {
				t.Fatalf("kv.get expected to fail but got payload: %+v", got)
			}

			if err.Error() != test.errorString {
				t.Fatalf("kv.get got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
