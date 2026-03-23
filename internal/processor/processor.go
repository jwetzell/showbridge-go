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

func GetProcessorsSchema() *jsonschema.Schema {
	processorRegistryMu.RLock()
	defer processorRegistryMu.RUnlock()

	schema := &jsonschema.Schema{
		Schema:      "https://json-schema.org/draft/2020-12/schema",
		ID:          "https://showbridge.io/processors.schema.json",
		Title:       "Processors",
		Description: "processor configurations",
		Type:        "array",
	}

	processorDefinitionSchemas := []*jsonschema.Schema{}
	for _, proc := range ProcessorRegistry {
		processorSchema := &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"type": {
					Const: jsonschema.Ptr[any](proc.Type),
				},
			},
			Required:             []string{"type"},
			AdditionalProperties: nil,
		}
		if proc.Title != "" {
			processorSchema.Title = proc.Title
		}
		if proc.Description != "" {
			processorSchema.Description = proc.Description
		}
		if proc.ParamsSchema != nil {
			processorSchema.Properties["params"] = proc.ParamsSchema
			processorSchema.Required = append(processorSchema.Required, "params")
		}
		processorDefinitionSchemas = append(processorDefinitionSchemas, processorSchema)
	}
	schema.Items = &jsonschema.Schema{
		OneOf: processorDefinitionSchemas,
	}
	return schema
}
