package config

import (
	"encoding/json"

	"github.com/google/jsonschema-go/jsonschema"
)

type ApiConfig struct {
	Enabled bool `json:"enabled"`
	Port    int  `json:"port"`
}

var ApiConfigSchema = jsonschema.Schema{
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
		},
	},
	Required: []string{"port"},
}
