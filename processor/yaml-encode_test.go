package processor_test

import (
	"slices"
	"testing"

	"github.com/jwetzell/osc-go"
	"github.com/jwetzell/showbridge-go/payload"
	"github.com/jwetzell/showbridge-go/processor"
)

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
			got, err := jsonEncoder.Process(t.Context(), payload.WrappedPayload{Payload: test.payload}, nil, nil)
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
			got, err := stringEncoder.Process(t.Context(), payload.WrappedPayload{Payload: test.payload}, nil, nil)

			if err == nil {
				t.Fatalf("yaml.encode expected to fail but got payload: %+v", got)
			}
			if err.Error() != test.errorString {
				t.Fatalf("yaml.encode got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
