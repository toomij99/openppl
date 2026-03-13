package motd

import (
	"strings"
	"testing"
)

func TestStudyTipBySection(t *testing.T) {
	if got := studyTip("K"); !strings.Contains(got, "teach-back") {
		t.Fatalf("expected knowledge tip to mention teach-back, got %q", got)
	}
	if got := studyTip("R"); !strings.Contains(got, "PAVE") {
		t.Fatalf("expected risk tip to mention PAVE, got %q", got)
	}
	if got := studyTip("S"); !strings.Contains(got, "Chair-fly") {
		t.Fatalf("expected skill tip to mention chair-fly, got %q", got)
	}
}

func TestStudyInsightBySectionAndCategory(t *testing.T) {
	if got := studyInsight("R", "Theory"); !strings.Contains(got, "trigger point") {
		t.Fatalf("expected risk/theory insight to mention trigger point, got %q", got)
	}
	if got := studyInsight("S", "CFI Flights"); !strings.Contains(got, "stable setup") {
		t.Fatalf("expected skill insight to mention setup, got %q", got)
	}
}
