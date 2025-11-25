package processing

import (
	"context"
	"fmt"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/vm"
)

// NOTE(jwetzell): see language definition https://expr-lang.org/docs/language-definition
type ProgramExpr struct {
	config  ProcessorConfig
	Program *vm.Program
}

func (pe *ProgramExpr) Process(ctx context.Context, payload any) (any, error) {

	output, err := expr.Run(pe.Program, payload)
	if err != nil {
		return nil, err
	}

	return output, nil
}

func (pe *ProgramExpr) Type() string {
	return pe.config.Type
}

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type: "program.expr",
		New: func(config ProcessorConfig) (Processor, error) {
			params := config.Params

			expression, ok := params["expression"]

			if !ok {
				return nil, fmt.Errorf("program.expr requires an expression parameter")
			}

			expressionString, ok := expression.(string)

			if !ok {
				return nil, fmt.Errorf("program.expr expression must be a string")
			}

			program, err := expr.Compile(expressionString)
			if err != nil {
				return nil, err
			}

			return &ProgramExpr{config: config, Program: program}, nil
		},
	})
}
