package processor_test

import (
	"slices"
	"testing"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/processor"
)

func TestStringEncodeFromRegistry(t *testing.T) {
	registration, ok := processor.ProcessorRegistry["string.encode"]
	if !ok {
		t.Fatalf("string.encode processor not registered")
	}

	processorInstance, err := registration.New(config.ProcessorConfig{
		Type: "string.encode",
	})
	if err != nil {
		t.Fatalf("failed to create string.encode processor: %s", err)
	}

	if processorInstance.Type() != "string.encode" {
		t.Fatalf("string.encode processor has wrong type: %s", processorInstance.Type())
	}

	payload := "hello"
	expected := []byte{'h', 'e', 'l', 'l', 'o'}

	got, err := processorInstance.Process(t.Context(), payload)
	if err != nil {
		t.Fatalf("string.encode processing failed: %s", err)
	}

	gotBytes, ok := got.([]byte)

	if !ok {
		t.Fatalf("string.encode should return byte slice")
	}

	if !slices.Equal(gotBytes, expected) {
		t.Fatalf("string.encode got %+v, expected %+v", got, expected)
	}
}

func TestGoodStringEncode(t *testing.T) {
	stringEncoder := processor.StringEncode{}
	tests := []struct {
		name     string
		payload  any
		expected []byte
	}{
		{
			name:     "hello",
			payload:  "hello",
			expected: []byte{0x68, 0x65, 0x6c, 0x6c, 0x6f},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := stringEncoder.Process(t.Context(), test.payload)

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
			got, err := stringEncoder.Process(t.Context(), test.payload)

			if err == nil {
				t.Fatalf("string.encode expected to fail but got payload: %s", got)
			}
			if err.Error() != test.errorString {
				t.Fatalf("string.encode got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
