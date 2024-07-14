package types

import "fmt"

const MaxPacketSize = 512

type PacketRecords struct {
	Answers           []Record
	AuthorityRecords  []Record
	AdditionalRecords []Record
}

func (r *PacketRecords) Len() int {
	return len(r.Answers) + len(r.AuthorityRecords) + len(r.AdditionalRecords)
}

type Packet struct {
	Header    Header
	Questions []Question
	Records   PacketRecords
}

func (p Packet) String() string {
	var bytes []byte

	header := fmt.Sprintf(
		"Id: %d, Type: %v, RD: %t, RA: %t, Section sizes: [%d, %d, %d, %d]",
		p.Header.ID, p.Header.PacketType,
		p.Header.RecursionDesired, p.Header.RecursionAvailable,
		p.Header.QuestionSectionSize, p.Header.AnswerSectionSize,
		p.Header.AuthorityRecordsSectionSize, p.Header.AdditionalRecordsSectionSize,
	)
	bytes = append(bytes, header...)
	bytes = append(bytes, '\n')

	for _, question := range p.Questions {
		bytes = append(bytes, question.String()...)
		bytes = append(bytes, '\n')
	}

	for _, answer := range p.Records.Answers {
		bytes = append(bytes, answer.String()...)
		bytes = append(bytes, '\n')
	}

	for _, authorityRecord := range p.Records.AuthorityRecords {
		bytes = append(bytes, authorityRecord.String()...)
		bytes = append(bytes, '\n')
	}

	for _, additionalRecord := range p.Records.AdditionalRecords {
		bytes = append(bytes, additionalRecord.String()...)
		bytes = append(bytes, '\n')
	}

	return string(bytes)
}
