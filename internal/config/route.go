package config

type RouteConfig struct {
	Id         string            `json:"id"`
	Input      string            `json:"input"`
	Processors []ProcessorConfig `json:"processors"`
}

type RouteError struct {
	Index  int         `json:"index"`
	Config RouteConfig `json:"config"`
	Error  string      `json:"error"`
}
