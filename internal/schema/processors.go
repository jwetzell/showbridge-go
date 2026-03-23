package schema

import (
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/jwetzell/showbridge-go/internal/processor"
)

func GetProcessorsSchema() *jsonschema.Schema {

	schema := &jsonschema.Schema{
		Schema:      "https://json-schema.org/draft/2020-12/schema",
		ID:          "https://showbridge.io/processors.schema.json",
		Title:       "Processors",
		Description: "processor configurations",
		Type:        "array",
	}

	processorDefinitionSchemas := []*jsonschema.Schema{}
	for _, proc := range processor.ProcessorRegistry {
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
