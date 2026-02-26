package services

import (
	"errors"
	"fmt"
	"sort"
	"time"

	"gorm.io/gorm"

	"ppl-study-planner/internal/model"
)

func BuildAutomationStatus(database *gorm.DB, now time.Time) (AutomationStatusResponse, error) {
	if database == nil {
		return AutomationStatusResponse{}, newAutomationValidationError("status.db_required", errors.New("database is required"))
	}

	var tasks []model.DailyTask
	if err := database.Order("date asc").Order("id asc").Find(&tasks).Error; err != nil {
		return AutomationStatusResponse{}, newAutomationRuntimeError("status.tasks_query_failed", fmt.Errorf("query tasks: %w", err))
	}

	var plan model.StudyPlan
	planErr := database.Order("id asc").First(&plan).Error
	if planErr != nil && !errors.Is(planErr, gorm.ErrRecordNotFound) {
		return AutomationStatusResponse{}, newAutomationRuntimeError("status.plan_query_failed", fmt.Errorf("query study plan: %w", planErr))
	}

	summary := AutomationStatusSummary{TotalTasks: len(tasks)}
	nextTasks := make([]AutomationStatusTask, 0, 5)
	for _, task := range tasks {
		if task.Completed {
			summary.CompletedTasks++
			continue
		}
		summary.PendingTasks++
		if len(nextTasks) < 5 {
			nextTasks = append(nextTasks, AutomationStatusTask{
				Date:     task.Date.UTC().Format("2006-01-02"),
				Category: task.Category,
				Title:    task.Title,
			})
		}
	}

	sort.SliceStable(nextTasks, func(i, j int) bool {
		if nextTasks[i].Date == nextTasks[j].Date {
			if nextTasks[i].Category == nextTasks[j].Category {
				return nextTasks[i].Title < nextTasks[j].Title
			}
			return nextTasks[i].Category < nextTasks[j].Category
		}
		return nextTasks[i].Date < nextTasks[j].Date
	})

	payload := &AutomationStatusPayload{Summary: summary, NextTasks: nextTasks}
	if planErr == nil {
		payload.CheckrideDate = plan.CheckrideDate.UTC().Format("2006-01-02")
	}

	return AutomationStatusResponse{
		Version:     AutomationVersionV1,
		ResultState: AutomationResultStateOK,
		Timestamp:   utcTimestamp(now),
		Status:      payload,
	}, nil
}
