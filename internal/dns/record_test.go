package dns

import (
	"reflect"
	"slices"
	"testing"

	"github.com/SergeyCherepiuk/dns-go/internal/utils"
)

func TestRecordConstants(t *testing.T) {
	entries := make(utils.DiffEntries, 0)

	entries = append(entries, utils.Diff(RecordTypeA, 1)...)
	entries = append(entries, utils.Diff(RecordTypeNS, 2)...)
	entries = append(entries, utils.Diff(RecordTypeMD, 3)...)
	entries = append(entries, utils.Diff(RecordTypeMF, 4)...)
	entries = append(entries, utils.Diff(RecordTypeCNAME, 5)...)
	entries = append(entries, utils.Diff(RecordTypeSOA, 6)...)
	entries = append(entries, utils.Diff(RecordTypeMB, 7)...)
	entries = append(entries, utils.Diff(RecordTypeMG, 8)...)
	entries = append(entries, utils.Diff(RecordTypeMR, 9)...)
	entries = append(entries, utils.Diff(RecordTypeNULL, 10)...)
	entries = append(entries, utils.Diff(RecordTypeWKS, 11)...)
	entries = append(entries, utils.Diff(RecordTypePTR, 12)...)
	entries = append(entries, utils.Diff(RecordTypeHINFO, 13)...)
	entries = append(entries, utils.Diff(RecordTypeMINFO, 14)...)
	entries = append(entries, utils.Diff(RecordTypeMX, 15)...)
	entries = append(entries, utils.Diff(RecordTypeTXT, 16)...)

	entries = append(entries, utils.Diff(RecordClassIN, 1)...)
	entries = append(entries, utils.Diff(RecordClassCS, 2)...)
	entries = append(entries, utils.Diff(RecordClassCH, 3)...)
	entries = append(entries, utils.Diff(RecordClassHS, 4)...)

	if len(entries) > 0 {
		t.Fatal(entries.String())
	}
}

func TestMarshalRecordEmptyLookup(t *testing.T) {
	var (
		record = Record{
			Domain: "google.com.",
			Type:   RecordTypeA,
			Class:  RecordClassIN,
			Ttl:    86400,
			Data:   []byte{142, 251, 37, 110},
		}
		lookup = map[int]string{}
	)

	expectedBytes := []byte{
		0x06, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x03,
		0x63, 0x6f, 0x6d, 0x00, 0x00, 0x01, 0x00, 0x01,
		0x00, 0x01, 0x51, 0x80, 0x00, 0x04, 0x8e, 0xfb,
		0x25, 0x6e,
	}

	actualBytes := marshalRecord(record, lookup)

	if !slices.Equal(actualBytes, expectedBytes) {
		entries := utils.Diff(actualBytes, expectedBytes)
		t.Fatal(entries.String())
	}
}

func TestMarshalRecordExactDomainInLookup(t *testing.T) {
	var (
		record = Record{
			Domain: "google.com.",
			Type:   RecordTypeA,
			Class:  RecordClassIN,
			Ttl:    86400,
			Data:   []byte{142, 251, 37, 110},
		}
		lookup = map[int]string{12: "google.com."}
	)

	expectedBytes := []byte{
		0xc0, 0x0c, 0x00, 0x01, 0x00, 0x01, 0x00, 0x01,
		0x51, 0x80, 0x00, 0x04, 0x8e, 0xfb, 0x25, 0x6e,
	}

	actualBytes := marshalRecord(record, lookup)

	if !slices.Equal(actualBytes, expectedBytes) {
		entries := utils.Diff(actualBytes, expectedBytes)
		t.Fatal(entries.String())
	}
}

func TestMarshalRecordPartOfDomainInLookup(t *testing.T) {
	var (
		record = Record{
			Domain: "mx.google.com.",
			Type:   RecordTypeMX,
			Class:  RecordClassIN,
			Ttl:    86400,
			Data:   []byte{142, 251, 37, 110},
		}
		lookup = map[int]string{12: "google.com."}
	)

	expectedBytes := []byte{
		0x02, 0x6d, 0x78, 0xc0, 0x0c, 0x00, 0x0f, 0x00,
		0x01, 0x00, 0x01, 0x51, 0x80, 0x00, 0x04, 0x8e,
		0xfb, 0x25, 0x6e,
	}

	actualBytes := marshalRecord(record, lookup)

	if !slices.Equal(actualBytes, expectedBytes) {
		entries := utils.Diff(actualBytes, expectedBytes)
		t.Fatal(entries.String())
	}
}

func TestUnmarshalRecordEmptyLookup(t *testing.T) {
	var (
		bytes = []byte{
			0x06, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x03,
			0x63, 0x6f, 0x6d, 0x00, 0x00, 0x01, 0x00, 0x01,
			0x00, 0x01, 0x51, 0x80, 0x00, 0x04, 0x8e, 0xfb,
			0x25, 0x6e,
		}
		lookup = map[int]string{}
	)

	var (
		expectedRecord = Record{
			Domain: "google.com.",
			Type:   RecordTypeA,
			Class:  RecordClassIN,
			Ttl:    86400,
			Data:   []byte{142, 251, 37, 110},
		}
		expectedBytesRead = 26
	)

	actualRecord, actualBytesRead := unmarshalRecord(bytes, lookup)

	if !reflect.DeepEqual(actualRecord, expectedRecord) {
		entries := utils.Diff(actualRecord, expectedRecord)
		t.Fatal(entries.String())
	}

	if actualBytesRead != expectedBytesRead {
		entries := utils.Diff(actualBytesRead, expectedBytesRead)
		t.Fatal(entries.String())
	}
}

func TestUnmarshalRecordExactDomainInLookup(t *testing.T) {
	var (
		bytes = []byte{
			0xc0, 0x0c, 0x00, 0x01, 0x00, 0x01, 0x00, 0x01,
			0x51, 0x80, 0x00, 0x04, 0x8e, 0xfb, 0x25, 0x6e,
		}
		lookup = map[int]string{12: "google.com."}
	)

	var (
		expectedRecord = Record{
			Domain: "google.com.",
			Type:   RecordTypeA,
			Class:  RecordClassIN,
			Ttl:    86400,
			Data:   []byte{142, 251, 37, 110},
		}
		expectedBytesRead = 16
	)

	actualRecord, actualBytesRead := unmarshalRecord(bytes, lookup)

	if !reflect.DeepEqual(actualRecord, expectedRecord) {
		entries := utils.Diff(actualRecord, expectedRecord)
		t.Fatal(entries.String())
	}

	if actualBytesRead != expectedBytesRead {
		entries := utils.Diff(actualBytesRead, expectedBytesRead)
		t.Fatal(entries.String())
	}
}

func TestUnmarshalRecordPartOfDomainInLookup(t *testing.T) {
	var (
		bytes = []byte{
			0x02, 0x6d, 0x78, 0xc0, 0x0c, 0x00, 0x0f, 0x00,
			0x01, 0x00, 0x01, 0x51, 0x80, 0x00, 0x04, 0x8e,
			0xfb, 0x25, 0x6e,
		}
		lookup = map[int]string{12: "google.com."}
	)

	var (
		expectedRecord = Record{
			Domain: "mx.google.com.",
			Type:   RecordTypeMX,
			Class:  RecordClassIN,
			Ttl:    86400,
			Data:   []byte{142, 251, 37, 110},
		}
		expectedBytesRead = 19
	)

	actualRecord, actualBytesRead := unmarshalRecord(bytes, lookup)

	if !reflect.DeepEqual(actualRecord, expectedRecord) {
		entries := utils.Diff(actualRecord, expectedRecord)
		t.Fatal(entries.String())
	}

	if actualBytesRead != expectedBytesRead {
		entries := utils.Diff(actualBytesRead, expectedBytesRead)
		t.Fatal(entries.String())
	}
}
