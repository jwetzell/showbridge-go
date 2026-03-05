package config

import (
	"errors"
	"fmt"
	"math"
	"reflect"
)

type Config struct {
	Modules []ModuleConfig `json:"modules"`
	Routes  []RouteConfig  `json:"routes"`
}

type Params map[string]any

var (
	ErrParamNotFound   = errors.New("not found")
	ErrParamNotString  = errors.New("not a string")
	ErrParamNotNumber  = errors.New("not a number")
	ErrParamNotInteger = errors.New("not an integer")
	ErrParamNotBool    = errors.New("not a boolean")
	ErrParamNotSlice   = errors.New("not a slice")
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
	if ok {
		return intValue, nil
	}

	uintValue, ok := value.(uint)
	if ok {
		return int(uintValue), nil
	}

	byteValue, ok := value.(byte)
	if ok {
		return int(byteValue), nil
	}

	floatValue, ok := value.(float64)
	if ok {
		if floatValue != math.Floor(floatValue) {
			return 0, ErrParamNotInteger
		}
		return int(floatValue), nil
	}

	return 0, ErrParamNotNumber
}

func (p Params) GetFloat32(key string) (float32, error) {
	value, ok := p[key]
	if !ok {
		return 0, ErrParamNotFound
	}

	float32Value, ok := value.(float32)
	if ok {
		return float32Value, nil
	}

	float64Value, ok := value.(float64)
	if ok {
		return float32(float64Value), nil
	}

	intValue, ok := value.(int)
	if ok {
		return float32(intValue), nil
	}

	uintValue, ok := value.(uint)
	if ok {
		return float32(uintValue), nil
	}

	byteValue, ok := value.(byte)
	if ok {
		return float32(byteValue), nil
	}

	return 0, ErrParamNotNumber
}

func (p Params) GetFloat64(key string) (float64, error) {
	value, ok := p[key]
	if !ok {
		return 0, ErrParamNotFound
	}

	float64Value, ok := value.(float64)
	if ok {
		return float64Value, nil
	}

	float32Value, ok := value.(float32)
	if ok {
		return float64(float32Value), nil
	}

	intValue, ok := value.(int)
	if ok {
		return float64(intValue), nil
	}

	uintValue, ok := value.(uint)
	if ok {
		return float64(uintValue), nil
	}

	byteValue, ok := value.(byte)
	if ok {
		return float64(byteValue), nil
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

	v := reflect.ValueOf(value)
	if v.Kind() != reflect.Slice {
		return nil, ErrParamNotSlice
	}

	result := make([]int, v.Len())
	for i := 0; i < v.Len(); i++ {
		elem := v.Index(i).Interface()
		byteValue, ok := elem.(byte)
		if ok {
			result[i] = int(byteValue)
			continue
		}
		uintValue, ok := elem.(uint)
		if ok {
			result[i] = int(uintValue)
			continue
		}
		intValue, ok := elem.(int)
		if ok {
			result[i] = int(intValue)
			continue
		}
		floatValue, ok := elem.(float64)
		if ok {
			if floatValue != math.Floor(floatValue) {
				return nil, fmt.Errorf("element at index %d is not an integer", i)
			}
			result[i] = int(floatValue)
			continue
		}
		return nil, fmt.Errorf("element at index %d is not a number", i)
	}
	return result, nil
}

func (p Params) GetByteSlice(key string) ([]byte, error) {
	value, ok := p[key]
	if !ok {
		return nil, ErrParamNotFound
	}

	v := reflect.ValueOf(value)
	if v.Kind() != reflect.Slice {
		return nil, ErrParamNotSlice
	}

	result := make([]byte, v.Len())
	for i := 0; i < v.Len(); i++ {
		elem := v.Index(i).Interface()
		byteValue, ok := elem.(byte)
		if ok {
			result[i] = byteValue
			continue
		}
		uintValue, ok := elem.(uint)
		if ok {
			if uintValue > 255 {
				return nil, fmt.Errorf("element at index %d is out of byte range", i)
			}
			result[i] = byte(uintValue)
			continue
		}
		intValue, ok := elem.(int)
		if ok {
			if intValue < 0 || intValue > 255 {
				return nil, fmt.Errorf("element at index %d is out of byte range", i)
			}
			result[i] = byte(intValue)
			continue
		}
		floatValue, ok := elem.(float64)
		if ok {
			if floatValue != math.Floor(floatValue) {
				return nil, fmt.Errorf("element at index %d is not an integer", i)
			}
			if floatValue < 0 || floatValue > 255 {
				return nil, fmt.Errorf("element at index %d is out of byte range", i)
			}
			result[i] = byte(floatValue)
			continue
		}
		return nil, fmt.Errorf("element at index %d is not a number", i)
	}
	return result, nil
}

type ModuleConfig struct {
	Id     string `json:"id"`
	Type   string `json:"type"`
	Params Params `json:"params"`
}

type RouteConfig struct {
	Input      string            `json:"input"`
	Processors []ProcessorConfig `json:"processors"`
}

type ProcessorConfig struct {
	Type   string `json:"type"`
	Params Params `json:"params"`
}
