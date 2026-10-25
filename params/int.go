package params

import (
	"math"
	"reflect"
)

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

func GetAnyAsIntSlice(value any) ([]int, bool) {

	// already a []int
	intSlice, ok := value.([]int)
	if ok {
		return intSlice, true
	}

	// check for a slice
	v := reflect.ValueOf(value)
	if v.Kind() != reflect.Slice && v.Kind() != reflect.Array {
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
