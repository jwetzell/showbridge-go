package processor

import (
	"context"
	"fmt"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/vm"
	"github.com/jwetzell/showbridge-go/internal/config"
)

// NOTE(jwetzell): see language definition https://expr-lang.org/docs/language-definition
type ScriptExpr struct {
	config  config.ProcessorConfig
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
		New: func(config config.ProcessorConfig) (Processor, error) {
			params := config.Params

			expressionString, err := params.GetString("expression")
			if err != nil {
				return nil, fmt.Errorf("script.expr expression error: %w", err)
			}

			program, err := expr.Compile(expressionString)
			if err != nil {
				return nil, err
			}

			return &ScriptExpr{config: config, Program: program}, nil
		},
	})
}
