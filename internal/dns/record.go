package dns

import (
	"fmt"
	"net"
)

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

func marshalRecord(w *PacketWriter, record Record) error {
	err := w.WriteDomain(record.Domain)
	if err != nil {
		return err
	}

	err = w.WriteUint16(uint16(record.Type))
	if err != nil {
		return err
	}

	err = w.WriteUint16(uint16(record.Class))
	if err != nil {
		return err
	}

	err = w.WriteUint32(record.Ttl)
	if err != nil {
		return err
	}

	switch record.Type {
	case RecordTypeA, RecordTypeAAAA:
		bytes := []byte(record.Data.(net.IP))

		err = w.WriteUint16(uint16(len(bytes)))
		if err != nil {
			return err
		}

		err = w.WriteBytes(bytes)
		if err != nil {
			return err
		}

	case RecordTypeNS, RecordTypeCNAME:
		domain := record.Data.(string)

		err = w.WriteDomain(domain)
		if err != nil {
			return err
		}
	}

	return nil
}

func unmarshalRecord(r *PacketReader) (Record, error) {
	domain, err := r.ReadDomain()
	if err != nil {
		return Record{}, err
	}

	recordType, err := r.ReadUint16()
	if err != nil {
		return Record{}, err
	}

	recordClass, err := r.ReadUint16()
	if err != nil {
		return Record{}, err
	}

	ttl, err := r.ReadUint32()
	if err != nil {
		return Record{}, err
	}

	record := Record{
		Domain: domain,
		Type:   RecordType(recordType),
		Class:  RecordClass(recordClass),
		Ttl:    ttl,
	}

	length, err := r.ReadUint16()
	if err != nil {
		return Record{}, err
	}

	switch record.Type {
	case RecordTypeA, RecordTypeAAAA:
		ip, err := r.ReadBytes(int(length))
		if err != nil {
			return Record{}, err
		}

		record.Data = net.IP(ip)

	case RecordTypeNS, RecordTypeCNAME:
		record.Data, err = r.ReadDomain()
		if err != nil {
			return Record{}, err
		}
	}

	return record, nil
}
