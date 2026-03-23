package processor_test

import (
	"testing"

	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/processor"
	"github.com/jwetzell/showbridge-go/internal/test"
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
	testCases := []struct {
		name    string
		params  map[string]any
		payload any
		match   bool
	}{
		{
			name: "number",
			params: map[string]any{
				"expression": "Payload.Int > 0",
			},
			payload: test.TestStruct{
				Int: 1,
			},
			match: true,
		},
		{
			name: "string",
			params: map[string]any{
				"expression": "Payload.String == 'hello'",
			},
			payload: test.TestStruct{
				String: "hello",
			},
			match: true,
		},
		{
			name: "not matching",
			params: map[string]any{
				"expression": "Payload.Int > 0",
			},
			payload: test.TestStruct{
				Int: 0,
			},
			match: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			registration, ok := processor.ProcessorRegistry["filter.expr"]
			if !ok {
				t.Fatalf("filter.expr processor not registered")
			}

			processorInstance, err := registration.New(config.ProcessorConfig{
				Type:   "filter.expr",
				Params: testCase.params,
			})

			if err != nil {
				t.Fatalf("filter.expr failed to create processor: %s", err)
			}

			got, err := processorInstance.Process(t.Context(), common.GetWrappedPayload(t.Context(), testCase.payload))

			if err != nil {
				t.Fatalf("filter.expr processing failed: %s", err)
			}

			//TODO(jwetzell): work out better way to compare the any/any
			if got.End != !testCase.match {
				t.Fatalf("filter.expr did fitler properly %+v (%T), expected %+v (%T)", got, got, testCase.match, testCase.match)
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
			payload:     test.TestStruct{},
			errorString: "filter.expr expression error: not found",
		},
		{
			name: "non-string expression parameter",
			params: map[string]any{
				"expression": 12345,
			},
			payload:     test.TestStruct{},
			errorString: "filter.expr expression error: not a string",
		},
		{
			name: "invalid expression",
			params: map[string]any{
				"expression": "foo +",
			},
			payload:     test.TestStruct{},
			errorString: "unexpected token EOF (1:5)\n | foo +\n | ....^",
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
			got, err := processorInstance.Process(t.Context(), common.GetWrappedPayload(t.Context(), test.payload))

			if err == nil {
				t.Fatalf("filter.expr expected to fail but succeeded, got: %v", got)

			}
			if err.Error() != test.errorString {
				t.Fatalf("filter.expr got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
