package motd

import (
	"bytes"
	"strings"
	"testing"

	"ppl-study-planner/internal/services"
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

func TestConfigCommand_DisablesQuizMode(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	var buf bytes.Buffer
	code := Execute([]string{"config", "quiz", "off"}, strings.NewReader(""), &buf)
	if code != 0 {
		t.Fatalf("Execute(config quiz off) = %d; want 0", code)
	}
	if !strings.Contains(buf.String(), "MOTD quiz mode: off") {
		t.Fatalf("expected config output to confirm disabled mode, got %q", buf.String())
	}

	cfg, err := services.LoadMOTDConfig()
	if err != nil {
		t.Fatalf("LoadMOTDConfig returned error: %v", err)
	}
	if cfg.QuizMode {
		t.Fatalf("expected quiz mode to be disabled")
	}
}

func TestDisplay_UsesReviewLayoutWhenQuizDisabled(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	if err := services.SaveMOTDConfig(services.MOTDConfig{QuizMode: false}); err != nil {
		t.Fatalf("SaveMOTDConfig returned error: %v", err)
	}

	var buf bytes.Buffer
	code := Execute([]string{"display"}, strings.NewReader(""), &buf)
	if code != 0 {
		t.Fatalf("Execute(display) = %d; want 0", code)
	}
	out := buf.String()
	if !strings.Contains(out, "ACS Daily Review") {
		t.Fatalf("expected review header in output, got %q", out)
	}
	if !strings.Contains(out, "Question:") || !strings.Contains(out, "Answer:") {
		t.Fatalf("expected question/answer review layout, got %q", out)
	}
	if strings.Contains(out, "Run: openppl motd quiz") {
		t.Fatalf("expected quiz prompt to be hidden when quiz mode disabled, got %q", out)
	}
	if !strings.Contains(out, "openppl motd config quiz on") {
		t.Fatalf("expected hint to re-enable quiz mode, got %q", out)
	}
}
