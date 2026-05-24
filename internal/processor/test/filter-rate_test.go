package processor_test

import (
	"testing"

	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/processor"
	"github.com/jwetzell/showbridge-go/internal/test"
)

func TestFilterRateFromRegistry(t *testing.T) {
	registration, ok := processor.ProcessorRegistry["filter.rate"]
	if !ok {
		t.Fatalf("filter.rate processor not registered")
	}

	processorInstance, err := registration.New(config.ProcessorConfig{
		Type: "filter.rate",
		Params: map[string]any{
			"rate": 1,
		},
	})
	if err != nil {
		t.Fatalf("failed to create filter.rate processor: %s", err)
	}

	if processorInstance.Type() != "filter.rate" {
		t.Fatalf("filter.rate processor has wrong type: %s", processorInstance.Type())
	}
}

func TestGoodFilterRate(t *testing.T) {
	testCases := []struct {
		name    string
		params  map[string]any
		payload any
		match   bool
	}{}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			registration, ok := processor.ProcessorRegistry["filter.rate"]
			if !ok {
				t.Fatalf("filter.rate processor not registered")
			}

			processorInstance, err := registration.New(config.ProcessorConfig{
				Type:   "filter.rate",
				Params: testCase.params,
			})

			if err != nil {
				t.Fatalf("filter.rate failed to create processor: %s", err)
			}

			_, err = processorInstance.Process(t.Context(), common.WrappedPayload{Payload: testCase.payload})
			// TODO(jwetzell): figure out how to test the rate limiting behavior
		})
	}
}

func TestBadFilterRate(t *testing.T) {
	tests := []struct {
		name        string
		params      map[string]any
		payload     any
		errorString string
	}{
		{
			name:   "no rate parameter",
			params: map[string]any{
				// no rate parameter
			},
			payload:     test.TestStruct{},
			errorString: "filter.rate rate error: not found",
		},
		{
			name: "non-int rate parameter",
			params: map[string]any{
				"rate": "12345",
			},
			payload:     test.TestStruct{},
			errorString: "filter.rate rate error: not a number",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			registration, ok := processor.ProcessorRegistry["filter.rate"]
			if !ok {
				t.Fatalf("filter.rate processor not registered")
			}

			processorInstance, err := registration.New(config.ProcessorConfig{
				Type:   "filter.rate",
				Params: test.params,
			})
			if err != nil {
				if err.Error() != test.errorString {
					t.Fatalf("filter.rate got error '%s', expected '%s'", err.Error(), test.errorString)
				}
				return
			}
			got, err := processorInstance.Process(t.Context(), common.WrappedPayload{Payload: test.payload})

			if err == nil {
				t.Fatalf("filter.rate expected to fail but succeeded, got: %v", got)

			}
			if err.Error() != test.errorString {
				t.Fatalf("filter.rate got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
