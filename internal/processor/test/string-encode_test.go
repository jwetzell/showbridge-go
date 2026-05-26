package processor_test

import (
	"fmt"
	"slices"
	"testing"

	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/processor"
)

func TestStringEncodeFromRegistry(t *testing.T) {
	registration, ok := processor.GetProcessorRegistration("string.encode")
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

	got, err := processorInstance.Process(t.Context(), common.WrappedPayload{Payload: payload})
	if err != nil {
		t.Fatalf("string.encode processing failed: %s", err)
	}

	gotBytes, ok := got.Payload.([]byte)

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
			name:     "basic string",
			payload:  "hello",
			expected: []byte{0x68, 0x65, 0x6c, 0x6c, 0x6f},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := stringEncoder.Process(t.Context(), common.WrappedPayload{Payload: test.payload})
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
			got, err := stringEncoder.Process(t.Context(), common.WrappedPayload{Payload: test.payload})

			if err == nil {
				t.Fatalf("string.encode expected to fail but got payload: %+v", got)
			}
			if err.Error() != test.errorString {
				t.Fatalf("string.encode got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}

func BenchmarkStringEncode(b *testing.B) {
	registration, ok := processor.GetProcessorRegistration("string.encode")
	if !ok {
		b.Fatalf("string.encode processor not registered")
	}

	processorInstance, err := registration.New(config.ProcessorConfig{
		Type: "string.encode",
	})

	if err != nil {
		b.Fatalf("string.encode failed to create processor: %s", err)
	}

	count := 0
	for b.Loop() {
		_, err := processorInstance.Process(b.Context(), common.WrappedPayload{Payload: fmt.Sprintf("%d", count)})
		if err != nil {
			b.Fatalf("string.encode processing failed: %s", err)
		}
		count++
	}
}
