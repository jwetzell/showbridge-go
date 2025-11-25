package processing

import (
	"context"
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

	output, err := vm.Eval(sj.Program, quickjs.EvalGlobal)
	if err != nil {
		return nil, err
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
