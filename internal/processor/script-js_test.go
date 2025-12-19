package processor_test

import (
	"testing"

	"github.com/jwetzell/showbridge-go/internal/processor"
)

func TestGoodScriptJS(t *testing.T) {
	tests := []struct {
		processor processor.Processor
		name      string
		payload   any
		expected  any
	}{
		{
			processor: &processor.ScriptJS{Program: `
			payload = payload + 1
			`},
			name:     "number",
			payload:  1,
			expected: 2,
		},
		{
			processor: &processor.ScriptJS{Program: `
			payload = payload + "1"
			`},
			name:     "string",
			payload:  "1",
			expected: "11",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.processor.Process(t.Context(), test.payload)

			if err != nil {
				t.Errorf("script.js failed: %s", err)
			}

			//TODO(jwetzell): work out better way to compare the any/any
			if got != test.expected {
				t.Errorf("script.js got %+v, expected %+v", got, test.expected)
			}
		})
	}
}
