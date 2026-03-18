package config

type Config struct {
	Api     ApiConfig      `json:"api"`
	Modules []ModuleConfig `json:"modules"`
	Routes  []RouteConfig  `json:"routes"`
}

type ApiConfig struct {
	Enabled bool `json:"enabled"`
	Port    int  `json:"port"`
}
type ModuleConfig struct {
	Id     string `json:"id"`
	Type   string `json:"type"`
	Params Params `json:"params,omitempty"`
}

type RouteConfig struct {
	Input      string            `json:"input"`
	Processors []ProcessorConfig `json:"processors"`
}

type ProcessorConfig struct {
	Type   string `json:"type"`
	Params Params `json:"params,omitempty"`
}
