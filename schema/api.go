package schema

import (
	"encoding/json"

	"github.com/google/jsonschema-go/jsonschema"
)

var ApiConfigSchema = jsonschema.Schema{
	ID:   "https://showbridge.io/api.schema.json",
	Type: "object",
	Properties: map[string]*jsonschema.Schema{
		"enabled": {
			Type:        "boolean",
			Description: "Whether the API server is enabled",
			Default:     json.RawMessage(`false`),
		},
		"port": {
			Type:        "integer",
			Description: "Port for the API server to listen on",
			Minimum:     jsonschema.Ptr[float64](1024),
			Maximum:     jsonschema.Ptr[float64](65535),
			Default:     json.RawMessage(`8080`),
		},
	},
	Required:             []string{"port"},
	Default:              json.RawMessage(`{"enabled": false, "port": 8080}`),
	AdditionalProperties: &jsonschema.Schema{Not: &jsonschema.Schema{}},
}
