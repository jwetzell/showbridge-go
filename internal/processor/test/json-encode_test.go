package processor_test

import (
	"slices"
	"testing"

	"github.com/jwetzell/osc-go"
	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/processor"
)

func TestJsonEncodeFromRegistry(t *testing.T) {
	registration, ok := processor.ProcessorRegistry["json.encode"]
	if !ok {
		t.Fatalf("json.encode processor not registered")
	}

	processorInstance, err := registration.New(config.ProcessorConfig{
		Type: "json.encode",
	})
	if err != nil {
		t.Fatalf("failed to create json.encode processor: %s", err)
	}

	if processorInstance.Type() != "json.encode" {
		t.Fatalf("json.encode processor has wrong type: %s", processorInstance.Type())
	}

	payload := struct {
		Property string `json:"property"`
	}{
		Property: "hello",
	}

	expected := []byte("{\"property\":\"hello\"}")

	got, err := processorInstance.Process(t.Context(), payload)
	if err != nil {
		t.Fatalf("json.encode processing failed: %s", err)
	}

	gotBytes, ok := got.([]byte)

	if !ok {
		t.Fatalf("json.encode should return byte slice")
	}

	if !slices.Equal(gotBytes, expected) {
		t.Fatalf("json.encode got %+v, expected %+v", got, expected)
	}
}

func TestGoodJsonEncode(t *testing.T) {
	jsonEncoder := processor.JsonEncode{}
	tests := []struct {
		name     string
		payload  any
		expected []byte
	}{
		{
			name: "basic struct",
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
