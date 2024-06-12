package utils

func BytesToUint16(bytes [2]byte) uint16 {
	var u16 uint16
	u16 += uint16(bytes[0]) << 8
	u16 += uint16(bytes[1]) << 0
	return u16
}

func BytesToUint32(bytes [4]byte) uint32 {
	var u32 uint32
	u32 += uint32(bytes[0]) << 24
	u32 += uint32(bytes[1]) << 16
	u32 += uint32(bytes[2]) << 8
	u32 += uint32(bytes[3]) << 0
	return u32
}

func Uint16ToBytes(number uint16) [2]byte {
	return [2]byte{
		byte(number>>8) & 0xFF,
		byte(number>>0) & 0xFF,
	}
}

func Uint32ToBytes(number uint32) [4]byte {
	return [4]byte{
		byte(number>>24) & 0xFF,
		byte(number>>16) & 0xFF,
		byte(number>>8) & 0xFF,
		byte(number>>0) & 0xFF,
	}
}
