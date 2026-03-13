package processor

import (
	"context"
	"encoding/json"
	"fmt"

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

func (sj *ScriptJS) Process(ctx context.Context, payload any) (any, error) {

	//NOTE(jwetzell): some weird conversion going on with these types
	_, isUint8Slice := common.GetAnyAs[[]uint8](payload)
	_, isbyteSlice := common.GetAnyAs[[]byte](payload)

	if isUint8Slice || isbyteSlice {
		intSlice, ok := common.GetAnyAsIntSlice(payload)

		if ok {
			payload = intSlice
		}
	}

	sj.vm.SetProperty(sj.vm.GlobalObject(), sj.payloadAtom, payload)

	sender := ctx.Value(common.SenderContextKey)
	sj.vm.SetProperty(sj.vm.GlobalObject(), sj.senderAtom, sender)

	_, err := sj.vm.Eval(sj.Program, quickjs.EvalGlobal)

	if err != nil {
		return nil, err
	}

	output, err := sj.vm.GetProperty(sj.vm.GlobalObject(), sj.payloadAtom)

	if err != nil {
		return nil, err
	}

	// NOTE(jwetzell): turn undefined into nil
	_, ok := output.(quickjs.Undefined)

	if ok {
		return nil, nil
	}

	// NOTE(jwetzell): turn object into map[string]interface{}
	outputObject, ok := output.(*quickjs.Object)

	if ok {
		var outputMap map[string]interface{}
		fmt.Println(outputObject.String())
		err := json.Unmarshal([]byte(outputObject.String()), &outputMap)
		return outputMap, err
	}

	return output, nil
}

func (sj *ScriptJS) Type() string {
	return sj.config.Type
}

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type: "script.js",
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
