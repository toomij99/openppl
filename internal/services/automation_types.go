package services

import (
	"fmt"
	"time"
)

const (
	AutomationVersionV1 = "v1"

	AutomationResultStateOK       = "ok"
	AutomationResultStateError    = "error"
	AutomationResultStateExecuted = "executed"
	AutomationResultStateReplayed = "replayed"
	AutomationResultStateRejected = "rejected"
)

type AutomationError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type AutomationStatusSummary struct {
	TotalTasks     int `json:"total_tasks"`
	CompletedTasks int `json:"completed_tasks"`
	PendingTasks   int `json:"pending_tasks"`
}

type AutomationStatusTask struct {
	Date     string `json:"date"`
	Category string `json:"category"`
	Title    string `json:"title"`
}

type AutomationStatusPayload struct {
	CheckrideDate string                  `json:"checkride_date,omitempty"`
	Summary       AutomationStatusSummary `json:"summary"`
	NextTasks     []AutomationStatusTask  `json:"next_tasks"`
}

type AutomationStatusResponse struct {
	Version     string                   `json:"version"`
	ResultState string                   `json:"result_state"`
	Timestamp   string                   `json:"timestamp"`
	Status      *AutomationStatusPayload `json:"status,omitempty"`
	Error       *AutomationError         `json:"error,omitempty"`
}

type AutomationActionRequest struct {
	Name       string            `json:"name"`
	RequestID  string            `json:"request_id"`
	ActorScope string            `json:"actor_scope"`
	Args       map[string]string `json:"args,omitempty"`
}

type AutomationActionPayload struct {
	ActionName      string `json:"action_name"`
	RequestID       string `json:"request_id"`
	ActorScope      string `json:"actor_scope"`
	ReminderTitle   string `json:"reminder_title,omitempty"`
	ReminderDueDate string `json:"reminder_due_date,omitempty"`
	TargetList      string `json:"target_list,omitempty"`
	CreatedCount    int    `json:"created_count"`
}

type AutomationActionResponse struct {
	Version     string                   `json:"version"`
	ResultState string                   `json:"result_state"`
	Timestamp   string                   `json:"timestamp"`
	Action      *AutomationActionPayload `json:"action,omitempty"`
	Error       *AutomationError         `json:"error,omitempty"`
}

type AutomationCommandError struct {
	Kind string
	Code string
	Err  error
}

func (e *AutomationCommandError) Error() string {
	if e == nil {
		return "automation command failed"
	}
	if e.Err == nil {
		return fmt.Sprintf("automation command failed: %s", e.Code)
	}
	return fmt.Sprintf("automation command failed: %s: %v", e.Code, e.Err)
}

func (e *AutomationCommandError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

func newAutomationValidationError(code string, err error) error {
	return &AutomationCommandError{Kind: "validation", Code: code, Err: err}
}

func newAutomationRuntimeError(code string, err error) error {
	return &AutomationCommandError{Kind: "runtime", Code: code, Err: err}
}

func utcTimestamp(now time.Time) string {
	return now.UTC().Format(time.RFC3339)
}
