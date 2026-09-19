package processor_test

import (
	"reflect"
	"testing"

	"github.com/jwetzell/showbridge-go/config"
	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/processor"
	"github.com/jwetzell/showbridge-go/internal/test"
)

func TestOsWriteFileFromRegistry(t *testing.T) {
	registration, ok := processor.GetProcessorRegistration("os.writefile")
	if !ok {
		t.Fatalf("os.writefile processor not registered")
	}

	processorInstance, err := registration.New(config.ProcessorConfig{
		Id:   "test-id",
		Type: "os.writefile",
		Params: map[string]any{
			"path": "/tmp/test.txt",
		},
	})
	if err != nil {
		t.Fatalf("failed to create os.writefile processor: %s", err)
	}

	if processorInstance.Id() != "test-id" {
		t.Fatalf("os.writefile processor has wrong id: %s", processorInstance.Id())
	}

	if processorInstance.Type() != "os.writefile" {
		t.Fatalf("os.writefile processor has wrong type: %s", processorInstance.Type())
	}
}

func TestGoodOsWriteFile(t *testing.T) {
	testCases := []struct {
		name     string
		params   map[string]any
		payload  any
		expected []byte
	}{}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			registration, ok := processor.GetProcessorRegistration("os.writefile")
			if !ok {
				t.Fatalf("os.writefile processor not registered")
			}

			processorInstance, err := registration.New(config.ProcessorConfig{
				Type:   "os.writefile",
				Params: testCase.params,
			})

			if err != nil {
				t.Fatalf("os.writefile failed to create processor: %s", err)
			}

			got, err := processorInstance.Process(t.Context(), common.WrappedPayload{Payload: testCase.payload})

			if err != nil {
				t.Fatalf("os.writefile processing failed: %s", err)
			}

			if !reflect.DeepEqual(got.Payload, testCase.expected) {
				t.Fatalf("os.writefile got payload '%v', expected '%v'", got.Payload, testCase.expected)
			}
		})
	}
}

func TestBadOsWriteFile(t *testing.T) {
	tests := []struct {
		name        string
		params      map[string]any
		payload     any
		errorString string
	}{
		{
			name:        "no path parameter",
			params:      map[string]any{},
			payload:     test.TestStruct{},
			errorString: "os.writefile path error: not found",
		},
		{
			name: "non-string path parameter",
			params: map[string]any{
				"path": 12345,
			},
			payload:     test.TestStruct{},
			errorString: "os.writefile path error: not a string",
		},
		{
			name: "path template syntax error",
			params: map[string]any{
				"path": "{{",
			},
			payload:     test.TestStruct{},
			errorString: "template: path:1: unclosed action",
		},
		{
			name: "path templating error",
			params: map[string]any{
				"path": "{{.NonExistentField}}",
			},
			payload:     test.TestStruct{},
			errorString: "template: path:1:2: executing \"path\" at <.NonExistentField>: can't evaluate field NonExistentField in type common.WrappedPayload",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			registration, ok := processor.GetProcessorRegistration("os.writefile")
			if !ok {
				t.Fatalf("os.writefile processor not registered")
			}

			processorInstance, err := registration.New(config.ProcessorConfig{
				Type:   "os.writefile",
				Params: test.params,
			})
			if err != nil {
				if err.Error() != test.errorString {
					t.Fatalf("os.writefile got error '%s', expected '%s'", err.Error(), test.errorString)
				}
				return
			}
			got, err := processorInstance.Process(t.Context(), common.WrappedPayload{Payload: test.payload})

			if err == nil {
				t.Fatalf("os.writefile expected to fail but succeeded, got: %v", got)

			}
			if err.Error() != test.errorString {
				t.Fatalf("os.writefile got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
