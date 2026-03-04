package processor

import (
	"context"
	"fmt"
	"sync"

	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/config"
)

type Processor interface {
	Type() string
	Process(context.Context, any) (any, error)
}

type ProcessorRegistration struct {
	Type string `json:"type"`
	New  func(config.ProcessorConfig) (Processor, error)
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

func GetAnyAs[T any](p any) (T, bool) {
	typed, ok := p.(T)
	return typed, ok
}

type TemplateData struct {
	Payload any
	Modules any
}

func GetTemplateData(ctx context.Context, payload any) TemplateData {
	templateData := TemplateData{Payload: payload}
	modules := ctx.Value(common.ModulesContextKey)
	if modules != nil {
		templateData.Modules = modules
	}
	return templateData
}
