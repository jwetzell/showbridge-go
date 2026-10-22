package processor_test

import (
	"reflect"
	"testing"

	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/processor"
)

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
