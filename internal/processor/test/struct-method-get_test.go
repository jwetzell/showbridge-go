package processor_test

import (
	"testing"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/processor"
)

func TestStructMethodGetFromRegistry(t *testing.T) {
	registration, ok := processor.ProcessorRegistry["struct.method.get"]
	if !ok {
		t.Fatalf("struct.method.get processor not registered")
	}

	processorInstance, err := registration.New(config.ProcessorConfig{
		Type: "struct.method.get",
		Params: map[string]any{
			"name": "GetData",
		},
	})
	if err != nil {
		t.Fatalf("failed to create struct.method.get processor: %s", err)
	}

	if processorInstance.Type() != "struct.method.get" {
		t.Fatalf("struct.method.get processor has wrong type: %s", processorInstance.Type())
	}

	payload := TestStruct{Data: "hello"}
	expected := "hello"

	got, err := processorInstance.Process(t.Context(), payload)
	if err != nil {
		t.Fatalf("struct.method.get processing failed: %s", err)
	}

	if got != expected {
		t.Fatalf("struct.method.get got %+v, expected %+v", got, expected)
	}
}

func TestGoodStructMethodGet(t *testing.T) {

	tests := []struct {
		name     string
		params   map[string]any
		payload  any
		expected any
	}{
		{
			name:     "string field",
			params:   map[string]any{"name": "GetString"},
			payload:  TestStruct{String: "hello"},
			expected: "hello",
		},
		{
			name:     "int field",
			params:   map[string]any{"name": "GetInt"},
			payload:  TestStruct{Int: 42},
			expected: 42,
		},
		{
			name:     "float field",
			params:   map[string]any{"name": "GetFloat"},
			payload:  TestStruct{Float: 3.14},
			expected: 3.14,
		},
		{
			name:     "bool field",
			params:   map[string]any{"name": "GetBool"},
			payload:  TestStruct{Bool: true},
			expected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			registration, ok := processor.ProcessorRegistry["struct.method.get"]
			if !ok {
				t.Fatalf("struct.method.get processor not registered")
			}

			processorInstance, err := registration.New(config.ProcessorConfig{
				Type:   "struct.method.get",
				Params: test.params,
			})

			if err != nil {
				t.Fatalf("struct.method.get failed to create processor: %s", err)
			}

			got, err := processorInstance.Process(t.Context(), test.payload)

			if err != nil {
				t.Fatalf("struct.method.get failed: %s", err)
			}

			if got != test.expected {
				t.Fatalf("struct.method.get got %s, expected %s", got, test.expected)
			}
		})
	}
}

func TestBadStructMethodGet(t *testing.T) {
	tests := []struct {
		name        string
		params      map[string]any
		payload     any
		errorString string
	}{
		{
			name:        "no name param",
			payload:     TestStruct{Data: "hello"},
			params:      map[string]any{},
			errorString: "struct.method.get name error: not found",
		},
		{
			name:    "non string name",
			payload: TestStruct{Data: "hello"},
			params: map[string]any{
				"name": 1,
			},
			errorString: "struct.method.get name error: not a string",
		},
		{
			name:    "missing method",
			payload: TestStruct{String: "hello"},
			params: map[string]any{
				"name": "NonExistentMethod",
			},
			errorString: "struct.method.get method 'NonExistentMethod' does not exist",
		},
		{
			name:    "not a struct payload",
			payload: "not a struct",
			params: map[string]any{
				"name": "NonExistentMethod",
			},
			errorString: "struct.method.get processor only accepts a struct payload",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			registration, ok := processor.ProcessorRegistry["struct.method.get"]
			if !ok {
				t.Fatalf("struct.method.get processor not registered")
			}

			processorInstance, err := registration.New(config.ProcessorConfig{
				Type:   "struct.method.get",
				Params: test.params,
			})

			if err != nil {
				if test.errorString != err.Error() {
					t.Fatalf("struct.method.get got error '%s', expected '%s'", err.Error(), test.errorString)
				}
				return
			}

			got, err := processorInstance.Process(t.Context(), test.payload)

			if err == nil {
				t.Fatalf("struct.method.get expected to fail but got payload: %s", got)
			}

			if err.Error() != test.errorString {
				t.Fatalf("struct.method.get got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
