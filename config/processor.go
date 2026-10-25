package config

import "github.com/jwetzell/showbridge-go/params"

type ProcessorConfig struct {
	Id     string        `json:"id"`
	Type   string        `json:"type"`
	Params params.Params `json:"params,omitempty"`
}
