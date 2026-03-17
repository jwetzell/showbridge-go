package processor_test

import (
	"reflect"
	"testing"

	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/processor"
)

func TestSipResponseDTMFCreateFromRegistry(t *testing.T) {
	registration, ok := processor.ProcessorRegistry["sip.response.dtmf.create"]
	if !ok {
		t.Fatalf("sip.response.dtmf.create processor not registered")
	}

	processorInstance, err := registration.New(config.ProcessorConfig{
		Type: "sip.response.dtmf.create",
		Params: map[string]any{
			"preWait":  0,
			"digits":   "good.wav",
			"postWait": 0,
		},
	})

	if err != nil {
		t.Fatalf("failed to filter sip.response.dtmf.create processor: %s", err)
	}

	if processorInstance.Type() != "sip.response.dtmf.create" {
		t.Fatalf("sip.response.dtmf.create processor has wrong type: %s", processorInstance.Type())
	}
}

func TestGoodSipResponseDTMFCreate(t *testing.T) {

	tests := []struct {
		name     string
		params   map[string]any
		payload  any
		expected any
	}{
		{
			name: "basic",
			params: map[string]any{
				"preWait":  0,
				"digits":   "12345",
				"postWait": 0,
			},
			payload: nil,
			expected: processor.SipDTMFResponse{
				PreWait:  0,
				PostWait: 0,
				Digits:   "12345",
			},
		},
		{
			name: "template digits",
			params: map[string]any{
				"preWait":  0,
				"digits":   "{{.Payload}}",
				"postWait": 0,
			},
			payload: "67890",
			expected: processor.SipDTMFResponse{
				PreWait:  0,
				PostWait: 0,
				Digits:   "67890",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			registration, ok := processor.ProcessorRegistry["sip.response.dtmf.create"]
			if !ok {
				t.Fatalf("sip.response.dtmf.create processor not registered")
			}

			processorInstance, err := registration.New(config.ProcessorConfig{
				Type:   "sip.response.dtmf.create",
				Params: test.params,
			})

			if err != nil {
				t.Fatalf("sip.response.dtmf.create failed to create processor: %s", err)
			}

			got, err := processorInstance.Process(t.Context(), common.GetWrappedPayload(t.Context(), test.payload))
			if err != nil {
				t.Fatalf("sip.response.dtmf.create processing failed: %s", err)
			}

			if !reflect.DeepEqual(got.Payload, test.expected) {
				t.Fatalf("sip.response.dtmf.create got %+v (%T), expected %+v (%T)", got.Payload, got.Payload, test.expected, test.expected)
			}
		})
	}
}

func TestBadSipResponseDTMFCreate(t *testing.T) {
	tests := []struct {
		name        string
		params      map[string]any
		payload     any
		errorString string
	}{
		{
			name: "missing preWait param",
			params: map[string]any{
				"digits":   "good.wav",
				"postWait": 0,
			},
			errorString: "sip.response.dtmf.create preWait error: not found",
		},
		{
			name: "non-numeric preWait param",
			params: map[string]any{
				"preWait":  "not a number",
				"digits":   "good.wav",
				"postWait": 0,
			},
			errorString: "sip.response.dtmf.create preWait error: not a number",
		},
		{
			name: "missing digits param",
			params: map[string]any{
				"preWait":  0,
				"postWait": 0,
			},
			errorString: "sip.response.dtmf.create digits error: not found",
		},
		{
			name: "non-string digits param",
			params: map[string]any{
				"preWait":  0,
				"digits":   12345,
				"postWait": 0,
			},
			errorString: "sip.response.dtmf.create digits error: not a string",
		},
		{
			name: "digits template syntax error",
			params: map[string]any{
				"preWait":  0,
				"digits":   "{{.Unclosed",
				"postWait": 0,
			},
			errorString: "template: digits:1: unclosed action",
		},
		{
			name: "digits template error",
			params: map[string]any{
				"preWait":  0,
				"digits":   "{{.NonExistentField}}	",
				"postWait": 0,
			},
			errorString: "template: digits:1:2: executing \"digits\" at <.NonExistentField>: can't evaluate field NonExistentField in type common.WrappedPayload",
		},
		{
			name:    "invalid digits template result",
			payload: "nhf",
			params: map[string]any{
				"preWait":  0,
				"digits":   "{{.Payload}}",
				"postWait": 0,
			},
			errorString: "sip.response.dtmf.create result of digits template contains invalid characters",
		},
		{
			name: "missing postWait param",
			params: map[string]any{
				"preWait": 0,
				"digits":  "good.wav",
			},
			errorString: "sip.response.dtmf.create postWait error: not found",
		},
		{
			name: "non-numeric postWait param",
			params: map[string]any{
				"preWait":  0,
				"digits":   "good.wav",
				"postWait": "not a number",
			},
			errorString: "sip.response.dtmf.create postWait error: not a number",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			registration, ok := processor.ProcessorRegistry["sip.response.dtmf.create"]
			if !ok {
				t.Fatalf("sip.response.dtmf.create processor not registered")
			}

			processorInstance, err := registration.New(config.ProcessorConfig{
				Type:   "sip.response.dtmf.create",
				Params: test.params,
			})

			if err != nil {
				if test.errorString != err.Error() {
					t.Fatalf("sip.response.dtmf.create got error '%s', expected '%s'", err.Error(), test.errorString)
				}
				return
			}

			got, err := processorInstance.Process(t.Context(), common.GetWrappedPayload(t.Context(), test.payload))

			if err == nil {
				t.Fatalf("sip.response.dtmf.create expected to fail but succeeded, got: %v", got)
			}

			if err.Error() != test.errorString {
				t.Fatalf("sip.response.dtmf.create got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
