package processor_test

import (
	"testing"

	"github.com/jwetzell/showbridge-go/config"
	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/processor"
)

func TestGoodFilterChange(t *testing.T) {
	tests := []struct {
		name    string
		params  map[string]any
		payload string
		match   bool
	}{
		{
			name:    "basic",
			payload: "hello",
			params:  nil,
			match:   true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			registration, ok := processor.GetProcessorRegistration("filter.change")
			if !ok {
				t.Fatalf("filter.change processor not registered")
			}

			processorInstance, err := registration.New(config.ProcessorConfig{
				Type:   "filter.change",
				Params: test.params,
			})

			if err != nil {
				t.Fatalf("filter.change failed to create processor: %s", err)
			}

			got, err := processorInstance.Process(t.Context(), common.WrappedPayload{Payload: test.payload})

			if err != nil {
				t.Fatalf("filter.change processing failed: %s", err)
			}

			if got.End != !test.match {
				t.Fatalf("filter.change did not filter properly %+v, expected %+v", got, test.match)
			}
		})
	}
}

func TestBadFilterChange(t *testing.T) {
	tests := []struct {
		name        string
		params      map[string]any
		payload     any
		errorString string
	}{}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			registration, ok := processor.GetProcessorRegistration("filter.change")
			if !ok {
				t.Fatalf("filter.change processor not registered")
			}

			processorInstance, err := registration.New(config.ProcessorConfig{
				Type:   "filter.change",
				Params: test.params,
			})

			if err != nil {
				if test.errorString != err.Error() {
					t.Fatalf("filter.change got error '%s', expected '%s'", err.Error(), test.errorString)
				}
				return
			}

			got, err := processorInstance.Process(t.Context(), common.WrappedPayload{Payload: test.payload})

			if err == nil {
				t.Fatalf("filter.change expected to fail but got payload: %+v", got)
			}

			if err.Error() != test.errorString {
				t.Fatalf("filter.change got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
