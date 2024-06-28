package dns

import (
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
