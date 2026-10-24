package processor_test

import (
	"reflect"
	"testing"

	"github.com/jwetzell/showbridge-go/config"
	"github.com/jwetzell/showbridge-go/payload"
	"github.com/jwetzell/showbridge-go/processor"
)

func TestGoodScriptWASM(t *testing.T) {
	tests := []struct {
		name     string
		payload  any
		params   map[string]any
		expected any
	}{
		{
			name: "string input, default process function with wasi",
			params: map[string]any{
				"path": "test/good.wasm",
			},
			payload:  []byte("hello"),
			expected: []byte("Processed: hello"),
		},
		{
			name: "string input, specified function with wasi",
			params: map[string]any{
				"path":     "test/good.wasm",
				"function": "greet",
			},
			payload:  []byte("world"),
			expected: []byte("Hello, world"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			registration, ok := processor.GetProcessorRegistration("script.wasm")
			if !ok {
				t.Fatalf("script.wasm processor not registered")
			}

			processorInstance, err := registration.New(config.ProcessorConfig{
				Type:   "script.wasm",
				Params: test.params,
			})

			if err != nil {
				t.Fatalf("script.wasm failed to create processor: %s", err)
			}

			got, err := processorInstance.Process(t.Context(), payload.WrappedPayload{Payload: test.payload}, nil, nil)

			if err != nil {
				t.Fatalf("script.wasm processing failed: %s", err)
			}

			gotBytes, ok := got.Payload.([]byte)
			if !ok {
				t.Fatalf("script.wasm returned a %T payload: %+v", got, got)
			}

			if !reflect.DeepEqual(gotBytes, test.expected) {
				t.Fatalf("script.wasm got %+v, expected %+v", gotBytes, test.expected)
			}
		})
	}
}

func TestBadScriptWASM(t *testing.T) {
	tests := []struct {
		name        string
		params      map[string]any
		payload     any
		errorString string
	}{
		{
			name:        "no path parameter",
			params:      map[string]any{},
			payload:     []byte("hello"),
			errorString: "script.wasm path error: not found",
		},
		{
			name: "non-string path parameter",
			params: map[string]any{
				"path": 12345,
			},
			payload:     []byte("hello"),
			errorString: "script.wasm path error: not a string",
		},
		{
			name: "non-string function",
			params: map[string]any{
				"path":     "test/good.wasm",
				"function": 12345,
			},
			payload:     []byte("hello"),
			errorString: "script.wasm function error: not a string",
		},
		{
			name: "non-boolean enableWasi",
			params: map[string]any{
				"path":       "test/good.wasm",
				"enableWasi": "true",
			},
			payload:     []byte("hello"),
			errorString: "script.wasm enableWasi error: not a boolean",
		},
		{
			name: "non-byte slice input",
			params: map[string]any{
				"path": "test/good.wasm",
			},
			payload:     123,
			errorString: "script.wasm can only process a byte array",
		},
		{
			name: "function not found in module",
			params: map[string]any{
				"path":     "test/good.wasm",
				"function": "asdf",
			},
			payload:     []byte("hello"),
			errorString: "unknown function: asdf",
		},
		{
			name: "path doesn't exist",
			params: map[string]any{
				"path": "asdf.wasm",
			},
			payload:     []byte("hello"),
			errorString: "script.wasm file not found: asdf.wasm",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			registration, ok := processor.GetProcessorRegistration("script.wasm")
			if !ok {
				t.Fatalf("script.wasm processor not registered")
			}

			processorInstance, err := registration.New(config.ProcessorConfig{
				Type:   "script.wasm",
				Params: test.params,
			})

			if err != nil {
				if test.errorString != err.Error() {
					t.Fatalf("script.wasm got error '%s', expected '%s'", err.Error(), test.errorString)
				}
				return
			}

			got, err := processorInstance.Process(t.Context(), payload.WrappedPayload{Payload: test.payload}, nil, nil)

			if err == nil {
				t.Fatalf("script.wasm expected to fail but succeeded, got: %v", got)
			}

			if err.Error() != test.errorString {
				t.Fatalf("script.wasm got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
