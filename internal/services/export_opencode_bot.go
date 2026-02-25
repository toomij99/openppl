package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"ppl-study-planner/internal/model"
)

const defaultOpenCodeBotSchemaVersion = "v1"

var opencodeNow = time.Now

type OpenCodeBotExportOptions struct {
	OutputDir string
	Version   string
}

type OpenCodeBotExportResult struct {
	Path      string
	Version   string
	TaskCount int
}

type openCodeBotPayload struct {
	Version     string            `json:"version"`
	GeneratedAt string            `json:"generated_at"`
	Tasks       []openCodeBotTask `json:"tasks"`
}

type openCodeBotTask struct {
	TaskID      uint            `json:"task_id"`
	Title       string          `json:"title"`
	Category    string          `json:"category"`
	Description string          `json:"description"`
	DueAt       string          `json:"due_at"`
	Completed   bool            `json:"completed"`
	Metadata    openCodeBotMeta `json:"metadata"`
}

type openCodeBotMeta struct {
	Source   string `json:"source"`
	Identity string `json:"identity"`
}

func ExportOpenCodeBotTasks(tasks []model.DailyTask, opts OpenCodeBotExportOptions) (OpenCodeBotExportResult, error) {
	if len(tasks) == 0 {
		return OpenCodeBotExportResult{}, errors.New("no tasks available for OpenCode bot export")
	}

	resolvedDir, err := ResolveArtifactOutputDir(opts.OutputDir)
	if err != nil {
		return OpenCodeBotExportResult{}, fmt.Errorf("resolve bot export directory: %w", err)
	}

	if err := os.MkdirAll(resolvedDir, 0o755); err != nil {
		return OpenCodeBotExportResult{}, fmt.Errorf("create bot export directory: %w", err)
	}

	version := strings.TrimSpace(opts.Version)
	if version == "" {
		version = defaultOpenCodeBotSchemaVersion
	}

	sortedTasks := append([]model.DailyTask(nil), tasks...)
	sort.Slice(sortedTasks, func(i, j int) bool {
		if sortedTasks[i].Date.Equal(sortedTasks[j].Date) {
			if sortedTasks[i].ID == sortedTasks[j].ID {
				return sortedTasks[i].Title < sortedTasks[j].Title
			}
			return sortedTasks[i].ID < sortedTasks[j].ID
		}
		return sortedTasks[i].Date.Before(sortedTasks[j].Date)
	})

	payload := openCodeBotPayload{
		Version:     version,
		GeneratedAt: opencodeNow().UTC().Format(time.RFC3339),
		Tasks:       make([]openCodeBotTask, 0, len(sortedTasks)),
	}

	for _, task := range sortedTasks {
		title := strings.TrimSpace(task.Title)
		if title == "" {
			title = "Study Task"
		}

		description := strings.TrimSpace(task.Description)
		if description == "" {
			description = strings.TrimSpace(task.Category)
		}

		payload.Tasks = append(payload.Tasks, openCodeBotTask{
			TaskID:      task.ID,
			Title:       title,
			Category:    strings.TrimSpace(task.Category),
			Description: description,
			DueAt:       task.Date.UTC().Format(time.RFC3339),
			Completed:   task.Completed,
			Metadata: openCodeBotMeta{
				Source:   "openppl",
				Identity: deterministicGoogleTaskIdentity(task),
			},
		})
	}

	b, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return OpenCodeBotExportResult{}, fmt.Errorf("marshal bot export payload: %w", err)
	}

	filename := fmt.Sprintf("opencode-bot-%s-%s.json", version, opencodeNow().UTC().Format("20060102-150405"))
	outputPath := filepath.Join(resolvedDir, filename)
	if err := os.WriteFile(outputPath, b, 0o644); err != nil {
		return OpenCodeBotExportResult{}, fmt.Errorf("write bot export payload: %w", err)
	}

	return OpenCodeBotExportResult{
		Path:      outputPath,
		Version:   version,
		TaskCount: len(payload.Tasks),
	}, nil
}
