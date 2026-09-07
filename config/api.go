package config

type ApiConfig struct {
	Enabled bool `json:"enabled"`
	Port    int  `json:"port"`
}
