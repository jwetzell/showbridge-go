package processor_test

import (
	"reflect"
	"testing"

	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/processor"
	"github.com/jwetzell/showbridge-go/internal/test"
)

func TestModuleOutputFromRegistry(t *testing.T) {
	registration, ok := processor.ProcessorRegistry["module.output"]
	if !ok {
		t.Fatalf("module.output processor not registered")
	}

	processorInstance, err := registration.New(config.ProcessorConfig{
		Type: "module.output",
		Params: config.Params{
			"module": "test",
		},
	})

	if err != nil {
		t.Fatalf("failed to create module.output processor: %s", err)
	}

	if processorInstance.Type() != "module.output" {
		t.Fatalf("module.output processor has wrong type: %s", processorInstance.Type())
	}

	payload := "test"
	expected := "test"

	got, err := processorInstance.Process(t.Context(), common.WrappedPayload{
		Router:  test.GetNewTestRouter(),
		Modules: map[string]common.Module{"test": &test.TestOutputModule{}},
		Payload: payload,
	})
	if err != nil {
		t.Fatalf("module.output processing failed: %s", err)
	}

	if got.Payload != expected {
		t.Fatalf("module.output got %+v, expected %+v", got, expected)
	}
}

func TestGoodModuleOutput(t *testing.T) {

	testCases := []struct {
		name     string
		params   map[string]any
		payload  any
		expected any
	}{}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {

			registration, ok := processor.ProcessorRegistry["module.output"]
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

			got, err := processorInstance.Process(t.Context(), common.WrappedPayload{Payload: testCase.payload})
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
		name        string
		params      map[string]any
		payload     any
		modules     map[string]common.Module
		errorString string
	}{
		{
			name:        "no module param",
			params:      map[string]any{},
			payload:     "test",
			modules:     map[string]common.Module{"test": &test.TestModule{}},
			errorString: "module.output module error: not found",
		},
		{
			name: "non-string module",
			params: map[string]any{
				"module": 123,
			},
			payload:     "test",
			modules:     map[string]common.Module{"test": &test.TestModule{}},
			errorString: "module.output module error: not a string",
		},
		{
			name: "modules not found in context",
			params: map[string]any{
				"module": "test",
			},
			payload:     "test",
			modules:     nil,
			errorString: "module.output wrapped payload has no modules",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {

			registration, ok := processor.ProcessorRegistry["module.output"]
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

			got, err := processorInstance.Process(t.Context(), common.WrappedPayload{Modules: testCase.modules, Payload: testCase.payload})

			if err == nil {
				t.Fatalf("module.output expected to fail but succeeded, got: %v", got)
			}

			if err.Error() != testCase.errorString {
				t.Fatalf("module.output got error '%s', expected '%s'", err.Error(), testCase.errorString)
			}
		})
	}
}
