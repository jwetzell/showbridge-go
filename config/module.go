package config

import "github.com/jwetzell/showbridge-go/params"

type ModuleConfig struct {
	Id     string        `json:"id"`
	Type   string        `json:"type"`
	Params params.Params `json:"params,omitempty"`
}

type ModuleError struct {
	Index  int          `json:"index"`
	Config ModuleConfig `json:"config"`
	Error  string       `json:"error"`
}
