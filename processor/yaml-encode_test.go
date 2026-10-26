package processor_test

import (
	"slices"
	"testing"

	"github.com/jwetzell/osc-go"
	"github.com/jwetzell/showbridge-go/config"
	"github.com/jwetzell/showbridge-go/processor"
	"github.com/jwetzell/showbridge-go/routing"
)

func TestGoodYamlEncode(t *testing.T) {
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
			registration, ok := processor.GetProcessorRegistration("yaml.encode")
			if !ok {
				t.Fatalf("yaml.encode processor not registered")
			}

			processorInstance, err := registration.New(config.ProcessorConfig{
				Type:   "yaml.encode",
				Params: nil,
			})
			got, err := processorInstance.Process(t.Context(), routing.WrappedPayload{Payload: test.payload}, nil, nil)
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
			registration, ok := processor.GetProcessorRegistration("yaml.encode")
			if !ok {
				t.Fatalf("yaml.encode processor not registered")
			}

			processorInstance, err := registration.New(config.ProcessorConfig{
				Type:   "yaml.encode",
				Params: nil,
			})
			if err != nil {
				if err.Error() == test.errorString {
					return
				}
				t.Fatalf("yaml.encode got error '%s', expected '%s'", err.Error(), test.errorString)
			}
			got, err := processorInstance.Process(t.Context(), routing.WrappedPayload{Payload: test.payload}, nil, nil)

			if err == nil {
				t.Fatalf("yaml.encode expected to fail but got payload: %+v", got)
			}
			if err.Error() != test.errorString {
				t.Fatalf("yaml.encode got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
