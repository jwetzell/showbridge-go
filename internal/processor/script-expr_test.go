package processor_test

import (
	"testing"

	"github.com/expr-lang/expr"
	"github.com/jwetzell/showbridge-go/internal/processor"
)

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
