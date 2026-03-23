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

	if _, ok := ProcessorRegistry[string(processor.Type)]; ok {
		panic(fmt.Sprintf("processor already registered: %s", processor.Type))
	}
	ProcessorRegistry[string(processor.Type)] = processor
}

var (
	processorRegistryMu sync.RWMutex
	ProcessorRegistry   = make(map[string]ProcessorRegistration)
)
