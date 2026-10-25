package processor_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/jwetzell/showbridge-go/common"
	"github.com/jwetzell/showbridge-go/config"
	"github.com/jwetzell/showbridge-go/payload"
	"github.com/jwetzell/showbridge-go/processor"
)

func NewTestModule(id string) *TestModule {
	return &TestModule{
		id: id,
	}
}

type TestModule struct {
	id string
}

func (m *TestModule) Start(ctx context.Context, inputHandler common.InputHandler) error {
	<-ctx.Done()
	return nil
}

func (m *TestModule) Stop() {}

func TestGoodModuleOutput(t *testing.T) {

	testCases := []struct {
		name     string
		params   map[string]any
		payload  any
		expected any
	}{}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {

			registration, ok := processor.GetProcessorRegistration("module.output")
			if !ok {
				t.Fatalf("module.output processor not registered")
			}

			processorInstance, err := registration.New(config.ProcessorConfig{
				Type:   "module.output",
				Params: testCase.params,
			})

			if err != nil {
				t.Fatalf("module.output failed to create processor: %s", err)
			}

			got, err := processorInstance.Process(t.Context(), payload.WrappedPayload{Payload: testCase.payload}, nil, nil)
			if err != nil {
				t.Fatalf("module.output processing failed: %s", err)
			}

			if !reflect.DeepEqual(got.Payload, testCase.expected) {
				t.Fatalf("module.output got %+v (%T), expected %+v (%T)", got.Payload, got.Payload, testCase.expected, testCase.expected)
			}
		})
	}
}

func TestBadModuleOutput(t *testing.T) {
	testCases := []struct {
		name          string
		params        map[string]any
		payload       any
		moduleManager common.ModuleManager
		errorString   string
	}{
		{
			name:          "no module param",
			params:        map[string]any{},
			payload:       "test",
			moduleManager: testModuleManager(map[string]common.Module{"test": NewTestModule("test")}),
			errorString:   "module.output module error: not found",
		},
		{
			name: "non-string module",
			params: map[string]any{
				"module": 123,
			},
			payload:       "test",
			moduleManager: testModuleManager(map[string]common.Module{"test": NewTestModule("test")}),
			errorString:   "module.output module error: not a string",
		},
		{
			name: "modules does not exist",
			params: map[string]any{
				"module": "test",
			},
			payload:       "test",
			moduleManager: nil,
			errorString:   "module.output module manager is nil",
		},
		{
			name: "module not found",
			params: map[string]any{
				"module": "test",
			},
			payload:       "test",
			moduleManager: testModuleManager(map[string]common.Module{}),
			errorString:   "module.output unable to find module with id: test",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {

			registration, ok := processor.GetProcessorRegistration("module.output")
			if !ok {
				t.Fatalf("module.output processor not registered")
			}

			processorInstance, err := registration.New(config.ProcessorConfig{
				Type:   "module.output",
				Params: testCase.params,
			})

			if err != nil {
				if testCase.errorString != err.Error() {
					t.Fatalf("module.output got error '%s', expected '%s'", err.Error(), testCase.errorString)
				}
				return
			}

			got, err := processorInstance.Process(t.Context(), payload.WrappedPayload{Payload: testCase.payload}, testCase.moduleManager, nil)

			if err == nil {
				t.Fatalf("module.output expected to fail but succeeded, got: %v", got)
			}

			if err.Error() != testCase.errorString {
				t.Fatalf("module.output got error '%s', expected '%s'", err.Error(), testCase.errorString)
			}
		})
	}
}
