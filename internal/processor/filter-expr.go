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

type PayloadStruct struct {
	Payload any
}

func (se *FilterExpr) Process(ctx context.Context, payload any) (any, error) {
	payloadType := fmt.Sprintf("%T", payload)

	exprEnv := payload

	switch payloadType {
	case "uint", "uint8", "uint16", "uint32", "uint64":
		exprEnv = PayloadStruct{Payload: payload}
	case "int", "int8", "int16", "int32", "int64":
		exprEnv = PayloadStruct{Payload: payload}
	case "float32", "float64":
		exprEnv = PayloadStruct{Payload: payload}
	case "string":
		exprEnv = PayloadStruct{Payload: payload}
	case "bool":
		exprEnv = PayloadStruct{Payload: payload}
	}

	output, err := expr.Run(se.Program, exprEnv)
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
