package services

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"ppl-study-planner/internal/model"
)

func TestICSExport_GeneratesEnvelopeAndEvents(t *testing.T) {
	tasks := []model.DailyTask{
		{
			ID:          101,
			Date:        time.Date(2026, 3, 10, 0, 0, 0, 0, time.UTC),
			Category:    "Theory",
			Title:       "Area 1 Review",
			Description: "Review aerodynamics",
		},
		{
			ID:          102,
			Date:        time.Date(2026, 3, 11, 0, 0, 0, 0, time.UTC),
			Category:    "CFI Flights",
			Title:       "Maneuvers prep",
			Description: "Chair-fly emergency operations",
		},
	}

	result, err := ExportICS(tasks, testICSOutputDir(t))
	if err != nil {
		t.Fatalf("ExportICS failed: %v", err)
	}

	if result.EventCount != 2 {
		t.Fatalf("expected 2 events, got %d", result.EventCount)
	}

	b, err := os.ReadFile(result.Path)
	if err != nil {
		t.Fatalf("reading ICS output failed: %v", err)
	}

	content := string(b)
	mustContain := []string{
		"BEGIN:VCALENDAR",
		"END:VCALENDAR",
		"BEGIN:VEVENT",
		"END:VEVENT",
		"UID:",
		"DTSTAMP:",
		"DTSTART:",
		"SUMMARY:",
	}

	for _, token := range mustContain {
		if !strings.Contains(content, token) {
			t.Fatalf("expected ICS output to contain %q", token)
		}
	}

	if filepath.Ext(result.Path) != ".ics" {
		t.Fatalf("expected .ics output, got %s", result.Path)
	}
}

func TestICSExport_UsesUTCDateTimeFields(t *testing.T) {
	tasks := []model.DailyTask{
		{
			ID:          2,
			Date:        time.Date(2026, 5, 1, 23, 59, 0, 0, time.FixedZone("UTC+9", 9*3600)),
			Title:       "UTC export check",
			Description: "Ensure Z suffix",
		},
	}

	result, err := ExportICS(tasks, testICSOutputDir(t))
	if err != nil {
		t.Fatalf("ExportICS failed: %v", err)
	}

	b, err := os.ReadFile(result.Path)
	if err != nil {
		t.Fatalf("reading ICS output failed: %v", err)
	}

	content := string(b)
	startUTC := regexp.MustCompile(`DTSTART:\d{8}T\d{6}Z`)
	endUTC := regexp.MustCompile(`DTEND:\d{8}T\d{6}Z`)
	stampUTC := regexp.MustCompile(`DTSTAMP:\d{8}T\d{6}Z`)

	if !startUTC.MatchString(content) {
		t.Fatalf("expected DTSTART with UTC Z suffix")
	}
	if !endUTC.MatchString(content) {
		t.Fatalf("expected DTEND with UTC Z suffix")
	}
	if !stampUTC.MatchString(content) {
		t.Fatalf("expected DTSTAMP with UTC Z suffix")
	}
}

func TestICSDeterministicUID(t *testing.T) {
	task := model.DailyTask{
		ID:    88,
		Date:  time.Date(2026, 6, 12, 0, 0, 0, 0, time.UTC),
		Title: "Navigation prep",
	}

	uidA := deterministicTaskUID(task)
	uidB := deterministicTaskUID(task)

	if uidA != uidB {
		t.Fatalf("deterministicTaskUID should return same UID for same task, got %q and %q", uidA, uidB)
	}

	if !strings.Contains(uidA, "task-88-20260612@openppl") {
		t.Fatalf("unexpected deterministic UID format: %q", uidA)
	}
}

func testICSOutputDir(t *testing.T) string {
	t.Helper()

	dir := filepath.Join("icss", "test-artifacts", strings.ReplaceAll(t.Name(), "/", "-"))
	resolvedDir, err := ResolveArtifactOutputDir(dir)
	if err == nil {
		t.Cleanup(func() {
			_ = os.RemoveAll(resolvedDir)
		})
	}

	return dir
}
