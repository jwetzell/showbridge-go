package processing

import (
	"context"
	"fmt"

	"modernc.org/quickjs"
)

type ProgramJS struct {
	config  ProcessorConfig
	Program string
}

func (pj *ProgramJS) Process(ctx context.Context, payload any) (any, error) {

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

	output, err := vm.Eval(pj.Program, quickjs.EvalGlobal)
	if err != nil {
		return nil, err
	}

	return output, nil
}

func (pj *ProgramJS) Type() string {
	return pj.config.Type
}

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type: "program.js",
		New: func(config ProcessorConfig) (Processor, error) {
			params := config.Params

			program, ok := params["program"]

			if !ok {
				return nil, fmt.Errorf("program.js requires a program parameter")
			}

			programString, ok := program.(string)

			if !ok {
				return nil, fmt.Errorf("program.js program must be a string")
			}

			return &ProgramJS{config: config, Program: programString}, nil
		},
	})
}
