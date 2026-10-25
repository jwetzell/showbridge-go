package processor_test

import (
	"reflect"
	"testing"

	"github.com/jwetzell/showbridge-go/config"
	"github.com/jwetzell/showbridge-go/processor"
	"github.com/jwetzell/showbridge-go/routing"
)

// TODO(jwetzell): figure out how to test the actual HTTP response writing
func TestGoodHTTPResponseCreate(t *testing.T) {

	tests := []struct {
		name     string
		params   map[string]any
		expected any
		payload  any
	}{}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			registration, ok := processor.GetProcessorRegistration("http.response")
			if !ok {
				t.Fatalf("http.response processor not registered")
			}

			processorInstance, err := registration.New(config.ProcessorConfig{
				Type:   "http.response",
				Params: test.params,
			})

			if err != nil {
				t.Fatalf("http.response failed to create processor: %s", err)
			}

			got, err := processorInstance.Process(t.Context(), routing.WrappedPayload{Payload: test.payload}, nil, nil)
			if err != nil {
				t.Fatalf("http.response processing failed: %s", err)
			}

			if !reflect.DeepEqual(got.Payload, test.expected) {
				t.Fatalf("http.response got %+v (%T), expected %+v (%T)", got.Payload, got.Payload, test.expected, test.expected)
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
			errorString: "http.response status error: not found",
		},
		{
			name:        "non-number status",
			params:      map[string]any{"status": "200", "bodyTemplate": "Hello, World!"},
			payload:     nil,
			errorString: "http.response status error: not a number",
		},
		{
			name:        "missing bodyTemplate",
			params:      map[string]any{"status": 200},
			payload:     nil,
			errorString: "http.response bodyTemplate error: not found",
		},
		{
			name:        "non-string bodyTemplate",
			params:      map[string]any{"status": 200, "bodyTemplate": 123},
			payload:     nil,
			errorString: "http.response bodyTemplate error: not a string",
		},
		{
			name:        "bodyTemplate template error",
			params:      map[string]any{"status": 200, "bodyTemplate": "{{.MissingField}}"},
			payload:     nil,
			errorString: "template: body:1:2: executing \"body\" at <.MissingField>: can't evaluate field MissingField in type routing.WrappedPayload",
		},
		{
			name:        "bodyTemplate template syntax error",
			params:      map[string]any{"status": 200, "bodyTemplate": "{{.MissingField"},
			payload:     nil,
			errorString: "template: body:1: unclosed action",
		},
		{
			name:        "no response writer in context",
			params:      map[string]any{"status": 200, "bodyTemplate": "Hello, World!"},
			payload:     nil,
			errorString: "http.response input must come from a http.server",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			registration, ok := processor.GetProcessorRegistration("http.response")
			if !ok {
				t.Fatalf("http.response processor not registered")
			}

			processorInstance, err := registration.New(config.ProcessorConfig{
				Type:   "http.response",
				Params: test.params,
			})

			if err != nil {
				if test.errorString != err.Error() {
					t.Fatalf("http.response got error '%s', expected '%s'", err.Error(), test.errorString)
				}
				return
			}

			got, err := processorInstance.Process(t.Context(), routing.WrappedPayload{Payload: test.payload}, nil, nil)

			if err == nil {
				t.Fatalf("http.response expected to fail but succeeded, got: %v", got)
			}

			if err.Error() != test.errorString {
				t.Fatalf("http.response got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
