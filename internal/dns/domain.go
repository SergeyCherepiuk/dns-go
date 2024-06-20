package dns

import (
	"strings"

	"github.com/SergeyCherepiuk/dns-go/internal/utils"
)

func MarshalDomain(domain string, lookup map[int]string) []byte {
	var bytes []byte

	subdomains := strings.Split(domain, ".")
	for i, subdomain := range subdomains {
		fullDomain := strings.Join(subdomains[i:], ".")
		index, ok := utils.KeyByValue(lookup, fullDomain)

		if ok {
			pointer := 0b11000000_00000000 | uint16(index)
			pointerBytes := utils.Uint16ToBytes(pointer)
			bytes = append(bytes, pointerBytes[:]...)
			break
		}

		size := byte(len(subdomain))
		bytes = append(bytes, size)
		bytes = append(bytes, subdomain...)
	}

	return bytes
}

func UnmarshalDomain(bytes []byte, lookup map[int]string) (string, int) {
	var (
		domain    []byte
		bytesRead int
	)

	for {
		size := int(bytes[bytesRead])
		bytesRead += 1

		if size == 0 {
			break
		}

		if size&0b11000000 == 0b11000000 {
			pointerBytes := [2]byte{byte(size) & 0b00111111, bytes[bytesRead]}
			bytesRead += 1

			pointer := int(utils.BytesToUint16(pointerBytes))
			lookedupDomain := []byte(lookup[pointer])

			domain = append(domain, lookedupDomain...)
			break
		}

		domain = append(domain, bytes[bytesRead:bytesRead+size]...)
		domain = append(domain, '.')
		bytesRead += size
	}

	return string(domain), bytesRead
}
