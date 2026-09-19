package processor_test

import (
	"fmt"
	"slices"
	"testing"

	"github.com/jwetzell/osc-go"
	"github.com/jwetzell/showbridge-go/config"
	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/processor"
)

func TestYamlEncodeFromRegistry(t *testing.T) {
	registration, ok := processor.GetProcessorRegistration("yaml.encode")
	if !ok {
		t.Fatalf("yaml.encode processor not registered")
	}

	processorInstance, err := registration.New(config.ProcessorConfig{
		Id:   "test-id",
		Type: "yaml.encode",
	})
	if err != nil {
		t.Fatalf("failed to create yaml.encode processor: %s", err)
	}

	if processorInstance.Id() != "test-id" {
		t.Fatalf("yaml.encode processor has wrong id: %s", processorInstance.Id())
	}

	if processorInstance.Type() != "yaml.encode" {
		t.Fatalf("yaml.encode processor has wrong type: %s", processorInstance.Type())
	}

	payload := struct {
		Property string `json:"property"`
	}{
		Property: "hello",
	}

	expected := []byte("property: hello\n")

	got, err := processorInstance.Process(t.Context(), common.WrappedPayload{Payload: payload})
	if err != nil {
		t.Fatalf("yaml.encode processing failed: %s", err)
	}

	gotBytes, ok := got.Payload.([]byte)

	if !ok {
		t.Fatalf("yaml.encode should return byte slice got %T: %+v", got.Payload, got.Payload)
	}

	if !slices.Equal(gotBytes, expected) {
		t.Fatalf("yaml.encode got %+v, expected %+v", got, expected)
	}
}

func TestGoodYamlEncode(t *testing.T) {
	jsonEncoder := processor.YamlEncode{}
	tests := []struct {
		name     string
		payload  any
		expected []byte
	}{
		{
			name: "basic struct",
			payload: osc.Message{
				Address: "/hello",
			},
			expected: []byte("address: /hello\nargs: null\n"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := jsonEncoder.Process(t.Context(), common.WrappedPayload{Payload: test.payload})
			if err != nil {
				t.Fatalf("yaml.encode processing failed: %s", err)
			}

			gotBytes, ok := got.Payload.([]byte)
			if !ok {
				t.Fatalf("yaml.encode returned a %T payload: %+v", got.Payload, got.Payload)
			}

			if !slices.Equal(gotBytes, test.expected) {
				t.Fatalf("yaml.encode got %+v, expected %s", got, test.expected)
			}
		})
	}
}

func TestBadYamlEncode(t *testing.T) {
	stringEncoder := processor.YamlEncode{}
	tests := []struct {
		name        string
		payload     any
		errorString string
	}{
		{
			name:        "unencodable type",
			payload:     make(chan int),
			errorString: "error marshaling into JSON: json: unsupported type: chan int",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := stringEncoder.Process(t.Context(), common.WrappedPayload{Payload: test.payload})

			if err == nil {
				t.Fatalf("yaml.encode expected to fail but got payload: %+v", got)
			}
			if err.Error() != test.errorString {
				t.Fatalf("yaml.encode got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}

func BenchmarkYamlEncode(b *testing.B) {
	registration, ok := processor.GetProcessorRegistration("yaml.encode")
	if !ok {
		b.Fatalf("yaml.encode processor not registered")
	}

	processorInstance, err := registration.New(config.ProcessorConfig{
		Type: "yaml.encode",
	})

	if err != nil {
		b.Fatalf("yaml.encode failed to create processor: %s", err)
	}

	count := 0
	for b.Loop() {
		_, err := processorInstance.Process(b.Context(), common.WrappedPayload{Payload: fmt.Sprintf("{\"key\":%d}", count)})
		if err != nil {
			b.Fatalf("yaml.encode processing failed: %s", err)
		}
		count++
	}
}
