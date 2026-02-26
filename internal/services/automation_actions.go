package services

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"

	"ppl-study-planner/internal/model"
)

type reminderExporter func(tasks []model.DailyTask, opts RemindersExportOptions) (RemindersExportResult, error)

var automationReminderExporter reminderExporter = ExportAppleReminders

type AutomationActionService struct {
	db       *gorm.DB
	exporter reminderExporter
	now      func() time.Time
}

func NewAutomationActionService(database *gorm.DB) *AutomationActionService {
	return &AutomationActionService{
		db:       database,
		exporter: automationReminderExporter,
		now:      time.Now,
	}
}

// SetAutomationReminderExporterForTest allows tests to replace the reminder side effect.
func SetAutomationReminderExporterForTest(exporter func(tasks []model.DailyTask, opts RemindersExportOptions) (RemindersExportResult, error)) func() {
	previous := automationReminderExporter
	if exporter != nil {
		automationReminderExporter = exporter
	}
	return func() {
		automationReminderExporter = previous
	}
}

func (s *AutomationActionService) WithExporter(exporter reminderExporter) *AutomationActionService {
	if exporter != nil {
		s.exporter = exporter
	}
	return s
}

func (s *AutomationActionService) WithClock(clock func() time.Time) *AutomationActionService {
	if clock != nil {
		s.now = clock
	}
	return s
}

func (s *AutomationActionService) RunAutomationAction(req AutomationActionRequest) (AutomationActionResponse, error) {
	if s == nil || s.db == nil {
		return AutomationActionResponse{}, newAutomationValidationError("action.db_required", errors.New("database is required"))
	}

	name := strings.TrimSpace(strings.ToLower(req.Name))
	if name == "" {
		return AutomationActionResponse{}, newAutomationValidationError("action.name_required", errors.New("action name is required"))
	}
	if req.RequestID == "" {
		return AutomationActionResponse{}, newAutomationValidationError("action.request_id_required", errors.New("request_id is required"))
	}
	if name != "remind" {
		return AutomationActionResponse{}, newAutomationValidationError("action.not_allowlisted", fmt.Errorf("unsupported action %q", req.Name))
	}

	if req.ActorScope == "" {
		req.ActorScope = "default"
	}

	hash, err := normalizedArgsHash(req.Args)
	if err != nil {
		return AutomationActionResponse{}, newAutomationValidationError("action.invalid_args", err)
	}

	var existing model.AutomationIdempotency
	err = s.db.Where(
		"action_name = ? AND request_id = ? AND args_hash = ? AND actor_scope = ?",
		name,
		req.RequestID,
		hash,
		req.ActorScope,
	).First(&existing).Error
	if err == nil {
		var prior AutomationActionResponse
		if unmarshalErr := json.Unmarshal([]byte(existing.ResponseJSON), &prior); unmarshalErr != nil {
			return AutomationActionResponse{}, newAutomationRuntimeError("action.replay_decode_failed", unmarshalErr)
		}
		prior.ResultState = AutomationResultStateReplayed
		prior.Timestamp = utcTimestamp(s.now())
		return prior, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return AutomationActionResponse{}, newAutomationRuntimeError("action.idempotency_lookup_failed", err)
	}

	response, err := s.executeRemind(req)
	if err != nil {
		return AutomationActionResponse{}, err
	}

	encoded, err := json.Marshal(response)
	if err != nil {
		return AutomationActionResponse{}, newAutomationRuntimeError("action.encode_failed", err)
	}

	record := model.AutomationIdempotency{
		ActionName:   name,
		RequestID:    req.RequestID,
		ArgsHash:     hash,
		ActorScope:   req.ActorScope,
		ResultState:  response.ResultState,
		ResponseJSON: string(encoded),
	}
	if err := s.db.Create(&record).Error; err != nil {
		return AutomationActionResponse{}, newAutomationRuntimeError("action.idempotency_write_failed", err)
	}

	return response, nil
}

func (s *AutomationActionService) executeRemind(req AutomationActionRequest) (AutomationActionResponse, error) {
	var task model.DailyTask
	err := s.db.Where("completed = ?", false).Order("date asc").Order("id asc").First(&task).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return AutomationActionResponse{}, newAutomationValidationError("action.no_pending_tasks", errors.New("no pending study task to remind"))
	}
	if err != nil {
		return AutomationActionResponse{}, newAutomationRuntimeError("action.pending_task_lookup_failed", err)
	}

	result, err := s.exporter([]model.DailyTask{task}, RemindersExportOptions{})
	if err != nil {
		return AutomationActionResponse{}, newAutomationRuntimeError("action.reminder_export_failed", err)
	}

	return AutomationActionResponse{
		Version:     AutomationVersionV1,
		ResultState: AutomationResultStateExecuted,
		Timestamp:   utcTimestamp(s.now()),
		Action: &AutomationActionPayload{
			ActionName:      "remind",
			RequestID:       req.RequestID,
			ActorScope:      req.ActorScope,
			ReminderTitle:   task.Title,
			ReminderDueDate: task.Date.UTC().Format("2006-01-02"),
			TargetList:      result.ListName,
			CreatedCount:    result.Created,
		},
	}, nil
}

func RunAutomationAction(database *gorm.DB, req AutomationActionRequest) (AutomationActionResponse, error) {
	service := NewAutomationActionService(database)
	return service.RunAutomationAction(req)
}

func normalizedArgsHash(args map[string]string) (string, error) {
	if len(args) == 0 {
		hash := sha256.Sum256([]byte("{}"))
		return hex.EncodeToString(hash[:]), nil
	}

	encoded, err := json.Marshal(args)
	if err != nil {
		return "", err
	}
	hash := sha256.Sum256(encoded)
	return hex.EncodeToString(hash[:]), nil
}
