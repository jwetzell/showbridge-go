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
	tests := []struct {
		name        string
		params      map[string]any
		payload     any
		errorString string
	}{
		{
			name: "missing id",
			params: map[string]any{
				"pan":   "180",
				"tilt":  "90",
				"roll":  "-180",
				"posX":  "131069",
				"posY":  "131070",
				"posZ":  "131071",
				"zoom":  "66051",
				"focus": "263430",
			},
			errorString: "freed.create id error: not found",
		},
		{
			name: "missing pan",
			params: map[string]any{
				"id":    "1",
				"tilt":  "90",
				"roll":  "-180",
				"posX":  "131069",
				"posY":  "131070",
				"posZ":  "131071",
				"zoom":  "66051",
				"focus": "263430",
			},
			errorString: "freed.create pan error: not found",
		},
		{
			name: "missing tilt",
			params: map[string]any{
				"id":    "1",
				"pan":   "180",
				"roll":  "-180",
				"posX":  "131069",
				"posY":  "131070",
				"posZ":  "131071",
				"zoom":  "66051",
				"focus": "263430",
			},
			errorString: "freed.create tilt error: not found",
		},
		{
			name: "missing roll",
			params: map[string]any{
				"id":    "1",
				"pan":   "180",
				"tilt":  "90",
				"posX":  "131069",
				"posY":  "131070",
				"posZ":  "131071",
				"zoom":  "66051",
				"focus": "263430",
			},
			errorString: "freed.create roll error: not found",
		},
		{
			name: "missing posX",
			params: map[string]any{
				"id":    "1",
				"pan":   "180",
				"tilt":  "90",
				"roll":  "-180",
				"posY":  "131070",
				"posZ":  "131071",
				"zoom":  "66051",
				"focus": "263430",
			},
			errorString: "freed.create posX error: not found",
		},
		{
			name: "missing posY",
			params: map[string]any{
				"id":    "1",
				"pan":   "180",
				"tilt":  "90",
				"roll":  "-180",
				"posX":  "131069",
				"posZ":  "131071",
				"zoom":  "66051",
				"focus": "263430",
			},
			errorString: "freed.create posY error: not found",
		},
		{
			name: "missing posZ",
			params: map[string]any{
				"id":    "1",
				"pan":   "180",
				"tilt":  "90",
				"roll":  "-180",
				"posX":  "131069",
				"posY":  "131070",
				"zoom":  "66051",
				"focus": "263430",
			},
			errorString: "freed.create posZ error: not found",
		},
		{
			name: "missing zoom",
			params: map[string]any{
				"id":    "1",
				"pan":   "180",
				"tilt":  "90",
				"roll":  "-180",
				"posX":  "131069",
				"posY":  "131070",
				"posZ":  "131071",
				"focus": "263430",
			},
			errorString: "freed.create zoom error: not found",
		},
		{
			name: "missing focus",
			params: map[string]any{
				"id":   "1",
				"pan":  "180",
				"tilt": "90",
				"roll": "-180",
				"posX": "131069",
				"posY": "131070",
				"posZ": "131071",
				"zoom": "66051",
			},
			errorString: "freed.create focus error: not found",
		},
		{
			name: "id not string",
			params: map[string]any{
				"id":    1,
				"pan":   "180",
				"tilt":  "90",
				"roll":  "-180",
				"posX":  "131069",
				"posY":  "131070",
				"posZ":  "131071",
				"zoom":  "66051",
				"focus": "263430",
			},
			errorString: "freed.create id error: not a string",
			payload:     nil,
		},
		{
			name: "pan not string",
			params: map[string]any{
				"id":    "1",
				"pan":   180,
				"tilt":  "90",
				"roll":  "-180",
				"posX":  "131069",
				"posY":  "131070",
				"posZ":  "131071",
				"zoom":  "66051",
				"focus": "263430",
			},
			errorString: "freed.create pan error: not a string",
			payload:     nil,
		},
		{
			name: "tilt not string",
			params: map[string]any{
				"id":    "1",
				"pan":   "180",
				"tilt":  90,
				"roll":  "-180",
				"posX":  "131069",
				"posY":  "131070",
				"posZ":  "131071",
				"zoom":  "66051",
				"focus": "263430",
			},
			errorString: "freed.create tilt error: not a string",
			payload:     nil,
		},
		{
			name: "roll not string",
			params: map[string]any{
				"id":    "1",
				"pan":   "180",
				"tilt":  "90",
				"roll":  -180,
				"posX":  "131069",
				"posY":  "131070",
				"posZ":  "131071",
				"zoom":  "66051",
				"focus": "263430",
			},
			errorString: "freed.create roll error: not a string",
			payload:     nil,
		},
		{
			name: "posX not string",
			params: map[string]any{
				"id":    "1",
				"pan":   "180",
				"tilt":  "90",
				"roll":  "-180",
				"posX":  131069,
				"posY":  "131070",
				"posZ":  "131071",
				"zoom":  "66051",
				"focus": "263430",
			},
			errorString: "freed.create posX error: not a string",
			payload:     nil,
		},
		{
			name: "posY not string",
			params: map[string]any{
				"id":    "1",
				"pan":   "180",
				"tilt":  "90",
				"roll":  "-180",
				"posX":  "131069",
				"posY":  131070,
				"posZ":  "131071",
				"zoom":  "66051",
				"focus": "263430",
			},
			errorString: "freed.create posY error: not a string",
			payload:     nil,
		},
		{
			name: "posZ not string",
			params: map[string]any{
				"id":    "1",
				"pan":   "180",
				"tilt":  "90",
				"roll":  "-180",
				"posX":  "131069",
				"posY":  "131070",
				"posZ":  131071,
				"zoom":  "66051",
				"focus": "263430",
			},
			errorString: "freed.create posZ error: not a string",
			payload:     nil,
		},
		{
			name: "zoom not string",
			params: map[string]any{
				"id":    "1",
				"pan":   "180",
				"tilt":  "90",
				"roll":  "-180",
				"posX":  "131069",
				"posY":  "131070",
				"posZ":  "131071",
				"zoom":  66051,
				"focus": "263430",
			},
			errorString: "freed.create zoom error: not a string",
			payload:     nil,
		},
		{
			name: "focus not string",
			params: map[string]any{
				"id":    "1",
				"pan":   "180",
				"tilt":  "90",
				"roll":  "-180",
				"posX":  "131069",
				"posY":  "131070",
				"posZ":  "131071",
				"zoom":  "66051",
				"focus": 263430,
			},
			errorString: "freed.create focus error: not a string",
			payload:     nil,
		},

		{
			name: "id template syntax error",
			params: map[string]any{
				"id":    "{{",
				"pan":   "180",
				"tilt":  "90",
				"roll":  "-180",
				"posX":  "131069",
				"posY":  "131070",
				"posZ":  "131071",
				"zoom":  "66051",
				"focus": "263430",
			},
			errorString: "template: id:1: unclosed action",
			payload:     nil,
		},
		{
			name: "pan template syntax error",
			params: map[string]any{
				"id":    "1",
				"pan":   "{{",
				"tilt":  "90",
				"roll":  "-180",
				"posX":  "131069",
				"posY":  "131070",
				"posZ":  "131071",
				"zoom":  "66051",
				"focus": "263430",
			},
			errorString: "template: pan:1: unclosed action",
			payload:     nil,
		},
		{
			name: "tilt template syntax error",
			params: map[string]any{
				"id":    "1",
				"pan":   "180",
				"tilt":  "{{",
				"roll":  "-180",
				"posX":  "131069",
				"posY":  "131070",
				"posZ":  "131071",
				"zoom":  "66051",
				"focus": "263430",
			},
			errorString: "template: tilt:1: unclosed action",
			payload:     nil,
		},
		{
			name: "roll template syntax error",
			params: map[string]any{
				"id":    "1",
				"pan":   "180",
				"tilt":  "90",
				"roll":  "{{",
				"posX":  "131069",
				"posY":  "131070",
				"posZ":  "131071",
				"zoom":  "66051",
				"focus": "263430",
			},
			errorString: "template: roll:1: unclosed action",
			payload:     nil,
		},
		{
			name: "posX template syntax error",
			params: map[string]any{
				"id":    "1",
				"pan":   "180",
				"tilt":  "90",
				"roll":  "-180",
				"posX":  "{{",
				"posY":  "131070",
				"posZ":  "131071",
				"zoom":  "66051",
				"focus": "263430",
			},
			errorString: "template: posX:1: unclosed action",
			payload:     nil,
		},
		{
			name: "posY template syntax error",
			params: map[string]any{
				"id":    "1",
				"pan":   "180",
				"tilt":  "90",
				"roll":  "-180",
				"posX":  "131069",
				"posY":  "{{",
				"posZ":  "131071",
				"zoom":  "66051",
				"focus": "263430",
			},
			errorString: "template: posY:1: unclosed action",
			payload:     nil,
		},
		{
			name: "posZ template syntax error",
			params: map[string]any{
				"id":    "1",
				"pan":   "180",
				"tilt":  "90",
				"roll":  "-180",
				"posX":  "131069",
				"posY":  "131070",
				"posZ":  "{{",
				"zoom":  "66051",
				"focus": "263430",
			},
			errorString: "template: posZ:1: unclosed action",
			payload:     nil,
		},
		{
			name: "zoom template syntax error",
			params: map[string]any{
				"id":    "1",
				"pan":   "180",
				"tilt":  "90",
				"roll":  "-180",
				"posX":  "131069",
				"posY":  "131070",
				"posZ":  "131071",
				"zoom":  "{{",
				"focus": "263430",
			},
			errorString: "template: zoom:1: unclosed action",
			payload:     nil,
		},
		{
			name: "focus template syntax error",
			params: map[string]any{
				"id":    "1",
				"pan":   "180",
				"tilt":  "90",
				"roll":  "-180",
				"posX":  "131069",
				"posY":  "131070",
				"posZ":  "131071",
				"zoom":  "66051",
				"focus": "{{",
			},
			errorString: "template: focus:1: unclosed action",
			payload:     nil,
		},
		{
			name: "id template error",
			params: map[string]any{
				"id":    "{{.Unknown}}",
				"pan":   "180",
				"tilt":  "90",
				"roll":  "-180",
				"posX":  "131069",
				"posY":  "131070",
				"posZ":  "131071",
				"zoom":  "66051",
				"focus": "263430",
			},
			errorString: "template: id:1:2: executing \"id\" at <.Unknown>: can't evaluate field Unknown in type common.WrappedPayload",
			payload:     nil,
		},
		{
			name: "pan template error",
			params: map[string]any{
				"id":    "1",
				"pan":   "{{.Unknown}}",
				"tilt":  "90",
				"roll":  "-180",
				"posX":  "131069",
				"posY":  "131070",
				"posZ":  "131071",
				"zoom":  "66051",
				"focus": "263430",
			},
			errorString: "template: pan:1:2: executing \"pan\" at <.Unknown>: can't evaluate field Unknown in type common.WrappedPayload",
			payload:     nil,
		},
		{
			name: "tilt template error",
			params: map[string]any{
				"id":    "1",
				"pan":   "180",
				"tilt":  "{{.Unknown}}",
				"roll":  "-180",
				"posX":  "131069",
				"posY":  "131070",
				"posZ":  "131071",
				"zoom":  "66051",
				"focus": "263430",
			},
			errorString: "template: tilt:1:2: executing \"tilt\" at <.Unknown>: can't evaluate field Unknown in type common.WrappedPayload",
			payload:     nil,
		},
		{
			name: "roll template error",
			params: map[string]any{
				"id":    "1",
				"pan":   "180",
				"tilt":  "90",
				"roll":  "{{.Unknown}}",
				"posX":  "131069",
				"posY":  "131070",
				"posZ":  "131071",
				"zoom":  "66051",
				"focus": "263430",
			},
			errorString: "template: roll:1:2: executing \"roll\" at <.Unknown>: can't evaluate field Unknown in type common.WrappedPayload",
			payload:     nil,
		},
		{
			name: "posX template error",
			params: map[string]any{
				"id":    "1",
				"pan":   "180",
				"tilt":  "90",
				"roll":  "-180",
				"posX":  "{{.Unknown}}",
				"posY":  "131070",
				"posZ":  "131071",
				"zoom":  "66051",
				"focus": "263430",
			},
			errorString: "template: posX:1:2: executing \"posX\" at <.Unknown>: can't evaluate field Unknown in type common.WrappedPayload",
			payload:     nil,
		},
		{
			name: "posY template error",
			params: map[string]any{
				"id":    "1",
				"pan":   "180",
				"tilt":  "90",
				"roll":  "-180",
				"posX":  "131069",
				"posY":  "{{.Unknown}}",
				"posZ":  "131071",
				"zoom":  "66051",
				"focus": "263430",
			},
			errorString: "template: posY:1:2: executing \"posY\" at <.Unknown>: can't evaluate field Unknown in type common.WrappedPayload",
			payload:     nil,
		},
		{
			name: "posZ template error",
			params: map[string]any{
				"id":    "1",
				"pan":   "180",
				"tilt":  "90",
				"roll":  "-180",
				"posX":  "131069",
				"posY":  "131070",
				"posZ":  "{{.Unknown}}",
				"zoom":  "66051",
				"focus": "263430",
			},
			errorString: "template: posZ:1:2: executing \"posZ\" at <.Unknown>: can't evaluate field Unknown in type common.WrappedPayload",
			payload:     nil,
		},
		{
			name: "zoom template error",
			params: map[string]any{
				"id":    "1",
				"pan":   "180",
				"tilt":  "90",
				"roll":  "-180",
				"posX":  "131069",
				"posY":  "131070",
				"posZ":  "131071",
				"zoom":  "{{.Unknown}}",
				"focus": "263430",
			},
			errorString: "template: zoom:1:2: executing \"zoom\" at <.Unknown>: can't evaluate field Unknown in type common.WrappedPayload",
			payload:     nil,
		},
		{
			name: "focus template error",
			params: map[string]any{
				"id":    "1",
				"pan":   "180",
				"tilt":  "90",
				"roll":  "-180",
				"posX":  "131069",
				"posY":  "131070",
				"posZ":  "131071",
				"zoom":  "66051",
				"focus": "{{.Unknown}}",
			},
			errorString: "template: focus:1:2: executing \"focus\" at <.Unknown>: can't evaluate field Unknown in type common.WrappedPayload",
			payload:     nil,
		},
		{
			name: "id number parsing error",
			params: map[string]any{
				"id":    "{{.Payload.id}}",
				"pan":   "180",
				"tilt":  "90",
				"roll":  "-180",
				"posX":  "131069",
				"posY":  "131070",
				"posZ":  "131071",
				"zoom":  "66051",
				"focus": "263430",
			},
			payload: map[string]any{
				"id": "not a number",
			},
			errorString: "strconv.ParseUint: parsing \"not a number\": invalid syntax",
		},
		{
			name: "pan number parsing error",
			params: map[string]any{
				"id":    "1",
				"pan":   "{{.Payload.pan}}",
				"tilt":  "90",
				"roll":  "-180",
				"posX":  "131069",
				"posY":  "131070",
				"posZ":  "131071",
				"zoom":  "66051",
				"focus": "263430",
			},
			payload: map[string]any{
				"pan": "not a number",
			},
			errorString: "strconv.ParseFloat: parsing \"not a number\": invalid syntax",
		},
		{
			name: "tilt number parsing error",
			params: map[string]any{
				"id":    "1",
				"pan":   "180",
				"tilt":  "{{.Payload.tilt}}",
				"roll":  "-180",
				"posX":  "131069",
				"posY":  "131070",
				"posZ":  "131071",
				"zoom":  "66051",
				"focus": "263430",
			},
			payload: map[string]any{
				"tilt": "not a number",
			},
			errorString: "strconv.ParseFloat: parsing \"not a number\": invalid syntax",
		},
		{
			name: "roll number parsing error",
			params: map[string]any{
				"id":    "1",
				"pan":   "180",
				"tilt":  "90",
				"roll":  "{{.Payload.roll}}",
				"posX":  "131069",
				"posY":  "131070",
				"posZ":  "131071",
				"zoom":  "66051",
				"focus": "263430",
			},
			payload: map[string]any{
				"roll": "not a number",
			},
			errorString: "strconv.ParseFloat: parsing \"not a number\": invalid syntax",
		},
		{
			name: "posX number parsing error",
			params: map[string]any{
				"id":    "1",
				"pan":   "180",
				"tilt":  "90",
				"roll":  "-180",
				"posX":  "{{.Payload.posX}}",
				"posY":  "131070",
				"posZ":  "131071",
				"zoom":  "66051",
				"focus": "263430",
			},
			payload: map[string]any{
				"posX": "not a number",
			},
			errorString: "strconv.ParseFloat: parsing \"not a number\": invalid syntax",
		},
		{
			name: "posY number parsing error",
			params: map[string]any{
				"id":    "1",
				"pan":   "180",
				"tilt":  "90",
				"roll":  "-180",
				"posX":  "131069",
				"posY":  "{{.Payload.posY}}",
				"posZ":  "131071",
				"zoom":  "66051",
				"focus": "263430",
			},
			payload: map[string]any{
				"posY": "not a number",
			},
			errorString: "strconv.ParseFloat: parsing \"not a number\": invalid syntax",
		},
		{
			name: "posZ number parsing error",
			params: map[string]any{
				"id":    "1",
				"pan":   "180",
				"tilt":  "90",
				"roll":  "-180",
				"posX":  "131069",
				"posY":  "131070",
				"posZ":  "{{.Payload.posZ}}",
				"zoom":  "66051",
				"focus": "263430",
			},
			payload: map[string]any{
				"posZ": "not a number",
			},
			errorString: "strconv.ParseFloat: parsing \"not a number\": invalid syntax",
		},
		{
			name: "zoom number parsing error",
			params: map[string]any{
				"id":    "1",
				"pan":   "180",
				"tilt":  "90",
				"roll":  "-180",
				"posX":  "131069",
				"posY":  "131070",
				"posZ":  "131071",
				"zoom":  "{{.Payload.zoom}}",
				"focus": "263430",
			},
			payload: map[string]any{
				"zoom": "not a number",
			},
			errorString: "strconv.ParseInt: parsing \"not a number\": invalid syntax",
		},
		{
			name: "focus number parsing error",
			params: map[string]any{
				"id":    "1",
				"pan":   "180",
				"tilt":  "90",
				"roll":  "-180",
				"posX":  "131069",
				"posY":  "131070",
				"posZ":  "131071",
				"zoom":  "66051",
				"focus": "{{.Payload.focus}}",
			},
			payload: map[string]any{
				"focus": "not a number",
			},
			errorString: "strconv.ParseInt: parsing \"not a number\": invalid syntax",
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
				if test.errorString != err.Error() {
					t.Fatalf("freed.create got error '%s', expected '%s'", err.Error(), test.errorString)
				}
				return
			}

			got, err := processorInstance.Process(t.Context(), common.GetWrappedPayload(t.Context(), test.payload))

			if err == nil {
				t.Fatalf("freed.create expected to fail but succeeded, got: %v", got)
			}

			if err.Error() != test.errorString {
				t.Fatalf("freed.create got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
