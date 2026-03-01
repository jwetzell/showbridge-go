package processor

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jwetzell/showbridge-go/internal/config"
	"modernc.org/quickjs"
)

type ScriptJS struct {
	config  config.ProcessorConfig
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

	if err != nil {
		return nil, err
	}

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
		New: func(config config.ProcessorConfig) (Processor, error) {
			params := config.Params

			programString, err := params.GetString("program")
			if err != nil {
				return nil, fmt.Errorf("script.js program error: %w", err)
			}

			return &ScriptJS{config: config, Program: programString}, nil
		},
	})
}
