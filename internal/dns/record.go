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
	Data   []byte
}

func marshalRecord(record Record, lookup map[int]string) []byte {
	var bytes []byte

	domainBytes := marshalDomain(record.Domain, lookup)
	bytes = append(bytes, domainBytes...)

	typeBytes := utils.Uint16ToBytes(uint16(record.Type))
	bytes = append(bytes, typeBytes[:]...)

	classBytes := utils.Uint16ToBytes(uint16(record.Class))
	bytes = append(bytes, classBytes[:]...)

	ttlBytes := utils.Uint32ToBytes(record.Ttl)
	bytes = append(bytes, ttlBytes[:]...)

	length := uint16(len(record.Data))
	lengthBytes := utils.Uint16ToBytes(length)
	bytes = append(bytes, lengthBytes[:]...)

	bytes = append(bytes, record.Data...)

	return bytes
}

func unmarshalRecord(bytes []byte, lookup map[int]string) (Record, int) {
	domain, bytesRead := unmarshalDomain(bytes, lookup)

	typeBytes := [2]byte(bytes[bytesRead : bytesRead+2])
	bytesRead += 2

	classBytes := [2]byte(bytes[bytesRead : bytesRead+2])
	bytesRead += 2

	ttlBytes := [4]byte(bytes[bytesRead : bytesRead+4])
	bytesRead += 4

	lengthBytes := [2]byte(bytes[bytesRead : bytesRead+2])
	length := int(utils.BytesToUint16(lengthBytes))
	bytesRead += 2

	dataBytes := bytes[bytesRead : bytesRead+length]
	bytesRead += length

	record := Record{
		Domain: domain,
		Type:   RecordType(utils.BytesToUint16(typeBytes)),
		Class:  RecordClass(utils.BytesToUint16(classBytes)),
		Ttl:    utils.BytesToUint32(ttlBytes),
		Data:   dataBytes,
	}

	return record, bytesRead
}
