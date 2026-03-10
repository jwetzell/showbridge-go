package common

import (
	"math"
	"reflect"
)

func GetAnyAs[T any](value any) (T, bool) {
	typed, ok := value.(T)
	return typed, ok
}

func GetAnyAsInt(value any) (int, bool) {

	intValue, ok := value.(int)
	if ok {
		return intValue, true
	}

	uintValue, ok := value.(uint)
	if ok {
		return int(uintValue), true
	}

	byteValue, ok := value.(byte)
	if ok {
		return int(byteValue), true
	}

	floatValue, ok := value.(float64)
	if ok {
		if floatValue != math.Floor(floatValue) {
			return 0, false
		}
		return int(floatValue), true
	}
	return 0, false
}

func GetAnyAsByteSlice(value any) ([]byte, bool) {
	v := reflect.ValueOf(value)
	if v.Kind() != reflect.Slice {
		return nil, false
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
				return nil, false
			}
			result[i] = byte(uintValue)
			continue
		}
		intValue, ok := elem.(int)
		if ok {
			if intValue < 0 || intValue > 255 {
				return nil, false
			}
			result[i] = byte(intValue)
			continue
		}
		floatValue, ok := elem.(float64)
		if ok {
			if floatValue != math.Floor(floatValue) {
				return nil, false
			}
			if floatValue < 0 || floatValue > 255 {
				return nil, false
			}
			result[i] = byte(floatValue)
			continue
		}
		return nil, false
	}
	return result, true
}

func GetAnyAsIntSlice(value any) ([]int, bool) {
	v := reflect.ValueOf(value)
	if v.Kind() != reflect.Slice {
		return nil, false
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
				return nil, false
			}
			result[i] = int(floatValue)
			continue
		}
		return nil, false
	}
	return result, true
}

func GetAnyAsFloat32(value any) (float32, bool) {
	float32Value, ok := value.(float32)
	if ok {
		return float32Value, true
	}

	float64Value, ok := value.(float64)
	if ok {
		return float32(float64Value), true
	}

	intValue, ok := value.(int)
	if ok {
		return float32(intValue), true
	}

	uintValue, ok := value.(uint)
	if ok {
		return float32(uintValue), true
	}

	byteValue, ok := value.(byte)
	if ok {
		return float32(byteValue), true
	}

	return 0, false
}

func GetAnyAsFloat64(value any) (float64, bool) {
	float64Value, ok := value.(float64)
	if ok {
		return float64Value, true
	}

	float32Value, ok := value.(float32)
	if ok {
		return float64(float32Value), true
	}

	intValue, ok := value.(int)
	if ok {
		return float64(intValue), true
	}

	uintValue, ok := value.(uint)
	if ok {
		return float64(uintValue), true
	}

	byteValue, ok := value.(byte)
	if ok {
		return float64(byteValue), true
	}

	return 0, false
}
