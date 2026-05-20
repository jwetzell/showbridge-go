package processor_test

import (
	"reflect"
	"testing"

	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/processor"
	"github.com/jwetzell/showbridge-go/internal/test"
	_ "modernc.org/sqlite"
)

func TestPubSubPublishFromRegistry(t *testing.T) {
	registration, ok := processor.ProcessorRegistry["pubsub.publish"]
	if !ok {
		t.Fatalf("pubsub.publish processor not registered")
	}

	processorInstance, err := registration.New(config.ProcessorConfig{
		Type: "pubsub.publish",
		Params: map[string]any{
			"module": "test",
			"topic":  "test",
		},
	})
	if err != nil {
		t.Fatalf("failed to create pubsub.publish processor: %s", err)
	}

	if processorInstance.Type() != "pubsub.publish" {
		t.Fatalf("pubsub.publish processor has wrong type: %s", processorInstance.Type())
	}

	payload := "hello"
	expected := "hello"

	got, err := processorInstance.Process(t.Context(), common.WrappedPayload{
		Payload: payload,
		Modules: map[string]common.Module{
			"test": test.NewTestPubSubModule("test"),
		},
	})
	if err != nil {
		t.Fatalf("pubsub.publish processing failed: %s", err)
	}

	if !reflect.DeepEqual(got.Payload, expected) {
		t.Fatalf("pubsub.publish got %+v, expected %+v", got.Payload, expected)
	}
}

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
			registration, ok := processor.ProcessorRegistry["pubsub.publish"]
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
				Modules: map[string]common.Module{
					"test": test.NewTestPubSubModule("test"),
				},
				Payload: testCase.payload,
			})

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
		name                  string
		params                map[string]any
		payload               any
		wrappedPayloadModules map[string]common.Module
		errorString           string
	}{
		{
			name:    "no module param",
			payload: test.TestStruct{Data: "hello"},
			params: map[string]any{
				"topic": "test",
			},
			wrappedPayloadModules: map[string]common.Module{
				"test": test.NewTestPubSubModule("test"),
			},
			errorString: "pubsub.publish module error: not found",
		},
		{
			name:    "non string module",
			payload: test.TestStruct{Data: "hello"},
			params: map[string]any{
				"module": 1,
				"topic":  "test",
			},
			wrappedPayloadModules: map[string]common.Module{
				"test": test.NewTestPubSubModule("test"),
			},
			errorString: "pubsub.publish module error: not a string",
		},
		{
			name:    "no topic param",
			payload: test.TestStruct{Data: "hello"},
			params: map[string]any{
				"module": "test",
			},
			wrappedPayloadModules: map[string]common.Module{
				"test": test.NewTestPubSubModule("test"),
			},
			errorString: "pubsub.publish topic error: not found",
		},
		{
			name:    "non string topic",
			payload: test.TestStruct{Data: "hello"},
			params: map[string]any{
				"module": "test",
				"topic":  1,
			},
			wrappedPayloadModules: map[string]common.Module{
				"test": test.NewTestPubSubModule("test"),
			},
			errorString: "pubsub.publish topic error: not a string",
		},
		{
			name:    "topic template syntax error",
			payload: test.TestStruct{Data: "hello"},
			params: map[string]any{
				"module": "test",
				"topic":  "{{",
			},
			wrappedPayloadModules: map[string]common.Module{
				"test": test.NewTestPubSubModule("test"),
			},
			errorString: "template: topic:1: unclosed action",
		},
		{
			name:    "topic template error",
			payload: test.TestStruct{Data: "hello"},
			params: map[string]any{
				"module": "test",
				"topic":  "{{.Data}}",
			},
			wrappedPayloadModules: map[string]common.Module{
				"test": test.NewTestPubSubModule("test"),
			},
			errorString: "template: topic:1:2: executing \"topic\" at <.Data>: can't evaluate field Data in type common.WrappedPayload",
		},
		{
			name:    "no modules in context",
			payload: test.TestStruct{Data: "hello"},
			params: map[string]any{
				"module": "test",
				"topic":  "test",
			},
			wrappedPayloadModules: nil,
			errorString:           "pubsub.publish wrapped payload has no modules",
		},
		{
			name:    "module not found in context",
			payload: test.TestStruct{Data: "hello"},
			params: map[string]any{
				"module": "test",
				"topic":  "test",
			},
			wrappedPayloadModules: map[string]common.Module{},
			errorString:           "pubsub.publish unable to find module with id: test",
		},
		{
			name:    "module not an OutputModule",
			payload: test.TestStruct{Data: "hello"},
			params: map[string]any{
				"module": "test",
				"topic":  "test",
			},
			wrappedPayloadModules: map[string]common.Module{
				"test": test.NewTestKVModule("test", nil),
			},
			errorString: "pubsub.publish module with id test is not an OutputModule",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			registration, ok := processor.ProcessorRegistry["pubsub.publish"]
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
				Modules: test.wrappedPayloadModules,
			})

			if err == nil {
				t.Fatalf("pubsub.publish expected to fail but got payload: %+v", got)
			}

			if err.Error() != test.errorString {
				t.Fatalf("pubsub.publish got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
