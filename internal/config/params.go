package config

import (
	"errors"
	"fmt"

	"github.com/jwetzell/showbridge-go/internal/common"
)

type Params map[string]any

var (
	ErrParamNotFound     = errors.New("not found")
	ErrParamNotString    = errors.New("not a string")
	ErrParamNotNumber    = errors.New("not a number")
	ErrParamNotInteger   = errors.New("not an integer")
	ErrParamNotBool      = errors.New("not a boolean")
	ErrParamNotSlice     = errors.New("not a slice")
	ErrParamNotByteSlice = errors.New("not a byte slice")
	ErrParamNotIntSlice  = errors.New("not an int slice")
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

	intValue, ok := common.GetAnyAsInt(value)
	if ok {
		return intValue, nil
	}

	return 0, ErrParamNotNumber
}

func (p Params) GetFloat32(key string) (float32, error) {
	value, ok := p[key]
	if !ok {
		return 0, ErrParamNotFound
	}

	floatValue, ok := common.GetAnyAsFloat32(value)
	if ok {
		return floatValue, nil
	}

	return 0, ErrParamNotNumber
}

func (p Params) GetFloat64(key string) (float64, error) {
	value, ok := p[key]
	if !ok {
		return 0, ErrParamNotFound
	}

	floatValue, ok := common.GetAnyAsFloat64(value)
	if ok {
		return floatValue, nil
	}

	return 0, ErrParamNotNumber
}

func (p Params) GetBool(key string) (bool, error) {
	value, ok := p[key]
	if !ok {
		return false, ErrParamNotFound
	}

	boolValue, ok := value.(bool)
	if !ok {
		return false, ErrParamNotBool
	}
	return boolValue, nil
}

func (p Params) GetStringSlice(key string) ([]string, error) {
	value, ok := p[key]
	if !ok {
		return nil, ErrParamNotFound
	}

	interfaceSlice, ok := value.([]any)
	if !ok {
		return nil, ErrParamNotSlice
	}

	stringSlice := make([]string, len(interfaceSlice))
	for i, v := range interfaceSlice {
		str, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("element at index %d is not a string", i)
		}
		stringSlice[i] = str
	}
	return stringSlice, nil
}

func (p Params) GetIntSlice(key string) ([]int, error) {
	value, ok := p[key]
	if !ok {
		return nil, ErrParamNotFound
	}

	intSlice, ok := common.GetAnyAsIntSlice(value)
	if !ok {
		return nil, ErrParamNotIntSlice
	}

	return intSlice, nil
}

func (p Params) GetByteSlice(key string) ([]byte, error) {
	value, ok := p[key]
	if !ok {
		return nil, ErrParamNotFound
	}

	byteSlice, ok := common.GetAnyAsByteSlice(value)

	if !ok {
		return nil, ErrParamNotByteSlice
	}

	return byteSlice, nil
}
