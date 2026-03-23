package config

import (
	"github.com/google/jsonschema-go/jsonschema"
)

type Config struct {
	Api     ApiConfig      `json:"api"`
	Modules []ModuleConfig `json:"modules"`
	Routes  []RouteConfig  `json:"routes"`
}

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
