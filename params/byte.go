package params

import (
	"encoding"
	"math"
	"reflect"
)

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
	// already a []byte
	byteSlice, ok := value.([]byte)
	if ok {
		return byteSlice, true
	}

	// check for a slice
	v := reflect.ValueOf(value)
	if v.Kind() != reflect.Slice && v.Kind() != reflect.Array {

		valueString, ok := value.(string)
		if ok {
			return []byte(valueString), true
		}
		// something that can be marshaled to []byte
		valueMarshalBinary, ok := value.(encoding.BinaryMarshaler)
		if ok {
			bytes, err := valueMarshalBinary.MarshalBinary()
			if err != nil {
				return nil, false
			}
			return bytes, true
		}

		return nil, false
	}

	// try to convert each element to a byte
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
