package processor_test

import (
	"slices"
	"testing"

	"github.com/jwetzell/showbridge-go/internal/processor"
)

func TestGoodStringEncode(t *testing.T) {
	stringEncoder := processor.StringEncode{}
	tests := []struct {
		processor processor.Processor
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
				t.Fatalf("string.encode returned a %T payload: %s", got, got)
			}
			if err != nil {
				t.Fatalf("string.encode failed: %s", err)
			}
			if !slices.Equal(gotBytes, test.expected) {
				t.Fatalf("string.encode got %s, expected %s", got, test.expected)
			}
		})
	}
}

func TestBadStringEncode(t *testing.T) {
	stringEncoder := processor.StringEncode{}
	tests := []struct {
		processor   processor.Processor
		name        string
		payload     any
		errorString string
	}{
		{
			processor:   &stringEncoder,
			name:        "non-string input",
			payload:     []byte{0x68, 0x65, 0x6c, 0x6c, 0x6f},
			errorString: "string.encode processor only accepts a string",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.processor.Process(t.Context(), test.payload)

			if err == nil {
				t.Fatalf("string.encode expected to fail but got payload: %s", got)
			}
			if err.Error() != test.errorString {
				t.Fatalf("string.encode got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
