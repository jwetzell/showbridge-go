package config

import "uuid"

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
	if config.Modules == nil {
		config.Modules = []ModuleConfig{}
	}
	if config.Routes == nil {
		config.Routes = []RouteConfig{}
	}

	for routeIndex := range config.Routes {
		if config.Routes[routeIndex].Id == "" {
			config.Routes[routeIndex].Id = uuid.NewV4().String()
		}

		for processorIndex := range config.Routes[routeIndex].Processors {
			if config.Routes[routeIndex].Processors[processorIndex].Id == "" {
				config.Routes[routeIndex].Processors[processorIndex].Id = uuid.NewV4().String()
			}
		}
	}
	return config
}
