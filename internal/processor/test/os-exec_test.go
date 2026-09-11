package processor_test

import (
	"reflect"
	"testing"

	"github.com/jwetzell/showbridge-go/config"
	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/processor"
	"github.com/jwetzell/showbridge-go/internal/test"
)

func TestOsExecFromRegistry(t *testing.T) {
	registration, ok := processor.GetProcessorRegistration("os.exec")
	if !ok {
		t.Fatalf("os.exec processor not registered")
	}

	processorInstance, err := registration.New(config.ProcessorConfig{
		Id:   "test-id",
		Type: "os.exec",
		Params: map[string]any{
			"command": "ls",
		},
	})
	if err != nil {
		t.Fatalf("failed to create os.exec processor: %s", err)
	}

	if processorInstance.Id() != "test-id" {
		t.Fatalf("os.exec processor has wrong id: %s", processorInstance.Id())
	}

	if processorInstance.Type() != "os.exec" {
		t.Fatalf("os.exec processor has wrong type: %s", processorInstance.Type())
	}
}

func TestGoodOsExec(t *testing.T) {
	testCases := []struct {
		name     string
		params   map[string]any
		payload  any
		expected []byte
	}{}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			registration, ok := processor.GetProcessorRegistration("os.exec")
			if !ok {
				t.Fatalf("os.exec processor not registered")
			}

			processorInstance, err := registration.New(config.ProcessorConfig{
				Type:   "os.exec",
				Params: testCase.params,
			})

			if err != nil {
				t.Fatalf("os.exec failed to create processor: %s", err)
			}

			got, err := processorInstance.Process(t.Context(), common.WrappedPayload{Payload: testCase.payload})

			if err != nil {
				t.Fatalf("os.exec processing failed: %s", err)
			}

			if !reflect.DeepEqual(got.Payload, testCase.expected) {
				t.Fatalf("os.exec got payload '%v', expected '%v'", got.Payload, testCase.expected)
			}
		})
	}
}

func TestBadOsExec(t *testing.T) {
	tests := []struct {
		name        string
		params      map[string]any
		payload     any
		errorString string
	}{
		{
			name:        "no command parameter",
			params:      map[string]any{},
			payload:     test.TestStruct{},
			errorString: "os.exec command error: not found",
		},
		{
			name: "non-string command parameter",
			params: map[string]any{
				"command": 12345,
			},
			payload:     test.TestStruct{},
			errorString: "os.exec command error: not a string",
		},
		{
			name: "command template syntax error",
			params: map[string]any{
				"command": "{{",
			},
			payload:     test.TestStruct{},
			errorString: "template: command:1: unclosed action",
		},
		{
			name: "command templating error",
			params: map[string]any{
				"command": "{{.NonExistentField}}",
			},
			payload:     test.TestStruct{},
			errorString: "template: command:1:2: executing \"command\" at <.NonExistentField>: can't evaluate field NonExistentField in type common.WrappedPayload",
		},
		{
			name: "non-string in args",
			params: map[string]any{
				"command": "echo",
				"args":    []any{12345},
			},
			payload:     test.TestStruct{},
			errorString: "os.exec args error: not a string slice",
		},
		{
			name: "args template syntax error",
			params: map[string]any{
				"command": "echo",
				"args":    []any{"{{"},
			},
			payload:     test.TestStruct{},
			errorString: "template: arg-0:1: unclosed action",
		},
		{
			name: "args templating error",
			params: map[string]any{
				"command": "echo",
				"args":    []any{"{{.Unknown}}"},
			},
			payload:     test.TestStruct{},
			errorString: "template: arg-0:1:2: executing \"arg-0\" at <.Unknown>: can't evaluate field Unknown in type common.WrappedPayload",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			registration, ok := processor.GetProcessorRegistration("os.exec")
			if !ok {
				t.Fatalf("os.exec processor not registered")
			}

			processorInstance, err := registration.New(config.ProcessorConfig{
				Type:   "os.exec",
				Params: test.params,
			})
			if err != nil {
				if err.Error() != test.errorString {
					t.Fatalf("os.exec got error '%s', expected '%s'", err.Error(), test.errorString)
				}
				return
			}
			got, err := processorInstance.Process(t.Context(), common.WrappedPayload{Payload: test.payload})

			if err == nil {
				t.Fatalf("os.exec expected to fail but succeeded, got: %v", got)

			}
			if err.Error() != test.errorString {
				t.Fatalf("os.exec got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
