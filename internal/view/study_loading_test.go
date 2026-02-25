package view

import (
	"errors"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"ppl-study-planner/internal/model"
	"ppl-study-planner/internal/services"
)

func TestStudyLoadingStateTransitions(t *testing.T) {
	sv := &StudyView{
		tasks: []model.DailyTask{
			{ID: 1, Category: "Theory", Title: "Review regulations"},
		},
		filteredTasks: []model.DailyTask{
			{ID: 1, Category: "Theory", Title: "Review regulations"},
		},
	}

	cmd := sv.exportICS()
	if cmd == nil {
		t.Fatal("expected exportICS to return async command")
	}
	if !sv.operation.loading {
		t.Fatal("expected loading state after export trigger")
	}
	if sv.operation.label != "ICS export" {
		t.Fatalf("expected operation label %q, got %q", "ICS export", sv.operation.label)
	}
	if sv.status.severity != studyStatusSeverityInfo {
		t.Fatalf("expected info severity while loading, got %q", sv.status.severity)
	}

	_, _ = sv.Update(icsExportDoneMsg{result: services.ICSExportResult{Path: "exports/ppl.ics", EventCount: 1}})
	if sv.operation.loading {
		t.Fatal("expected loading state to clear after completion")
	}
	if sv.status.severity != studyStatusSeveritySuccess {
		t.Fatalf("expected success status after completion, got %q", sv.status.severity)
	}
	if !strings.Contains(sv.status.message, "ICS export complete") {
		t.Fatalf("expected completion message, got %q", sv.status.message)
	}
}

func TestStudyLoadingIndicatorRendering(t *testing.T) {
	sv := &StudyView{
		tasks:         []model.DailyTask{{ID: 1, Category: "Theory", Title: "Review regulations"}},
		filteredTasks: []model.DailyTask{{ID: 1, Category: "Theory", Title: "Review regulations"}},
	}

	_ = sv.exportReminders()
	renderedLoading := sv.View()
	if !strings.Contains(renderedLoading, "Operation in progress: Reminders export") {
		t.Fatalf("expected loading hint in view, got %q", renderedLoading)
	}
	if !strings.Contains(renderedLoading, "repeat e/r/g/o is ignored") {
		t.Fatalf("expected duplicate-key hint in view, got %q", renderedLoading)
	}
	if !strings.Contains(renderedLoading, "Exporting Apple Reminders") {
		t.Fatalf("expected loading status text in view, got %q", renderedLoading)
	}

	_, _ = sv.Update(remindersExportDoneMsg{err: errors.New("boom")})
	renderedDone := sv.View()
	if strings.Contains(renderedDone, "Operation in progress") {
		t.Fatalf("expected loading hint to clear after completion, got %q", renderedDone)
	}
	if !strings.Contains(renderedDone, "Reminders export failed") {
		t.Fatalf("expected completion status after loading, got %q", renderedDone)
	}
}

func TestStudyDuplicateOperationSuppression(t *testing.T) {
	sv := &StudyView{
		tasks: []model.DailyTask{{ID: 1, Category: "Theory", Title: "Review regulations"}},
	}

	firstCmd := sv.exportICS()
	if firstCmd == nil {
		t.Fatal("expected first command to start operation")
	}

	for _, key := range []string{"e", "r", "g", "o"} {
		_, cmd := sv.handleNav(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(key)})
		if cmd != nil {
			t.Fatalf("expected duplicate key %q to be suppressed", key)
		}
		if sv.status.severity != studyStatusSeverityWarning {
			t.Fatalf("expected warning status for duplicate key %q, got %q", key, sv.status.severity)
		}
		if !strings.Contains(sv.status.message, "already in progress") {
			t.Fatalf("expected duplicate warning for key %q, got %q", key, sv.status.message)
		}
	}

	_, _ = sv.Update(icsExportDoneMsg{result: services.ICSExportResult{Path: "exports/ppl.ics", EventCount: 1}})
	if sv.operation.loading {
		t.Fatal("expected loading to clear after completion")
	}
	if sv.operation.label != "" {
		t.Fatalf("expected operation label to reset, got %q", sv.operation.label)
	}
}
