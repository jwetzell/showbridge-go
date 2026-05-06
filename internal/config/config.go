package config

type Config struct {
	Api     ApiConfig      `json:"api"`
	Modules []ModuleConfig `json:"modules"`
	Routes  []RouteConfig  `json:"routes"`
}

type Configurable interface {
	UpdateConfig(newConfig Config, triggerChangeChannel bool) (error, []ModuleError, []RouteError)
	GetRunningConfig() Config
}
