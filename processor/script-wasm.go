package processor

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	extism "github.com/extism/go-sdk"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/jwetzell/showbridge-go/config"
	"github.com/jwetzell/showbridge-go/params"
	"github.com/jwetzell/showbridge-go/routing"
)

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type:        "script.wasm",
		Title:       "Run WASM Plugin",
		Description: "Runs an Extism WASM plugin with the specified function and options",
		ParamsSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"path": {
					Title:       "Path",
					Description: "path to the WASM plugin to load",
					Type:        "string",
				},
				"function": {
					Title:       "Function",
					Description: "exported function to call on the WASM plugin",
					Type:        "string",
					Default:     json.RawMessage(`"process"`),
				},
				"enableWasi": {
					Title:       "Enable WASI",
					Description: "whether to enable WASI when running the WASM plugin",
					Type:        "boolean",
					Default:     json.RawMessage("false"),
				},
			},
			Required:             []string{"path"},
			AdditionalProperties: &jsonschema.Schema{Not: &jsonschema.Schema{}},
		},
		New: func(processorConfig config.ProcessorConfig) (routing.Processor, error) {
			processorParams := processorConfig.Params

			pathString, err := processorParams.GetString("path")
			if err != nil {
				return nil, fmt.Errorf("script.wasm path error: %w", err)
			}

			functionString, err := processorParams.GetString("function")
			if err != nil {
				if errors.Is(err, params.ErrParamNotFound) {
					functionString = "process"
				} else {
					return nil, fmt.Errorf("script.wasm function error: %w", err)
				}
			}

			enableWasiBool, err := processorParams.GetBool("enableWasi")
			if err != nil {
				if errors.Is(err, params.ErrParamNotFound) {
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
				if errors.Is(err, os.ErrNotExist) {
					return nil, fmt.Errorf("script.wasm file not found: %s", pathString)
				}
				return nil, err
			}

			programInstance, err := program.Instance(context.Background(), extism.PluginInstanceConfig{})

			if err != nil {
				return nil, err
			}

			return &ScriptWASM{config: processorConfig, Program: programInstance, Function: functionString}, nil
		},
	})
}

type ScriptWASM struct {
	config   config.ProcessorConfig
	Program  *extism.Plugin
	Function string
}

func (sw *ScriptWASM) Process(ctx context.Context, wrappedPayload routing.WrappedPayload, moduleManager routing.ModuleManager, router routing.Router) (routing.WrappedPayload, error) {

	payload := wrappedPayload.Payload
	payloadBytes, ok := params.GetAnyAsByteSlice(payload)

	if !ok {
		wrappedPayload.End = true
		return wrappedPayload, fmt.Errorf("script.wasm can only process a byte array")
	}

	_, output, err := sw.Program.Call(sw.Function, payloadBytes)

	if err != nil {
		wrappedPayload.End = true
		return wrappedPayload, err
	}
	wrappedPayload.Payload = output

	return wrappedPayload, nil
}
