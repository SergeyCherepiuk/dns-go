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

	headerBytes := marshalHeader(packet.Header)
	bytes = append(bytes, headerBytes[:]...)
	bytesWritten += HeaderSize

	lookup := make(map[int]string, 0)
	bytesWritten += marshalQuestions(packet.Questions, &bytes, bytesWritten, lookup)
	bytesWritten += marshalRecords(packet.Answers, &bytes, bytesWritten, lookup)
	bytesWritten += marshalRecords(packet.AuthorityRecords, &bytes, bytesWritten, lookup)
	bytesWritten += marshalRecords(packet.AdditionalRecords, &bytes, bytesWritten, lookup)

	return bytes
}

func UnmarshalPacket(bytes []byte) Packet {
	bytesRead := 0

	var (
		headerBytes = [HeaderSize]byte(bytes[bytesRead : bytesRead+HeaderSize])
		header      = unmarshalHeader(headerBytes)
		packet      = Packet{Header: header}
	)
	bytesRead += HeaderSize

	lookup := make(map[int]string, 0)
	packet.Questions = unmarshalQuestions(bytes, &bytesRead, header.QuestionSectionSize, lookup)
	packet.Answers = unmarshalRecords(bytes, &bytesRead, header.AnswerSectionSize, lookup)
	packet.AuthorityRecords = unmarshalRecords(bytes, &bytesRead, header.AuthorityRecordsSectionSize, lookup)
	packet.AdditionalRecords = unmarshalRecords(bytes, &bytesRead, header.AdditionalRecordsSectionSize, lookup)

	return packet
}

func marshalQuestions(questions []Question, bytes *[]byte, offset int, lookup map[int]string) int {
	var bytesWritten int

	for _, question := range questions {
		questionBytes := marshalQuestion(question, lookup)
		*bytes = append(*bytes, questionBytes...)

		cacheDomain(question.Domain, offset+bytesWritten, lookup)

		bytesWritten += len(questionBytes)
	}

	return bytesWritten
}

func marshalRecords(records []Record, bytes *[]byte, offset int, lookup map[int]string) int {
	var bytesWritten int

	for _, record := range records {
		recordBytes := marshalRecord(record, lookup)
		*bytes = append(*bytes, recordBytes...)

		cacheDomain(record.Domain, offset+bytesWritten, lookup)

		bytesWritten += len(recordBytes)
	}

	return bytesWritten
}

func unmarshalQuestions(bytes []byte, offset *int, count uint16, lookup map[int]string) []Question {
	var (
		questions = make([]Question, count)
		bytesRead int
	)

	for i := range count {
		question, n := unmarshalQuestion(bytes[*offset+bytesRead:], lookup)
		questions[i] = question

		cacheDomain(question.Domain, *offset+bytesRead, lookup)

		bytesRead += n
	}

	*offset += bytesRead

	return questions
}

func unmarshalRecords(bytes []byte, offset *int, count uint16, lookup map[int]string) []Record {
	var (
		records   = make([]Record, count)
		bytesRead int
	)

	for i := range count {
		record, n := unmarshalRecord(bytes[*offset+bytesRead:], lookup)
		records[i] = record

		cacheDomain(record.Domain, *offset+bytesRead, lookup)

		bytesRead += n

		// TODO: Implement data parsing of the essential question/record types
		// (A, AAAA, NS, CNAME, MX). Move this logic to the Marshal-/Unmarshal-
		// functions. Think about creating separate strutures for each type.
		if record.Type == RecordTypeCNAME {
			canonicalDomain, _ := unmarshalDomain(record.Data, lookup)
			cacheDomain(canonicalDomain, *offset+bytesRead-len(record.Data), lookup)
		}
	}

	*offset += bytesRead

	return records
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
