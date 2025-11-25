package processing

import (
	"context"
	"fmt"

	"modernc.org/quickjs"
)

// NOTE(jwetzell): see language definition https://expr-lang.org/docs/language-definition
type DebugJS struct {
	config  ProcessorConfig
	Program string
}

func (dl *DebugJS) Process(ctx context.Context, payload any) (any, error) {

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

	output, err := vm.Eval(dl.Program, quickjs.EvalGlobal)
	if err != nil {
		return nil, err
	}

	return output, nil
}

func (dl *DebugJS) Type() string {
	return dl.config.Type
}

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type: "debug.js",
		New: func(config ProcessorConfig) (Processor, error) {
			params := config.Params

			program, ok := params["program"]

			if !ok {
				return nil, fmt.Errorf("debug.js requires a program parameter")
			}

			programString, ok := program.(string)

			if !ok {
				return nil, fmt.Errorf("debug.js program must be a string")
			}

			return &DebugJS{config: config, Program: programString}, nil
		},
	})
}
