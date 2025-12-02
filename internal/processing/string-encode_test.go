package processing_test

import (
	"slices"
	"testing"

	"github.com/jwetzell/showbridge-go/internal/processing"
)

func TestGoodStringEncode(t *testing.T) {
	stringEncoder := processing.StringEncode{}
	tests := []struct {
		processor processing.Processor
		name      string
		payload   any
		expected  []byte
	}{
		{
			processor: &stringEncoder,
			name:      "hello",
			payload:   "hello",
			expected:  []byte{0x68, 0x65, 0x6c, 0x6c, 0x6f},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.processor.Process(t.Context(), test.payload)

			gotBytes, ok := got.([]byte)
			if !ok {
				t.Errorf("string.encode returned a %T payload: %s", got, got)
			}
			if err != nil {
				t.Errorf("string.encode failed: %s", err)
			}
			if !slices.Equal(gotBytes, test.expected) {
				t.Errorf("string.encode got %s, expected %s", got, test.expected)
			}
		})
	}
}
