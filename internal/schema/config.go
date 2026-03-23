package schema

import (
	"encoding/json"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/jwetzell/showbridge-go/internal/config"
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

func ValidateConfig(config config.Config) error {
	resolvedSchema, err := GetResolvedConfigSchema()
	if err != nil {
		return err
	}

	jsonBytes, err := json.Marshal(config)
	if err != nil {
		return err
	}

	jsonMap := make(map[string]any)
	err = json.Unmarshal(jsonBytes, &jsonMap)
	if err != nil {
		return err
	}

	return resolvedSchema.Validate(jsonMap)
}
