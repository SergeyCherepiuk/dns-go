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
