package processor_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/jwetzell/showbridge-go/config"
	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/processor"
)

func TestYamlDecodeFromRegistry(t *testing.T) {
	registration, ok := processor.GetProcessorRegistration("yaml.decode")
	if !ok {
		t.Fatalf("yaml.decode processor not registered")
	}

	processorInstance, err := registration.New(config.ProcessorConfig{
		Id:   "test-id",
		Type: "yaml.decode",
	})
	if err != nil {
		t.Fatalf("failed to create yaml.decode processor: %s", err)
	}

	if processorInstance.Id() != "test-id" {
		t.Fatalf("yaml.decode processor has wrong id: %s", processorInstance.Id())
	}

	if processorInstance.Type() != "yaml.decode" {
		t.Fatalf("yaml.decode processor has wrong type: %s", processorInstance.Type())
	}

	payload := "{\"property\":\"hello\"}"

	expected := map[string]any{
		"property": "hello",
	}

	got, err := processorInstance.Process(t.Context(), common.WrappedPayload{Payload: payload})
	if err != nil {
		t.Fatalf("yaml.decode processing failed: %s", err)
	}

	gotMap, ok := got.Payload.(map[string]any)

	if !ok {
		t.Fatalf("yaml.decode should return byte slice")
	}

	if !reflect.DeepEqual(gotMap, expected) {
		t.Fatalf("yaml.decode got %+v, expected %+v", got, expected)
	}
}

func TestGoodYamlDecode(t *testing.T) {
	jsonDecoder := processor.YamlDecode{}
	tests := []struct {
		name     string
		payload  string
		expected map[string]any
	}{
		{
			name:    "basic json",
			payload: "address: /hello\nargs: null\n",
			expected: map[string]any{
				"address": "/hello",
				"args":    nil,
			},
		},
		{
			name:    "array",
			payload: "address: /hello\nargs:\n- 1\n- 2\n- 3\n",
			expected: map[string]any{
				"address": "/hello",
				"args":    []any{1.0, 2.0, 3.0},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := jsonDecoder.Process(t.Context(), common.WrappedPayload{Payload: test.payload})
			if err != nil {
				t.Fatalf("yaml.decode processing failed: %s", err)
			}

			gotMap, ok := got.Payload.(map[string]any)
			if !ok {
				t.Fatalf("yaml.decode returned a %T payload: %+v", got, got)
			}

			if !reflect.DeepEqual(gotMap, test.expected) {
				t.Fatalf("yaml.decode got %+v, expected %s", got, test.expected)
			}
		})
	}
}

func TestBadYamlDecode(t *testing.T) {
	stringEncoder := processor.YamlDecode{}
	tests := []struct {
		name        string
		payload     any
		errorString string
	}{
		{
			name:        "non-string or byte input",
			payload:     123,
			errorString: "yaml.decode can only process a string or []byte",
		},
		{
			name:        "invalid json",
			payload:     "address:property",
			errorString: "error unmarshaling JSON: while decoding JSON: json: cannot unmarshal string into Go value of type map[string]interface {}",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := stringEncoder.Process(t.Context(), common.WrappedPayload{Payload: test.payload})

			if err == nil {
				t.Fatalf("yaml.decode expected to fail but got payload: %+v", got)
			}
			if err.Error() != test.errorString {
				t.Fatalf("yaml.decode got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}

func BenchmarkYamlDecode(b *testing.B) {
	registration, ok := processor.GetProcessorRegistration("yaml.decode")
	if !ok {
		b.Fatalf("yaml.decode processor not registered")
	}

	processorInstance, err := registration.New(config.ProcessorConfig{
		Type: "yaml.decode",
	})

	if err != nil {
		b.Fatalf("yaml.decode failed to create processor: %s", err)
	}

	count := 0
	for b.Loop() {
		_, err := processorInstance.Process(b.Context(), common.WrappedPayload{Payload: fmt.Appendf(nil, "{\"key\":%d}", count)})
		if err != nil {
			b.Fatalf("yaml.decode processing failed: %s", err)
		}
		count++
	}
}
