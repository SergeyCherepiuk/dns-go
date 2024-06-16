package dns

import (
	"testing"

	"github.com/SergeyCherepiuk/dns-go/internal/utils"
)

func TestQuestionConstants(t *testing.T) {
	entries := make(utils.DiffEntries, 0)

	entries = append(entries, utils.Diff(QuestionTypeA, 1)...)
	entries = append(entries, utils.Diff(QuestionTypeNS, 2)...)
	entries = append(entries, utils.Diff(QuestionTypeMD, 3)...)
	entries = append(entries, utils.Diff(QuestionTypeMF, 4)...)
	entries = append(entries, utils.Diff(QuestionTypeCNAME, 5)...)
	entries = append(entries, utils.Diff(QuestionTypeSOA, 6)...)
	entries = append(entries, utils.Diff(QuestionTypeMB, 7)...)
	entries = append(entries, utils.Diff(QuestionTypeMG, 8)...)
	entries = append(entries, utils.Diff(QuestionTypeMR, 9)...)
	entries = append(entries, utils.Diff(QuestionTypeNULL, 10)...)
	entries = append(entries, utils.Diff(QuestionTypeWKS, 11)...)
	entries = append(entries, utils.Diff(QuestionTypePTR, 12)...)
	entries = append(entries, utils.Diff(QuestionTypeHINFO, 13)...)
	entries = append(entries, utils.Diff(QuestionTypeMINFO, 14)...)
	entries = append(entries, utils.Diff(QuestionTypeMX, 15)...)
	entries = append(entries, utils.Diff(QuestionTypeTXT, 16)...)
	entries = append(entries, utils.Diff(QuestionTypeAXFR, 252)...)
	entries = append(entries, utils.Diff(QuestionTypeMAILB, 253)...)
	entries = append(entries, utils.Diff(QuestionTypeMAILA, 254)...)
	entries = append(entries, utils.Diff(QuestionTypeALL, 255)...)

	entries = append(entries, utils.Diff(QuestionClassIN, 1)...)
	entries = append(entries, utils.Diff(QuestionClassCS, 2)...)
	entries = append(entries, utils.Diff(QuestionClassCH, 3)...)
	entries = append(entries, utils.Diff(QuestionClassHS, 4)...)
	entries = append(entries, utils.Diff(QuestionClassALL, 255)...)

	if len(entries) > 0 {
		t.Fatal(entries.String())
	}
}

func TestUnmarshalQuestionEmptyLookup(t *testing.T) {
	var (
		bytes = []byte{
			0x06, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x03,
			0x63, 0x6f, 0x6d, 0x00, 0x00, 0x01, 0x00, 0x01,
		}
		lookup = map[int]string{}
	)

	var (
		expectedQuestion = Question{
			Domain: "google.com.",
			Type:   QuestionTypeA,
			Class:  QuestionClassIN,
		}
		expectedBytesRead = 16
	)

	actualQuestion, actualBytesRead := UnmarshalQuestion(bytes, lookup)

	if actualQuestion != expectedQuestion {
		entries := utils.Diff(actualQuestion, expectedQuestion)
		t.Fatal(entries.String())
	}

	if actualBytesRead != expectedBytesRead {
		entries := utils.Diff(actualBytesRead, expectedBytesRead)
		t.Fatal(entries.String())
	}
}

func TestUnmarshalQuestionExactDomainInLookup(t *testing.T) {
	var (
		bytes = []byte{
			0xc0, 0x0c, 0x00, 0x00, 0x01, 0x00, 0x01,
		}
		lookup = map[int]string{12: "google.com."}
	)

	var (
		expectedQuestion = Question{
			Domain: "google.com.",
			Type:   QuestionTypeA,
			Class:  QuestionClassIN,
		}
		expectedBytesRead = 7
	)

	actualQuestion, actualBytesRead := UnmarshalQuestion(bytes, lookup)

	if actualQuestion != expectedQuestion {
		entries := utils.Diff(actualQuestion, expectedQuestion)
		t.Fatal(entries.String())
	}

	if actualBytesRead != expectedBytesRead {
		entries := utils.Diff(actualBytesRead, expectedBytesRead)
		t.Fatal(entries.String())
	}
}

func TestUnmarshalQuestionPartOfDomainInLookup(t *testing.T) {
	var (
		bytes = []byte{
			0x02, 0x6d, 0x78, 0xc0, 0x0c, 0x00, 0x00, 0x0f,
			0x00, 0x01,
		}
		lookup = map[int]string{12: "google.com."}
	)

	var (
		expectedQuestion = Question{
			Domain: "mx.google.com.",
			Type:   QuestionTypeMX,
			Class:  QuestionClassIN,
		}
		expectedBytesRead = 10
	)

	actualQuestion, actualBytesRead := UnmarshalQuestion(bytes, lookup)

	if actualQuestion != expectedQuestion {
		entries := utils.Diff(actualQuestion, expectedQuestion)
		t.Fatal(entries.String())
	}

	if actualBytesRead != expectedBytesRead {
		entries := utils.Diff(actualBytesRead, expectedBytesRead)
		t.Fatal(entries.String())
	}
}
