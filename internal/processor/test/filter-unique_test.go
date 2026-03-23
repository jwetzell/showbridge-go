package processor_test

import (
	"reflect"
	"testing"

	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/processor"
)

func TestFilterUniqueFromRegistry(t *testing.T) {
	registration, ok := processor.ProcessorRegistry["filter.unique"]
	if !ok {
		t.Fatalf("filter.unique processor not registered")
	}

	processorInstance, err := registration.New(config.ProcessorConfig{
		Type: "filter.unique",
	})
	if err != nil {
		t.Fatalf("failed to create filter.unique processor: %s", err)
	}

	if processorInstance.Type() != "filter.unique" {
		t.Fatalf("filter.unique processor has wrong type: %s", processorInstance.Type())
	}

	payload := "hello"
	expected := "hello"

	got, err := processorInstance.Process(t.Context(), common.GetWrappedPayload(t.Context(), payload))
	if err != nil {
		t.Fatalf("filter.unique processing failed: %s", err)
	}

	if !reflect.DeepEqual(got.Payload, expected) {
		t.Fatalf("filter.unique got %+v, expected %+v", got.Payload, expected)
	}
}

func TestGoodFilterUnique(t *testing.T) {
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
			registration, ok := processor.ProcessorRegistry["filter.unique"]
			if !ok {
				t.Fatalf("filter.unique processor not registered")
			}

			processorInstance, err := registration.New(config.ProcessorConfig{
				Type:   "filter.unique",
				Params: test.params,
			})

			if err != nil {
				t.Fatalf("filter.unique failed to create processor: %s", err)
			}

			got, err := processorInstance.Process(t.Context(), common.GetWrappedPayload(t.Context(), test.payload))

			if err != nil {
				t.Fatalf("filter.unique processing failed: %s", err)
			}

			if got.End != !test.match {
				t.Fatalf("filter.unique did not filter properly %+v, expected %+v", got, test.match)
			}
		})
	}
}

func TestBadFilterUnique(t *testing.T) {
	tests := []struct {
		name        string
		params      map[string]any
		payload     any
		errorString string
	}{}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			registration, ok := processor.ProcessorRegistry["filter.unique"]
			if !ok {
				t.Fatalf("filter.unique processor not registered")
			}

			processorInstance, err := registration.New(config.ProcessorConfig{
				Type:   "filter.unique",
				Params: test.params,
			})

			if err != nil {
				if test.errorString != err.Error() {
					t.Fatalf("filter.unique got error '%s', expected '%s'", err.Error(), test.errorString)
				}
				return
			}

			got, err := processorInstance.Process(t.Context(), common.GetWrappedPayload(t.Context(), test.payload))

			if err == nil {
				t.Fatalf("filter.unique expected to fail but got payload: %+v", got)
			}

			if err.Error() != test.errorString {
				t.Fatalf("filter.unique got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
