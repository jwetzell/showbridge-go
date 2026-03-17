package processor_test

import (
	"reflect"
	"testing"

	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/processor"
)

func TestDebugLogFromRegistry(t *testing.T) {
	registration, ok := processor.ProcessorRegistry["debug.log"]
	if !ok {
		t.Fatalf("debug.log processor not registered")
	}

	processorInstance, err := registration.New(config.ProcessorConfig{
		Type: "debug.log",
	})

	if err != nil {
		t.Fatalf("failed to create debug.log processor: %s", err)
	}

	if processorInstance.Type() != "debug.log" {
		t.Fatalf("debug.log processor has wrong type: %s", processorInstance.Type())
	}

	payload := "test"
	expected := "test"

	got, err := processorInstance.Process(t.Context(), common.GetWrappedPayload(t.Context(), payload))
	if err != nil {
		t.Fatalf("debug.log processing failed: %s", err)
	}

	if got.Payload != expected {
		t.Fatalf("debug.log got %+v, expected %+v", got, expected)
	}
}

func TestGoodDebugLog(t *testing.T) {

	tests := []struct {
		name     string
		params   map[string]any
		payload  any
		expected any
	}{}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			registration, ok := processor.ProcessorRegistry["debug.log"]
			if !ok {
				t.Fatalf("debug.log processor not registered")
			}

			processorInstance, err := registration.New(config.ProcessorConfig{
				Type:   "debug.log",
				Params: test.params,
			})

			if err != nil {
				t.Fatalf("debug.log failed to create processor: %s", err)
			}

			got, err := processorInstance.Process(t.Context(), common.GetWrappedPayload(t.Context(), test.payload))
			if err != nil {
				t.Fatalf("debug.log processing failed: %s", err)
			}

			if !reflect.DeepEqual(got.Payload, test.expected) {
				t.Fatalf("debug.log got %+v (%T), expected %+v (%T)", got.Payload, got.Payload, test.expected, test.expected)
			}
		})
	}
}

func TestBadDebugLog(t *testing.T) {
	tests := []struct {
		name        string
		params      map[string]any
		payload     any
		errorString string
	}{}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			registration, ok := processor.ProcessorRegistry["debug.log"]
			if !ok {
				t.Fatalf("debug.log processor not registered")
			}

			processorInstance, err := registration.New(config.ProcessorConfig{
				Type:   "debug.log",
				Params: test.params,
			})

			if err != nil {
				if test.errorString != err.Error() {
					t.Fatalf("debug.log got error '%s', expected '%s'", err.Error(), test.errorString)
				}
				return
			}

			got, err := processorInstance.Process(t.Context(), common.GetWrappedPayload(t.Context(), test.payload))

			if err == nil {
				t.Fatalf("debug.log expected to fail but succeeded, got: %v", got)
			}

			if err.Error() != test.errorString {
				t.Fatalf("debug.log got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
