package processor

import (
	"context"
	"fmt"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/vm"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/config"
)

// NOTE(jwetzell): see language definition https://expr-lang.org/docs/language-definition
type ScriptExpr struct {
	config  config.ProcessorConfig
	Program *vm.Program
}

func (se *ScriptExpr) Process(ctx context.Context, wrappedPayload common.WrappedPayload) (common.WrappedPayload, error) {

	exprEnv := wrappedPayload

	output, err := expr.Run(se.Program, exprEnv)
	if err != nil {
		wrappedPayload.End = true
		return wrappedPayload, err
	}

	wrappedPayload.Payload = output
	return wrappedPayload, nil
}

func (se *ScriptExpr) Type() string {
	return se.config.Type
}

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type:  "script.expr",
		Title: "Evaluate Expr Expression",
		ParamsSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"expression": {
					Title: "Expression",
					Type:  "string",
				},
			},
			Required:             []string{"expression"},
			AdditionalProperties: nil,
		},
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
