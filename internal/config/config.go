package config

import "github.com/jwetzell/showbridge-go/internal/processing"

type Config struct {
	Modules []ModuleConfig `json:"modules"`
	Routes  []RouteConfig  `json:"routes"`
}

type ModuleConfig struct {
	Id     string         `json:"id"`
	Type   string         `json:"type"`
	Params map[string]any `json:"params"`
}

type RouteConfig struct {
	Input      string                       `json:"input"`
	Processors []processing.ProcessorConfig `json:"processors"`
	Output     string                       `json:"output"`
}
