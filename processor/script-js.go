package processor

import (
	"context"
	"fmt"

	"github.com/dop251/goja"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/jwetzell/showbridge-go/config"
	"github.com/jwetzell/showbridge-go/internal/types"
	"github.com/jwetzell/showbridge-go/routing"
)

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type:        "script.js",
		Title:       "Run JavaScript",
		Description: "Runs a JavaScript program with access to the incoming payload",
		ParamsSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"program": {
					Title:       "Program",
					Description: "JavaScript program to run",
					Type:        "string",
				},
			},
			Required:             []string{"program"},
			AdditionalProperties: &jsonschema.Schema{Not: &jsonschema.Schema{}},
		},
		New: func(config config.ProcessorConfig) (routing.Processor, error) {
			processorParams := config.Params

			programString, err := processorParams.GetString("program")
			if err != nil {
				return nil, fmt.Errorf("script.js program error: %w", err)
			}

			vm := goja.New()

			return &ScriptJS{config: config, Program: programString, vm: vm}, nil
		},
	})
}

type ScriptJS struct {
	config  config.ProcessorConfig
	vm      *goja.Runtime
	Program string
}

func (sj *ScriptJS) Process(next routing.ProcessFunc) routing.ProcessFunc {
	return func(ctx context.Context, wrappedPayload routing.WrappedPayload, moduleManager routing.ModuleManager, router routing.Router) (routing.WrappedPayload, error) {

		// NOTE(jwetzell): some weird conversion going on with these types
		_, isUint8Slice := wrappedPayload.Payload.([]uint8)
		_, isByteSlice := wrappedPayload.Payload.([]byte)

		if isUint8Slice || isByteSlice {
			intSlice, ok := types.GetAnyAsIntSlice(wrappedPayload.Payload)

			if ok {
				wrappedPayload.Payload = intSlice
			}
		}

		err := sj.vm.Set("payload", wrappedPayload.Payload)
		if err != nil {
			return wrappedPayload, err
		}

		err = sj.vm.Set("source", wrappedPayload.Source)
		if err != nil {
			return wrappedPayload, err
		}

		err = sj.vm.Set("kv", kvGetter{
			ctx: ctx,
			mm:  moduleManager,
		})
		if err != nil {
			return wrappedPayload, err
		}

		_, err = sj.vm.RunString(sj.Program)

		if err != nil {
			return wrappedPayload, err
		}

		jsPayload := sj.vm.Get("payload").Export()
		// TODO(jwetzell): handle any potential post-processing for specific "troublesome" types

		wrappedPayload.Payload = jsPayload
		return next(ctx, wrappedPayload, moduleManager, router)
	}
}
