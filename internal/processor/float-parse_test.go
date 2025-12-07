package processor_test

import (
	"testing"

	"github.com/jwetzell/showbridge-go/internal/processor"
)

func TestGoodFloatParse(t *testing.T) {
	floatParser := processor.FloatParse{}
	tests := []struct {
		processor processor.Processor
		name      string
		payload   any
		expected  float64
	}{
		{
			name:     "positive number",
			payload:  "12345.67",
			expected: 12345.67,
		},
		{
			name:     "negative number",
			payload:  "-12345.67",
			expected: -12345.67,
		},
		{
			name:     "zero",
			payload:  "0",
			expected: 0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := floatParser.Process(t.Context(), test.payload)

			gotFloat, ok := got.(float64)
			if !ok {
				t.Errorf("float.parse returned a %T payload: %s", got, got)
			}
			if err != nil {
				t.Errorf("float.parse failed: %s", err)
			}
			if gotFloat != test.expected {
				t.Errorf("float.parse got %f, expected %f", gotFloat, test.expected)
			}
		})
	}
}

func TestBadFloatParse(t *testing.T) {
	floatParser := processor.FloatParse{}
	tests := []struct {
		processor   processor.Processor
		name        string
		payload     any
		errorString string
	}{
		{
			name:        "non-string input",
			payload:     []byte{0x01},
			errorString: "float.parse processor only accepts a string",
		},
		{
			name:        "not float string",
			payload:     "abcd",
			errorString: "strconv.ParseFloat: parsing \"abcd\": invalid syntax",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := floatParser.Process(t.Context(), test.payload)

			if err == nil {
				t.Errorf("float.parse expected to fail but succeeded, got: %v", got)

			}
			if err.Error() != test.errorString {
				t.Errorf("float.parse got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
