package utils

// TODO: Unit test
func BytesToUint16(bytes [2]byte) uint16 {
	var u16 uint16
	u16 += uint16(bytes[0]) << 8
	u16 += uint16(bytes[1])
	return u16
}

// TODO: Unit test
func BytesToUint32(bytes [4]byte) uint32 {
	var u32 uint32
	u32 += uint32(bytes[0]) << 24
	u32 += uint32(bytes[1]) << 16
	u32 += uint32(bytes[2]) << 8
	u32 += uint32(bytes[3])
	return u32
}
