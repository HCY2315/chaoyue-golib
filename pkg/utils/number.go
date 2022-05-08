package utils

func IsThisFloatAInt(val float64) bool {
	return val == float64(int(val))
}
