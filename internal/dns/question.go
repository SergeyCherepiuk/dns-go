package dns

import (
	"github.com/SergeyCherepiuk/dns-go/internal/utils"
)

type QuestionType uint16

const (
	_ = QuestionType(iota)
	QuestionTypeA
	QuestionTypeNS
	QuestionTypeMD
	QuestionTypeMF
	QuestionTypeCNAME
	QuestionTypeSOA
	QuestionTypeMB
	QuestionTypeMG
	QuestionTypeMR
	QuestionTypeNULL
	QuestionTypeWKS
	QuestionTypePTR
	QuestionTypeHINFO
	QuestionTypeMINFO
	QuestionTypeMX
	QuestionTypeTXT
	QuestionTypeAXFR = QuestionType(iota + 235)
	QuestionTypeMAILB
	QuestionTypeMAILA
	QuestionTypeALL
)

type QuestionClass uint16

const (
	_ = QuestionClass(iota)
	QuestionClassIN
	QuestionClassCS
	QuestionClassCH
	QuestionClassHS
	QuestionClassALL = QuestionClass(iota + 250)
)

type Question struct {
	Domain string
	Type   QuestionType
	Class  QuestionClass
}

func MarshalQuestion(question Question, lookup map[string]int) []byte {
	var bytes []byte

	domainBytes := MarshalDomain(question.Domain, lookup)
	bytes = append(bytes, domainBytes...)

	questionTypeBytes := utils.Uint16ToBytes(uint16(question.Type))
	bytes = append(bytes, questionTypeBytes[:]...)

	questionClassBytes := utils.Uint16ToBytes(uint16(question.Class))
	bytes = append(bytes, questionClassBytes[:]...)

	return bytes
}

func UnmarshalQuestion(bytes []byte, lookup map[int]string) (Question, int) {
	domain, bytesRead := UnmarshalDomain(bytes, lookup)

	questionTypeBytes := [2]byte{bytes[bytesRead], bytes[bytesRead+1]}
	bytesRead += 2

	questionClassBytes := [2]byte{bytes[bytesRead], bytes[bytesRead+1]}
	bytesRead += 2

	question := Question{
		Domain: domain,
		Type:   QuestionType(utils.BytesToUint16(questionTypeBytes)),
		Class:  QuestionClass(utils.BytesToUint16(questionClassBytes)),
	}

	return question, bytesRead
}
