package dns

import "fmt"

const MaxPacketSize = 512

type Packet struct {
	Header            Header
	Questions         []Question
	Answers           []Record
	AuthorityRecords  []Record
	AdditionalRecords []Record
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

	for _, answer := range p.Answers {
		bytes = append(bytes, answer.String()...)
		bytes = append(bytes, '\n')
	}

	for _, authorityRecord := range p.AuthorityRecords {
		bytes = append(bytes, authorityRecord.String()...)
		bytes = append(bytes, '\n')
	}

	for _, additionalRecord := range p.AdditionalRecords {
		bytes = append(bytes, additionalRecord.String()...)
		bytes = append(bytes, '\n')
	}

	return string(bytes)
}

func MarshalPacket(packet Packet) ([]byte, error) {
	writer := newPacketWriter()
	err := marshalHeader(writer, packet.Header)
	if err != nil {
		return nil, err
	}

	for _, question := range packet.Questions {
		err := marshalQuestion(writer, question)
		if err != nil {
			return nil, err
		}
	}

	for _, answer := range packet.Answers {
		err := marshalRecord(writer, answer)
		if err != nil {
			return nil, err
		}
	}

	for _, authorityRecord := range packet.AuthorityRecords {
		err := marshalRecord(writer, authorityRecord)
		if err != nil {
			return nil, err
		}
	}

	for _, additionalRecord := range packet.AdditionalRecords {
		err := marshalRecord(writer, additionalRecord)
		if err != nil {
			return nil, err
		}
	}

	return writer.Bytes(), nil
}

func UnmarshalPacket(bytes []byte) (Packet, error) {
	reader, err := newPacketReader(bytes)
	if err != nil {
		return Packet{}, err
	}

	header, err := unmarshalHeader(reader)
	if err != nil {
		return Packet{}, err
	}

	packet := Packet{Header: header}

	for range packet.Header.QuestionSectionSize {
		question, err := unmarshalQuestion(reader)
		if err != nil {
			return Packet{}, err
		}

		packet.Questions = append(packet.Questions, question)
	}

	for range packet.Header.AnswerSectionSize {
		answer, err := unmarshalRecord(reader)
		if err != nil {
			return Packet{}, err
		}

		packet.Answers = append(packet.Answers, answer)
	}

	for range packet.Header.AuthorityRecordsSectionSize {
		authorityRecord, err := unmarshalRecord(reader)
		if err != nil {
			return Packet{}, err
		}

		packet.AuthorityRecords = append(packet.AuthorityRecords, authorityRecord)
	}

	for range packet.Header.AdditionalRecordsSectionSize {
		additionalRecord, err := unmarshalRecord(reader)
		if err != nil {
			return Packet{}, err
		}

		packet.AdditionalRecords = append(packet.AdditionalRecords, additionalRecord)
	}

	return packet, nil
}
