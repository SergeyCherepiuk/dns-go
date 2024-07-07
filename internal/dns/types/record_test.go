package types

import (
	"testing"

	"github.com/SergeyCherepiuk/dns-go/internal/utils"
)

func TestRecordConstants(t *testing.T) {
	entries := make(utils.DiffEntries, 0)

	entries = append(entries, utils.Diff(RecordTypeA, 1)...)
	entries = append(entries, utils.Diff(RecordTypeNS, 2)...)
	entries = append(entries, utils.Diff(RecordTypeCNAME, 5)...)
	entries = append(entries, utils.Diff(RecordTypeMX, 15)...)
	entries = append(entries, utils.Diff(RecordTypeAAAA, 28)...)

	entries = append(entries, utils.Diff(RecordClassIN, 1)...)

	if len(entries) > 0 {
		t.Fatal(entries.String())
	}
}
