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

func TestGoodRouterInput(t *testing.T) {

	testCases := []struct {
		name     string
		params   map[string]any
		payload  any
		expected any
	}{}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {

			registration, ok := processor.GetProcessorRegistration("router.input")
			if !ok {
				t.Fatalf("router.input processor not registered")
			}

			processorInstance, err := registration.New(config.ProcessorConfig{
				Type:   "router.input",
				Params: testCase.params,
			})

			if err != nil {
				t.Fatalf("router.input failed to create processor: %s", err)
			}

			got, err := processorInstance.Process(t.Context(), payload.WrappedPayload{Payload: testCase.payload}, nil, nil)
			if err != nil {
				t.Fatalf("router.input processing failed: %s", err)
			}

			if !reflect.DeepEqual(got.Payload, testCase.expected) {
				t.Fatalf("router.input got %+v (%T), expected %+v (%T)", got.Payload, got.Payload, testCase.expected, testCase.expected)
			}
		})
	}
}

func TestBadRouterInput(t *testing.T) {
	testCases := []struct {
		name         string
		params       map[string]any
		payload      any
		inputHandler common.InputHandler
		errorString  string
	}{
		{
			name:    "no source param",
			params:  map[string]any{},
			payload: "test",
			inputHandler: func(ctx context.Context, sourceId string, payload any) (bool, []common.RouteIOError) {
				return false, nil
			},
			errorString: "router.input source error: not found",
		},
		{
			name: "non-string source",
			params: map[string]any{
				"source": 123,
			},
			payload: "test",
			inputHandler: func(ctx context.Context, sourceId string, payload any) (bool, []common.RouteIOError) {
				return false, nil
			},
			errorString: "router.input source error: not a string",
		},
		{
			name: "router not found in context",
			params: map[string]any{
				"source": "test",
			},
			payload:      "test",
			inputHandler: nil,
			errorString:  "router.input no input handler found",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {

			registration, ok := processor.GetProcessorRegistration("router.input")
			if !ok {
				t.Fatalf("router.input processor not registered")
			}

			processorInstance, err := registration.New(config.ProcessorConfig{
				Type:   "router.input",
				Params: testCase.params,
			})

			if err != nil {
				if testCase.errorString != err.Error() {
					t.Fatalf("router.input got error '%s', expected '%s'", err.Error(), testCase.errorString)
				}
				return
			}

			got, err := processorInstance.Process(t.Context(), payload.WrappedPayload{Payload: testCase.payload}, nil, testCase.inputHandler)

			if err == nil {
				t.Fatalf("router.input expected to fail but succeeded, got: %v", got)
			}

			if err.Error() != testCase.errorString {
				t.Fatalf("router.input got error '%s', expected '%s'", err.Error(), testCase.errorString)
			}
		})
	}
}
