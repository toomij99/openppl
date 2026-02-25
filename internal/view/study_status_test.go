package view

import (
	"errors"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"ppl-study-planner/internal/model"
	"ppl-study-planner/internal/services"
)

func TestStudyStatusTranslations(t *testing.T) {
	tests := []struct {
		name      string
		err       error
		expected  string
		operation string
	}{
		{
			name:      "google missing credentials",
			err:       &services.GoogleAuthError{Kind: services.GoogleAuthErrorMissingCredentials, Err: errors.New("missing")},
			expected:  "GOOGLE_OAUTH_CREDENTIALS_PATH",
			operation: "Google sync",
		},
		{
			name:      "google token permissions",
			err:       &services.GoogleAuthError{Kind: services.GoogleAuthErrorTokenPermissions, Err: errors.New("bad perms")},
			expected:  "0600",
			operation: "Google sync",
		},
		{
			name: "google auth wrapped by calendar error",
			err: &services.GoogleCalendarError{
				Kind: services.GoogleCalendarErrorAuth,
				Err:  &services.GoogleAuthError{Kind: services.GoogleAuthErrorInvalidCredentials, Err: errors.New("bad json")},
			},
			expected:  "valid OAuth client file",
			operation: "Google sync",
		},
		{
			name:      "reminders permission",
			err:       &services.RemindersExportError{Kind: "permission", Err: errors.New("denied")},
			expected:  "Allow automation permissions",
			operation: "Reminders export",
		},
		{
			name:      "unknown errors are user-safe",
			err:       errors.New("raw internal stack trace"),
			expected:  "Check configuration and try again",
			operation: "ICS export",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			status := newStudyStatusFromError(tc.operation, tc.err)
			if status.severity != studyStatusSeverityError {
				t.Fatalf("expected severity=%q, got %q", studyStatusSeverityError, status.severity)
			}
			if !strings.Contains(status.message, tc.expected) {
				t.Fatalf("expected message %q to contain %q", status.message, tc.expected)
			}
			if strings.Contains(status.message, "stack trace") {
				t.Fatalf("expected user-safe message, got %q", status.message)
			}
		})
	}
}

func TestStudyInputValidation(t *testing.T) {
	t.Run("invalid date keeps input mode active with warning", func(t *testing.T) {
		sv := &StudyView{inputMode: true, dateInput: "13/40/2026"}

		_, _ = sv.handleInput(tea.KeyMsg{Type: tea.KeyEnter})

		if !sv.inputMode {
			t.Fatal("expected input mode to remain active")
		}
		if sv.status.severity != studyStatusSeverityWarning {
			t.Fatalf("expected warning severity, got %q", sv.status.severity)
		}
		if !strings.Contains(sv.status.message, "MM/DD/YYYY") {
			t.Fatalf("expected format guidance, got %q", sv.status.message)
		}
	})

	t.Run("valid date exits input mode and saves", func(t *testing.T) {
		sv := newStudyViewForTest(t)
		sv.inputMode = true
		sv.dateInput = "12/31/2026"

		_, _ = sv.handleInput(tea.KeyMsg{Type: tea.KeyEnter})

		if sv.inputMode {
			t.Fatal("expected input mode to be closed")
		}
		if sv.dateInput != "" {
			t.Fatalf("expected cleared input, got %q", sv.dateInput)
		}
		if sv.status.severity != studyStatusSeveritySuccess {
			t.Fatalf("expected success severity, got %q", sv.status.severity)
		}
		if !sv.hasCheckride {
			t.Fatal("expected checkride date to be set")
		}
	})
}

func TestStudyStatusRendering(t *testing.T) {
	sv := &StudyView{}

	_ = sv.exportICS()
	if sv.status.severity != studyStatusSeverityWarning {
		t.Fatalf("expected ICS no-task warning, got %q", sv.status.severity)
	}

	_ = sv.exportReminders()
	if sv.status.severity != studyStatusSeverityWarning {
		t.Fatalf("expected Reminders no-task warning, got %q", sv.status.severity)
	}

	_ = sv.syncGoogleCalendar()
	if sv.status.severity != studyStatusSeverityWarning {
		t.Fatalf("expected Google no-task warning, got %q", sv.status.severity)
	}

	_ = sv.exportOpenCodeBot()
	if sv.status.severity != studyStatusSeverityWarning {
		t.Fatalf("expected OpenCode no-task warning, got %q", sv.status.severity)
	}

	errStatus := newStudyStatusFromError("ICS export", errors.New("sensitive internal error"))
	renderedErr := renderStudyStatus(errStatus)
	if !strings.Contains(renderedErr, "ICS export failed") {
		t.Fatalf("expected rendered error message, got %q", renderedErr)
	}
	if strings.Contains(renderedErr, "sensitive internal error") {
		t.Fatalf("expected raw error details to be hidden, got %q", renderedErr)
	}
}

func newStudyViewForTest(t *testing.T) *StudyView {
	t.Helper()

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}

	if err := db.AutoMigrate(&model.StudyPlan{}, &model.DailyTask{}); err != nil {
		t.Fatalf("automigrate: %v", err)
	}

	sv := NewStudyView(db)
	return sv
}
