package processing

import (
	"context"
	"fmt"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/vm"
)

// NOTE(jwetzell): see language definition https://expr-lang.org/docs/language-definition
type DebugExpr struct {
	config  ProcessorConfig
	Program *vm.Program
}

func (dl *DebugExpr) Process(ctx context.Context, payload any) (any, error) {

	output, err := expr.Run(dl.Program, payload)
	if err != nil {
		return nil, err
	}

	return output, nil
}

func (dl *DebugExpr) Type() string {
	return dl.config.Type
}

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type: "debug.expr",
		New: func(config ProcessorConfig) (Processor, error) {
			params := config.Params

			expression, ok := params["expression"]

			if !ok {
				return nil, fmt.Errorf("debug.expr requires an expression parameter")
			}

			expressionString, ok := expression.(string)

			if !ok {
				return nil, fmt.Errorf("debug.expr url must be a string")
			}

			program, err := expr.Compile(expressionString)
			if err != nil {
				return nil, err
			}

			return &DebugExpr{config: config, Program: program}, nil
		},
	})
}
