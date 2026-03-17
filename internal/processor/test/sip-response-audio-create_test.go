package processor_test

import (
	"reflect"
	"testing"

	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/processor"
)

func TestSipResponseAudioCreateFromRegistry(t *testing.T) {
	registration, ok := processor.ProcessorRegistry["sip.response.audio.create"]
	if !ok {
		t.Fatalf("sip.response.audio.create processor not registered")
	}

	processorInstance, err := registration.New(config.ProcessorConfig{
		Type: "sip.response.audio.create",
		Params: map[string]any{
			"preWait":   0,
			"audioFile": "good.wav",
			"postWait":  0,
		},
	})

	if err != nil {
		t.Fatalf("failed to filter sip.response.audio.create processor: %s", err)
	}

	if processorInstance.Type() != "sip.response.audio.create" {
		t.Fatalf("sip.response.audio.create processor has wrong type: %s", processorInstance.Type())
	}
}

func TestGoodSipResponseAudioCreate(t *testing.T) {

	tests := []struct {
		name     string
		params   map[string]any
		payload  any
		expected any
	}{}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			registration, ok := processor.ProcessorRegistry["sip.response.audio.create"]
			if !ok {
				t.Fatalf("sip.response.audio.create processor not registered")
			}

			processorInstance, err := registration.New(config.ProcessorConfig{
				Type:   "sip.response.audio.create",
				Params: test.params,
			})

			if err != nil {
				t.Fatalf("sip.response.audio.create failed to create processor: %s", err)
			}

			got, err := processorInstance.Process(t.Context(), common.GetWrappedPayload(t.Context(), test.payload))
			if err != nil {
				t.Fatalf("sip.response.audio.create processing failed: %s", err)
			}

			if !reflect.DeepEqual(got.Payload, test.expected) {
				t.Fatalf("sip.response.audio.create got %+v (%T), expected %+v (%T)", got.Payload, got.Payload, test.expected, test.expected)
			}
		})
	}
}

func TestBadSipResponseAudioCreate(t *testing.T) {
	tests := []struct {
		name        string
		params      map[string]any
		payload     any
		errorString string
	}{}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			registration, ok := processor.ProcessorRegistry["sip.response.audio.create"]
			if !ok {
				t.Fatalf("sip.response.audio.create processor not registered")
			}

			processorInstance, err := registration.New(config.ProcessorConfig{
				Type:   "sip.response.audio.create",
				Params: test.params,
			})

			if err != nil {
				if test.errorString != err.Error() {
					t.Fatalf("sip.response.audio.create got error '%s', expected '%s'", err.Error(), test.errorString)
				}
				return
			}

			got, err := processorInstance.Process(t.Context(), common.GetWrappedPayload(t.Context(), test.payload))

			if err == nil {
				t.Fatalf("sip.response.audio.create expected to fail but succeeded, got: %v", got)
			}

			if err.Error() != test.errorString {
				t.Fatalf("sip.response.audio.create got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
