package processor_test

import (
	"testing"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/processor"
)

func TestScriptExprFromRegistry(t *testing.T) {
	registration, ok := processor.ProcessorRegistry["script.expr"]
	if !ok {
		t.Fatalf("script.expr processor not registered")
	}

	processorInstance, err := registration.New(config.ProcessorConfig{
		Type: "script.expr",
		Params: map[string]any{
			"expression": "foo + bar",
		},
	})
	if err != nil {
		t.Fatalf("failed to create script.expr processor: %s", err)
	}

	if processorInstance.Type() != "script.expr" {
		t.Fatalf("script.expr processor has wrong type: %s", processorInstance.Type())
	}
}

func TestGoodScriptExpr(t *testing.T) {
	tests := []struct {
		name     string
		params   map[string]any
		payload  map[string]any
		expected any
	}{
		{
			name: "number",
			params: map[string]any{
				"expression": "Payload.foo + Payload.bar",
			},
			payload: map[string]any{
				"foo": 1,
				"bar": 1,
			},
			expected: 2,
		},
		{
			name: "string",
			params: map[string]any{
				"expression": "Payload.foo + Payload.bar",
			},
			payload: map[string]any{
				"foo": "1",
				"bar": "1",
			},
			expected: "11",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			registration, ok := processor.ProcessorRegistry["script.expr"]
			if !ok {
				t.Fatalf("script.expr processor not registered")
			}

			processorInstance, err := registration.New(config.ProcessorConfig{
				Type:   "script.expr",
				Params: test.params,
			})

			if err != nil {
				t.Fatalf("script.expr failed to create processor: %s", err)
			}

			got, err := processorInstance.Process(t.Context(), test.payload)

			if err != nil {
				t.Fatalf("script.expr processing failed: %s", err)
			}

			//TODO(jwetzell): work out better way to compare the any/any
			if got != test.expected {
				t.Fatalf("script.expr got %+v (%T), expected %+v (%T)", got, got, test.expected, test.expected)
			}
		})
	}
}

func TestBadScriptExpr(t *testing.T) {
	tests := []struct {
		name        string
		params      map[string]any
		payload     any
		errorString string
	}{
		{
			name:        "no expression parameter",
			params:      map[string]any{},
			payload:     map[string]any{"foo": 1, "bar": 1},
			errorString: "script.expr expression error: not found",
		},
		{
			name: "accessing missing field",
			params: map[string]any{
				"expression": "Payload.foo + Payload.bar",
			},
			payload: map[string]any{
				"foo": 1,
			},
			errorString: "invalid operation: int + <nil> (1:13)\n | Payload.foo + Payload.bar\n | ............^",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			registration, ok := processor.ProcessorRegistry["script.expr"]
			if !ok {
				t.Fatalf("script.expr processor not registered")
			}

			processorInstance, err := registration.New(config.ProcessorConfig{
				Type:   "script.expr",
				Params: test.params,
			})

			if err != nil {
				if err.Error() != test.errorString {
					t.Fatalf("script.expr got error '%s', expected '%s'", err.Error(), test.errorString)
				}
				return
			}

			got, err := processorInstance.Process(t.Context(), test.payload)

			if err == nil {
				t.Fatalf("script.expr expected to fail but succeeded, got: %v", got)

			}
			if err.Error() != test.errorString {
				t.Fatalf("script.expr got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
