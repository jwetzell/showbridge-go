package processor_test

import (
	"slices"
	"testing"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/processor"
)

func TestStringSplitFromRegistry(t *testing.T) {
	registration, ok := processor.ProcessorRegistry["string.split"]
	if !ok {
		t.Fatalf("string.split processor not registered")
	}

	processorInstance, err := registration.New(config.ProcessorConfig{
		Type: "string.split",
		Params: map[string]any{
			"separator": ",",
		},
	})
	if err != nil {
		t.Fatalf("failed to create string.split processor: %s", err)
	}

	if processorInstance.Type() != "string.split" {
		t.Fatalf("string.split processor has wrong type: %s", processorInstance.Type())
	}

	payload := "part1,part2,part3"
	expected := []string{"part1", "part2", "part3"}

	got, err := processorInstance.Process(t.Context(), payload)
	if err != nil {
		t.Fatalf("string.split processing failed: %s", err)
	}

	gotStrings, ok := got.([]string)

	if !ok {
		t.Fatalf("string.split should return a slice of strings")
	}

	if !slices.Equal(gotStrings, expected) {
		t.Fatalf("string.split got %+v, expected %+v", got, expected)
	}
}

func TestGoodStringSplit(t *testing.T) {
	tests := []struct {
		name     string
		params   map[string]any
		payload  any
		expected []string
	}{
		{
			name:     "comma separated",
			params:   map[string]any{"separator": ","},
			payload:  "part1,part2,part3",
			expected: []string{"part1", "part2", "part3"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			registration, ok := processor.ProcessorRegistry["string.split"]
			if !ok {
				t.Fatalf("string.split processor not registered")
			}

			processorInstance, err := registration.New(config.ProcessorConfig{
				Type:   "string.split",
				Params: test.params,
			})

			if err != nil {
				t.Fatalf("string.split failed to create processor: %s", err)
			}

			got, err := processorInstance.Process(t.Context(), test.payload)

			gotStrings, ok := got.([]string)
			if !ok {
				t.Fatalf("string.split returned a %T payload: %s", got, got)
			}
			if err != nil {
				t.Fatalf("string.split failed: %s", err)
			}
			if !slices.Equal(gotStrings, test.expected) {
				t.Fatalf("string.split got %s, expected %s", got, test.expected)
			}
		})
	}
}

func TestBadStringSplit(t *testing.T) {
	tests := []struct {
		name        string
		payload     any
		params      map[string]any
		errorString string
	}{
		{
			name:        "non-string input",
			payload:     []byte{0x68, 0x65, 0x6c, 0x6c, 0x6f},
			params:      map[string]any{"separator": ","},
			errorString: "string.split only accepts a string",
		},
		{
			name:        "missing separator param",
			payload:     "part1,part2,part3",
			params:      map[string]any{},
			errorString: "string.split requires a separator",
		},
		{
			name:        "non-string separator param",
			payload:     "part1,part2,part3",
			params:      map[string]any{"separator": 123},
			errorString: "string.split separator must be a string",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			registration, ok := processor.ProcessorRegistry["string.split"]
			if !ok {
				t.Fatalf("string.split processor not registered")
			}

			processorInstance, err := registration.New(config.ProcessorConfig{
				Type:   "string.split",
				Params: test.params,
			})

			if err != nil {
				if test.errorString != err.Error() {
					t.Fatalf("string.split got error '%s', expected '%s'", err.Error(), test.errorString)
				}
				return
			}

			got, err := processorInstance.Process(t.Context(), test.payload)

			if err == nil {
				t.Fatalf("string.split expected error but got none, payload: %s", got)
			}
			if err.Error() != test.errorString {
				t.Fatalf("string.split got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
