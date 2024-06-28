package dns

type QuestionType uint16

const (
	QuestionTypeA     = QuestionType(RecordTypeA)
	QuestionTypeNS    = QuestionType(RecordTypeNS)
	QuestionTypeMD    = QuestionType(RecordTypeMD)
	QuestionTypeMF    = QuestionType(RecordTypeMF)
	QuestionTypeCNAME = QuestionType(RecordTypeCNAME)
	QuestionTypeSOA   = QuestionType(RecordTypeSOA)
	QuestionTypeMB    = QuestionType(RecordTypeMB)
	QuestionTypeMG    = QuestionType(RecordTypeMG)
	QuestionTypeMR    = QuestionType(RecordTypeMR)
	QuestionTypeNULL  = QuestionType(RecordTypeNULL)
	QuestionTypeWKS   = QuestionType(RecordTypeWKS)
	QuestionTypePTR   = QuestionType(RecordTypePTR)
	QuestionTypeHINFO = QuestionType(RecordTypeHINFO)
	QuestionTypeMINFO = QuestionType(RecordTypeMINFO)
	QuestionTypeMX    = QuestionType(RecordTypeMX)
	QuestionTypeTXT   = QuestionType(RecordTypeTXT)

	QuestionTypeAXFR = QuestionType(iota + 236)
	QuestionTypeMAILB
	QuestionTypeMAILA
	QuestionTypeALL
)

type QuestionClass uint16

const (
	QuestionClassIN = QuestionClass(RecordClassIN)
	QuestionClassCS = QuestionClass(RecordClassCS)
	QuestionClassCH = QuestionClass(RecordClassCH)
	QuestionClassHS = QuestionClass(RecordClassHS)

	QuestionClassALL = QuestionClass(iota + 251)
)

type Question struct {
	Domain string
	Type   QuestionType
	Class  QuestionClass
}

func marshalQuestion(w *PacketWriter, question Question) error {
	err := w.WriteDomain(question.Domain)
	if err != nil {
		return err
	}

	err = w.WriteUint16(uint16(question.Type))
	if err != nil {
		return err
	}

	err = w.WriteUint16(uint16(question.Class))
	if err != nil {
		return err
	}

	return nil
}

func unmarshalQuestion(r *PacketReader) (Question, error) {
	domain, err := r.ReadDomain()
	if err != nil {
		return Question{}, err
	}

	questionType, err := r.ReadUint16()
	if err != nil {
		return Question{}, err
	}

	questionClass, err := r.ReadUint16()
	if err != nil {
		return Question{}, err
	}

	question := Question{
		Domain: domain,
		Type:   QuestionType(questionType),
		Class:  QuestionClass(questionClass),
	}

	return question, nil
}
