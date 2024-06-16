package utils

func MapContainsValue[K, V comparable](table map[K]V, target V) bool {
	for _, value := range table {
		if value == target {
			return true
		}
	}
	return false
}
