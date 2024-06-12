package utils

func BoolToUint8(value bool) uint8 {
	if value {
		return 1
	}
	return 0
}
