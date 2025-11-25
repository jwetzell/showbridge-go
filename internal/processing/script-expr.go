package processing

import (
	"context"
	"fmt"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/vm"
)

// NOTE(jwetzell): see language definition https://expr-lang.org/docs/language-definition
type ScriptExpr struct {
	config  ProcessorConfig
	Program *vm.Program
}

func (se *ScriptExpr) Process(ctx context.Context, payload any) (any, error) {

	output, err := expr.Run(se.Program, payload)
	if err != nil {
		return nil, err
	}

	return output, nil
}

func (se *ScriptExpr) Type() string {
	return se.config.Type
}

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type: "script.expr",
		New: func(config ProcessorConfig) (Processor, error) {
			params := config.Params

			expression, ok := params["expression"]

			if !ok {
				return nil, fmt.Errorf("script.expr requires an expression parameter")
			}

			expressionString, ok := expression.(string)

			if !ok {
				return nil, fmt.Errorf("script.expr expression must be a string")
			}

			program, err := expr.Compile(expressionString)
			if err != nil {
				return nil, err
			}

			return &ScriptExpr{config: config, Program: program}, nil
		},
	})
}
