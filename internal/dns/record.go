package dns

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

	err = w.WriteUint16(uint16(len(record.Data)))
	if err != nil {
		return err
	}

	err = w.WriteBytes(record.Data)
	if err != nil {
		return err
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

	length, err := r.ReadUint16()
	if err != nil {
		return Record{}, err
	}

	data, err := r.ReadBytes(int(length))
	if err != nil {
		return Record{}, err
	}

	record := Record{
		Domain: domain,
		Type:   RecordType(recordType),
		Class:  RecordClass(recordClass),
		Ttl:    ttl,
		Data:   data,
	}

	return record, nil
}
