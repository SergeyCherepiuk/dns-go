package serde

import (
	"net"

	"github.com/SergeyCherepiuk/dns-go/internal/dns/io"
	"github.com/SergeyCherepiuk/dns-go/internal/dns/types"
)

func marshalRecord(w *io.PacketWriter, record types.Record) error {
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
	case types.RecordTypeA, types.RecordTypeAAAA:
		bytes := []byte(record.Data.(net.IP))

		err = w.WriteUint16(uint16(len(bytes)))
		if err != nil {
			return err
		}

		err = w.WriteBytes(bytes)
		if err != nil {
			return err
		}

	case types.RecordTypeNS, types.RecordTypeCNAME:
		domain := record.Data.(string)

		err = w.WriteDomainWithLength(domain)
		if err != nil {
			return err
		}
	}

	return nil
}

func unmarshalRecord(r *io.PacketReader) (types.Record, error) {
	domain, err := r.ReadDomain()
	if err != nil {
		return types.Record{}, err
	}

	recordType, err := r.ReadUint16()
	if err != nil {
		return types.Record{}, err
	}

	recordClass, err := r.ReadUint16()
	if err != nil {
		return types.Record{}, err
	}

	ttl, err := r.ReadUint32()
	if err != nil {
		return types.Record{}, err
	}

	record := types.Record{
		Domain: domain,
		Type:   types.RecordType(recordType),
		Class:  types.RecordClass(recordClass),
		Ttl:    ttl,
	}

	length, err := r.ReadUint16()
	if err != nil {
		return types.Record{}, err
	}

	switch record.Type {
	case types.RecordTypeA, types.RecordTypeAAAA:
		ip, err := r.ReadBytes(int(length))
		if err != nil {
			return types.Record{}, err
		}

		record.Data = net.IP(ip)

	case types.RecordTypeNS, types.RecordTypeCNAME:
		record.Data, err = r.ReadDomain()
		if err != nil {
			return types.Record{}, err
		}
	}

	return record, nil
}
