package processor_test

import (
	"reflect"
	"testing"

	"github.com/jwetzell/showbridge-go/config"
	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/processor"
	"github.com/jwetzell/showbridge-go/internal/test"
	_ "modernc.org/sqlite"
)

func TestGoodPubSubPublish(t *testing.T) {

	testCases := []struct {
		name     string
		params   map[string]any
		payload  any
		expected any
	}{
		{
			name: "basic topic",
			params: map[string]any{
				"module": "test",
				"topic":  "test",
			},
			payload:  "",
			expected: "",
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			registration, ok := processor.GetProcessorRegistration("pubsub.publish")
			if !ok {
				t.Fatalf("pubsub.publish processor not registered")
			}

			processorInstance, err := registration.New(config.ProcessorConfig{
				Type:   "pubsub.publish",
				Params: testCase.params,
			})

			if err != nil {
				t.Fatalf("pubsub.publish failed to create processor: %s", err)
			}

			got, err := processorInstance.Process(t.Context(), common.WrappedPayload{
				Payload: testCase.payload,
			}, testModuleManager(map[string]common.Module{
				"test": test.NewTestPubSubModule("test"),
			}), nil)

			if err != nil {
				t.Fatalf("pubsub.publish processing failed: %s", err)
			}

			if !reflect.DeepEqual(got.Payload, testCase.expected) {
				t.Fatalf("pubsub.publish got payload: %+v, expected %+v", got.Payload, testCase.expected)
			}
		})
	}
}

func TestBadPubSubPublish(t *testing.T) {
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
				"topic": "test",
			},
			moduleManager: testModuleManager(map[string]common.Module{
				"test": test.NewTestPubSubModule("test"),
			}),
			errorString: "pubsub.publish module error: not found",
		},
		{
			name:    "non string module",
			payload: test.TestStruct{Data: "hello"},
			params: map[string]any{
				"module": 1,
				"topic":  "test",
			},
			moduleManager: testModuleManager(map[string]common.Module{
				"test": test.NewTestPubSubModule("test"),
			}),
			errorString: "pubsub.publish module error: not a string",
		},
		{
			name:    "no topic param",
			payload: test.TestStruct{Data: "hello"},
			params: map[string]any{
				"module": "test",
			},
			moduleManager: testModuleManager(map[string]common.Module{
				"test": test.NewTestPubSubModule("test"),
			}),
			errorString: "pubsub.publish topic error: not found",
		},
		{
			name:    "non string topic",
			payload: test.TestStruct{Data: "hello"},
			params: map[string]any{
				"module": "test",
				"topic":  1,
			},
			moduleManager: testModuleManager(map[string]common.Module{
				"test": test.NewTestPubSubModule("test"),
			}),
			errorString: "pubsub.publish topic error: not a string",
		},
		{
			name:    "topic template syntax error",
			payload: test.TestStruct{Data: "hello"},
			params: map[string]any{
				"module": "test",
				"topic":  "{{",
			},
			moduleManager: testModuleManager(map[string]common.Module{
				"test": test.NewTestPubSubModule("test"),
			}),
			errorString: "template: topic:1: unclosed action",
		},
		{
			name:    "topic template error",
			payload: test.TestStruct{Data: "hello"},
			params: map[string]any{
				"module": "test",
				"topic":  "{{.Data}}",
			},
			moduleManager: testModuleManager(map[string]common.Module{
				"test": test.NewTestPubSubModule("test"),
			}),
			errorString: "template: topic:1:2: executing \"topic\" at <.Data>: can't evaluate field Data in type common.WrappedPayload",
		},
		{
			name:    "no modules",
			payload: test.TestStruct{Data: "hello"},
			params: map[string]any{
				"module": "test",
				"topic":  "test",
			},
			moduleManager: nil,
			errorString:   "pubsub.publish module manager is nil",
		},
		{
			name:    "module not found",
			payload: test.TestStruct{Data: "hello"},
			params: map[string]any{
				"module": "test",
				"topic":  "test",
			},
			moduleManager: testModuleManager(map[string]common.Module{}),
			errorString:   "pubsub.publish unable to find module with id: test",
		},
		{
			name:    "module not an OutputModule",
			payload: test.TestStruct{Data: "hello"},
			params: map[string]any{
				"module": "test",
				"topic":  "test",
			},
			moduleManager: testModuleManager(map[string]common.Module{
				"test": test.NewTestKVModule("test", nil),
			}),
			errorString: "pubsub.publish module with id test is not a PubSubModule",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			registration, ok := processor.GetProcessorRegistration("pubsub.publish")
			if !ok {
				t.Fatalf("pubsub.publish processor not registered")
			}

			processorInstance, err := registration.New(config.ProcessorConfig{
				Type:   "pubsub.publish",
				Params: test.params,
			})

			if err != nil {
				if test.errorString != err.Error() {
					t.Fatalf("pubsub.publish got error '%s', expected '%s'", err.Error(), test.errorString)
				}
				return
			}

			got, err := processorInstance.Process(t.Context(), common.WrappedPayload{
				Payload: test.payload,
			}, test.moduleManager, nil)

			if err == nil {
				t.Fatalf("pubsub.publish expected to fail but got payload: %+v", got)
			}

			if err.Error() != test.errorString {
				t.Fatalf("pubsub.publish got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
