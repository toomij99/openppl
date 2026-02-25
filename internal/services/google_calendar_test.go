package services

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"

	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/googleapi"

	"ppl-study-planner/internal/model"
)

func TestGoogleCalendar_MapsTaskToEventPayload(t *testing.T) {
	task := model.DailyTask{
		ID:          42,
		Date:        time.Date(2026, 4, 22, 12, 0, 0, 0, time.FixedZone("UTC+5", 5*3600)),
		Category:    "Theory",
		Title:       "Weather briefing",
		Description: "Review winds aloft",
	}

	event := mapTaskToGoogleEvent(task)
	if event == nil {
		t.Fatal("expected event")
	}
	if event.Summary != task.Title {
		t.Fatalf("expected summary %q, got %q", task.Title, event.Summary)
	}
	if event.Description != task.Description {
		t.Fatalf("expected description %q, got %q", task.Description, event.Description)
	}
	if event.Start == nil || event.End == nil {
		t.Fatal("expected start and end to be set")
	}
	if event.Start.TimeZone != "UTC" || event.End.TimeZone != "UTC" {
		t.Fatal("expected UTC timezone on start/end")
	}
	if event.ExtendedProperties == nil || event.ExtendedProperties.Private == nil {
		t.Fatal("expected extended private properties")
	}
	if got := event.ExtendedProperties.Private["openppl_task_identity"]; got != "task-42-20260422" {
		t.Fatalf("unexpected deterministic identity: %q", got)
	}
}

func TestGoogleCalendar_RetriesRetryableFailures(t *testing.T) {
	tasks := []model.DailyTask{{ID: 1, Date: time.Date(2026, 3, 10, 0, 0, 0, 0, time.UTC), Title: "Task"}}
	writer := &fakeGoogleCalendarWriter{
		errorsByAttempt: []error{
			&googleapi.Error{Code: 500, Message: "temporary failure"},
			nil,
		},
	}

	restore := stubGoogleCalendarDependencies(t, writer)
	defer restore()

	result, err := SyncTasksToGoogleCalendar(context.Background(), tasks, GoogleCalendarSyncOptions{MaxRetries: 2, InitialBackoff: time.Millisecond})
	if err != nil {
		t.Fatalf("SyncTasksToGoogleCalendar returned error: %v", err)
	}

	if result.Created != 1 {
		t.Fatalf("expected 1 created event, got %d", result.Created)
	}
	if result.Failed != 0 {
		t.Fatalf("expected 0 failures, got %d", result.Failed)
	}
	if writer.insertCalls != 2 {
		t.Fatalf("expected 2 insert calls (retry), got %d", writer.insertCalls)
	}
}

func TestGoogleCalendar_DoesNotRetryPermanentFailure(t *testing.T) {
	tasks := []model.DailyTask{{ID: 9, Date: time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC), Title: "Bad task"}}
	writer := &fakeGoogleCalendarWriter{
		errorsByAttempt: []error{
			&googleapi.Error{Code: 400, Message: "invalid request"},
		},
	}

	restore := stubGoogleCalendarDependencies(t, writer)
	defer restore()

	result, err := SyncTasksToGoogleCalendar(context.Background(), tasks, GoogleCalendarSyncOptions{MaxRetries: 3, InitialBackoff: time.Millisecond})
	if err != nil {
		t.Fatalf("SyncTasksToGoogleCalendar returned error: %v", err)
	}

	if result.Created != 0 {
		t.Fatalf("expected 0 created events, got %d", result.Created)
	}
	if result.Failed != 1 || len(result.Failures) != 1 {
		t.Fatalf("expected exactly one failure entry, got failed=%d entries=%d", result.Failed, len(result.Failures))
	}
	if result.Failures[0].Retryable {
		t.Fatal("expected permanent failure to be non-retryable")
	}
	if writer.insertCalls != 1 {
		t.Fatalf("expected one insert attempt for permanent error, got %d", writer.insertCalls)
	}
}

func TestGoogleCalendar_ReportsPartialFailures(t *testing.T) {
	tasks := []model.DailyTask{
		{ID: 1, Date: time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC), Title: "Task A"},
		{ID: 2, Date: time.Date(2026, 4, 2, 0, 0, 0, 0, time.UTC), Title: "Task B"},
	}
	writer := &fakeGoogleCalendarWriter{
		errorsByAttempt: []error{nil, &googleapi.Error{Code: 429, Message: "quota exceeded"}, &googleapi.Error{Code: 429, Message: "quota exceeded"}},
	}

	restore := stubGoogleCalendarDependencies(t, writer)
	defer restore()

	result, err := SyncTasksToGoogleCalendar(context.Background(), tasks, GoogleCalendarSyncOptions{MaxRetries: 1, InitialBackoff: time.Millisecond})
	if err != nil {
		t.Fatalf("SyncTasksToGoogleCalendar returned error: %v", err)
	}

	if result.Created != 1 {
		t.Fatalf("expected 1 created event, got %d", result.Created)
	}
	if result.Failed != 1 || len(result.Failures) != 1 {
		t.Fatalf("expected one failed item, got failed=%d entries=%d", result.Failed, len(result.Failures))
	}
	failure := result.Failures[0]
	if failure.TaskID != 2 {
		t.Fatalf("expected failed task ID 2, got %d", failure.TaskID)
	}
	if failure.StatusCode != 429 {
		t.Fatalf("expected status code 429, got %d", failure.StatusCode)
	}
	if !failure.Retryable {
		t.Fatal("expected 429 failure to be marked retryable")
	}
}

func TestGoogleCalendar_MapsAuthFailureToTypedError(t *testing.T) {
	originalProvider := googleAuthClientProvider
	googleAuthClientProvider = func(ctx context.Context, opts GoogleAuthOptions) (GoogleAuthResult, error) {
		return GoogleAuthResult{}, errors.New("auth unavailable")
	}
	t.Cleanup(func() { googleAuthClientProvider = originalProvider })

	_, err := SyncTasksToGoogleCalendar(context.Background(), []model.DailyTask{{ID: 1, Date: time.Now()}}, GoogleCalendarSyncOptions{})
	if err == nil {
		t.Fatal("expected auth error")
	}

	var syncErr *GoogleCalendarError
	if !errors.As(err, &syncErr) {
		t.Fatalf("expected GoogleCalendarError, got %T", err)
	}
	if syncErr.Kind != GoogleCalendarErrorAuth {
		t.Fatalf("expected auth error kind, got %q", syncErr.Kind)
	}
}

type fakeGoogleCalendarWriter struct {
	insertCalls     int
	errorsByAttempt []error
	events          []*calendar.Event
}

func (f *fakeGoogleCalendarWriter) Insert(calendarID string, event *calendar.Event) (*calendar.Event, error) {
	f.insertCalls++
	f.events = append(f.events, event)

	idx := f.insertCalls - 1
	if idx < len(f.errorsByAttempt) && f.errorsByAttempt[idx] != nil {
		return nil, f.errorsByAttempt[idx]
	}

	return &calendar.Event{Id: fmt.Sprintf("event-%d", f.insertCalls)}, nil
}

func stubGoogleCalendarDependencies(t *testing.T, writer googleCalendarEventsWriter) func() {
	t.Helper()

	originalProvider := googleAuthClientProvider
	originalFactory := googleEventsWriterFactory
	originalSleep := googleRetrySleepWithContext

	googleAuthClientProvider = func(ctx context.Context, opts GoogleAuthOptions) (GoogleAuthResult, error) {
		return GoogleAuthResult{Client: &http.Client{}, UsedCachedToken: true}, nil
	}
	googleEventsWriterFactory = func(ctx context.Context, authResult GoogleAuthResult) (googleCalendarEventsWriter, error) {
		return writer, nil
	}
	googleRetrySleepWithContext = func(context.Context, time.Duration) error { return nil }

	return func() {
		googleAuthClientProvider = originalProvider
		googleEventsWriterFactory = originalFactory
		googleRetrySleepWithContext = originalSleep
	}
}
