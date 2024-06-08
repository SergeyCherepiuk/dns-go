package utils

// TODO: Unit test
func MaskBit(b byte, n uint) uint8 {
	return (b >> n) & 0x1
}
