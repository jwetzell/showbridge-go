package schema

import (
	"github.com/google/jsonschema-go/jsonschema"
)

var ConfigSchema = jsonschema.Schema{
	Schema:      "https://json-schema.org/draft/2020-12/schema",
	ID:          "https://showbridge.io/config.schema.json",
	Title:       "Config",
	Description: "showbridge configuration",
	Type:        "object",
	Properties: map[string]*jsonschema.Schema{
		"api": &ApiConfigSchema,
		"modules": {
			Ref: "https://showbridge.io/modules.schema.json",
		},
		"routes": {
			Ref: "https://showbridge.io/routes.schema.json",
		},
	},
}

func ApplyDefaults(cfg *map[string]any) error {
	resolvedSchema, err := GetResolvedConfigSchema()
	if err != nil {
		return err
	}

	return resolvedSchema.ApplyDefaults(cfg)
}

func ValidateConfig(cfg map[string]any) error {
	resolvedSchema, err := GetResolvedConfigSchema()
	if err != nil {
		return err
	}
	return resolvedSchema.Validate(cfg)
}
