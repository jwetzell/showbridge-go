//go:build !js

package processor

import (
	"context"
	"fmt"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/config"
	"modernc.org/quickjs"
)

type ScriptJS struct {
	config      config.ProcessorConfig
	vm          *quickjs.VM
	payloadAtom quickjs.Atom
	senderAtom  quickjs.Atom
	Program     string
}

func (sj *ScriptJS) Process(ctx context.Context, wrappedPayload common.WrappedPayload) (common.WrappedPayload, error) {

	//NOTE(jwetzell): some weird conversion going on with these types
	_, isUint8Slice := common.GetAnyAs[[]uint8](wrappedPayload.Payload)
	_, isByteSlice := common.GetAnyAs[[]byte](wrappedPayload.Payload)

	if isUint8Slice || isByteSlice {
		intSlice, ok := common.GetAnyAsIntSlice(wrappedPayload.Payload)

		if ok {
			wrappedPayload.Payload = intSlice
		}
	}

	sj.vm.SetProperty(sj.vm.GlobalObject(), sj.payloadAtom, wrappedPayload.Payload)

	sj.vm.SetProperty(sj.vm.GlobalObject(), sj.senderAtom, wrappedPayload.Sender)

	_, err := sj.vm.Eval(sj.Program, quickjs.EvalGlobal)

	if err != nil {
		wrappedPayload.End = true
		return wrappedPayload, err
	}

	output, err := sj.vm.GetProperty(sj.vm.GlobalObject(), sj.payloadAtom)

	if err != nil {
		wrappedPayload.End = true
		return wrappedPayload, err
	}

	// NOTE(jwetzell): turn undefined into nil
	_, ok := output.(quickjs.Undefined)

	if ok {
		wrappedPayload.End = true
		wrappedPayload.Payload = nil
		return wrappedPayload, nil
	}

	// NOTE(jwetzell): turn object into map[string]interface{}
	outputObject, ok := output.(*quickjs.Object)

	if ok {
		var outputSlice []interface{}

		err = outputObject.Into(&outputSlice)

		if err != nil {
			var outputMap map[string]interface{}
			err = outputObject.Into(&outputMap)
			if err != nil {
				wrappedPayload.End = true
				return wrappedPayload, err
			} else {
				wrappedPayload.Payload = outputMap
				return wrappedPayload, nil
			}

		} else {
			wrappedPayload.Payload = outputSlice
			return wrappedPayload, nil
		}
	}

	wrappedPayload.Payload = output
	return wrappedPayload, nil
}

func (sj *ScriptJS) Type() string {
	return sj.config.Type
}

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type:  "script.js",
		Title: "Run JavaScript",
		ParamsSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"program": {
					Title: "Program",
					Type:  "string",
				},
			},
			Required:             []string{"program"},
			AdditionalProperties: nil,
		},
		New: func(config config.ProcessorConfig) (Processor, error) {
			params := config.Params

			programString, err := params.GetString("program")
			if err != nil {
				return nil, fmt.Errorf("script.js program error: %w", err)
			}

			vm, err := quickjs.NewVM()

			if err != nil {
				return nil, err
			}

			payloadAtom, err := vm.NewAtom("payload")
			if err != nil {
				return nil, err
			}

			senderAtom, err := vm.NewAtom("sender")
			if err != nil {
				return nil, err
			}

			return &ScriptJS{config: config, Program: programString, vm: vm, payloadAtom: payloadAtom, senderAtom: senderAtom}, nil
		},
	})
}
