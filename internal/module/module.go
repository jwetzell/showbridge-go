package module

import (
	"fmt"
	"log/slog"
	"sync"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/config"
)

type ModuleError struct {
	Index  int                 `json:"index"`
	Config config.ModuleConfig `json:"config"`
	Error  string              `json:"error"`
}

type ModuleRegistration struct {
	Type         string             `json:"type"`
	Title        string             `json:"title,omitempty"`
	Description  string             `json:"description,omitempty"`
	ParamsSchema *jsonschema.Schema `json:"paramsSchema,omitempty"`
	New          func(config.ModuleConfig) (common.Module, error)
}

func RegisterModule(mod ModuleRegistration) {

	if mod.Type == "" {
		panic("module type is missing")
	}
	if mod.New == nil {
		panic("missing ModuleInfo.New")
	}

	moduleRegistryMu.Lock()
	defer moduleRegistryMu.Unlock()

	if _, ok := ModuleRegistry[string(mod.Type)]; ok {
		panic(fmt.Sprintf("module already registered: %s", mod.Type))
	}
	ModuleRegistry[string(mod.Type)] = mod
}

var (
	moduleRegistryMu sync.RWMutex
	ModuleRegistry   = make(map[string]ModuleRegistration)
)

func CreateLogger(config config.ModuleConfig) *slog.Logger {
	return slog.Default().With("component", "module", "id", config.Id, "type", config.Type)
}

func GetModulesSchema() *jsonschema.Schema {
	moduleRegistryMu.RLock()
	defer moduleRegistryMu.RUnlock()

	schema := &jsonschema.Schema{
		Schema:      "https://json-schema.org/draft/2020-12/schema",
		ID:          "https://showbridge.io/modules.schema.json",
		Title:       "Modules",
		Description: "module configurations",
		Type:        "array",
	}

	moduleDefinitionSchemas := []*jsonschema.Schema{}
	for _, mod := range ModuleRegistry {
		moduleSchema := &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"id": {
					Type:      "string",
					MinLength: jsonschema.Ptr(1),
				},
				"type": {
					Const: jsonschema.Ptr[any](mod.Type),
				},
			},
			Required:             []string{"id", "type"},
			AdditionalProperties: nil,
		}
		if mod.Title != "" {
			moduleSchema.Title = mod.Title
		}
		if mod.Description != "" {
			moduleSchema.Description = mod.Description
		}
		if mod.ParamsSchema != nil {
			moduleSchema.Properties["params"] = mod.ParamsSchema
			moduleSchema.Required = append(moduleSchema.Required, "params")
		}
		moduleDefinitionSchemas = append(moduleDefinitionSchemas, moduleSchema)
	}
	schema.Items = &jsonschema.Schema{
		OneOf: moduleDefinitionSchemas,
	}
	return schema
}
