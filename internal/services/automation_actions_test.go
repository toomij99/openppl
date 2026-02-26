package services

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"ppl-study-planner/internal/model"
)

func TestRunAutomationAction(t *testing.T) {
	db := setupAutomationActionsTestDB(t)
	now := time.Date(2026, 2, 26, 12, 0, 0, 0, time.UTC)

	plan := model.StudyPlan{CheckrideDate: time.Date(2026, 8, 10, 0, 0, 0, 0, time.UTC)}
	if err := db.Create(&plan).Error; err != nil {
		t.Fatalf("create plan: %v", err)
	}
	if err := db.Create(&model.DailyTask{StudyPlanID: plan.ID, Date: time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC), Title: "Area I review", Category: "Theory", Completed: false}).Error; err != nil {
		t.Fatalf("create task: %v", err)
	}

	callCount := 0
	service := NewAutomationActionService(db).
		WithClock(func() time.Time { return now }).
		WithExporter(func(tasks []model.DailyTask, opts RemindersExportOptions) (RemindersExportResult, error) {
			callCount++
			if len(tasks) != 1 {
				t.Fatalf("expected exactly one reminder task, got %d", len(tasks))
			}
			return RemindersExportResult{ListName: "OpenPPL Study Tasks", Created: 1}, nil
		})

	first, err := service.RunAutomationAction(AutomationActionRequest{Name: "remind", RequestID: "req-1", ActorScope: "telegram:user:42"})
	if err != nil {
		t.Fatalf("first run failed: %v", err)
	}
	if first.ResultState != AutomationResultStateExecuted {
		t.Fatalf("expected executed, got %q", first.ResultState)
	}
	if callCount != 1 {
		t.Fatalf("expected exporter to be called once, got %d", callCount)
	}

	replayed, err := service.RunAutomationAction(AutomationActionRequest{Name: "remind", RequestID: "req-1", ActorScope: "telegram:user:42"})
	if err != nil {
		t.Fatalf("replay run failed: %v", err)
	}
	if replayed.ResultState != AutomationResultStateReplayed {
		t.Fatalf("expected replayed, got %q", replayed.ResultState)
	}
	if callCount != 1 {
		t.Fatalf("expected no extra exporter call on replay, got %d", callCount)
	}

	if _, err := service.RunAutomationAction(AutomationActionRequest{Name: "unknown", RequestID: "req-2"}); err == nil {
		t.Fatal("expected allowlist rejection for unknown action")
	}
	if _, err := service.RunAutomationAction(AutomationActionRequest{Name: "remind"}); err == nil {
		t.Fatal("expected request_id validation failure")
	}

}

func TestRunAutomationAction_NoPendingTask(t *testing.T) {
	db := setupAutomationActionsTestDB(t)
	service := NewAutomationActionService(db).WithClock(time.Now).WithExporter(func(tasks []model.DailyTask, opts RemindersExportOptions) (RemindersExportResult, error) {
		return RemindersExportResult{}, errors.New("should not be called")
	})

	if _, err := service.RunAutomationAction(AutomationActionRequest{Name: "remind", RequestID: "r", ActorScope: "a"}); err == nil {
		t.Fatal("expected no pending tasks error")
	}
}

func setupAutomationActionsTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	dsn := fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&model.StudyPlan{}, &model.DailyTask{}, &model.AutomationIdempotency{}); err != nil {
		t.Fatalf("automigrate: %v", err)
	}
	return db
}
