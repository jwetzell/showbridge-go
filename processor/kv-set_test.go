package processor_test

import (
	"reflect"
	"testing"

	"github.com/jwetzell/showbridge-go/common"
	"github.com/jwetzell/showbridge-go/config"
	"github.com/jwetzell/showbridge-go/payload"
	"github.com/jwetzell/showbridge-go/processor"
)

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
			},
			payload:  "test",
			expected: "test",
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			registration, ok := processor.GetProcessorRegistration("kv.set")
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

			got, err := processorInstance.Process(t.Context(), payload.WrappedPayload{
				Payload: testCase.payload,
			}, testModuleManager(map[string]common.Module{
				"test": &TestKVModule{},
			}), nil)

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
		name          string
		params        map[string]any
		payload       any
		moduleManager common.ModuleManager
		errorString   string
	}{
		{
			name:    "no module param",
			payload: TestStruct{Data: "hello"},
			params: map[string]any{
				"key": "test",
			},
			moduleManager: testModuleManager(map[string]common.Module{
				"test": &TestKVModule{},
			}),
			errorString: "kv.set module error: not found",
		},
		{
			name:    "non string module",
			payload: TestStruct{Data: "hello"},
			params: map[string]any{
				"module": 1,
				"key":    "test",
			},
			moduleManager: testModuleManager(map[string]common.Module{
				"test": &TestKVModule{},
			}),
			errorString: "kv.set module error: not a string",
		},
		{
			name:    "no key param",
			payload: TestStruct{Data: "hello"},
			params: map[string]any{
				"module": "test",
			},
			moduleManager: testModuleManager(map[string]common.Module{
				"test": &TestKVModule{},
			}),
			errorString: "kv.set key error: not found",
		},
		{
			name:    "non string key",
			payload: TestStruct{Data: "hello"},
			params: map[string]any{
				"module": "test",
				"key":    1,
			},
			moduleManager: testModuleManager(map[string]common.Module{
				"test": &TestKVModule{},
			}),
			errorString: "kv.set key error: not a string",
		},
		{
			name:    "no modules",
			payload: TestStruct{Data: "hello"},
			params: map[string]any{
				"module": "test",
				"key":    "test",
			},
			moduleManager: nil,
			errorString:   "kv.set module manager is nil",
		},
		{
			name:    "module not found",
			payload: TestStruct{Data: "hello"},
			params: map[string]any{
				"module": "test",
				"key":    "test",
			},
			moduleManager: testModuleManager(map[string]common.Module{}),
			errorString:   "kv.set unable to find module with id: test",
		},
		{
			name:    "module not a kv module",
			payload: TestStruct{Data: "hello"},
			params: map[string]any{
				"module": "test",
				"key":    "test",
			},
			moduleManager: testModuleManager(map[string]common.Module{
				"test": NewTestDBModule("test"),
			}),
			errorString: "kv.set module with id test is not a KeyValueModule",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {

			registration, ok := processor.GetProcessorRegistration("kv.set")
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

			got, err := processorInstance.Process(t.Context(), payload.WrappedPayload{Payload: testCase.payload}, testCase.moduleManager, nil)

			if err == nil {
				t.Fatalf("kv.set expected to fail but got payload: %+v", got)
			}

			if err.Error() != testCase.errorString {
				t.Fatalf("kv.set got error '%s', expected '%s'", err.Error(), testCase.errorString)
			}
		})
	}
}
