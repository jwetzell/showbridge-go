package processing_test

import (
	"testing"

	"github.com/jwetzell/showbridge-go/internal/processing"
)

func TestGoodStringDecode(t *testing.T) {
	stringDecoder := processing.StringDecode{}
	tests := []struct {
		processor processing.Processor
		name      string
		payload   any
		expected  string
	}{
		{
			processor: &stringDecoder,
			name:      "hello",
			payload:   []byte{0x68, 0x65, 0x6c, 0x6c, 0x6f},
			expected:  "hello",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.processor.Process(t.Context(), test.payload)

			gotString, ok := got.(string)
			if !ok {
				t.Errorf("string.decode returned a %T payload: %s", got, got)
			}
			if err != nil {
				t.Errorf("string.decode failed: %s", err)
			}
			if gotString != test.expected {
				t.Errorf("string.decode got %s, expected %s", got, test.expected)
			}
		})
	}
}
