package processor_test

import (
	"reflect"
	"testing"

	"github.com/jwetzell/showbridge-go/config"
	"github.com/jwetzell/showbridge-go/processor"
	"github.com/jwetzell/showbridge-go/routing"
	"github.com/stretchr/testify/require"
)

func TestGoodJsonDecode(t *testing.T) {
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
			registration, ok := processor.GetProcessorRegistration("json.decode")
			if !ok {
				t.Fatalf("json.decode processor not registered")
			}

			processorInstance, err := registration.New(config.ProcessorConfig{
				Type:   "json.decode",
				Params: nil,
			})
			require.NoError(t, err)
			got, err := processorInstance.Process(t.Context(), routing.WrappedPayload{Payload: test.payload}, nil, nil)
			require.NoError(t, err)

			gotMap, ok := got.Payload.(map[string]any)
			if !ok {
				t.Fatalf("json.decode returned a %T payload: %+v", got, got)
			}

			if !reflect.DeepEqual(gotMap, test.expected) {
				t.Fatalf("json.decode got %+v, expected %s", got, test.expected)
			}
		})
	}
}

func TestBadJsonDecode(t *testing.T) {
	tests := []struct {
		name        string
		payload     any
		errorString string
	}{
		{
			name:        "non-string or byte input",
			payload:     123,
			errorString: "json.decode can only process a string or []byte",
		},
		{
			name:        "invalid json",
			payload:     "{\"address\":\"/hello\",\"args\":}",
			errorString: "invalid character '}' looking for beginning of value",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			registration, ok := processor.GetProcessorRegistration("json.decode")
			if !ok {
				t.Fatalf("json.decode processor not registered")
			}

			processorInstance, err := registration.New(config.ProcessorConfig{
				Type:   "json.decode",
				Params: nil,
			})
			if err != nil {
				if err.Error() == test.errorString {
					return
				}
				t.Fatalf("json.decode got error '%s', expected '%s'", err.Error(), test.errorString)
			}
			got, err := processorInstance.Process(t.Context(), routing.WrappedPayload{Payload: test.payload}, nil, nil)

			if err == nil {
				t.Fatalf("json.decode expected to fail but got payload: %+v", got)
			}
			if err.Error() != test.errorString {
				t.Fatalf("json.decode got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
