package processor_test

import (
	"testing"

	"github.com/jwetzell/showbridge-go/config"
	"github.com/jwetzell/showbridge-go/processor"
	"github.com/jwetzell/showbridge-go/routing"
)

func TestGoodFilterRate(t *testing.T) {
	testCases := []struct {
		name    string
		params  map[string]any
		payload any
		match   bool
	}{}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			registration, ok := processor.GetProcessorRegistration("filter.rate")
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

			_, err = processorInstance.Process(t.Context(), routing.WrappedPayload{Payload: testCase.payload}, nil, nil)
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
			payload:     TestStruct{},
			errorString: "filter.rate rate error: not found",
		},
		{
			name: "non-int rate parameter",
			params: map[string]any{
				"rate": "12345",
			},
			payload:     TestStruct{},
			errorString: "filter.rate rate error: not a number",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			registration, ok := processor.GetProcessorRegistration("filter.rate")
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
			got, err := processorInstance.Process(t.Context(), routing.WrappedPayload{Payload: test.payload}, nil, nil)

			if err == nil {
				t.Fatalf("filter.rate expected to fail but succeeded, got: %v", got)

			}
			if err.Error() != test.errorString {
				t.Fatalf("filter.rate got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
