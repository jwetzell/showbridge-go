package processor_test

import (
	"reflect"
	"testing"

	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/processor"
)

func TestHTTPRequestCreateFromRegistry(t *testing.T) {
	registration, ok := processor.ProcessorRegistry["http.request.do"]
	if !ok {
		t.Fatalf("http.request.do processor not registered")
	}

	processorInstance, err := registration.New(config.ProcessorConfig{
		Type: "http.request.do",
		Params: map[string]any{
			"method": "GET",
			"url":    "http://example.com",
		},
	})

	if err != nil {
		t.Fatalf("failed to create http.request.do processor: %s", err)
	}

	if processorInstance.Type() != "http.request.do" {
		t.Fatalf("http.request.do processor has wrong type: %s", processorInstance.Type())
	}
}

func TestGoodHTTPRequestDo(t *testing.T) {

	tests := []struct {
		name     string
		expected processor.NATSMessage
		params   map[string]any
		payload  any
	}{}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			registration, ok := processor.ProcessorRegistry["http.request.do"]
			if !ok {
				t.Fatalf("http.request.do processor not registered")
			}

			processorInstance, err := registration.New(config.ProcessorConfig{
				Type:   "http.request.do",
				Params: test.params,
			})

			if err != nil {
				t.Fatalf("http.request.do failed to create processor: %s", err)
			}

			got, err := processorInstance.Process(t.Context(), common.GetWrappedPayload(t.Context(), test.payload))
			if err != nil {
				t.Fatalf("http.request.do processing failed: %s", err)
			}

			if !reflect.DeepEqual(got.Payload, test.expected) {
				t.Fatalf("http.request.do got %+v (%T), expected %+v (%T)", got.Payload, got.Payload, test.expected, test.expected)
			}
		})
	}
}

func TestBadHTTPRequestDo(t *testing.T) {
	tests := []struct {
		name        string
		params      map[string]any
		payload     any
		errorString string
	}{
		{
			name:        "missing method",
			params:      map[string]any{"url": "http://example.com"},
			payload:     nil,
			errorString: "http.request.do method error: not found",
		},
		{
			name: "method not a string",
			params: map[string]any{
				"method": 123,
				"url":    "http://example.com",
			},
			payload:     nil,
			errorString: "http.request.do method error: not a string",
		},
		{
			name: "missing url",
			params: map[string]any{
				"method": "GET",
			},
			payload:     nil,
			errorString: "http.request.do url error: not found",
		},
		{
			name: "url not a string",
			params: map[string]any{
				"method": "GET",
				"url":    123,
			},
			payload:     nil,
			errorString: "http.request.do url error: not a string",
		},
		{
			name: "url template error",
			params: map[string]any{
				"method": "GET",
				"url":    "http://example.com/{{.Unknown}}",
			},
			payload:     nil,
			errorString: "template: url:1:21: executing \"url\" at <.Unknown>: can't evaluate field Unknown in type common.WrappedPayload",
		},
		{
			name: "url template syntax error",
			params: map[string]any{
				"method": "GET",
				"url":    "http://example.com/{{.MissingEndBrace",
			},
			payload:     nil,
			errorString: "template: url:1: unclosed action",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			registration, ok := processor.ProcessorRegistry["http.request.do"]
			if !ok {
				t.Fatalf("http.request.do processor not registered")
			}

			processorInstance, err := registration.New(config.ProcessorConfig{
				Type:   "http.request.do",
				Params: test.params,
			})

			if err != nil {
				if test.errorString != err.Error() {
					t.Fatalf("http.request.do got error '%s', expected '%s'", err.Error(), test.errorString)
				}
				return
			}

			got, err := processorInstance.Process(t.Context(), common.GetWrappedPayload(t.Context(), test.payload))

			if err == nil {
				t.Fatalf("http.request.do expected to fail but succeeded, got: %v", got)
			}

			if err.Error() != test.errorString {
				t.Fatalf("http.request.do got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
