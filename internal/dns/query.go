package dns

import (
	"strings"

	"github.com/SergeyCherepiuk/dns-go/internal/utils"
)

type QueryPacket struct {
	Header    Header
	Questions []Question
}

// TODO: Implement "MarshalQueryPacket"

func UnmarshalQueryPacket(bytes []byte) QueryPacket {
	bytesRead := 0

	var (
		headerBytes = [HeaderSize]byte(bytes[bytesRead : bytesRead+HeaderSize])
		header      = UnmarshalHeader(headerBytes)
	)
	bytesRead += HeaderSize

	queryPacket := QueryPacket{
		Header:    header,
		Questions: make([]Question, header.QuestionSectionSize),
	}

	lookup := make(map[int]string, 0)
	for i := range header.QuestionSectionSize {
		question, n := UnmarshalQuestion(bytes[bytesRead:], lookup)

		queryPacket.Questions[i] = question
		cacheDomain(question.Domain, bytesRead, lookup)

		bytesRead += n
	}

	return queryPacket
}

func cacheDomain(domain string, offset int, lookup map[int]string) {
	subdomains := strings.Split(domain, ".")

	for i := len(subdomains) - 1; i >= 0; i-- {
		subdomain := strings.Join(subdomains[i:], ".")
		if !utils.MapContainsValue(lookup, subdomain) {
			var (
				subdomainsBefore       = subdomains[:i]
				delimiters             = len(subdomainsBefore)
				subdomainsBeforeLength = lenSum(subdomainsBefore)
				bytesBefore            = subdomainsBeforeLength + delimiters
				startByte              = offset + bytesBefore
			)
			lookup[startByte] = subdomain
		}
	}
}

func lenSum(strings []string) int {
	sum := 0
	for _, s := range strings {
		sum += len(s)
	}
	return sum
}
