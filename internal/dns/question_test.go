package dns

import (
	"testing"

	"github.com/SergeyCherepiuk/dns-go/internal/utils"
)

func TestQuestionConstants(t *testing.T) {
	entries := make(utils.DiffEntries, 0)

	entries = append(entries, utils.Diff(QuestionTypeA, 1)...)
	entries = append(entries, utils.Diff(QuestionTypeNS, 2)...)
	entries = append(entries, utils.Diff(QuestionTypeCNAME, 5)...)
	entries = append(entries, utils.Diff(QuestionTypeMX, 15)...)
	entries = append(entries, utils.Diff(QuestionTypeAAAA, 28)...)

	entries = append(entries, utils.Diff(QuestionClassIN, 1)...)

	if len(entries) > 0 {
		t.Fatal(entries.String())
	}
}
