package dns

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
