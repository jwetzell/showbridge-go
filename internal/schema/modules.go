package schema

import (
	"encoding/json"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/jwetzell/showbridge-go/internal/module"
)

func GetModulesSchema() *jsonschema.Schema {

	schema := &jsonschema.Schema{
		Schema:      "https://json-schema.org/draft/2020-12/schema",
		ID:          "https://showbridge.io/modules.schema.json",
		Title:       "Modules",
		Description: "module configurations",
		Type:        "array",
		Default:     json.RawMessage(`[]`),
	}

	moduleDefinitionSchemas := []*jsonschema.Schema{}
	for _, mod := range module.ModuleRegistry {
		moduleSchema := &jsonschema.Schema{
			ID:   mod.Type,
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
			AdditionalProperties: &jsonschema.Schema{Not: &jsonschema.Schema{}},
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
