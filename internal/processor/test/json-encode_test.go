package processor_test

import (
	"slices"
	"testing"

	"github.com/jwetzell/osc-go"
	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/processor"
)

func TestGoodJsonEncode(t *testing.T) {
	jsonEncoder := processor.JsonEncode{}
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
			expected: []byte("{\"address\":\"/hello\",\"args\":null}"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := jsonEncoder.Process(t.Context(), common.WrappedPayload{Payload: test.payload}, nil, nil)
			if err != nil {
				t.Fatalf("json.encode processing failed: %s", err)
			}

			gotBytes, ok := got.Payload.([]byte)
			if !ok {
				t.Fatalf("json.encode returned a %T payload: %+v", got.Payload, got.Payload)
			}

			if !slices.Equal(gotBytes, test.expected) {
				t.Fatalf("json.encode got %+v, expected %s", got, test.expected)
			}
		})
	}
}

func TestBadJsonEncode(t *testing.T) {
	stringEncoder := processor.JsonEncode{}
	tests := []struct {
		name        string
		payload     any
		errorString string
	}{
		{
			name:        "unencodable type",
			payload:     make(chan int),
			errorString: "json: unsupported type: chan int",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := stringEncoder.Process(t.Context(), common.WrappedPayload{Payload: test.payload}, nil, nil)

			if err == nil {
				t.Fatalf("json.encode expected to fail but got payload: %+v", got)
			}
			if err.Error() != test.errorString {
				t.Fatalf("json.encode got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
