package processor_test

import (
	"reflect"
	"testing"

	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/processor"
)

func TestHTTPResponseCreateFromRegistry(t *testing.T) {
	registration, ok := processor.ProcessorRegistry["http.response.create"]
	if !ok {
		t.Fatalf("http.response.create processor not registered")
	}

	processorInstance, err := registration.New(config.ProcessorConfig{
		Type: "http.response.create",
		Params: map[string]any{
			"status":       200,
			"bodyTemplate": "Hello, World!",
		},
	})

	if err != nil {
		t.Fatalf("failed to create http.response.create processor: %s", err)
	}

	if processorInstance.Type() != "http.response.create" {
		t.Fatalf("http.response.create processor has wrong type: %s", processorInstance.Type())
	}
}

func TestGoodHTTPResponseCreate(t *testing.T) {

	tests := []struct {
		name     string
		expected processor.HTTPResponse
		params   map[string]any
		payload  any
	}{
		{
			name: "simple template",
			expected: processor.HTTPResponse{
				Status: 200,
				Body:   []byte("Hello, World!"),
			},
			params:  map[string]any{"status": 200, "bodyTemplate": "Hello, World!"},
			payload: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			registration, ok := processor.ProcessorRegistry["http.response.create"]
			if !ok {
				t.Fatalf("http.response.create processor not registered")
			}

			processorInstance, err := registration.New(config.ProcessorConfig{
				Type:   "http.response.create",
				Params: test.params,
			})

			if err != nil {
				t.Fatalf("http.response.create failed to create processor: %s", err)
			}

			got, err := processorInstance.Process(t.Context(), common.GetWrappedPayload(t.Context(), test.payload))
			if err != nil {
				t.Fatalf("http.response.create processing failed: %s", err)
			}

			if !reflect.DeepEqual(got.Payload, test.expected) {
				t.Fatalf("http.response.create got %+v (%T), expected %+v (%T)", got.Payload, got.Payload, test.expected, test.expected)
			}
		})
	}
}

func TestBadHTTPResponseCreate(t *testing.T) {
	tests := []struct {
		name        string
		params      map[string]any
		payload     any
		errorString string
	}{
		{
			name:        "missing status",
			params:      map[string]any{"bodyTemplate": "Hello, World!"},
			payload:     nil,
			errorString: "http.response.create status error: not found",
		},
		{
			name:        "non-number status",
			params:      map[string]any{"status": "200", "bodyTemplate": "Hello, World!"},
			payload:     nil,
			errorString: "http.response.create status error: not a number",
		},
		{
			name:        "missing bodyTemplate",
			params:      map[string]any{"status": 200},
			payload:     nil,
			errorString: "http.response.create bodyTemplate error: not found",
		},
		{
			name:        "non-string bodyTemplate",
			params:      map[string]any{"status": 200, "bodyTemplate": 123},
			payload:     nil,
			errorString: "http.response.create bodyTemplate error: not a string",
		},
		{
			name:        "bodyTemplate template error",
			params:      map[string]any{"status": 200, "bodyTemplate": "{{.MissingField}}"},
			payload:     nil,
			errorString: "template: body:1:2: executing \"body\" at <.MissingField>: can't evaluate field MissingField in type common.WrappedPayload",
		},
		{
			name:        "bodyTemplate template syntax error",
			params:      map[string]any{"status": 200, "bodyTemplate": "{{.MissingField"},
			payload:     nil,
			errorString: "template: body:1: unclosed action",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			registration, ok := processor.ProcessorRegistry["http.response.create"]
			if !ok {
				t.Fatalf("http.response.create processor not registered")
			}

			processorInstance, err := registration.New(config.ProcessorConfig{
				Type:   "http.response.create",
				Params: test.params,
			})

			if err != nil {
				if test.errorString != err.Error() {
					t.Fatalf("http.response.create got error '%s', expected '%s'", err.Error(), test.errorString)
				}
				return
			}

			got, err := processorInstance.Process(t.Context(), common.GetWrappedPayload(t.Context(), test.payload))

			if err == nil {
				t.Fatalf("http.response.create expected to fail but succeeded, got: %v", got)
			}

			if err.Error() != test.errorString {
				t.Fatalf("http.response.create got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
