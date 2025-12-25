package processor_test

import (
	"testing"

	"github.com/jwetzell/showbridge-go/internal/processor"
)

func TestGoodUintParse(t *testing.T) {
	uintParser := processor.UintParse{}
	tests := []struct {
		processor processor.Processor
		name      string
		payload   any
		expected  uint64
	}{
		{
			name:     "positive number",
			payload:  "12345",
			expected: 12345,
		},
		{
			name:     "zero",
			payload:  "0",
			expected: 0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := uintParser.Process(t.Context(), test.payload)

			gotUint, ok := got.(uint64)
			if !ok {
				t.Fatalf("uint.parse returned a %T payload: %s", got, got)
			}
			if err != nil {
				t.Fatalf("uint.parse failed: %s", err)
			}
			if gotUint != test.expected {
				t.Fatalf("uint.parse got %d, expected %d", gotUint, test.expected)
			}
		})
	}
}

func TestBadUintParse(t *testing.T) {
	uintParser := processor.UintParse{}
	tests := []struct {
		processor   processor.Processor
		name        string
		payload     any
		errorString string
	}{
		{
			name:        "non-string input",
			payload:     []byte{0x01},
			errorString: "uint.parse processor only accepts a string",
		},
		{
			name:        "not uint string",
			payload:     "-1234",
			errorString: "strconv.ParseUint: parsing \"-1234\": invalid syntax",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := uintParser.Process(t.Context(), test.payload)

			if err == nil {
				t.Fatalf("uint.parse expected to fail but succeeded, got: %v", got)

			}
			if err.Error() != test.errorString {
				t.Fatalf("uint.parse got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
