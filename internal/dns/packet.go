package dns

import (
	"strings"

	"github.com/SergeyCherepiuk/dns-go/internal/utils"
)

type Packet struct {
	Header            Header
	Questions         []Question
	Answers           []Record
	AuthorityRecords  []Record
	AdditionalRecords []Record
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

	bytesWritten += marshal(marshalInput[Question]{
		Items:       packet.Questions,
		MarshalFunc: MarshalQuestion,
		DomainFunc:  func(q Question) string { return q.Domain },
		Bytes:       &bytes,
		Offset:      bytesWritten,
		Lookup:      lookup,
	})

	bytesWritten += marshal(marshalInput[Record]{
		Items:       packet.Answers,
		MarshalFunc: MarshalRecord,
		DomainFunc:  func(r Record) string { return r.Domain },
		Bytes:       &bytes,
		Offset:      bytesWritten,
		Lookup:      lookup,
	})

	bytesWritten += marshal(marshalInput[Record]{
		Items:       packet.AuthorityRecords,
		MarshalFunc: MarshalRecord,
		DomainFunc:  func(r Record) string { return r.Domain },
		Bytes:       &bytes,
		Offset:      bytesWritten,
		Lookup:      lookup,
	})

	bytesWritten += marshal(marshalInput[Record]{
		Items:       packet.AdditionalRecords,
		MarshalFunc: MarshalRecord,
		DomainFunc:  func(r Record) string { return r.Domain },
		Bytes:       &bytes,
		Offset:      bytesWritten,
		Lookup:      lookup,
	})

	return bytes
}

func UnmarshalPacket(bytes []byte) Packet {
	bytesRead := 0

	var (
		headerBytes = [HeaderSize]byte(bytes[bytesRead : bytesRead+HeaderSize])
		header      = UnmarshalHeader(headerBytes)
		packet      = Packet{Header: header}
	)
	bytesRead += HeaderSize

	lookup := make(map[int]string, 0)

	packet.Questions = unmarshal(unmarshalInput[Question]{
		Bytes:         bytes,
		Offset:        &bytesRead,
		Count:         header.QuestionSectionSize,
		UnmarshalFunc: UnmarshalQuestion,
		DomainFunc:    func(q Question) string { return q.Domain },
		Lookup:        lookup,
	})

	packet.Answers = unmarshal(unmarshalInput[Record]{
		Bytes:         bytes,
		Offset:        &bytesRead,
		Count:         header.AnswerSectionSize,
		UnmarshalFunc: UnmarshalRecord,
		DomainFunc:    func(r Record) string { return r.Domain },
		Lookup:        lookup,
	})

	packet.AuthorityRecords = unmarshal(unmarshalInput[Record]{
		Bytes:         bytes,
		Offset:        &bytesRead,
		Count:         header.AuthorityRecordsSectionSize,
		UnmarshalFunc: UnmarshalRecord,
		DomainFunc:    func(r Record) string { return r.Domain },
		Lookup:        lookup,
	})

	packet.AdditionalRecords = unmarshal(unmarshalInput[Record]{
		Bytes:         bytes,
		Offset:        &bytesRead,
		Count:         header.AdditionalRecordsSectionSize,
		UnmarshalFunc: UnmarshalRecord,
		DomainFunc:    func(r Record) string { return r.Domain },
		Lookup:        lookup,
	})

	return packet
}

type marshalInput[T any] struct {
	Items       []T
	MarshalFunc func(T, map[int]string) []byte
	DomainFunc  func(T) string
	Bytes       *[]byte
	Offset      int
	Lookup      map[int]string
}

func marshal[T any](input marshalInput[T]) int {
	var bytesWritten int

	for _, item := range input.Items {
		itemBytes := input.MarshalFunc(item, input.Lookup)
		*input.Bytes = append(*input.Bytes, itemBytes...)

		cacheDomain(input.DomainFunc(item), input.Offset+bytesWritten, input.Lookup)

		bytesWritten += len(itemBytes)
	}

	return bytesWritten
}

type unmarshalInput[T any] struct {
	Bytes         []byte
	Offset        *int
	Count         uint16
	UnmarshalFunc func([]byte, map[int]string) (T, int)
	DomainFunc    func(T) string
	Lookup        map[int]string
}

func unmarshal[T any](input unmarshalInput[T]) []T {
	var (
		items     = make([]T, input.Count)
		bytesRead int
	)

	for i := range input.Count {
		item, n := input.UnmarshalFunc(input.Bytes[*input.Offset+bytesRead:], input.Lookup)
		items[i] = item

		cacheDomain(input.DomainFunc(item), *input.Offset+bytesRead, input.Lookup)

		bytesRead += n
	}

	*input.Offset += bytesRead

	return items
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
