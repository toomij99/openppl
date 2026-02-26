package automation

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"ppl-study-planner/internal/model"
	"ppl-study-planner/internal/services"
)

func TestAutomationStatusCLI(t *testing.T) {
	db := setupAutomationCLITestDB(t)
	plan := model.StudyPlan{CheckrideDate: time.Date(2026, 8, 10, 0, 0, 0, 0, time.UTC)}
	if err := db.Create(&plan).Error; err != nil {
		t.Fatalf("create plan: %v", err)
	}
	if err := db.Create(&model.DailyTask{StudyPlanID: plan.ID, Date: time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC), Title: "Task A", Category: "Theory", Completed: false}).Error; err != nil {
		t.Fatalf("create task: %v", err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Execute(db, []string{"status"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d, stderr: %s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), `"result_state":"ok"`) {
		t.Fatalf("expected status ok payload, got %s", stdout.String())
	}

	stdout.Reset()
	stderr.Reset()
	code = Execute(db, []string{"bogus"}, &stdout, &stderr)
	if code == 0 {
		t.Fatal("expected non-zero exit for unknown subcommand")
	}
	if !strings.Contains(stderr.String(), "automation.unknown_subcommand") {
		t.Fatalf("expected unknown subcommand error, got %s", stderr.String())
	}

	stdout.Reset()
	stderr.Reset()
	code = Execute(db, []string{"status", "--invalid"}, &stdout, &stderr)
	if code == 0 {
		t.Fatal("expected non-zero exit for invalid status arguments")
	}
	if !strings.Contains(stderr.String(), "status.invalid_arguments") {
		t.Fatalf("expected invalid arguments error, got %s", stderr.String())
	}
}

func TestAutomationAction(t *testing.T) {
	restore := services.SetAutomationReminderExporterForTest(func(tasks []model.DailyTask, opts services.RemindersExportOptions) (services.RemindersExportResult, error) {
		return services.RemindersExportResult{ListName: "OpenPPL Study Tasks", Created: len(tasks)}, nil
	})
	defer restore()

	db := setupAutomationCLITestDB(t)
	plan := model.StudyPlan{CheckrideDate: time.Date(2026, 8, 10, 0, 0, 0, 0, time.UTC)}
	if err := db.Create(&plan).Error; err != nil {
		t.Fatalf("create plan: %v", err)
	}
	if err := db.Create(&model.DailyTask{StudyPlanID: plan.ID, Date: time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC), Title: "Task A", Category: "Theory", Completed: false}).Error; err != nil {
		t.Fatalf("create task: %v", err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Execute(db, []string{"action", "--name", "remind", "--request-id", "req-1", "--actor-scope", "telegram:user:1"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("expected success for action, got %d, stderr=%s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), `"result_state":"executed"`) {
		t.Fatalf("expected executed result, got %s", stdout.String())
	}

	stdout.Reset()
	stderr.Reset()
	code = Execute(db, []string{"action", "--name", "remind", "--request-id", "req-1", "--actor-scope", "telegram:user:1"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("expected replay success, got %d, stderr=%s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), `"result_state":"replayed"`) {
		t.Fatalf("expected replayed result, got %s", stdout.String())
	}

	stdout.Reset()
	stderr.Reset()
	code = Execute(db, []string{"action", "--name", "remind"}, &stdout, &stderr)
	if code == 0 {
		t.Fatal("expected failure for missing request-id")
	}
	if !strings.Contains(stderr.String(), "action.request_id_required") {
		t.Fatalf("expected missing request_id error, got %s", stderr.String())
	}

	stdout.Reset()
	stderr.Reset()
	code = Execute(db, []string{"action", "--name", "unknown", "--request-id", "req-2"}, &stdout, &stderr)
	if code == 0 {
		t.Fatal("expected rejection for unsupported action")
	}
	if !strings.Contains(stderr.String(), "action.not_allowlisted") {
		t.Fatalf("expected allowlist rejection, got %s", stderr.String())
	}
}

func setupAutomationCLITestDB(t *testing.T) *gorm.DB {
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
