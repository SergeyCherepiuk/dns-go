package dns

import (
	"slices"
	"testing"

	"github.com/SergeyCherepiuk/dns-go/internal/utils"
)

func TestMarshalDomainEmptyLookup(t *testing.T) {
	var (
		domain = "google.com."
		lookup = map[int]string{}
	)

	expectedBytes := []byte{
		0x06, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x03,
		0x63, 0x6f, 0x6d, 0x00,
	}

	actualBytes := MarshalDomain(domain, lookup)

	if !slices.Equal(actualBytes, expectedBytes) {
		entries := utils.Diff(actualBytes, expectedBytes)
		t.Fatal(entries.String())
	}
}

func TestMarshalDomainExactDomainInLookup(t *testing.T) {
	var (
		domain = "google.com."
		lookup = map[int]string{12: "google.com."}
	)

	expectedBytes := []byte{
		0xc0, 0x0c,
	}

	actualBytes := MarshalDomain(domain, lookup)

	if !slices.Equal(actualBytes, expectedBytes) {
		entries := utils.Diff(actualBytes, expectedBytes)
		t.Fatal(entries.String())
	}
}

func TestMarshalDomainPartOfDomainInLookup(t *testing.T) {
	var (
		domain = "mx.google.com."
		lookup = map[int]string{12: "google.com."}
	)

	expectedBytes := []byte{
		0x02, 0x6d, 0x78, 0xc0, 0x0c,
	}

	actualBytes := MarshalDomain(domain, lookup)

	if !slices.Equal(actualBytes, expectedBytes) {
		entries := utils.Diff(actualBytes, expectedBytes)
		t.Fatal(entries.String())
	}
}

func TestUnmarshalDomainEmptyLookup(t *testing.T) {
	var (
		bytes = []byte{
			0x06, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x03,
			0x63, 0x6f, 0x6d, 0x00,
		}
		lookup = map[int]string{}
	)

	var (
		expectedDomain    = "google.com."
		expectedBytesRead = 12
	)

	actualDomain, actualBytesRead := UnmarshalDomain(bytes, lookup)

	if actualDomain != expectedDomain {
		entries := utils.Diff(actualDomain, expectedDomain)
		t.Fatal(entries.String())
	}

	if actualBytesRead != expectedBytesRead {
		entries := utils.Diff(actualBytesRead, expectedBytesRead)
		t.Fatal(entries.String())
	}
}

func TestUnmarshalDomainExactDomainInLookup(t *testing.T) {
	var (
		bytes = []byte{
			0xc0, 0x0c,
		}
		lookup = map[int]string{12: "google.com."}
	)

	var (
		expectedDomain    = "google.com."
		expectedBytesRead = 2
	)

	actualDomain, actualBytesRead := UnmarshalDomain(bytes, lookup)

	if actualDomain != expectedDomain {
		entries := utils.Diff(actualDomain, expectedDomain)
		t.Fatal(entries.String())
	}

	if actualBytesRead != expectedBytesRead {
		entries := utils.Diff(actualBytesRead, expectedBytesRead)
		t.Fatal(entries.String())
	}
}

func TestUnmarshalDomainPartOfDomainInLookup(t *testing.T) {
	var (
		bytes = []byte{
			0x02, 0x6d, 0x78, 0xc0, 0x0c,
		}
		lookup = map[int]string{12: "google.com."}
	)

	var (
		expectedDomain    = "mx.google.com."
		expectedBytesRead = 5
	)

	actualDomain, actualBytesRead := UnmarshalDomain(bytes, lookup)

	if actualDomain != expectedDomain {
		entries := utils.Diff(actualDomain, expectedDomain)
		t.Fatal(entries.String())
	}

	if actualBytesRead != expectedBytesRead {
		entries := utils.Diff(actualBytesRead, expectedBytesRead)
		t.Fatal(entries.String())
	}
}
