package processor_test

import (
	"slices"
	"testing"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/processor"
)

func TestScriptWASMFromRegistry(t *testing.T) {
	registration, ok := processor.ProcessorRegistry["script.wasm"]
	if !ok {
		t.Fatalf("script.wasm processor not registered")
	}

	processorInstance, err := registration.New(config.ProcessorConfig{
		Type: "script.wasm",
		Params: map[string]any{
			"path": "good.wasm",
		},
	})
	if err != nil {
		t.Fatalf("failed to create script.wasm processor: %s", err)
	}

	if processorInstance.Type() != "script.wasm" {
		t.Fatalf("script.wasm processor has wrong type: %s", processorInstance.Type())
	}

}

func TestGoodScriptWASM(t *testing.T) {
	tests := []struct {
		name     string
		payload  []byte
		params   map[string]any
		expected []byte
	}{
		{
			name: "string input, default process function with wasi",
			params: map[string]any{
				"path":       "good.wasm",
				"enableWasi": true,
			},
			payload:  []byte("hello"),
			expected: []byte("Processed: hello"),
		},
		{
			name: "string input, specified function with wasi",
			params: map[string]any{
				"path":       "good.wasm",
				"enableWasi": true,
				"function":   "greet",
			},
			payload:  []byte("world"),
			expected: []byte("Hello, world"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			registration, ok := processor.ProcessorRegistry["script.wasm"]
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

			got, err := processorInstance.Process(t.Context(), test.payload)

			if err != nil {
				t.Fatalf("script.wasm process failed: %s", err)
			}

			gotBytes, ok := got.([]byte)
			if !ok {
				t.Fatalf("script.wasm returned a %T payload: %s", got, got)
			}

			if !slices.Equal(gotBytes, test.expected) {
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
			errorString: "script.wasm requires a path parameter",
		},
		{
			name: "non-string path parameter",
			params: map[string]any{
				"path": 12345,
			},
			payload:     []byte("hello"),
			errorString: "script.wasm path must be a string",
		},
		{
			name: "non-string function",
			params: map[string]any{
				"path":       "good.wasm",
				"enableWasi": true,
				"function":   12345,
			},
			payload:     []byte("hello"),
			errorString: "script.wasm function must be a string",
		},
		{
			name: "non-boolean enableWasi",
			params: map[string]any{
				"path":       "good.wasm",
				"enableWasi": "true",
			},
			payload:     []byte("hello"),
			errorString: "script.wasm enableWasi must be a boolean",
		},
		{
			name: "non-byte slice input",
			params: map[string]any{
				"path":       "good.wasm",
				"enableWasi": true,
			},
			payload:     "hello",
			errorString: "script.wasm can only operator on byte array",
		},
		{
			name: "function not found in module",
			params: map[string]any{
				"path":       "good.wasm",
				"enableWasi": true,
				"function":   "asdf",
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
			errorString: "open asdf.wasm: no such file or directory",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			registration, ok := processor.ProcessorRegistry["script.wasm"]
			if !ok {
				t.Fatalf("script.wasm processor not registered")
			}

			processorInstance, err := registration.New(config.ProcessorConfig{
				Type:   "script.wasm",
				Params: test.params,
			})

			if err != nil {
				if test.errorString != err.Error() {
					t.Fatalf("string.create got error '%s', expected '%s'", err.Error(), test.errorString)
				}
				return
			}

			got, err := processorInstance.Process(t.Context(), test.payload)

			if err == nil {
				t.Fatalf("script.wasm expected to fail but succeeded, got: %v", got)
			}

			if err.Error() != test.errorString {
				t.Fatalf("script.wasm got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
