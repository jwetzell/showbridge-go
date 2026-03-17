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
	}{}

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
	}{}

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
