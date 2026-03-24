package schema

import (
	"fmt"
	"net/url"

	"github.com/google/jsonschema-go/jsonschema"
)

func GetResolvedConfigSchema() (*jsonschema.Resolved, error) {
	return ConfigSchema.Resolve(&jsonschema.ResolveOptions{
		Loader: func(uri *url.URL) (*jsonschema.Schema, error) {
			switch uri.String() {
			case "https://showbridge.io/modules.schema.json":
				return GetModulesSchema(), nil
			case "https://showbridge.io/processors.schema.json":
				return GetProcessorsSchema(), nil
			case "https://showbridge.io/routes.schema.json":
				return &RoutesConfigSchema, nil
			default:
				return nil, fmt.Errorf("unknown schema reference: %s", uri.String())
			}
		},
		ValidateDefaults: true,
	})
}
