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
				t.Errorf("int.parse returned a %T payload: %s", got, got)
			}
			if err != nil {
				t.Errorf("int.parse failed: %s", err)
			}
			if gotInt != test.expected {
				t.Errorf("int.parse got %d, expected %d", gotInt, test.expected)
			}
		})
	}
}
