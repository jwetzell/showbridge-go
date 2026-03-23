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

	float32Value, ok := value.(float32)
	if ok {
		if float64(float32Value) != math.Floor(float64(float32Value)) {
			return 0, false
		}
		return int(float32Value), true
	}

	float64Value, ok := value.(float64)
	if ok {
		if float64Value != math.Floor(float64Value) {
			return 0, false
		}
		return int(float64Value), true
	}
	return 0, false
}

func GetAnyAsByte(value any) (byte, bool) {

	byteValue, ok := value.(byte)
	if ok {
		return byte(byteValue), true
	}

	intValue, ok := value.(int)
	if ok {
		return byte(intValue), true
	}

	uintValue, ok := value.(uint)
	if ok {
		return byte(uintValue), true
	}

	float32Value, ok := value.(float32)
	if ok {
		if float64(float32Value) != math.Floor(float64(float32Value)) {
			return 0, false
		}
		return byte(float32Value), true
	}

	float64Value, ok := value.(float64)
	if ok {
		if float64Value != math.Floor(float64Value) {
			return 0, false
		}
		return byte(float64Value), true
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
		elemValue, ok := GetAnyAsByte(elem)
		if !ok {
			return nil, false
		}
		result[i] = elemValue
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
		elemInt, ok := GetAnyAsInt(elem)
		if !ok {
			return nil, false
		}
		result[i] = elemInt
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
