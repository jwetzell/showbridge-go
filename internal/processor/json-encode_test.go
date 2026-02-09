package processor_test

import (
	"slices"
	"testing"

	"github.com/jwetzell/osc-go"
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
			name: "hello",
			payload: osc.OSCMessage{
				Address: "/hello",
			},
			expected: []byte("{\"address\":\"/hello\",\"args\":null}"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := jsonEncoder.Process(t.Context(), test.payload)

			gotBytes, ok := got.([]byte)
			if !ok {
				t.Fatalf("json.encode returned a %T payload: %s", got, got)
			}
			if err != nil {
				t.Fatalf("json.encode failed: %s", err)
			}
			if !slices.Equal(gotBytes, test.expected) {
				t.Fatalf("json.encode got %x, expected %s", got, test.expected)
			}
		})
	}
}
