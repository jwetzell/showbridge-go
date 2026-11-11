package processor_test

import (
	"context"
	"testing"

	"github.com/jwetzell/showbridge-go/config"
	"github.com/jwetzell/showbridge-go/processor"
	"github.com/jwetzell/showbridge-go/routing"
	"github.com/stretchr/testify/require"
)

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
			payload: TestStruct{
				Int: 1,
			},
			match: true,
		},
		{
			name: "string",
			params: map[string]any{
				"expression": "Payload.String == 'hello'",
			},
			payload: TestStruct{
				String: "hello",
			},
			match: true,
		},
		{
			name: "not matching",
			params: map[string]any{
				"expression": "Payload.Int > 0",
			},
			payload: TestStruct{
				Int: 0,
			},
			match: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			registration, ok := processor.GetProcessorRegistration("filter.expr")
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

			filtered := false

			processFunc := processorInstance.Process(func(ctx context.Context, wp routing.WrappedPayload, mm routing.ModuleManager, r routing.Router) (routing.WrappedPayload, error) {
				filtered = true
				return wp, nil
			})

			_, err = processFunc(t.Context(), routing.WrappedPayload{Payload: testCase.payload}, nil, nil)
			require.NoError(t, err)
			require.Equal(t, testCase.match, filtered)
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
			registration, ok := processor.GetProcessorRegistration("filter.expr")
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
			got, err := EvaluateProcessFunc(processorInstance, t.Context(), routing.WrappedPayload{Payload: test.payload}, nil, nil)

			if err == nil {
				t.Fatalf("filter.expr expected to fail but succeeded, got: %v", got)

			}
			if err.Error() != test.errorString {
				t.Fatalf("filter.expr got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
