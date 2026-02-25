package services

import (
	"errors"
	"fmt"
	"time"
)

const (
	GoogleCalendarErrorValidation  = "validation"
	GoogleCalendarErrorAuth        = "auth"
	GoogleCalendarErrorServiceInit = "service_init"
)

type GoogleCalendarSyncOptions struct {
	AuthOptions    GoogleAuthOptions
	CalendarID     string
	MaxRetries     int
	InitialBackoff time.Duration
}

type GoogleCalendarTaskFailure struct {
	TaskID     uint
	TaskTitle  string
	StatusCode int
	Retryable  bool
	Message    string
}

type GoogleCalendarSyncResult struct {
	CalendarID     string
	Attempted      int
	Created        int
	Failed         int
	Failures       []GoogleCalendarTaskFailure
	UsedCachedAuth bool
}

type GoogleCalendarError struct {
	Kind string
	Err  error
}

func (e *GoogleCalendarError) Error() string {
	if e == nil {
		return "google calendar sync failed"
	}
	if e.Err == nil {
		return fmt.Sprintf("google calendar sync failed: %s", e.Kind)
	}
	return fmt.Sprintf("google calendar sync failed: %s: %v", e.Kind, e.Err)
}

func (e *GoogleCalendarError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

func normalizeGoogleCalendarOptions(opts GoogleCalendarSyncOptions) (GoogleCalendarSyncOptions, error) {
	if opts.MaxRetries < 0 {
		return GoogleCalendarSyncOptions{}, &GoogleCalendarError{Kind: GoogleCalendarErrorValidation, Err: errors.New("max retries cannot be negative")}
	}

	if opts.CalendarID == "" {
		opts.CalendarID = "primary"
	}

	if opts.InitialBackoff <= 0 {
		opts.InitialBackoff = 500 * time.Millisecond
	}

	return opts, nil
}
