package processor_test

import (
	"testing"

	"github.com/jwetzell/showbridge-go/internal/processor"
)

func TestGoodIntParse(t *testing.T) {
	intParser := processor.IntParse{}
	tests := []struct {
		processor processor.Processor
		name      string
		payload   any
		expected  int64
	}{
		{
			name:     "positive number",
			payload:  "12345",
			expected: 12345,
		},
		{
			name:     "negative number",
			payload:  "-12345",
			expected: -12345,
		},
		{
			name:     "zero",
			payload:  "0",
			expected: 0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := intParser.Process(t.Context(), test.payload)

			gotInt, ok := got.(int64)
			if !ok {
				t.Fatalf("int.parse returned a %T payload: %s", got, got)
			}
			if err != nil {
				t.Fatalf("int.parse failed: %s", err)
			}
			if gotInt != test.expected {
				t.Fatalf("int.parse got %d, expected %d", gotInt, test.expected)
			}
		})
	}
}

func TestBadIntParse(t *testing.T) {
	intParser := processor.IntParse{}
	tests := []struct {
		processor   processor.Processor
		name        string
		payload     any
		errorString string
	}{
		{
			name:        "non-string input",
			payload:     []byte{0x01},
			errorString: "int.parse processor only accepts a string",
		},
		{
			name:        "not int string",
			payload:     "123.46",
			errorString: "strconv.ParseInt: parsing \"123.46\": invalid syntax",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := intParser.Process(t.Context(), test.payload)

			if err == nil {
				t.Fatalf("int.parse expected to fail but succeeded, got: %v", got)

			}
			if err.Error() != test.errorString {
				t.Fatalf("int.parse got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
