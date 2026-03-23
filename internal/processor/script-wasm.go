package processor

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	extism "github.com/extism/go-sdk"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/config"
)

type ScriptWASM struct {
	config   config.ProcessorConfig
	Program  *extism.CompiledPlugin
	Function string
}

func (sw *ScriptWASM) Process(ctx context.Context, wrappedPayload common.WrappedPayload) (common.WrappedPayload, error) {

	payload := wrappedPayload.Payload
	payloadBytes, ok := common.GetAnyAsByteSlice(payload)

	if !ok {
		wrappedPayload.End = true
		return wrappedPayload, fmt.Errorf("script.wasm can only process a byte array")
	}

	program, err := sw.Program.Instance(ctx, extism.PluginInstanceConfig{})

	if err != nil {
		wrappedPayload.End = true
		return wrappedPayload, err
	}

	_, output, err := program.Call(sw.Function, payloadBytes)

	if err != nil {
		wrappedPayload.End = true
		return wrappedPayload, err
	}
	wrappedPayload.Payload = output

	return wrappedPayload, nil
}

func (sw *ScriptWASM) Type() string {
	return sw.config.Type
}

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type:  "script.wasm",
		Title: "Run WASM Plugin",
		ParamsSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"path": {
					Title: "Path",
					Type:  "string",
				},
				"function": {
					Title:   "Function",
					Type:    "string",
					Default: json.RawMessage(`"process"`),
				},
				"enableWasi": {
					Title:   "Enable WASI",
					Type:    "boolean",
					Default: json.RawMessage("false"),
				},
			},
			Required:             []string{"path"},
			AdditionalProperties: nil,
		},
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
