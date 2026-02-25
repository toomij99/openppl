package services

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	ics "github.com/arran4/golang-ical"

	"ppl-study-planner/internal/model"
)

// ICSExportResult contains metadata for a generated ICS file.
type ICSExportResult struct {
	Path       string
	EventCount int
}

// ExportICS writes study tasks into an RFC5545-compatible .ics file.
func ExportICS(tasks []model.DailyTask, outputDir string) (ICSExportResult, error) {
	if len(tasks) == 0 {
		return ICSExportResult{}, errors.New("no tasks available for ICS export")
	}

	if strings.TrimSpace(outputDir) == "" {
		outputDir = "exports"
	}

	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return ICSExportResult{}, fmt.Errorf("create export directory: %w", err)
	}

	cal := ics.NewCalendar()
	cal.SetMethod(ics.MethodPublish)
	cal.SetProductId("-//openppl//study-plan//EN")
	cal.SetVersion("2.0")

	nowUTC := time.Now().UTC()
	for _, task := range tasks {
		uid := deterministicTaskUID(task)
		event := cal.AddEvent(uid)
		event.SetDtStampTime(nowUTC)

		startUTC, endUTC := taskWindowUTC(task.Date)
		event.SetStartAt(startUTC)
		event.SetEndAt(endUTC)

		title := strings.TrimSpace(task.Title)
		if title == "" {
			title = "Study Task"
		}
		event.SetSummary(title)

		desc := strings.TrimSpace(task.Description)
		if desc == "" {
			desc = strings.TrimSpace(task.Category)
		}
		event.SetDescription(desc)
	}

	fileName := fmt.Sprintf("study-plan-%s.ics", nowUTC.Format("20060102-150405"))
	outputPath := filepath.Join(outputDir, fileName)

	if err := os.WriteFile(outputPath, []byte(cal.Serialize()), 0o644); err != nil {
		return ICSExportResult{}, fmt.Errorf("write ICS file: %w", err)
	}

	return ICSExportResult{
		Path:       outputPath,
		EventCount: len(tasks),
	}, nil
}

func deterministicTaskUID(task model.DailyTask) string {
	datePart := task.Date.UTC().Format("20060102")
	if task.ID != 0 {
		return fmt.Sprintf("task-%d-%s@openppl", task.ID, datePart)
	}

	titlePart := strings.ToLower(strings.TrimSpace(task.Title))
	titlePart = strings.ReplaceAll(titlePart, " ", "-")
	if titlePart == "" {
		titlePart = "untitled"
	}

	return fmt.Sprintf("task-%s-%s@openppl", titlePart, datePart)
}

func taskWindowUTC(taskDate time.Time) (time.Time, time.Time) {
	start := time.Date(taskDate.Year(), taskDate.Month(), taskDate.Day(), 9, 0, 0, 0, time.UTC)
	end := start.Add(30 * time.Minute)
	return start, end
}
