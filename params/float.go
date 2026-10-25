package params

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
