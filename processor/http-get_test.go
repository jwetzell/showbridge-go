package processor_test

import (
	"reflect"
	"testing"

	"github.com/jwetzell/showbridge-go/config"
	"github.com/jwetzell/showbridge-go/processor"
	"github.com/jwetzell/showbridge-go/routing"
)

// TODO(jwetzell): figure out how to test this
func TestGoodHTTPGet(t *testing.T) {

	tests := []struct {
		name     string
		expected any
		params   map[string]any
		payload  any
	}{}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			registration, ok := processor.GetProcessorRegistration("http.get")
			if !ok {
				t.Fatalf("http.get processor not registered")
			}

			processorInstance, err := registration.New(config.ProcessorConfig{
				Type:   "http.get",
				Params: test.params,
			})

			if err != nil {
				t.Fatalf("http.get failed to create processor: %s", err)
			}

			got, err := processorInstance.Process(t.Context(), routing.WrappedPayload{Payload: test.payload}, nil, nil)
			if err != nil {
				t.Fatalf("http.get processing failed: %s", err)
			}

			if !reflect.DeepEqual(got.Payload, test.expected) {
				t.Fatalf("http.get got %+v (%T), expected %+v (%T)", got.Payload, got.Payload, test.expected, test.expected)
			}
		})
	}
}

func TestBadHTTPGet(t *testing.T) {
	tests := []struct {
		name        string
		params      map[string]any
		payload     any
		errorString string
	}{
		{
			name:        "missing url",
			params:      map[string]any{},
			payload:     nil,
			errorString: "http.get url error: not found",
		},
		{
			name: "url not a string",
			params: map[string]any{
				"url": 123,
			},
			payload:     nil,
			errorString: "http.get url error: not a string",
		},
		{
			name: "url template error",
			params: map[string]any{
				"url": "http://example.com/{{.Unknown}}",
			},
			payload:     nil,
			errorString: "template: url:1:21: executing \"url\" at <.Unknown>: can't evaluate field Unknown in type routing.WrappedPayload",
		},
		{
			name: "url template syntax error",
			params: map[string]any{
				"url": "http://example.com/{{.MissingEndBrace",
			},
			payload:     nil,
			errorString: "template: url:1: unclosed action",
		},
		{
			name: "headers not a map",
			params: map[string]any{
				"url":     "http://example.com",
				"headers": "not a map",
			},
			payload:     nil,
			errorString: "http.get headers error: not a map",
		},
		{
			name: "headers not a string map",
			params: map[string]any{
				"url":     "http://example.com",
				"headers": map[string]any{"key": 123},
			},
			payload:     nil,
			errorString: "http.get headers error: not a string map",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			registration, ok := processor.GetProcessorRegistration("http.get")
			if !ok {
				t.Fatalf("http.get processor not registered")
			}

			processorInstance, err := registration.New(config.ProcessorConfig{
				Type:   "http.get",
				Params: test.params,
			})

			if err != nil {
				if test.errorString != err.Error() {
					t.Fatalf("http.get got error '%s', expected '%s'", err.Error(), test.errorString)
				}
				return
			}

			got, err := processorInstance.Process(t.Context(), routing.WrappedPayload{Payload: test.payload}, nil, nil)

			if err == nil {
				t.Fatalf("http.get expected to fail but succeeded, got: %v", got)
			}

			if err.Error() != test.errorString {
				t.Fatalf("http.get got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
