package utils

func KeyByValue[K, V comparable](table map[K]V, target V) (K, bool) {
	for key, value := range table {
		if value == target {
			return key, true
		}
	}

	var zero K
	return zero, false
}
