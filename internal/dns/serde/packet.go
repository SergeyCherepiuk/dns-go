package serde

import (
	"github.com/SergeyCherepiuk/dns-go/internal/dns/io"
	"github.com/SergeyCherepiuk/dns-go/internal/dns/types"
)

func MarshalPacket(packet types.Packet) ([]byte, error) {
	writer := io.NewPacketWriter()
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

	for _, answer := range packet.Records.Answers {
		err := marshalRecord(writer, answer)
		if err != nil {
			return nil, err
		}
	}

	for _, authorityRecord := range packet.Records.AuthorityRecords {
		err := marshalRecord(writer, authorityRecord)
		if err != nil {
			return nil, err
		}
	}

	for _, additionalRecord := range packet.Records.AdditionalRecords {
		err := marshalRecord(writer, additionalRecord)
		if err != nil {
			return nil, err
		}
	}

	return writer.Bytes(), nil
}

func UnmarshalPacket(bytes []byte) (types.Packet, error) {
	reader, err := io.NewPacketReader(bytes)
	if err != nil {
		return types.Packet{}, err
	}

	header, err := unmarshalHeader(reader)
	if err != nil {
		return types.Packet{}, err
	}

	packet := types.Packet{Header: header}

	for range packet.Header.QuestionSectionSize {
		question, err := unmarshalQuestion(reader)
		if err != nil {
			return types.Packet{}, err
		}

		packet.Questions = append(packet.Questions, question)
	}

	for range packet.Header.AnswerSectionSize {
		answer, err := unmarshalRecord(reader)
		if err != nil {
			return types.Packet{}, err
		}

		packet.Records.Answers = append(packet.Records.Answers, answer)
	}

	for range packet.Header.AuthorityRecordsSectionSize {
		authorityRecord, err := unmarshalRecord(reader)
		if err != nil {
			return types.Packet{}, err
		}

		packet.Records.AuthorityRecords = append(packet.Records.AuthorityRecords, authorityRecord)
	}

	for range packet.Header.AdditionalRecordsSectionSize {
		additionalRecord, err := unmarshalRecord(reader)
		if err != nil {
			return types.Packet{}, err
		}

		packet.Records.AdditionalRecords = append(packet.Records.AdditionalRecords, additionalRecord)
	}

	return packet, nil
}
