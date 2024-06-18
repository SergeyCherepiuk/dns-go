package dns

import (
	"strings"

	"github.com/SergeyCherepiuk/dns-go/internal/utils"
)

type Packet struct {
	Header    Header
	Questions []Question
	// TODO: Add "Answers" field of type []Answer
}

func MarshalPacket(packet Packet) []byte {
	var (
		bytes        []byte
		bytesWritten int
	)

	headerBytes := MarshalHeader(packet.Header)
	bytes = append(bytes, headerBytes[:]...)
	bytesWritten += HeaderSize

	lookup := make(map[int]string, 0)
	for _, question := range packet.Questions {
		questionBytes := MarshalQuestion(question, lookup)
		bytes = append(bytes, questionBytes...)

		cacheDomain(question.Domain, bytesWritten, lookup)

		bytesWritten += len(questionBytes)
	}

	// TODO: Marshal answers

	return bytes
}

func UnmarshalPacket(bytes []byte) Packet {
	bytesRead := 0

	var (
		headerBytes = [HeaderSize]byte(bytes[bytesRead : bytesRead+HeaderSize])
		header      = UnmarshalHeader(headerBytes)
	)
	bytesRead += HeaderSize

	packet := Packet{
		Header:    header,
		Questions: make([]Question, header.QuestionSectionSize),
	}

	lookup := make(map[int]string, 0)
	for i := range header.QuestionSectionSize {
		question, n := UnmarshalQuestion(bytes[bytesRead:], lookup)

		packet.Questions[i] = question
		cacheDomain(question.Domain, bytesRead, lookup)

		bytesRead += n
	}

	// TODO: Unmarshal answers

	return packet
}

// TODO (low priority): Optimize caching. Lookup table can potentially grow big.
func cacheDomain(domain string, offset int, lookup map[int]string) {
	subdomains := strings.Split(domain, ".")

	for i := 0; i < len(subdomains); i++ {
		subdomain := strings.Join(subdomains[i:], ".")

		if subdomain == "" {
			continue
		}

		var (
			subdomainsBefore       = subdomains[:i]
			delimiters             = len(subdomainsBefore)
			subdomainsBeforeLength = lenSum(subdomainsBefore)
			bytesBefore            = subdomainsBeforeLength + delimiters
			startByte              = offset + bytesBefore
		)

		_, ok := utils.KeyByValue(lookup, subdomain)

		if ok {
			lookup[startByte] = subdomain
			return
		}

		lookup[startByte] = subdomain
	}
}

func lenSum(strings []string) int {
	sum := 0
	for _, s := range strings {
		sum += len(s)
	}
	return sum
}
