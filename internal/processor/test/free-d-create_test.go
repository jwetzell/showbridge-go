package processor_test

import (
	"reflect"
	"testing"

	freeD "github.com/jwetzell/free-d-go"
	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/processor"
)

func TestFreeDCreateFromRegistry(t *testing.T) {
	registration, ok := processor.ProcessorRegistry["freed.create"]
	if !ok {
		t.Fatalf("freed.create processor not registered")
	}

	processorInstance, err := registration.New(config.ProcessorConfig{
		Type: "freed.create",
		Params: map[string]any{
			"id":    "0",
			"pan":   "0",
			"tilt":  "0",
			"roll":  "0",
			"posX":  "0",
			"posY":  "0",
			"posZ":  "0",
			"zoom":  "0",
			"focus": "0",
		},
	})

	if err != nil {
		t.Fatalf("failed to create freed.create processor: %s", err)
	}

	if processorInstance.Type() != "freed.create" {
		t.Fatalf("freed.create processor has wrong type: %s", processorInstance.Type())
	}
}

func TestGoodFreeDCreate(t *testing.T) {

	tests := []struct {
		name     string
		expected freeD.FreeDPosition
		params   map[string]any
		payload  any
	}{
		{
			name: "basic freed",
			params: map[string]any{
				"id":    "{{.Payload.id}}",
				"pan":   "{{.Payload.pan}}",
				"tilt":  "{{.Payload.tilt}}",
				"roll":  "{{.Payload.roll}}",
				"posX":  "{{.Payload.posX}}",
				"posY":  "{{.Payload.posY}}",
				"posZ":  "{{.Payload.posZ}}",
				"zoom":  "{{.Payload.zoom}}",
				"focus": "{{.Payload.focus}}",
			},
			payload: map[string]any{
				"id":    "1",
				"pan":   "180",
				"tilt":  "90",
				"roll":  "-180",
				"posX":  "131069",
				"posY":  "131070",
				"posZ":  "131071",
				"zoom":  "66051",
				"focus": "263430",
			},
			expected: freeD.FreeDPosition{
				ID:    1,
				Pan:   180,
				Tilt:  90,
				Roll:  -180,
				PosX:  131069,
				PosY:  131070,
				PosZ:  131071,
				Zoom:  66051,
				Focus: 263430,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			registration, ok := processor.ProcessorRegistry["freed.create"]
			if !ok {
				t.Fatalf("freed.create processor not registered")
			}

			processorInstance, err := registration.New(config.ProcessorConfig{
				Type:   "freed.create",
				Params: test.params,
			})

			if err != nil {
				t.Fatalf("freed.create failed to create processor: %s", err)
			}

			got, err := processorInstance.Process(t.Context(), common.GetWrappedPayload(t.Context(), test.payload))
			if err != nil {
				t.Fatalf("freed.create processing failed: %s", err)
			}

			if !reflect.DeepEqual(got.Payload, test.expected) {
				t.Fatalf("freed.create got %+v (%T), expected %+v (%T)", got.Payload, got.Payload, test.expected, test.expected)
			}
		})
	}
}

func TestBadFreeDCreate(t *testing.T) {
	packetEncoder := processor.FreeDCreate{}
	tests := []struct {
		name        string
		payload     any
		errorString string
	}{}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			got, err := packetEncoder.Process(t.Context(), common.GetWrappedPayload(t.Context(), test.payload))

			if err == nil {
				t.Fatalf("freed.create expected to fail but succeeded, got: %v", got)

			}
			if err.Error() != test.errorString {
				t.Fatalf("freed.create got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
