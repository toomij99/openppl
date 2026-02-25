package services

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"ppl-study-planner/internal/model"
)

func TestOpenCodeBotExport_WritesVersionedPayloadUnderICSS(t *testing.T) {
	fixedNow := time.Date(2026, 2, 25, 14, 0, 0, 0, time.UTC)
	originalNow := opencodeNow
	opencodeNow = func() time.Time { return fixedNow }
	t.Cleanup(func() { opencodeNow = originalNow })

	tasks := []model.DailyTask{
		{ID: 2, Date: time.Date(2026, 3, 11, 9, 0, 0, 0, time.UTC), Category: "Theory", Title: "Area 3", Description: "Weather"},
		{ID: 1, Date: time.Date(2026, 3, 10, 9, 0, 0, 0, time.UTC), Category: "CFI Flights", Title: "Maneuvers", Description: "Stalls"},
	}

	result, err := ExportOpenCodeBotTasks(tasks, OpenCodeBotExportOptions{OutputDir: "icss/test-bot-export"})
	if err != nil {
		t.Fatalf("ExportOpenCodeBotTasks failed: %v", err)
	}

	if result.Version != defaultOpenCodeBotSchemaVersion {
		t.Fatalf("expected default version %q, got %q", defaultOpenCodeBotSchemaVersion, result.Version)
	}
	if result.TaskCount != 2 {
		t.Fatalf("expected task count 2, got %d", result.TaskCount)
	}
	if !strings.Contains(result.Path, string(filepath.Separator)+"icss"+string(filepath.Separator)) {
		t.Fatalf("expected export path under icss, got %q", result.Path)
	}
	t.Cleanup(func() {
		_ = os.RemoveAll(filepath.Dir(result.Path))
	})

	b, err := os.ReadFile(result.Path)
	if err != nil {
		t.Fatalf("read export file: %v", err)
	}

	var payload openCodeBotPayload
	if err := json.Unmarshal(b, &payload); err != nil {
		t.Fatalf("unmarshal payload: %v", err)
	}

	if payload.Version != "v1" {
		t.Fatalf("expected version v1, got %q", payload.Version)
	}
	if payload.GeneratedAt != fixedNow.Format(time.RFC3339) {
		t.Fatalf("expected generated_at %q, got %q", fixedNow.Format(time.RFC3339), payload.GeneratedAt)
	}
	if len(payload.Tasks) != 2 {
		t.Fatalf("expected 2 payload tasks, got %d", len(payload.Tasks))
	}

	if payload.Tasks[0].TaskID != 1 {
		t.Fatalf("expected deterministic sorting by date/id, first task id should be 1, got %d", payload.Tasks[0].TaskID)
	}
	if payload.Tasks[0].Metadata.Source != "openppl" {
		t.Fatalf("expected metadata source openppl, got %q", payload.Tasks[0].Metadata.Source)
	}
}

func TestOpenCodeBotExport_GoldenSchemaContract(t *testing.T) {
	fixedNow := time.Date(2026, 2, 25, 14, 0, 0, 0, time.UTC)
	originalNow := opencodeNow
	opencodeNow = func() time.Time { return fixedNow }
	t.Cleanup(func() { opencodeNow = originalNow })

	tasks := []model.DailyTask{
		{ID: 10, Date: time.Date(2026, 3, 12, 18, 30, 0, 0, time.UTC), Category: "Theory", Title: "Area 4", Description: "Cross-country planning", Completed: true},
	}

	result, err := ExportOpenCodeBotTasks(tasks, OpenCodeBotExportOptions{OutputDir: "icss/test-bot-contract", Version: "v1"})
	if err != nil {
		t.Fatalf("ExportOpenCodeBotTasks failed: %v", err)
	}
	t.Cleanup(func() {
		_ = os.RemoveAll(filepath.Dir(result.Path))
	})

	b, err := os.ReadFile(result.Path)
	if err != nil {
		t.Fatalf("read export file: %v", err)
	}

	expected := `{
  "version": "v1",
  "generated_at": "2026-02-25T14:00:00Z",
  "tasks": [
    {
      "task_id": 10,
      "title": "Area 4",
      "category": "Theory",
      "description": "Cross-country planning",
      "due_at": "2026-03-12T18:30:00Z",
      "completed": true,
      "metadata": {
        "source": "openppl",
        "identity": "task-10-20260312"
      }
    }
  ]
}`

	if strings.TrimSpace(string(b)) != strings.TrimSpace(expected) {
		t.Fatalf("schema drift detected\nexpected:\n%s\n\nactual:\n%s", expected, string(b))
	}
}
