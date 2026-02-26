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

func TestAutomationStatusEnvelope(t *testing.T) {
	db := setupAutomationStatusTestDB(t)

	checkride := time.Date(2026, 8, 10, 0, 0, 0, 0, time.UTC)
	if err := db.Create(&model.StudyPlan{CheckrideDate: checkride}).Error; err != nil {
		t.Fatalf("create study plan: %v", err)
	}

	now := time.Date(2026, 2, 26, 12, 0, 0, 0, time.UTC)
	status, err := BuildAutomationStatus(db, now)
	if err != nil {
		t.Fatalf("BuildAutomationStatus failed: %v", err)
	}

	if status.Version != AutomationVersionV1 {
		t.Fatalf("expected version %q, got %q", AutomationVersionV1, status.Version)
	}
	if status.ResultState != AutomationResultStateOK {
		t.Fatalf("expected result_state ok, got %q", status.ResultState)
	}
	if status.Timestamp != now.Format(time.RFC3339) {
		t.Fatalf("unexpected timestamp %q", status.Timestamp)
	}
	if status.Status == nil {
		t.Fatal("expected non-nil status payload")
	}
	if status.Status.CheckrideDate != "2026-08-10" {
		t.Fatalf("unexpected checkride date: %q", status.Status.CheckrideDate)
	}
}

func TestAutomationTypesStatusEnvelope(t *testing.T) {
	if _, err := BuildAutomationStatus(nil, time.Now()); err == nil {
		t.Fatal("expected validation error for nil database")
	} else {
		var commandErr *AutomationCommandError
		if !errors.As(err, &commandErr) {
			t.Fatalf("expected AutomationCommandError, got %T", err)
		}
		if commandErr.Kind != "validation" {
			t.Fatalf("expected validation kind, got %q", commandErr.Kind)
		}
	}
}

func TestBuildAutomationStatus(t *testing.T) {
	db := setupAutomationStatusTestDB(t)

	plan := model.StudyPlan{CheckrideDate: time.Date(2026, 8, 10, 0, 0, 0, 0, time.UTC)}
	if err := db.Create(&plan).Error; err != nil {
		t.Fatalf("create plan: %v", err)
	}

	tasks := []model.DailyTask{
		{StudyPlanID: plan.ID, Date: time.Date(2026, 3, 3, 0, 0, 0, 0, time.UTC), Category: "Theory", Title: "Bravo", Completed: false},
		{StudyPlanID: plan.ID, Date: time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC), Category: "Theory", Title: "Alpha", Completed: false},
		{StudyPlanID: plan.ID, Date: time.Date(2026, 3, 2, 0, 0, 0, 0, time.UTC), Category: "Ground", Title: "Zulu", Completed: true},
	}
	for _, task := range tasks {
		if err := db.Create(&task).Error; err != nil {
			t.Fatalf("create task: %v", err)
		}
	}

	now := time.Date(2026, 2, 26, 12, 0, 0, 0, time.UTC)
	statusA, err := BuildAutomationStatus(db, now)
	if err != nil {
		t.Fatalf("BuildAutomationStatus failed: %v", err)
	}
	statusB, err := BuildAutomationStatus(db, now)
	if err != nil {
		t.Fatalf("BuildAutomationStatus second call failed: %v", err)
	}

	if statusA.Status.Summary.TotalTasks != 3 || statusA.Status.Summary.CompletedTasks != 1 || statusA.Status.Summary.PendingTasks != 2 {
		t.Fatalf("unexpected summary: %+v", statusA.Status.Summary)
	}
	if len(statusA.Status.NextTasks) != 2 {
		t.Fatalf("expected 2 next tasks, got %d", len(statusA.Status.NextTasks))
	}
	if statusA.Status.NextTasks[0].Title != "Alpha" || statusA.Status.NextTasks[1].Title != "Bravo" {
		t.Fatalf("expected deterministic ordering Alpha/Bravo, got %+v", statusA.Status.NextTasks)
	}
	if statusA.Status.NextTasks[0] != statusB.Status.NextTasks[0] || statusA.Status.NextTasks[1] != statusB.Status.NextTasks[1] {
		t.Fatal("expected deterministic ordering across repeated calls")
	}
}

func setupAutomationStatusTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	dsn := fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&model.StudyPlan{}, &model.DailyTask{}); err != nil {
		t.Fatalf("automigrate: %v", err)
	}
	return db
}
