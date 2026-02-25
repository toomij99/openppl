package services

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"

	"ppl-study-planner/internal/model"
)

type googleCalendarEventsWriter interface {
	Insert(calendarID string, event *calendar.Event) (*calendar.Event, error)
}

type calendarServiceEventsWriter struct {
	service *calendar.Service
}

func (w *calendarServiceEventsWriter) Insert(calendarID string, event *calendar.Event) (*calendar.Event, error) {
	return w.service.Events.Insert(calendarID, event).Do()
}

var (
	googleAuthClientProvider  = EnsureGoogleAuthClient
	googleEventsWriterFactory = func(ctx context.Context, authResult GoogleAuthResult) (googleCalendarEventsWriter, error) {
		svc, err := calendar.NewService(ctx, option.WithHTTPClient(authResult.Client))
		if err != nil {
			return nil, err
		}
		return &calendarServiceEventsWriter{service: svc}, nil
	}
	googleRetrySleepWithContext = sleepWithContext
)

func SyncTasksToGoogleCalendar(ctx context.Context, tasks []model.DailyTask, opts GoogleCalendarSyncOptions) (GoogleCalendarSyncResult, error) {
	if len(tasks) == 0 {
		return GoogleCalendarSyncResult{}, &GoogleCalendarError{Kind: GoogleCalendarErrorValidation, Err: errors.New("no tasks available for google calendar sync")}
	}

	resolvedOpts, err := normalizeGoogleCalendarOptions(opts)
	if err != nil {
		return GoogleCalendarSyncResult{}, err
	}

	authResult, err := googleAuthClientProvider(ctx, resolvedOpts.AuthOptions)
	if err != nil {
		return GoogleCalendarSyncResult{}, &GoogleCalendarError{Kind: GoogleCalendarErrorAuth, Err: err}
	}

	writer, err := googleEventsWriterFactory(ctx, authResult)
	if err != nil {
		return GoogleCalendarSyncResult{}, &GoogleCalendarError{Kind: GoogleCalendarErrorServiceInit, Err: err}
	}

	result := GoogleCalendarSyncResult{
		CalendarID:     resolvedOpts.CalendarID,
		Attempted:      len(tasks),
		UsedCachedAuth: authResult.UsedCachedToken,
	}

	for _, task := range tasks {
		event := mapTaskToGoogleEvent(task)
		if err := insertGoogleEventWithRetry(ctx, writer, resolvedOpts.CalendarID, event, resolvedOpts.MaxRetries, resolvedOpts.InitialBackoff); err != nil {
			statusCode := googleStatusCode(err)
			result.Failures = append(result.Failures, GoogleCalendarTaskFailure{
				TaskID:     task.ID,
				TaskTitle:  strings.TrimSpace(task.Title),
				StatusCode: statusCode,
				Retryable:  isRetryableGoogleStatus(statusCode),
				Message:    err.Error(),
			})
			continue
		}

		result.Created++
	}

	result.Failed = len(result.Failures)
	return result, nil
}

func mapTaskToGoogleEvent(task model.DailyTask) *calendar.Event {
	title := strings.TrimSpace(task.Title)
	if title == "" {
		title = "Study Task"
	}

	description := strings.TrimSpace(task.Description)
	if description == "" {
		description = strings.TrimSpace(task.Category)
	}

	startUTC, endUTC := taskWindowUTC(task.Date)
	identity := deterministicGoogleTaskIdentity(task)

	return &calendar.Event{
		Summary:     title,
		Description: description,
		Start: &calendar.EventDateTime{
			DateTime: startUTC.Format(time.RFC3339),
			TimeZone: "UTC",
		},
		End: &calendar.EventDateTime{
			DateTime: endUTC.Format(time.RFC3339),
			TimeZone: "UTC",
		},
		ExtendedProperties: &calendar.EventExtendedProperties{
			Private: map[string]string{
				"openppl_task_identity": identity,
				"openppl_task_id":       fmt.Sprintf("%d", task.ID),
			},
		},
	}
}

func deterministicGoogleTaskIdentity(task model.DailyTask) string {
	datePart := task.Date.UTC().Format("20060102")
	if task.ID != 0 {
		return fmt.Sprintf("task-%d-%s", task.ID, datePart)
	}

	titlePart := strings.ToLower(strings.TrimSpace(task.Title))
	titlePart = strings.ReplaceAll(titlePart, " ", "-")
	if titlePart == "" {
		titlePart = "untitled"
	}

	return fmt.Sprintf("task-%s-%s", titlePart, datePart)
}

func insertGoogleEventWithRetry(ctx context.Context, writer googleCalendarEventsWriter, calendarID string, event *calendar.Event, maxRetries int, initialBackoff time.Duration) error {
	var lastErr error

	for attempt := 0; attempt <= maxRetries; attempt++ {
		_, err := writer.Insert(calendarID, event)
		if err == nil {
			return nil
		}

		lastErr = err
		statusCode := googleStatusCode(err)
		if !isRetryableGoogleStatus(statusCode) || attempt == maxRetries {
			return err
		}

		if sleepErr := googleRetrySleepWithContext(ctx, initialBackoff*time.Duration(1<<attempt)); sleepErr != nil {
			return sleepErr
		}
	}

	if lastErr == nil {
		return errors.New("google calendar insert failed")
	}

	return lastErr
}

func sleepWithContext(ctx context.Context, d time.Duration) error {
	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

func googleStatusCode(err error) int {
	var apiErr *googleapi.Error
	if errors.As(err, &apiErr) {
		return apiErr.Code
	}
	return 0
}

func isRetryableGoogleStatus(statusCode int) bool {
	if statusCode == 401 || statusCode == 403 || statusCode == 429 {
		return true
	}
	return statusCode >= 500 && statusCode <= 599
}
