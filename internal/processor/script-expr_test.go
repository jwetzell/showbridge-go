package processor_test

import (
	"testing"

	"github.com/expr-lang/expr"
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

func TestScriptExprNoProgram(t *testing.T) {
	registration, ok := processor.ProcessorRegistry["script.expr"]
	if !ok {
		t.Fatalf("script.expr processor not registered")
	}

	_, err := registration.New(config.ProcessorConfig{
		Type:   "script.expr",
		Params: map[string]any{},
	})

	if err == nil {
		t.Fatalf("script.expr processor should have thrown an error when creating")
	}
}

func TestScriptExprBadConfigWrongExpressionType(t *testing.T) {
	registration, ok := processor.ProcessorRegistry["script.expr"]
	if !ok {
		t.Fatalf("script.expr processor not registered")
	}

	_, err := registration.New(config.ProcessorConfig{
		Type: "script.expr",
		Params: map[string]any{
			"expression": 12345,
		},
	})

	if err == nil {
		t.Fatalf("script.expr processor should have thrown an error when creating with non-string expression")
	}
}

func TestScriptExprBadConfigNonCompilingExpression(t *testing.T) {
	registration, ok := processor.ProcessorRegistry["script.expr"]
	if !ok {
		t.Fatalf("script.expr processor not registered")
	}

	_, err := registration.New(config.ProcessorConfig{
		Type: "script.expr",
		Params: map[string]any{
			"expression": "foo + ",
		},
	})

	if err == nil {
		t.Fatalf("script.expr processor should have thrown an error when creating with non-compiling expression")
	}
}

func TestGoodScriptExpr(t *testing.T) {
	tests := []struct {
		program  string
		name     string
		payload  map[string]any
		expected any
	}{
		{
			program: "foo + bar",
			name:    "number",
			payload: map[string]any{
				"foo": 1,
				"bar": 1,
			},
			expected: 2,
		},
		{
			program: "foo + bar",
			name:    "string",
			payload: map[string]any{
				"foo": "1",
				"bar": "1",
			},
			expected: "11",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			program, err := expr.Compile(test.program)
			if err != nil {
				t.Fatalf("script.expr failed to compile program: %s", err)
			}

			exprProcessor := &processor.ScriptExpr{Program: program}

			got, err := exprProcessor.Process(t.Context(), test.payload)

			if err != nil {
				t.Fatalf("script.expr failed: %s", err)
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
		program     string
		name        string
		payload     map[string]any
		errorString string
	}{
		{
			name:    "accessing missing field",
			program: "foo + bar",
			payload: map[string]any{
				"foo": 1,
			},
			errorString: "invalid operation: int + <nil> (1:5)\n | foo + bar\n | ....^",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			program, err := expr.Compile(test.program)
			if err != nil {
				t.Fatalf("script.expr failed to compile program: %s", err)
			}

			exprProcessor := &processor.ScriptExpr{Program: program}

			got, err := exprProcessor.Process(t.Context(), test.payload)

			if err == nil {
				t.Fatalf("script.expr expected to fail but succeeded, got: %v", got)

			}
			if err.Error() != test.errorString {
				t.Fatalf("script.expr got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
