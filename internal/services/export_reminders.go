package services

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"ppl-study-planner/internal/model"
)

const defaultRemindersList = "OpenPPL Study Tasks"

const remindersAppleScript = `on run argv
  set listName to item 1 of argv
  set secondsFromMidnight to (item 2 of argv) as integer
  set y to (item 3 of argv) as integer
  set m to (item 4 of argv) as integer
  set d to (item 5 of argv) as integer
  set theTitle to item 6 of argv
  set theNotes to item 7 of argv
  tell application "Reminders"
    if not (exists list listName) then
      make new list with properties {name:listName}
    end if
    set theList to list listName
    set theDate to current date
    set year of theDate to y
    set month of theDate to (item m of {January, February, March, April, May, June, July, August, September, October, November, December})
    set day of theDate to d
    set time of theDate to secondsFromMidnight
    make new reminder at end of reminders of theList with properties {name:theTitle, body:theNotes, due date:theDate}
  end tell
end run`

// RemindersExportOptions configures reminders export behavior.
type RemindersExportOptions struct {
	ListName string
	Timeout  time.Duration
}

// RemindersExportResult contains reminder export metadata.
type RemindersExportResult struct {
	ListName string
	Created  int
}

// RemindersExportError classifies reminder export failures.
type RemindersExportError struct {
	Kind   string
	Output string
	Err    error
}

func (e *RemindersExportError) Error() string {
	base := "apple reminders export failed"
	if e.Kind != "" {
		base = base + ": " + e.Kind
	}
	if e.Err != nil {
		base = base + ": " + e.Err.Error()
	}
	if strings.TrimSpace(e.Output) != "" {
		base = base + " (" + strings.TrimSpace(e.Output) + ")"
	}
	return base
}

func (e *RemindersExportError) Unwrap() error {
	return e.Err
}

type commandRunner func(ctx context.Context, name string, args ...string) ([]byte, error)

func defaultCommandRunner(ctx context.Context, name string, args ...string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, name, args...)
	return cmd.CombinedOutput()
}

// ExportAppleReminders creates reminders through osascript using argv-safe argument passing.
func ExportAppleReminders(tasks []model.DailyTask, opts RemindersExportOptions) (RemindersExportResult, error) {
	return exportAppleRemindersWithRunner(tasks, opts, defaultCommandRunner)
}

func exportAppleRemindersWithRunner(tasks []model.DailyTask, opts RemindersExportOptions, run commandRunner) (RemindersExportResult, error) {
	if len(tasks) == 0 {
		return RemindersExportResult{}, &RemindersExportError{Kind: "validation", Err: errors.New("no tasks available for reminders export")}
	}
	if run == nil {
		return RemindersExportResult{}, &RemindersExportError{Kind: "validation", Err: errors.New("nil command runner")}
	}

	artifactDir, err := ResolveArtifactOutputDir("")
	if err != nil {
		return RemindersExportResult{}, &RemindersExportError{Kind: "validation", Err: fmt.Errorf("resolve reminders artifact directory: %w", err)}
	}
	if err := os.MkdirAll(artifactDir, 0o755); err != nil {
		return RemindersExportResult{}, &RemindersExportError{Kind: "script_failure", Err: fmt.Errorf("create reminders artifact directory: %w", err)}
	}

	listName := strings.TrimSpace(opts.ListName)
	if listName == "" {
		listName = defaultRemindersList
	}

	timeout := opts.Timeout
	if timeout <= 0 {
		timeout = 12 * time.Second
	}

	created := 0
	for _, task := range tasks {
		taskDate := task.Date.UTC()
		title := strings.TrimSpace(task.Title)
		if title == "" {
			title = "Study Task"
		}
		notes := strings.TrimSpace(task.Description)
		if notes == "" {
			notes = strings.TrimSpace(task.Category)
		}

		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		output, err := run(
			ctx,
			"osascript",
			"-e", remindersAppleScript,
			"--",
			listName,
			"32400", // 09:00 local due time
			fmt.Sprintf("%d", taskDate.Year()),
			fmt.Sprintf("%d", int(taskDate.Month())),
			fmt.Sprintf("%d", taskDate.Day()),
			title,
			notes,
		)
		cancel()

		if err != nil {
			outputStr := string(output)
			if errors.Is(ctx.Err(), context.DeadlineExceeded) || errors.Is(err, context.DeadlineExceeded) {
				return RemindersExportResult{}, &RemindersExportError{Kind: "timeout", Output: outputStr, Err: err}
			}
			if isPermissionError(outputStr) {
				return RemindersExportResult{}, &RemindersExportError{Kind: "permission", Output: outputStr, Err: err}
			}
			return RemindersExportResult{}, &RemindersExportError{Kind: "script_failure", Output: outputStr, Err: err}
		}

		created++
	}

	return RemindersExportResult{ListName: listName, Created: created}, nil
}

func isPermissionError(output string) bool {
	output = strings.ToLower(output)
	return strings.Contains(output, "not authorized") ||
		strings.Contains(output, "not permitted") ||
		strings.Contains(output, "-1743")
}
