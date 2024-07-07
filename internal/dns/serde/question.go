package serde

import (
	"github.com/SergeyCherepiuk/dns-go/internal/dns/io"
	"github.com/SergeyCherepiuk/dns-go/internal/dns/types"
)

func marshalQuestion(w *io.PacketWriter, question types.Question) error {
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

func unmarshalQuestion(r *io.PacketReader) (types.Question, error) {
	domain, err := r.ReadDomain()
	if err != nil {
		return types.Question{}, err
	}

	questionType, err := r.ReadUint16()
	if err != nil {
		return types.Question{}, err
	}

	questionClass, err := r.ReadUint16()
	if err != nil {
		return types.Question{}, err
	}

	question := types.Question{
		Domain: domain,
		Type:   types.QuestionType(questionType),
		Class:  types.QuestionClass(questionClass),
	}

	return question, nil
}
