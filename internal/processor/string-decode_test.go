package processor_test

import (
	"testing"

	"github.com/jwetzell/showbridge-go/internal/processor"
)

func TestGoodStringDecode(t *testing.T) {
	stringDecoder := processor.StringDecode{}
	tests := []struct {
		processor processor.Processor
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
				t.Fatalf("string.decode returned a %T payload: %s", got, got)
			}
			if err != nil {
				t.Fatalf("string.decode failed: %s", err)
			}
			if gotString != test.expected {
				t.Fatalf("string.decode got %s, expected %s", got, test.expected)
			}
		})
	}
}

func TestBadStringDecode(t *testing.T) {
	stringDecoder := processor.StringDecode{}
	tests := []struct {
		processor   processor.Processor
		name        string
		payload     any
		errorString string
	}{
		{
			processor:   &stringDecoder,
			name:        "non-[]byte input",
			payload:     "hello",
			errorString: "string.decode processor only accepts a []byte",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.processor.Process(t.Context(), test.payload)

			if err == nil {
				t.Fatalf("string.decode expected to fail but got payload: %s", got)
			}
			if err.Error() != test.errorString {
				t.Fatalf("string.decode got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
