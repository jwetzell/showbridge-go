package processor_test

import (
	"reflect"
	"testing"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/processor"
)

func TestStringFilterFromRegistry(t *testing.T) {
	registration, ok := processor.ProcessorRegistry["string.filter"]
	if !ok {
		t.Fatalf("string.filter processor not registered")
	}

	processorInstance, err := registration.New(config.ProcessorConfig{
		Type: "string.filter",
		Params: map[string]any{
			"pattern": "hello",
		},
	})
	if err != nil {
		t.Fatalf("failed to create string.filter processor: %s", err)
	}

	if processorInstance.Type() != "string.filter" {
		t.Fatalf("string.filter processor has wrong type: %s", processorInstance.Type())
	}

	payload := "hello"
	expected := "hello"

	got, err := processorInstance.Process(t.Context(), payload)
	if err != nil {
		t.Fatalf("string.filter processing failed: %s", err)
	}

	gotString, ok := got.(string)

	if !ok {
		t.Fatalf("string.filter should return byte slice")
	}

	if gotString != expected {
		t.Fatalf("string.filter got %+v, expected %+v", got, expected)
	}
}

func TestGoodStringFilter(t *testing.T) {
	tests := []struct {
		name     string
		params   map[string]any
		payload  string
		expected any
	}{
		{
			name:     "matches pattern",
			payload:  "hello",
			params:   map[string]any{"pattern": "hello"},
			expected: "hello",
		},
		{
			name:     "does not match pattern",
			payload:  "hello",
			params:   map[string]any{"pattern": "world"},
			expected: nil,
		},
		{
			name:     "basic regex",
			payload:  "hello world",
			params:   map[string]any{"pattern": ".* world"},
			expected: "hello world",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			registration, ok := processor.ProcessorRegistry["string.filter"]
			if !ok {
				t.Fatalf("string.filter processor not registered")
			}

			processorInstance, err := registration.New(config.ProcessorConfig{
				Type:   "string.filter",
				Params: test.params,
			})

			if err != nil {
				t.Fatalf("string.filter failed to create processor: %s", err)
			}

			got, err := processorInstance.Process(t.Context(), test.payload)

			if err != nil {
				t.Fatalf("string.filter failed: %s", err)
			}

			if test.expected == nil {
				if got != nil {
					t.Fatalf("string.filter got %+v, expected nil", got)
				}
				return
			}

			gotString, ok := got.(string)
			if !ok {
				t.Fatalf("string.filter returned a %T payload: %s", got, got)
			}

			if !reflect.DeepEqual(gotString, test.expected) {
				t.Fatalf("string.filter got %+v, expected %+v", gotString, test.expected)
			}
		})
	}
}

func TestBadStringFilter(t *testing.T) {
	tests := []struct {
		name        string
		payload     any
		params      map[string]any
		errorString string
	}{
		{
			name:        "no pattern param",
			payload:     "hello",
			params:      map[string]any{},
			errorString: "string.filter requires a pattern parameter",
		},
		{
			name:    "non-string input",
			payload: []byte{0x68, 0x65, 0x6c, 0x6c, 0x6f},
			params: map[string]any{
				"pattern": "hello",
			},
			errorString: "string.filter processor only accepts a string",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			registration, ok := processor.ProcessorRegistry["string.filter"]
			if !ok {
				t.Fatalf("string.filter processor not registered")
			}

			processorInstance, err := registration.New(config.ProcessorConfig{
				Type:   "string.filter",
				Params: test.params,
			})

			if err != nil {
				if test.errorString != err.Error() {
					t.Fatalf("string.filter got error '%s', expected '%s'", err.Error(), test.errorString)
				}
				return
			}

			got, err := processorInstance.Process(t.Context(), test.payload)

			if err == nil {
				t.Fatalf("string.filter expected to fail but got payload: %s", got)
			}

			if err.Error() != test.errorString {
				t.Fatalf("string.filter got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
