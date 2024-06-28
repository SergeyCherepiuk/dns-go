package dns

const MaxPacketSize = 512

type Packet struct {
	Header            Header
	Questions         []Question
	Answers           []Record
	AuthorityRecords  []Record
	AdditionalRecords []Record
}

func MarshalPacket(w *PacketWriter, packet Packet) error {
	err := marshalHeader(w, packet.Header)
	if err != nil {
		return err
	}

	for _, question := range packet.Questions {
		err := marshalQuestion(w, question)
		if err != nil {
			return err
		}
	}

	for _, answer := range packet.Answers {
		err := marshalRecord(w, answer)
		if err != nil {
			return err
		}
	}

	for _, authorityRecord := range packet.AuthorityRecords {
		err := marshalRecord(w, authorityRecord)
		if err != nil {
			return err
		}
	}

	for _, additionalRecord := range packet.AdditionalRecords {
		err := marshalRecord(w, additionalRecord)
		if err != nil {
			return err
		}
	}

	return nil
}

func UnmarshalPacket(r *PacketReader) (Packet, error) {
	header, err := unmarshalHeader(r)
	if err != nil {
		return Packet{}, err
	}

	packet := Packet{Header: header}

	for range packet.Header.QuestionSectionSize {
		question, err := unmarshalQuestion(r)
		if err != nil {
			return Packet{}, err
		}

		packet.Questions = append(packet.Questions, question)
	}

	for range packet.Header.AnswerSectionSize {
		answer, err := unmarshalRecord(r)
		if err != nil {
			return Packet{}, err
		}

		packet.Answers = append(packet.Answers, answer)
	}

	for range packet.Header.AuthorityRecordsSectionSize {
		authorityRecord, err := unmarshalRecord(r)
		if err != nil {
			return Packet{}, err
		}

		packet.AuthorityRecords = append(packet.AuthorityRecords, authorityRecord)
	}

	for range packet.Header.AdditionalRecordsSectionSize {
		additionalRecord, err := unmarshalRecord(r)
		if err != nil {
			return Packet{}, err
		}

		packet.AdditionalRecords = append(packet.AdditionalRecords, additionalRecord)
	}

	return packet, nil
}
