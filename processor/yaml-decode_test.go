package processor_test

import (
	"reflect"
	"testing"

	"github.com/jwetzell/showbridge-go/config"
	"github.com/jwetzell/showbridge-go/processor"
	"github.com/jwetzell/showbridge-go/routing"
	"github.com/stretchr/testify/require"
)

func TestGoodYamlDecode(t *testing.T) {
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
			registration, ok := processor.GetProcessorRegistration("yaml.decode")
			if !ok {
				t.Fatalf("yaml.decode processor not registered")
			}

			processorInstance, err := registration.New(config.ProcessorConfig{
				Type:   "yaml.decode",
				Params: nil,
			})
			require.NoError(t, err)
			got, err := processorInstance.Process(t.Context(), routing.WrappedPayload{Payload: test.payload}, nil, nil)
			require.NoError(t, err)

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
			registration, ok := processor.GetProcessorRegistration("yaml.decode")
			if !ok {
				t.Fatalf("yaml.decode processor not registered")
			}

			processorInstance, err := registration.New(config.ProcessorConfig{
				Type:   "yaml.decode",
				Params: nil,
			})
			if err != nil {
				if err.Error() == test.errorString {
					return
				}
				t.Fatalf("yaml.decode got error '%s', expected '%s'", err.Error(), test.errorString)
			}
			got, err := processorInstance.Process(t.Context(), routing.WrappedPayload{Payload: test.payload}, nil, nil)

			if err == nil {
				t.Fatalf("yaml.decode expected to fail but got payload: %+v", got)
			}
			if err.Error() != test.errorString {
				t.Fatalf("yaml.decode got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
