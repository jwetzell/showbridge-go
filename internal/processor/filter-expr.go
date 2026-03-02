package processor

import (
	"context"
	"fmt"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/vm"
	"github.com/jwetzell/showbridge-go/internal/config"
)

// NOTE(jwetzell): see language definition https://expr-lang.org/docs/language-definition
type FilterExpr struct {
	config  config.ProcessorConfig
	Program *vm.Program
}

func (se *FilterExpr) Process(ctx context.Context, payload any) (any, error) {

	output, err := expr.Run(se.Program, payload)
	if err != nil {
		return nil, err
	}

	outputBool, ok := output.(bool)
	if !ok {
		return nil, fmt.Errorf("filter.expr expression did not return a boolean")
	}
	if !outputBool {
		return nil, nil
	}

	return payload, nil
}

func (se *FilterExpr) Type() string {
	return se.config.Type
}

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type: "filter.expr",
		New: func(config config.ProcessorConfig) (Processor, error) {
			params := config.Params

			expressionString, err := params.GetString("expression")
			if err != nil {
				return nil, fmt.Errorf("filter.expr expression error: %w", err)
			}

			program, err := expr.Compile(expressionString)
			if err != nil {
				return nil, err
			}

			return &FilterExpr{config: config, Program: program}, nil
		},
	})
}
