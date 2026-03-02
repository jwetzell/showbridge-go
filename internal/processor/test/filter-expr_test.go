package processor_test

import (
	"reflect"
	"testing"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/processor"
)

func TestFilterExprFromRegistry(t *testing.T) {
	registration, ok := processor.ProcessorRegistry["filter.expr"]
	if !ok {
		t.Fatalf("filter.expr processor not registered")
	}

	processorInstance, err := registration.New(config.ProcessorConfig{
		Type: "filter.expr",
		Params: map[string]any{
			"expression": "foo + bar",
		},
	})
	if err != nil {
		t.Fatalf("failed to create filter.expr processor: %s", err)
	}

	if processorInstance.Type() != "filter.expr" {
		t.Fatalf("filter.expr processor has wrong type: %s", processorInstance.Type())
	}
}

func TestGoodFilterExpr(t *testing.T) {
	tests := []struct {
		name     string
		params   map[string]any
		payload  any
		expected any
	}{
		{
			name: "number",
			params: map[string]any{
				"expression": "Int > 0",
			},
			payload: TestStruct{
				Int: 1,
			},
			expected: TestStruct{
				Int: 1,
			},
		},
		{
			name: "string",
			params: map[string]any{
				"expression": "String == 'hello'",
			},
			payload: TestStruct{
				String: "hello",
			},
			expected: TestStruct{
				String: "hello",
			},
		},
		{
			name: "not matching",
			params: map[string]any{
				"expression": "Int > 0",
			},
			payload: TestStruct{
				Int: 0,
			},
			expected: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			registration, ok := processor.ProcessorRegistry["filter.expr"]
			if !ok {
				t.Fatalf("filter.expr processor not registered")
			}

			processorInstance, err := registration.New(config.ProcessorConfig{
				Type:   "filter.expr",
				Params: test.params,
			})

			if err != nil {
				t.Fatalf("filter.expr failed to create processor: %s", err)
			}

			got, err := processorInstance.Process(t.Context(), test.payload)

			if err != nil {
				t.Fatalf("filter.expr failed: %s", err)
			}

			//TODO(jwetzell): work out better way to compare the any/any
			if !reflect.DeepEqual(got, test.expected) {
				t.Fatalf("filter.expr got %+v (%T), expected %+v (%T)", got, got, test.expected, test.expected)
			}
		})
	}
}

func TestBadFilterExpr(t *testing.T) {
	tests := []struct {
		name        string
		params      map[string]any
		payload     any
		errorString string
	}{
		{
			name:   "no expression parameter",
			params: map[string]any{
				// no expression parameter
			},
			payload:     TestStruct{},
			errorString: "filter.expr expression error: not found",
		},
		{
			name: "non-string expression parameter",
			params: map[string]any{
				"expression": 12345,
			},
			payload:     TestStruct{},
			errorString: "filter.expr expression error: not a string",
		},
		{
			name: "invalid expression",
			params: map[string]any{
				"expression": "foo +",
			},
			payload:     TestStruct{},
			errorString: "unexpected token EOF (1:5)\n | foo +\n | ....^",
		},
		{
			name: "accessing missing field",
			params: map[string]any{
				"expression": "foo + bar",
			},
			payload: map[string]any{
				"foo": 1,
			},
			errorString: "invalid operation: int + <nil> (1:5)\n | foo + bar\n | ....^",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			registration, ok := processor.ProcessorRegistry["filter.expr"]
			if !ok {
				t.Fatalf("filter.expr processor not registered")
			}

			processorInstance, err := registration.New(config.ProcessorConfig{
				Type:   "filter.expr",
				Params: test.params,
			})
			if err != nil {
				if err.Error() != test.errorString {
					t.Fatalf("filter.expr got error '%s', expected '%s'", err.Error(), test.errorString)
				}
				return
			}
			got, err := processorInstance.Process(t.Context(), test.payload)

			if err == nil {
				t.Fatalf("filter.expr expected to fail but succeeded, got: %v", got)

			}
			if err.Error() != test.errorString {
				t.Fatalf("filter.expr got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
