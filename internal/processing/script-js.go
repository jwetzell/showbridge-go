package processing

import (
	"context"
	"encoding/json"
	"fmt"

	"modernc.org/quickjs"
)

type ScriptJS struct {
	config  ProcessorConfig
	Program string
}

func (sj *ScriptJS) Process(ctx context.Context, payload any) (any, error) {

	vm, err := quickjs.NewVM()

	if err != nil {
		return nil, err
	}
	defer vm.Close()
	payloadAtom, err := vm.NewAtom("payload")

	if err != nil {
		return nil, err
	}

	vm.SetProperty(vm.GlobalObject(), payloadAtom, payload)

	_, err = vm.Eval(sj.Program, quickjs.EvalGlobal)

	output, err := vm.GetProperty(vm.GlobalObject(), payloadAtom)

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
		New: func(config ProcessorConfig) (Processor, error) {
			params := config.Params

			program, ok := params["program"]

			if !ok {
				return nil, fmt.Errorf("script.js requires a program parameter")
			}

			programString, ok := program.(string)

			if !ok {
				return nil, fmt.Errorf("script.js program must be a string")
			}

			return &ScriptJS{config: config, Program: programString}, nil
		},
	})
}
