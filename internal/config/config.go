package config

import (
	"errors"
)

type Config struct {
	Modules []ModuleConfig `json:"modules"`
	Routes  []RouteConfig  `json:"routes"`
}

type Params map[string]any

var (
	ErrParamNotFound  = errors.New("not found")
	ErrParamNotString = errors.New("not a string")
	ErrParamNotNumber = errors.New("not a number")
)

func (p Params) GetString(key string) (string, error) {
	value, ok := p[key]
	if !ok {
		return "", ErrParamNotFound
	}

	stringValue, ok := value.(string)
	if !ok {
		return "", ErrParamNotString
	}
	return stringValue, nil
}

func (p Params) GetInt(key string) (int, error) {
	value, ok := p[key]
	if !ok {
		return 0, ErrParamNotFound
	}

	intValue, ok := value.(int)
	if !ok {
		floatValue, ok := value.(float64)
		if !ok {
			return 0, ErrParamNotNumber
		}
		intValue = int(floatValue)
	}
	return intValue, nil
}

func (p Params) GetBool(key string) (bool, error) {
	value, ok := p[key]
	if !ok {
		return false, ErrParamNotFound
	}

	boolValue, ok := value.(bool)
	if !ok {
		return false, errors.New("not a boolean")
	}
	return boolValue, nil
}

type ModuleConfig struct {
	Id     string `json:"id"`
	Type   string `json:"type"`
	Params Params `json:"params"`
}

type RouteConfig struct {
	Input      string            `json:"input"`
	Processors []ProcessorConfig `json:"processors"`
	Output     string            `json:"output"`
}

type ProcessorConfig struct {
	Type   string `json:"type"`
	Params Params `json:"params"`
}
