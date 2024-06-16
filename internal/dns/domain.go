package dns

import "github.com/SergeyCherepiuk/dns-go/internal/utils"

// TODO: Implement "MarshalDomain"

func UnmarshalDomain(bytes []byte, lookup map[int]string) (string, int) {
	var (
		domain    []byte
		bytesRead uint16
	)

	for {
		size := uint16(bytes[bytesRead])
		bytesRead += 1

		if size == 0 {
			break
		}

		if size&0b11000000 == 0b11000000 {
			pointerBytes := [2]byte{byte(size) & 0b00111111, bytes[bytesRead]}
			pointer := utils.BytesToUint16(pointerBytes)
			bytesRead += 1

			lookedupDomain := []byte(lookup[int(pointer)])
			domain = append(domain, lookedupDomain...)

			continue
		}

		domain = append(domain, bytes[bytesRead:bytesRead+size]...)
		domain = append(domain, '.')
		bytesRead += size
	}

	return string(domain), int(bytesRead)
}
