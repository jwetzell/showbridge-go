package config

type Config struct {
	Api     ApiConfig      `json:"api"`
	Modules []ModuleConfig `json:"modules"`
	Routes  []RouteConfig  `json:"routes"`
}
