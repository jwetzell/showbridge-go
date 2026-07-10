package processor

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/config"
)

type Processor interface {
	Type() string
	Process(context.Context, common.WrappedPayload) (common.WrappedPayload, error)
}

type ProcessorRegistration struct {
	Type         string             `json:"type"`
	Title        string             `json:"title,omitempty"`
	Description  string             `json:"description,omitempty"`
	ParamsSchema *jsonschema.Schema `json:"paramsSchema,omitempty"`
	New          func(config.ProcessorConfig) (Processor, error)
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
