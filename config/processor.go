package config

type ProcessorConfig struct {
	Id     string `json:"id,omitempty"`
	Type   string `json:"type"`
	Params Params `json:"params,omitempty"`
}
