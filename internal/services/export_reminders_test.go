package services

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"ppl-study-planner/internal/model"
)

func TestRemindersExport_BuildsSafeArgsAndCountsCreated(t *testing.T) {
	tasks := []model.DailyTask{
		{
			Date:        time.Date(2026, 3, 15, 0, 0, 0, 0, time.UTC),
			Title:       "Area 3 weather",
			Description: "Review METAR/TAF",
		},
	}

	var gotName string
	var gotArgs []string
	runner := func(ctx context.Context, name string, args ...string) ([]byte, error) {
		gotName = name
		gotArgs = append([]string{}, args...)
		return []byte("ok"), nil
	}

	result, err := exportAppleRemindersWithRunner(tasks, RemindersExportOptions{ListName: "Flight Training", Timeout: 3 * time.Second}, runner)
	if err != nil {
		t.Fatalf("exportAppleRemindersWithRunner failed: %v", err)
	}

	if result.Created != 1 {
		t.Fatalf("expected created=1, got %d", result.Created)
	}
	if result.ListName != "Flight Training" {
		t.Fatalf("unexpected list name: %q", result.ListName)
	}

	if gotName != "osascript" {
		t.Fatalf("expected osascript command, got %q", gotName)
	}
	if len(gotArgs) < 10 {
		t.Fatalf("expected at least 10 args, got %d", len(gotArgs))
	}
	if gotArgs[0] != "-e" || gotArgs[2] != "--" {
		t.Fatalf("unexpected osascript argument shape: %v", gotArgs)
	}
	if strings.Contains(gotArgs[len(gotArgs)-1], "osascript") || strings.Contains(gotArgs[len(gotArgs)-1], ";") {
		t.Fatalf("expected notes/title to be discrete args, got suspicious payload: %q", gotArgs[len(gotArgs)-1])
	}
}

func TestRemindersExport_MapsTimeoutError(t *testing.T) {
	tasks := []model.DailyTask{{Date: time.Now(), Title: "Task"}}

	runner := func(ctx context.Context, name string, args ...string) ([]byte, error) {
		return []byte(""), context.DeadlineExceeded
	}

	_, err := exportAppleRemindersWithRunner(tasks, RemindersExportOptions{Timeout: time.Second}, runner)
	if err == nil {
		t.Fatal("expected timeout error")
	}

	var exportErr *RemindersExportError
	if !errors.As(err, &exportErr) {
		t.Fatalf("expected RemindersExportError, got %T", err)
	}
	if exportErr.Kind != "timeout" {
		t.Fatalf("expected timeout kind, got %q", exportErr.Kind)
	}
}

func TestRemindersExport_MapsPermissionError(t *testing.T) {
	tasks := []model.DailyTask{{Date: time.Now(), Title: "Task"}}

	runner := func(ctx context.Context, name string, args ...string) ([]byte, error) {
		return []byte("Not authorized to send Apple events to Reminders (-1743)"), errors.New("exit status 1")
	}

	_, err := exportAppleRemindersWithRunner(tasks, RemindersExportOptions{Timeout: time.Second}, runner)
	if err == nil {
		t.Fatal("expected permission error")
	}

	var exportErr *RemindersExportError
	if !errors.As(err, &exportErr) {
		t.Fatalf("expected RemindersExportError, got %T", err)
	}
	if exportErr.Kind != "permission" {
		t.Fatalf("expected permission kind, got %q", exportErr.Kind)
	}
}
