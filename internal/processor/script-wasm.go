package processor

import (
	"context"
	"errors"
	"fmt"

	extism "github.com/extism/go-sdk"
	"github.com/jwetzell/showbridge-go/internal/config"
)

type ScriptWASM struct {
	config   config.ProcessorConfig
	Program  *extism.CompiledPlugin
	Function string
}

func (se *ScriptWASM) Process(ctx context.Context, payload any) (any, error) {

	payloadBytes, ok := payload.([]byte)

	if !ok {
		return nil, fmt.Errorf("script.wasm can only operator on byte array")
	}

	program, err := se.Program.Instance(ctx, extism.PluginInstanceConfig{})

	if err != nil {
		return nil, err
	}

	_, output, err := program.Call(se.Function, payloadBytes)

	if err != nil {
		return nil, err
	}

	return output, nil
}

func (se *ScriptWASM) Type() string {
	return se.config.Type
}

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type: "script.wasm",
		New: func(processorConfig config.ProcessorConfig) (Processor, error) {
			params := processorConfig.Params

			pathString, err := params.GetString("path")
			if err != nil {
				return nil, fmt.Errorf("script.wasm path error: %w", err)
			}

			functionString, err := params.GetString("function")
			if err != nil {
				if errors.Is(err, config.ErrParamNotFound) {
					functionString = "process"
				} else {
					return nil, fmt.Errorf("script.wasm function error: %w", err)
				}
			}

			enableWasiBool, err := params.GetBool("enableWasi")
			if err != nil {
				if errors.Is(err, config.ErrParamNotFound) {
					enableWasiBool = false
				} else {
					return nil, fmt.Errorf("script.wasm enableWasi error: %w", err)
				}
			}

			manifest := extism.Manifest{
				Wasm: []extism.Wasm{
					extism.WasmFile{
						Path: pathString,
					},
				},
			}

			program, err := extism.NewCompiledPlugin(context.Background(), manifest, extism.PluginConfig{
				EnableWasi: enableWasiBool,
			}, []extism.HostFunction{})

			if err != nil {
				return nil, err
			}

			return &ScriptWASM{config: processorConfig, Program: program, Function: functionString}, nil
		},
	})
}
