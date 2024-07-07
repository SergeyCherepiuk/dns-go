package types

import "fmt"

type QuestionType uint16

const (
	QuestionTypeA     = QuestionType(1)
	QuestionTypeNS    = QuestionType(2)
	QuestionTypeCNAME = QuestionType(5)
	QuestionTypeMX    = QuestionType(15)
	QuestionTypeAAAA  = QuestionType(28)
)

type QuestionClass uint16

const QuestionClassIN = QuestionClass(1)

type Question struct {
	Domain string
	Type   QuestionType
	Class  QuestionClass
}

func (q Question) String() string {
	return fmt.Sprintf("%s, %v, %v", q.Domain, q.Type, q.Class)
}
