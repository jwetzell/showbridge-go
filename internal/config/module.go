package config

type ModuleConfig struct {
	Id     string `json:"id"`
	Type   string `json:"type"`
	Params Params `json:"params,omitempty"`
}

type ModuleError struct {
	Index  int          `json:"index"`
	Config ModuleConfig `json:"config"`
	Error  string       `json:"error"`
}
