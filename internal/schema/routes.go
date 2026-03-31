package schema

import (
	"encoding/json"

	"github.com/google/jsonschema-go/jsonschema"
)

var RoutesConfigSchema = jsonschema.Schema{
	Schema:      "https://json-schema.org/draft/2020-12/schema",
	ID:          "https://showbridge.io/routes.schema.json",
	Title:       "Routes",
	Description: "route configurations",
	Type:        "array",
	Items: &jsonschema.Schema{
		Type: "object",
		Properties: map[string]*jsonschema.Schema{
			"input": {
				Type:      "string",
				MinLength: jsonschema.Ptr(1),
			},
			"processors": {
				Ref: "https://showbridge.io/processors.schema.json",
			},
		},
		Required:             []string{"input"},
		AdditionalProperties: &jsonschema.Schema{Not: &jsonschema.Schema{}},
	},
	Default: json.RawMessage(`[]`),
}
