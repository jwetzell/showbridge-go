package processor

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/jwetzell/showbridge-go/config"
	"github.com/jwetzell/showbridge-go/routing"
)

type ProcessorRegistration struct {
	Type         string             `json:"type"`
	Title        string             `json:"title,omitempty"`
	Description  string             `json:"description,omitempty"`
	ParamsSchema *jsonschema.Schema `json:"paramsSchema,omitempty"`
	New          func(config.ProcessorConfig) (routing.Processor, error)
}

func RegisterProcessor(processor ProcessorRegistration) {

	if processor.Type == "" {
		panic("processor type is missing")
	}
	if processor.New == nil {
		panic("missing ProcessorRegistration.New")
	}

	processorRegistryMu.Lock()
	defer processorRegistryMu.Unlock()

	_, exists := processorRegistry[string(processor.Type)]
	if exists {
		panic(fmt.Sprintf("processor already registered: %s", processor.Type))
	}
	processorRegistry[string(processor.Type)] = processor
}

type ProcessorRegistry map[string]ProcessorRegistration

func GetProcessorRegistration(processorType string) (ProcessorRegistration, bool) {
	processorRegistryMu.RLock()
	defer processorRegistryMu.RUnlock()
	processor, ok := processorRegistry[processorType]
	return processor, ok
}

func GetProcessorRegistrations() []ProcessorRegistration {
	processorRegistryMu.RLock()
	defer processorRegistryMu.RUnlock()

	registrations := make([]ProcessorRegistration, 0, len(processorRegistry))
	for _, processor := range processorRegistry {
		registrations = append(registrations, processor)
	}
	return registrations
}

var (
	processorRegistryMu sync.RWMutex
	processorRegistry   = make(map[string]ProcessorRegistration)
)

type kvGetter struct {
	ctx context.Context
	mm  routing.ModuleManager
}

type templateData struct {
	routing.WrappedPayload
	KV kvGetter
}

func (kvg kvGetter) Get(store string, key string) (any, error) {
	module, ok := kvg.mm.GetModule(store)
	if !ok {
		return nil, fmt.Errorf("module not found: %s", store)
	}
	kvModule, ok := module.(routing.KeyValueModule)
	if !ok {
		return nil, fmt.Errorf("module is not a key-value module: %s", store)
	}
	value, err := kvModule.Get(kvg.ctx, key)
	if err != nil {
		return nil, err
	}
	return value, nil
}

func newTemplateData(ctx context.Context, payload routing.WrappedPayload, mm routing.ModuleManager) templateData {
	return templateData{
		WrappedPayload: payload,
		KV: kvGetter{
			ctx: ctx,
			mm:  mm,
		},
	}
}
