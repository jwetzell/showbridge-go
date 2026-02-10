package processor_test

import (
	"maps"
	"testing"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/processor"
)

func TestScriptJSFromRegistry(t *testing.T) {
	registration, ok := processor.ProcessorRegistry["script.js"]
	if !ok {
		t.Fatalf("script.js processor not registered")
	}

	processorInstance, err := registration.New(config.ProcessorConfig{
		Type: "script.js",
		Params: map[string]any{
			"program": `
			payload = payload + 1
			`,
		},
	})
	if err != nil {
		t.Fatalf("failed to create script.js processor: %s", err)
	}

	if processorInstance.Type() != "script.js" {
		t.Fatalf("script.js processor has wrong type: %s", processorInstance.Type())
	}

	payload := 1
	expected := 2

	got, err := processorInstance.Process(t.Context(), payload)
	if err != nil {
		t.Fatalf("script.js processing failed: %s", err)
	}

	if got != expected {
		t.Fatalf("script.js got %+v, expected %+v", got, expected)
	}
}

func TestScriptJSNoProgram(t *testing.T) {
	registration, ok := processor.ProcessorRegistry["script.js"]
	if !ok {
		t.Fatalf("script.js processor not registered")
	}

	_, err := registration.New(config.ProcessorConfig{
		Type:   "script.js",
		Params: map[string]any{},
	})

	if err == nil {
		t.Fatalf("script.js processor should have thrown an error when creating")
	}
}

func TestScriptJSBadConfigWrongProgramType(t *testing.T) {
	registration, ok := processor.ProcessorRegistry["script.js"]
	if !ok {
		t.Fatalf("script.js processor not registered")
	}

	_, err := registration.New(config.ProcessorConfig{
		Type: "script.js",
		Params: map[string]any{
			"program": 12345,
		},
	})

	if err == nil {
		t.Fatalf("script.js processor should have thrown an error when creating with non-string program")
	}
}

func TestGoodScriptJS(t *testing.T) {
	tests := []struct {
		name     string
		params   map[string]any
		payload  any
		expected any
	}{
		{
			name: "number",
			params: map[string]any{
				"program": `
				payload = payload + 1
				`,
			},
			payload:  1,
			expected: 2,
		},
		{
			name: "string",
			params: map[string]any{
				"program": `
				payload = payload + "1"
				`,
			},
			payload:  "1",
			expected: "11",
		},
		{
			name: "object",
			params: map[string]any{
				"program": `
				payload = { key: payload }
				`,
			},
			payload:  "1",
			expected: map[string]any{"key": "1"},
		},
		{
			name: "nil",
			params: map[string]any{
				"program": `
				payload = undefined
				`,
			},
			payload:  "1",
			expected: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			registration, ok := processor.ProcessorRegistry["script.js"]
			if !ok {
				t.Fatalf("script.js processor not registered")
			}

			processorInstance, err := registration.New(config.ProcessorConfig{
				Type:   "script.js",
				Params: test.params,
			})

			if err != nil {
				t.Fatalf("script.js failed to create processor: %s", err)
			}

			got, err := processorInstance.Process(t.Context(), test.payload)

			if err != nil {
				t.Fatalf("script.js process failed: %s", err)
			}

			//TODO(jwetzell): work out better way to compare the any/any

			gotMap, ok := got.(map[string]interface{})
			if ok {
				// got a map
				expectedMap, ok := test.expected.(map[string]interface{})
				if ok {
					if !maps.Equal(gotMap, expectedMap) {
						t.Fatalf("script.js got %+v, expected %+v", got, test.expected)
					}
				} else {
					t.Fatalf("script.js got %+v, expected %+v", got, test.expected)
				}
			} else {
				if got != test.expected {
					t.Fatalf("script.js got %+v, expected %+v", got, test.expected)
				}
			}
		})
	}
}

func TestBadScriptJS(t *testing.T) {
	tests := []struct {
		name        string
		params      map[string]any
		payload     any
		errorString string
	}{
		{
			name: "accessing not defined variable",
			params: map[string]any{
				"program": `
				paylod = foo
				`,
			},
			payload:     0,
			errorString: "ReferenceError: 'foo' is not defined",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			registration, ok := processor.ProcessorRegistry["script.js"]
			if !ok {
				t.Fatalf("script.js processor not registered")
			}

			processorInstance, err := registration.New(config.ProcessorConfig{
				Type:   "script.js",
				Params: test.params,
			})

			got, err := processorInstance.Process(t.Context(), test.payload)

			if err == nil {
				t.Fatalf("script.js expected to fail but succeeded, got: %v", got)
			}

			if err.Error() != test.errorString {
				t.Fatalf("script.js got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
