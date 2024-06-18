package dns

import (
	"github.com/SergeyCherepiuk/dns-go/internal/utils"
)

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

func MarshalQuestion(question Question, lookup map[int]string) []byte {
	var bytes []byte

	domainBytes := MarshalDomain(question.Domain, lookup)
	bytes = append(bytes, domainBytes...)

	typeBytes := utils.Uint16ToBytes(uint16(question.Type))
	bytes = append(bytes, typeBytes[:]...)

	classBytes := utils.Uint16ToBytes(uint16(question.Class))
	bytes = append(bytes, classBytes[:]...)

	return bytes
}

func UnmarshalQuestion(bytes []byte, lookup map[int]string) (Question, int) {
	domain, bytesRead := UnmarshalDomain(bytes, lookup)

	typeBytes := [2]byte(bytes[bytesRead : bytesRead+2])
	bytesRead += 2

	classBytes := [2]byte(bytes[bytesRead : bytesRead+2])
	bytesRead += 2

	question := Question{
		Domain: domain,
		Type:   QuestionType(utils.BytesToUint16(typeBytes)),
		Class:  QuestionClass(utils.BytesToUint16(classBytes)),
	}

	return question, bytesRead
}
