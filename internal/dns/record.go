package dns

import "github.com/SergeyCherepiuk/dns-go/internal/utils"

type RecordType uint16

const (
	_ = RecordType(iota)
	RecordTypeA
	RecordTypeNS
	RecordTypeMD
	RecordTypeMF
	RecordTypeCNAME
	RecordTypeSOA
	RecordTypeMB
	RecordTypeMG
	RecordTypeMR
	RecordTypeNULL
	RecordTypeWKS
	RecordTypePTR
	RecordTypeHINFO
	RecordTypeMINFO
	RecordTypeMX
	RecordTypeTXT
)

type RecordClass uint16

const (
	_ = RecordClass(iota)
	RecordClassIN
	RecordClassCS
	RecordClassCH
	RecordClassHS
)

type Record struct {
	Domain string
	Type   RecordType
	Class  RecordClass
	Ttl    uint32
	Length uint16
	Data   []byte
}

func MarshalRecord(record Record, lookup map[int]string) []byte {
	var bytes []byte

	domainBytes := MarshalDomain(record.Domain, lookup)
	bytes = append(bytes, domainBytes...)

	typeBytes := utils.Uint16ToBytes(uint16(record.Type))
	bytes = append(bytes, typeBytes[:]...)

	classBytes := utils.Uint16ToBytes(uint16(record.Class))
	bytes = append(bytes, classBytes[:]...)

	ttlBytes := utils.Uint32ToBytes(record.Ttl)
	bytes = append(bytes, ttlBytes[:]...)

	lengthBytes := utils.Uint16ToBytes(record.Length)
	bytes = append(bytes, lengthBytes[:]...)

	bytes = append(bytes, record.Data...)

	return bytes
}

func UnmarshalRecord(bytes []byte, lookup map[int]string) (Record, int) {
	domain, bytesRead := UnmarshalDomain(bytes, lookup)

	typeBytes := [2]byte(bytes[bytesRead : bytesRead+2])
	bytesRead += 2

	classBytes := [2]byte(bytes[bytesRead : bytesRead+2])
	bytesRead += 2

	ttlBytes := [4]byte(bytes[bytesRead : bytesRead+4])
	bytesRead += 4

	lengthBytes := [2]byte(bytes[bytesRead : bytesRead+2])
	bytesRead += 2

	record := Record{
		Domain: domain,
		Type:   RecordType(utils.BytesToUint16(typeBytes)),
		Class:  RecordClass(utils.BytesToUint16(classBytes)),
		Ttl:    utils.BytesToUint32(ttlBytes),
		Length: utils.BytesToUint16(lengthBytes),
	}

	record.Data = bytes[bytesRead : bytesRead+int(record.Length)]
	bytesRead += int(record.Length)

	return record, bytesRead
}
