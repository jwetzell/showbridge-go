package config

type ModuleConfig struct {
	Id     string `json:"id"`
	Type   string `json:"type"`
	Params Params `json:"params,omitempty"`
}
