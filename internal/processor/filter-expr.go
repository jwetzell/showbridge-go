package processor

import (
	"context"
	"fmt"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/vm"
	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/config"
)

// NOTE(jwetzell): see language definition https://expr-lang.org/docs/language-definition
type FilterExpr struct {
	config  config.ProcessorConfig
	Program *vm.Program
}

func (fe *FilterExpr) Process(ctx context.Context, wrappedPayload common.WrappedPayload) (common.WrappedPayload, error) {

	exprEnv := wrappedPayload

	output, err := expr.Run(fe.Program, exprEnv)
	if err != nil {
		wrappedPayload.End = true
		return wrappedPayload, err
	}

	outputBool, ok := output.(bool)
	if !ok {
		wrappedPayload.End = true
		return wrappedPayload, fmt.Errorf("filter.expr expression did not return a boolean")
	}
	if !outputBool {
		wrappedPayload.End = true
		return wrappedPayload, nil
	}

	wrappedPayload.Payload = exprEnv.Payload

	return wrappedPayload, nil
}

func (fe *FilterExpr) Type() string {
	return fe.config.Type
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
