package processor_test

import (
	"testing"

	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/processor"
)

func TestStringFilterFromRegistry(t *testing.T) {
	registration, ok := processor.ProcessorRegistry["filter.regex"]
	if !ok {
		t.Fatalf("filter.regex processor not registered")
	}

	processorInstance, err := registration.New(config.ProcessorConfig{
		Type: "filter.regex",
		Params: map[string]any{
			"pattern": "hello",
		},
	})
	if err != nil {
		t.Fatalf("failed to create filter.regex processor: %s", err)
	}

	if processorInstance.Type() != "filter.regex" {
		t.Fatalf("filter.regex processor has wrong type: %s", processorInstance.Type())
	}

	payload := "hello"
	expected := "hello"

	got, err := processorInstance.Process(t.Context(), common.GetWrappedPayload(t.Context(), payload))
	if err != nil {
		t.Fatalf("filter.regex processing failed: %s", err)
	}

	gotString, ok := got.Payload.(string)

	if !ok {
		t.Fatalf("filter.regex should return byte slice")
	}

	if gotString != expected {
		t.Fatalf("filter.regex got %+v, expected %+v", got, expected)
	}
}

func TestGoodStringFilter(t *testing.T) {
	tests := []struct {
		name    string
		params  map[string]any
		payload string
		match   bool
	}{
		{
			name:    "matches pattern",
			payload: "hello",
			params:  map[string]any{"pattern": "hello"},
			match:   true,
		},
		{
			name:    "does not match pattern",
			payload: "hello",
			params:  map[string]any{"pattern": "world"},
			match:   false,
		},
		{
			name:    "basic regex",
			payload: "hello world",
			params:  map[string]any{"pattern": ".* world"},
			match:   true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			registration, ok := processor.ProcessorRegistry["filter.regex"]
			if !ok {
				t.Fatalf("filter.regex processor not registered")
			}

			processorInstance, err := registration.New(config.ProcessorConfig{
				Type:   "filter.regex",
				Params: test.params,
			})

			if err != nil {
				t.Fatalf("filter.regex failed to create processor: %s", err)
			}

			got, err := processorInstance.Process(t.Context(), common.GetWrappedPayload(t.Context(), test.payload))

			if err != nil {
				t.Fatalf("filter.regex processing failed: %s", err)
			}

			if got.End != !test.match {
				t.Fatalf("filter.regex did not filter properly %+v, expected %+v", got, test.match)
			}
		})
	}
}

func TestBadStringFilter(t *testing.T) {
	tests := []struct {
		name        string
		params      map[string]any
		payload     any
		errorString string
	}{
		{
			name:        "no pattern param",
			payload:     "hello",
			params:      map[string]any{},
			errorString: "filter.regex pattern error: not found",
		},
		{
			name:    "non-string input",
			payload: []byte{0x68, 0x65, 0x6c, 0x6c, 0x6f},
			params: map[string]any{
				"pattern": "hello",
			},
			errorString: "filter.regex processor only accepts a string",
		},
		{
			name:    "non-string pattern param",
			payload: "hello",
			params: map[string]any{
				"pattern": 123,
			},
			errorString: "filter.regex pattern error: not a string",
		},
		{
			name:    "invalid regex pattern",
			payload: "hello",
			params: map[string]any{
				"pattern": "*invalid",
			},
			errorString: "error parsing regexp: missing argument to repetition operator: `*`",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			registration, ok := processor.ProcessorRegistry["filter.regex"]
			if !ok {
				t.Fatalf("filter.regex processor not registered")
			}

			processorInstance, err := registration.New(config.ProcessorConfig{
				Type:   "filter.regex",
				Params: test.params,
			})

			if err != nil {
				if test.errorString != err.Error() {
					t.Fatalf("filter.regex got error '%s', expected '%s'", err.Error(), test.errorString)
				}
				return
			}

			got, err := processorInstance.Process(t.Context(), common.GetWrappedPayload(t.Context(), test.payload))

			if err == nil {
				t.Fatalf("filter.regex expected to fail but got payload: %+v", got)
			}

			if err.Error() != test.errorString {
				t.Fatalf("filter.regex got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
