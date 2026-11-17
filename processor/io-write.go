package processor

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/jwetzell/showbridge-go/config"
	"github.com/jwetzell/showbridge-go/internal/types"
	"github.com/jwetzell/showbridge-go/routing"
)

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type:        "io.write",
		Title:       "Write IO",
		Description: "Sends the incoming payload to the specified module output",
		ParamsSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"module": {
					Title:       "Module ID",
					Description: "ID of module to send output to",
					Type:        "string",
				},
			},
			Required:             []string{"module"},
			AdditionalProperties: &jsonschema.Schema{Not: &jsonschema.Schema{}},
		},
		New: func(config config.ProcessorConfig) (routing.Processor, error) {

			processorParams := config.Params

			moduleId, err := processorParams.GetString("module")

			if err != nil {
				return nil, fmt.Errorf("io.write module error: %w", err)
			}

			return &IOWrite{config: config, ModuleId: moduleId, logger: slog.Default().With("component", "processor", "type", config.Type)}, nil
		},
	})
}

type IOWrite struct {
	config   config.ProcessorConfig
	ModuleId string
	logger   *slog.Logger
	module   routing.WriterModule
	mu       sync.Mutex
}

func (mo *IOWrite) Process(next routing.ProcessFunc) routing.ProcessFunc {
	return func(ctx context.Context, wrappedPayload routing.WrappedPayload, moduleManager routing.ModuleManager, router routing.Router) (routing.WrappedPayload, error) {

		mo.mu.Lock()
		defer mo.mu.Unlock()

		if mo.module == nil {
			if moduleManager == nil {
				return wrappedPayload, errors.New("io.write module manager is nil")
			}

			module, ok := moduleManager.GetModule(mo.ModuleId)
			if !ok {
				return wrappedPayload, fmt.Errorf("io.write unable to find module with id: %s", mo.ModuleId)
			}

			writerModule, ok := module.(routing.WriterModule)
			if !ok {
				return wrappedPayload, fmt.Errorf("io.write module with id %s is not a WriterModule", mo.ModuleId)
			}
			mo.module = writerModule
		}

		data, ok := types.GetAnyAsByteSlice(wrappedPayload.Payload)
		if !ok {
			return wrappedPayload, fmt.Errorf("io.write failed to convert payload to []byte")
		}

		err := mo.module.Write(ctx, data)

		if err != nil {
			return wrappedPayload, fmt.Errorf("io.write failed to send output: %w", err)
		}

		return next(ctx, wrappedPayload, moduleManager, router)
	}
}
