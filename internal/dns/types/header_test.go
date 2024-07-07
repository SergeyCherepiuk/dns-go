package types

import (
	"testing"

	"github.com/SergeyCherepiuk/dns-go/internal/utils"
)

func TestHeaderConstants(t *testing.T) {
	entries := make(utils.DiffEntries, 0)

	entries = append(entries, utils.Diff(HeaderSize, 12)...)

	entries = append(entries, utils.Diff(PacketTypeQuery, 0)...)
	entries = append(entries, utils.Diff(PacketTypeResponse, 1)...)

	entries = append(entries, utils.Diff(OpcodeQuery, 0)...)
	entries = append(entries, utils.Diff(OpcodeIQuery, 1)...)
	entries = append(entries, utils.Diff(OpcodeStatus, 2)...)

	entries = append(entries, utils.Diff(ResponseCodeNoError, 0)...)
	entries = append(entries, utils.Diff(ResponseCodeFormatError, 1)...)
	entries = append(entries, utils.Diff(ResponseCodeServerFailure, 2)...)
	entries = append(entries, utils.Diff(ResponseCodeNameError, 3)...)
	entries = append(entries, utils.Diff(ResponseCodeNotImplemented, 4)...)
	entries = append(entries, utils.Diff(ResponseCodeRefused, 5)...)

	if len(entries) > 0 {
		t.Fatal(entries.String())
	}
}
