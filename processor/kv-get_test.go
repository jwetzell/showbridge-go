package processor_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/jwetzell/showbridge-go/config"
	"github.com/jwetzell/showbridge-go/processor"
	"github.com/jwetzell/showbridge-go/routing"
)

func NewTestKVModule(id string, presetValues map[string]any) *TestKVModule {
	return &TestKVModule{
		id:     id,
		kvData: presetValues,
	}
}

type TestKVModule struct {
	id     string
	kvData map[string]any
}

func (m *TestKVModule) Start(ctx context.Context, inputHandler routing.InputHandler) error {
	<-ctx.Done()
	return nil
}

func (m *TestKVModule) Stop() {}

func (m *TestKVModule) Get(ctx context.Context, key string) (any, error) {
	if m.kvData == nil {
		return nil, nil
	}
	value, ok := m.kvData[key]
	if !ok {
		return nil, nil
	}
	return value, nil
}

func (m *TestKVModule) Set(ctx context.Context, key string, value any) error {
	if m.kvData == nil {
		m.kvData = make(map[string]any)
	}
	m.kvData[key] = value
	return nil
}

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

			got, err := processorInstance.Process(t.Context(), routing.WrappedPayload{
				Payload: testCase.payload,
			}, testModuleManager(map[string]routing.Module{
				"test": NewTestKVModule("test", testCase.presetValues),
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
		moduleManager routing.ModuleManager
		errorString   string
	}{
		{
			name:    "no module param",
			payload: TestStruct{Data: "hello"},
			params: map[string]any{
				"key": "test",
			},
			moduleManager: testModuleManager(map[string]routing.Module{
				"test": NewTestKVModule("test", nil),
			}),
			errorString: "kv.get module error: not found",
		},
		{
			name:    "non string module",
			payload: TestStruct{Data: "hello"},
			params: map[string]any{
				"module": 1,
				"key":    "test",
			},
			moduleManager: testModuleManager(map[string]routing.Module{
				"test": NewTestKVModule("test", nil),
			}),
			errorString: "kv.get module error: not a string",
		},
		{
			name:    "no key param",
			payload: TestStruct{Data: "hello"},
			params: map[string]any{
				"module": "test",
			},
			moduleManager: testModuleManager(map[string]routing.Module{
				"test": NewTestKVModule("test", nil),
			}),
			errorString: "kv.get key error: not found",
		},
		{
			name:    "non string key",
			payload: TestStruct{Data: "hello"},
			params: map[string]any{
				"module": "test",
				"key":    1,
			},
			moduleManager: testModuleManager(map[string]routing.Module{
				"test": NewTestKVModule("test", nil),
			}),
			errorString: "kv.get key error: not a string",
		},
		{
			name:    "no modules",
			payload: TestStruct{Data: "hello"},
			params: map[string]any{
				"module": "test",
				"key":    "test",
			},
			moduleManager: nil,
			errorString:   "kv.get module manager is nil",
		},
		{
			name:    "module not found",
			payload: TestStruct{Data: "hello"},
			params: map[string]any{
				"module": "test",
				"key":    "test",
			},
			moduleManager: testModuleManager(map[string]routing.Module{}),
			errorString:   "kv.get unable to find module with id: test",
		},
		{
			name:    "module not a kv module",
			payload: TestStruct{Data: "hello"},
			params: map[string]any{
				"module": "test",
				"key":    "test",
			},
			moduleManager: testModuleManager(map[string]routing.Module{
				"test": NewTestDBModule("test"),
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

			got, err := processorInstance.Process(t.Context(), routing.WrappedPayload{Payload: test.payload}, test.moduleManager, nil)

			if err == nil {
				t.Fatalf("kv.get expected to fail but got payload: %+v", got)
			}

			if err.Error() != test.errorString {
				t.Fatalf("kv.get got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
