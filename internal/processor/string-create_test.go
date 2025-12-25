package processor_test

import (
	"testing"
	"text/template"

	"github.com/jwetzell/showbridge-go/internal/processor"
)

type TestStruct struct {
	Data string
}

func (t TestStruct) GetData() string {
	return t.Data
}

func TestGoodStringCreate(t *testing.T) {

	tests := []struct {
		name     string
		template string
		payload  any
		expected string
	}{
		{
			name:     "string payload",
			template: "{{.}}",
			payload:  "hello",
			expected: "hello",
		},
		{
			name:     "number payload",
			template: "{{.}}",
			payload:  4,
			expected: "4",
		},
		{
			name:     "boolean payload",
			template: "{{.}}",
			payload:  true,
			expected: "true",
		},
		{
			name:     "struct payload - field",
			template: "{{.Data}}",
			payload:  TestStruct{Data: "test"},
			expected: "test",
		},
		{
			name:     "struct payload - method",
			template: "{{.GetData}}",
			payload:  TestStruct{Data: "test"},
			expected: "test",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			template, err := template.New("template").Parse(test.template)
			if err != nil {
				t.Fatalf("string.create template parsing failed: %s", err)
			}

			processor := &processor.StringCreate{Template: template}

			got, err := processor.Process(t.Context(), test.payload)

			gotStrings, ok := got.(string)
			if !ok {
				t.Fatalf("string.create returned a %T payload: %s", got, got)
			}
			if err != nil {
				t.Fatalf("string.create failed: %s", err)
			}
			if gotStrings != test.expected {
				t.Fatalf("string.create got %s, expected %s", got, test.expected)
			}
		})
	}
}
