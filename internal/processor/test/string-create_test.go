package processor_test

import (
	"testing"
	"text/template"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/processor"
)

type TestStruct struct {
	Data string
}

func (t TestStruct) GetData() string {
	return t.Data
}

func TestStringCreateFromRegistry(t *testing.T) {
	registration, ok := processor.ProcessorRegistry["string.create"]
	if !ok {
		t.Fatalf("string.create processor not registered")
	}

	processorInstance, err := registration.New(config.ProcessorConfig{
		Type: "string.create",
		Params: map[string]any{
			"template": "{{.}}",
		},
	})
	if err != nil {
		t.Fatalf("failed to create string.create processor: %s", err)
	}

	if processorInstance.Type() != "string.create" {
		t.Fatalf("string.create processor has wrong type: %s", processorInstance.Type())
	}

	payload := "hello"
	expected := "hello"

	got, err := processorInstance.Process(t.Context(), payload)
	if err != nil {
		t.Fatalf("string.create processing failed: %s", err)
	}

	if got != expected {
		t.Fatalf("string.create got %+v, expected %+v", got, expected)
	}
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

func TestBadStringCreate(t *testing.T) {
	tests := []struct {
		name        string
		params      map[string]any
		payload     any
		errorString string
	}{
		{
			name:        "no template param",
			payload:     "hello",
			params:      map[string]any{},
			errorString: "string.create requires a template parameter",
		},
		{
			name:    "non string template",
			payload: "hello",
			params: map[string]any{
				"template": 1,
			},
			errorString: "string.create template must be a string",
		},
		{
			name:    "invalid template",
			payload: "hello",
			params: map[string]any{
				"template": "{{.",
			},
			errorString: "template: template:1: illegal number syntax: \".\"",
		},
		{
			name:    "bad property in template",
			payload: "hello",
			params: map[string]any{
				"template": "{{.Invalid}}",
			},
			errorString: "template: template:1:2: executing \"template\" at <.Invalid>: can't evaluate field Invalid in type string",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			registration, ok := processor.ProcessorRegistry["string.create"]
			if !ok {
				t.Fatalf("string.create processor not registered")
			}

			processorInstance, err := registration.New(config.ProcessorConfig{
				Type:   "string.create",
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
				t.Fatalf("string.create expected to fail but got payload: %s", got)
			}

			if err.Error() != test.errorString {
				t.Fatalf("string.create got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
