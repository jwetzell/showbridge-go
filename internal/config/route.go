package config

type RouteConfig struct {
	Input      string            `json:"input"`
	Processors []ProcessorConfig `json:"processors"`
}

type RouteError struct {
	Index  int         `json:"index"`
	Config RouteConfig `json:"config"`
	Error  string      `json:"error"`
}
