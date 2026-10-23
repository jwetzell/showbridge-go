package processor_test

import (
	"slices"
	"testing"

	"github.com/jwetzell/showbridge-go/common"
	"github.com/jwetzell/showbridge-go/processor"
)

func TestGoodStringEncode(t *testing.T) {
	stringEncoder := processor.StringEncode{}
	tests := []struct {
		name     string
		payload  any
		expected []byte
	}{
		{
			name:     "basic string",
			payload:  "hello",
			expected: []byte{0x68, 0x65, 0x6c, 0x6c, 0x6f},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := stringEncoder.Process(t.Context(), common.WrappedPayload{Payload: test.payload}, nil, nil)
			if err != nil {
				t.Fatalf("string.encode processing failed: %s", err)
			}

			gotBytes, ok := got.Payload.([]byte)
			if !ok {
				t.Fatalf("string.encode returned a %T payload: %+v", got, got)
			}
			if !slices.Equal(gotBytes, test.expected) {
				t.Fatalf("string.encode got %+v, expected %+v", gotBytes, test.expected)
			}
		})
	}
}

func TestBadStringEncode(t *testing.T) {
	stringEncoder := processor.StringEncode{}
	tests := []struct {
		name        string
		payload     any
		errorString string
	}{
		{
			name:        "non-string input",
			payload:     []byte{0x68, 0x65, 0x6c, 0x6c, 0x6f},
			errorString: "string.encode processor only accepts a string",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := stringEncoder.Process(t.Context(), common.WrappedPayload{Payload: test.payload}, nil, nil)

			if err == nil {
				t.Fatalf("string.encode expected to fail but got payload: %+v", got)
			}
			if err.Error() != test.errorString {
				t.Fatalf("string.encode got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
