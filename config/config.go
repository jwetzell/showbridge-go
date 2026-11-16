package config

type Config struct {
	Api     ApiConfig      `json:"api"`
	Modules []ModuleConfig `json:"modules"`
	Routes  []RouteConfig  `json:"routes"`
}

type Configurable interface {
	UpdateConfig(newConfig Config, triggerChangeChannel bool) ([]ModuleError, []RouteError, error)
	GetRunningConfig() Config
}

func CleanConfig(config Config) Config {
	// NOTE(jwetzell): this should be done by the schema apply defaults but there is a bug in the jsonschema library
	if config.Modules == nil {
		config.Modules = []ModuleConfig{}
	}
	if config.Routes == nil {
		config.Routes = []RouteConfig{}
	}
	return config
}
