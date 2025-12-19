package processor

import (
	"context"
	"fmt"

	extism "github.com/extism/go-sdk"
	"github.com/jwetzell/showbridge-go/internal/config"
)

type ScriptWASM struct {
	config  config.ProcessorConfig
	Program *extism.CompiledPlugin
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

	_, output, err := program.Call("process", payloadBytes)

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
		New: func(config config.ProcessorConfig) (Processor, error) {
			params := config.Params

			path, ok := params["path"]

			if !ok {
				return nil, fmt.Errorf("script.wasm requires a path parameter")
			}

			pathString, ok := path.(string)

			if !ok {
				return nil, fmt.Errorf("script.wasm path must be a string")
			}

			manifest := extism.Manifest{
				Wasm: []extism.Wasm{
					extism.WasmFile{
						Path: pathString,
					},
				},
			}

			program, err := extism.NewCompiledPlugin(context.Background(), manifest, extism.PluginConfig{}, []extism.HostFunction{})

			if err != nil {
				return nil, err
			}

			return &ScriptWASM{config: config, Program: program}, nil
		},
	})
}
