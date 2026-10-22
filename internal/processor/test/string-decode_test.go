package processor_test

import (
	"testing"

	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/processor"
)

func TestGoodStringDecode(t *testing.T) {
	stringDecoder := processor.StringDecode{}
	tests := []struct {
		name     string
		payload  any
		expected string
	}{
		{
			name:     "basic string",
			payload:  []byte{0x68, 0x65, 0x6c, 0x6c, 0x6f},
			expected: "hello",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := stringDecoder.Process(t.Context(), common.WrappedPayload{Payload: test.payload})
			if err != nil {
				t.Fatalf("string.decode processing failed: %s", err)
			}

			gotString, ok := got.Payload.(string)
			if !ok {
				t.Fatalf("string.decode returned a %T payload: %+v", got, got)
			}

			if gotString != test.expected {
				t.Fatalf("string.decode got %+v, expected %s", got, test.expected)
			}
		})
	}
}

func TestBadStringDecode(t *testing.T) {
	stringDecoder := processor.StringDecode{}
	tests := []struct {
		name        string
		payload     any
		errorString string
	}{
		{
			name:        "non-[]byte input",
			payload:     123,
			errorString: "string.decode processor only accepts a []byte",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := stringDecoder.Process(t.Context(), common.WrappedPayload{Payload: test.payload})

			if err == nil {
				t.Fatalf("string.decode expected to fail but got payload: %+v", got)
			}
			if err.Error() != test.errorString {
				t.Fatalf("string.decode got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
