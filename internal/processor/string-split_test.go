package processor_test

import (
	"slices"
	"testing"

	"github.com/jwetzell/showbridge-go/internal/processor"
)

func TestGoodStringSplit(t *testing.T) {
	tests := []struct {
		processor processor.Processor
		name      string
		payload   any
		expected  []string
	}{
		{
			processor: &processor.StringSplit{Separator: ","},
			name:      "comma separated",
			payload:   "part1,part2,part3",
			expected:  []string{"part1", "part2", "part3"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.processor.Process(t.Context(), test.payload)

			gotStrings, ok := got.([]string)
			if !ok {
				t.Errorf("string.split returned a %T payload: %s", got, got)
			}
			if err != nil {
				t.Errorf("string.split failed: %s", err)
			}
			if !slices.Equal(gotStrings, test.expected) {
				t.Errorf("string.split got %s, expected %s", got, test.expected)
			}
		})
	}
}

func TestBasStringSplit(t *testing.T) {
	tests := []struct {
		processor   processor.Processor
		name        string
		payload     any
		errorString string
	}{
		{
			processor:   &processor.StringSplit{Separator: ","},
			name:        "hello",
			payload:     []byte{0x68, 0x65, 0x6c, 0x6c, 0x6f},
			errorString: "string.split only accepts a string",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.processor.Process(t.Context(), test.payload)

			if err == nil {
				t.Errorf("string.split expected error but got none, payload: %s", got)
			}
			if err.Error() != test.errorString {
				t.Errorf("string.split got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
