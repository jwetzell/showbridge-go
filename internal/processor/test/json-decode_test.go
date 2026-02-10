package processor_test

import (
	"reflect"
	"testing"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/processor"
)

func TestJsonDecodeFromRegistry(t *testing.T) {
	registration, ok := processor.ProcessorRegistry["json.decode"]
	if !ok {
		t.Fatalf("json.decode processor not registered")
	}

	processorInstance, err := registration.New(config.ProcessorConfig{
		Type: "json.decode",
	})
	if err != nil {
		t.Fatalf("failed to create json.decode processor: %s", err)
	}

	if processorInstance.Type() != "json.decode" {
		t.Fatalf("json.decode processor has wrong type: %s", processorInstance.Type())
	}

	payload := "{\"property\":\"hello\"}"

	expected := map[string]any{
		"property": "hello",
	}

	got, err := processorInstance.Process(t.Context(), payload)
	if err != nil {
		t.Fatalf("json.decode processing failed: %s", err)
	}

	gotMap, ok := got.(map[string]any)

	if !ok {
		t.Fatalf("json.decode should return byte slice")
	}

	if !reflect.DeepEqual(gotMap, expected) {
		t.Fatalf("json.decode got %+v, expected %+v", got, expected)
	}
}

func TestGoodJsonDecode(t *testing.T) {
	jsonDecoder := processor.JsonDecode{}
	tests := []struct {
		name     string
		payload  string
		expected map[string]any
	}{
		{
			name:    "basic json",
			payload: "{\"address\":\"/hello\",\"args\":null}",
			expected: map[string]any{
				"address": "/hello",
				"args":    nil,
			},
		},
		{
			name:    "array",
			payload: "{\"address\":\"/hello\",\"args\":[1,2,3]}",
			expected: map[string]any{
				"address": "/hello",
				"args":    []any{1.0, 2.0, 3.0},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := jsonDecoder.Process(t.Context(), test.payload)

			gotMap, ok := got.(map[string]any)
			if !ok {
				t.Fatalf("json.decode returned a %T payload: %s", got, got)
			}
			if err != nil {
				t.Fatalf("json.decode failed: %s", err)
			}
			if !reflect.DeepEqual(gotMap, test.expected) {
				t.Fatalf("json.decode got %x, expected %s", got, test.expected)
			}
		})
	}
}

func TestBadJsonDecode(t *testing.T) {
	stringEncoder := processor.JsonDecode{}
	tests := []struct {
		name        string
		payload     any
		errorString string
	}{
		{
			name:        "non-string input",
			payload:     []byte("hello"),
			errorString: "json.decode processor only accepts a string",
		},
		{
			name:        "invalid json",
			payload:     "{\"address\":\"/hello\",\"args\":}",
			errorString: "invalid character '}' looking for beginning of value",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := stringEncoder.Process(t.Context(), test.payload)

			if err == nil {
				t.Fatalf("json.decode expected to fail but got payload: %+v", got)
			}
			if err.Error() != test.errorString {
				t.Fatalf("json.decode got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
