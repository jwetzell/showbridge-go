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

func TestRouterInputFromRegistry(t *testing.T) {
	registration, ok := processor.ProcessorRegistry["router.input"]
	if !ok {
		t.Fatalf("router.input processor not registered")
	}

	processorInstance, err := registration.New(config.ProcessorConfig{
		Type: "router.input",
		Params: config.Params{
			"source": "test",
		},
	})

	if err != nil {
		t.Fatalf("failed to create router.input processor: %s", err)
	}

	if processorInstance.Type() != "router.input" {
		t.Fatalf("router.input processor has wrong type: %s", processorInstance.Type())
	}

	payload := "test"
	expected := "test"

	got, err := processorInstance.Process(test.GetContextWithRouter(t.Context()), common.GetWrappedPayload(t.Context(), payload))
	if err != nil {
		t.Fatalf("router.input processing failed: %s", err)
	}

	if got.Payload != expected {
		t.Fatalf("router.input got %+v, expected %+v", got, expected)
	}
}

func TestGoodRouterInput(t *testing.T) {

	testCases := []struct {
		name     string
		params   map[string]any
		payload  any
		expected any
	}{}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {

			registration, ok := processor.ProcessorRegistry["router.input"]
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

			got, err := processorInstance.Process(t.Context(), common.GetWrappedPayload(test.GetContextWithRouter(t.Context()), testCase.payload))
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
		name              string
		params            map[string]any
		payload           any
		processCtx        context.Context
		wrappedPayloadCtx context.Context
		errorString       string
	}{
		{
			name:              "no source param",
			params:            map[string]any{},
			payload:           "test",
			processCtx:        test.GetContextWithRouter(t.Context()),
			wrappedPayloadCtx: t.Context(),
			errorString:       "router.input source error: not found",
		},
		{
			name: "non-string source",
			params: map[string]any{
				"source": 123,
			},
			payload:           "test",
			processCtx:        test.GetContextWithRouter(t.Context()),
			wrappedPayloadCtx: t.Context(),
			errorString:       "router.input source error: not a string",
		},
		{
			name: "router not found in context",
			params: map[string]any{
				"source": "test",
			},
			payload:           "test",
			processCtx:        t.Context(),
			wrappedPayloadCtx: t.Context(),
			errorString:       "router.input no router found",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {

			registration, ok := processor.ProcessorRegistry["router.input"]
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

			got, err := processorInstance.Process(testCase.processCtx, common.GetWrappedPayload(testCase.wrappedPayloadCtx, testCase.payload))

			if err == nil {
				t.Fatalf("router.input expected to fail but succeeded, got: %v", got)
			}

			if err.Error() != testCase.errorString {
				t.Fatalf("router.input got error '%s', expected '%s'", err.Error(), testCase.errorString)
			}
		})
	}
}
