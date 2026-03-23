package processor_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/processor"
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

	got, err := processorInstance.Process(GetContextWithRouter(t.Context()), common.GetWrappedPayload(t.Context(), payload))
	if err != nil {
		t.Fatalf("router.input processing failed: %s", err)
	}

	if got.Payload != expected {
		t.Fatalf("router.input got %+v, expected %+v", got, expected)
	}
}

func TestGoodRouterInput(t *testing.T) {

	tests := []struct {
		name     string
		params   map[string]any
		payload  any
		expected any
	}{}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			registration, ok := processor.ProcessorRegistry["router.input"]
			if !ok {
				t.Fatalf("router.input processor not registered")
			}

			processorInstance, err := registration.New(config.ProcessorConfig{
				Type:   "router.input",
				Params: test.params,
			})

			if err != nil {
				t.Fatalf("router.input failed to create processor: %s", err)
			}

			got, err := processorInstance.Process(t.Context(), common.GetWrappedPayload(GetContextWithRouter(t.Context()), test.payload))
			if err != nil {
				t.Fatalf("router.input processing failed: %s", err)
			}

			if !reflect.DeepEqual(got.Payload, test.expected) {
				t.Fatalf("router.input got %+v (%T), expected %+v (%T)", got.Payload, got.Payload, test.expected, test.expected)
			}
		})
	}
}

func TestBadRouterInput(t *testing.T) {
	tests := []struct {
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
			processCtx:        GetContextWithRouter(t.Context()),
			wrappedPayloadCtx: t.Context(),
			errorString:       "router.input source error: not found",
		},
		{
			name: "non-string source",
			params: map[string]any{
				"source": 123,
			},
			payload:           "test",
			processCtx:        GetContextWithRouter(t.Context()),
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

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			registration, ok := processor.ProcessorRegistry["router.input"]
			if !ok {
				t.Fatalf("router.input processor not registered")
			}

			processorInstance, err := registration.New(config.ProcessorConfig{
				Type:   "router.input",
				Params: test.params,
			})

			if err != nil {
				if test.errorString != err.Error() {
					t.Fatalf("router.input got error '%s', expected '%s'", err.Error(), test.errorString)
				}
				return
			}

			got, err := processorInstance.Process(test.processCtx, common.GetWrappedPayload(test.wrappedPayloadCtx, test.payload))

			if err == nil {
				t.Fatalf("router.input expected to fail but succeeded, got: %v", got)
			}

			if err.Error() != test.errorString {
				t.Fatalf("router.input got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
