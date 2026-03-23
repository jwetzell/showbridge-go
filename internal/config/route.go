package config

type RouteConfig struct {
	Input      string            `json:"input"`
	Processors []ProcessorConfig `json:"processors"`
}
