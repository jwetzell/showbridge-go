package config

type ProcessorConfig struct {
	Type   string `json:"type"`
	Params Params `json:"params,omitempty"`
}
