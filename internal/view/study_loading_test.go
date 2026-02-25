package view

import (
	"strings"
	"testing"

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
