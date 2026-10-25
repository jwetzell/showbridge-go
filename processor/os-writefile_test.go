package processor_test

import (
	"reflect"
	"testing"

	"github.com/jwetzell/showbridge-go/config"
	"github.com/jwetzell/showbridge-go/payload"
	"github.com/jwetzell/showbridge-go/processor"
)

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

			got, err := processorInstance.Process(t.Context(), payload.WrappedPayload{Payload: testCase.payload}, nil, nil)

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
			payload:     TestStruct{},
			errorString: "os.writefile path error: not found",
		},
		{
			name: "non-string path parameter",
			params: map[string]any{
				"path": 12345,
			},
			payload:     TestStruct{},
			errorString: "os.writefile path error: not a string",
		},
		{
			name: "path template syntax error",
			params: map[string]any{
				"path": "{{",
			},
			payload:     TestStruct{},
			errorString: "template: path:1: unclosed action",
		},
		{
			name: "path templating error",
			params: map[string]any{
				"path": "{{.NonExistentField}}",
			},
			payload:     TestStruct{},
			errorString: "template: path:1:2: executing \"path\" at <.NonExistentField>: can't evaluate field NonExistentField in type payload.WrappedPayload",
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
			got, err := processorInstance.Process(t.Context(), payload.WrappedPayload{Payload: test.payload}, nil, nil)

			if err == nil {
				t.Fatalf("os.writefile expected to fail but succeeded, got: %v", got)

			}
			if err.Error() != test.errorString {
				t.Fatalf("os.writefile got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
