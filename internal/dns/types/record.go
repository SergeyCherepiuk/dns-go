package types

import "fmt"

type RecordType uint16

const (
	RecordTypeA     = RecordType(1)
	RecordTypeNS    = RecordType(2)
	RecordTypeCNAME = RecordType(5)
	RecordTypeMX    = RecordType(15)
	RecordTypeAAAA  = RecordType(28)
)

type RecordClass uint16

const RecordClassIN = RecordClass(1)

type Record struct {
	Domain string
	Type   RecordType
	Class  RecordClass
	Ttl    uint32
	Data   any
}

func (r Record) String() string {
	return fmt.Sprintf(
		"%s, %v, %v, %d, %v",
		r.Domain, r.Type, r.Class, r.Ttl, r.Data,
	)
}
