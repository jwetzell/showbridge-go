package processing

import (
	"context"
	"fmt"
	"sync"
)

type Processor interface {
	Type() string
	Process(context.Context, any) (any, error)
}

type ProcessorConfig struct {
	Type   string         `json:"type"`
	Params map[string]any `json:"params"`
}

type ProcessorRegistration struct {
	Type string `json:"type"`
	New  func(ProcessorConfig) (Processor, error)
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
