package processor_test

import (
	"testing"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/processor"
)

func TestStringDecodeFromRegistry(t *testing.T) {
	registration, ok := processor.ProcessorRegistry["string.decode"]
	if !ok {
		t.Fatalf("string.decode processor not registered")
	}

	processorInstance, err := registration.New(config.ProcessorConfig{
		Type: "string.decode",
	})
	if err != nil {
		t.Fatalf("failed to create string.decode processor: %s", err)
	}

	if processorInstance.Type() != "string.decode" {
		t.Fatalf("string.decode processor has wrong type: %s", processorInstance.Type())
	}

	payload := []byte{'h', 'e', 'l', 'l', 'o'}
	expected := "hello"

	got, err := processorInstance.Process(t.Context(), payload)
	if err != nil {
		t.Fatalf("string.decode processing failed: %s", err)
	}

	if got != expected {
		t.Fatalf("string.decode got %+v, expected %+v", got, expected)
	}
}

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
			got, err := stringDecoder.Process(t.Context(), test.payload)

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
		name        string
		payload     any
		errorString string
	}{
		{
			name:        "non-[]byte input",
			payload:     "hello",
			errorString: "string.decode processor only accepts a []byte",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := stringDecoder.Process(t.Context(), test.payload)

			if err == nil {
				t.Fatalf("string.decode expected to fail but got payload: %s", got)
			}
			if err.Error() != test.errorString {
				t.Fatalf("string.decode got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
