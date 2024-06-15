package dns

import (
	"github.com/SergeyCherepiuk/dns-go/internal/utils"
)

type QueryPacket struct {
	Header    Header
	Questions []Question
}

// TODO: Not implemented
func MarshalQueryPacket(packet QueryPacket) []byte {
	bytes := make([]byte, 0)

	headerBytes := MarshalHeader(packet.Header)
	bytes = append(bytes, headerBytes[:]...)

	for _, question := range packet.Questions {
		_ = question
		// TODO
	}

	return bytes
}

// TODO: Refactor this nightmare
func UnmarshalQueryPacket(bytes []byte) QueryPacket {
	j := uint16(0)

	var (
		headerBytes = [HeaderSize]byte(bytes[j : j+HeaderSize])
		header      = UnmarshalHeader(headerBytes)
	)
	j += HeaderSize

	queryPacket := QueryPacket{
		Header:    header,
		Questions: make([]Question, header.QuestionSectionSize),
	}

	for i := range header.QuestionSectionSize {
		var (
			domain      = make([]byte, 0)
			domainIndex = j
		)

		for {
			size := uint16(bytes[domainIndex])
			if domainIndex == j {
				j += 1
			}
			domainIndex += 1

			if size&0b11000000 == 0b11000000 {
				pointerBytes := [2]byte{byte(size) & 0b00111111, bytes[domainIndex]}
				pointer := utils.BytesToUint16(pointerBytes)
				domainIndex = pointer

				j += 1

				continue
			}

			if size == 0 {
				break
			}

			domain = append(domain, bytes[domainIndex:domainIndex+size]...)
			domain = append(domain, '.')

			if domainIndex == j {
				j += size
			}
			domainIndex += size
		}

		var (
			questionTypeBytes = [2]byte{bytes[j], bytes[j+1]}
			questionType      = QuestionType(utils.BytesToUint16(questionTypeBytes))
		)
		j += 2

		var (
			questionClassBytes = [2]byte{bytes[j], bytes[j+1]}
			questionClass      = QuestionClass(utils.BytesToUint16(questionClassBytes))
		)
		j += 2

		queryPacket.Questions[i] = Question{
			Domain: string(domain),
			Type:   questionType,
			Class:  questionClass,
		}
	}

	return queryPacket
}
