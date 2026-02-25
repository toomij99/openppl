package view

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"ppl-study-planner/internal/model"
	"ppl-study-planner/internal/services"
	"ppl-study-planner/internal/styles"

	"gorm.io/gorm"
)

// ProgressView shows task completion progress
type ProgressView struct {
	db             *gorm.DB
	tasks          []model.DailyTask
	overallPercent float64
	byCategory     map[string]services.ProgressStats
}

// NewProgressView creates a new progress view
func NewProgressView(db *gorm.DB) *ProgressView {
	pv := &ProgressView{db: db}
	pv.loadData()
	return pv
}

// loadData loads progress data from database
func (pv *ProgressView) loadData() {
	var plan model.StudyPlan
	if err := pv.db.Preload("DailyTasks").Last(&plan).Error; err == nil {
		pv.tasks = plan.DailyTasks
		_, _, pv.overallPercent = services.CalculateProgress(plan.DailyTasks)
		pv.byCategory = services.GetProgressByCategory(plan.DailyTasks)
	}
}

// Init implements tea.Model
func (pv *ProgressView) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (pv *ProgressView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case tea.KeyMsg:
		// Reload data on keypress
		pv.loadData()
	}
	return pv, nil
}

// View implements tea.Model
func (pv *ProgressView) View() string {
	var b strings.Builder

	b.WriteString(styles.Title.Render("Progress"))
	b.WriteString("\n\n")

	// Overall progress
	b.WriteString(styles.Normal.Render("Overall Progress:"))
	b.WriteString("\n")
	b.WriteString(pv.renderProgressBar(pv.overallPercent))
	b.WriteString(fmt.Sprintf(" %.0f%%\n\n", pv.overallPercent))

	// Progress by category
	b.WriteString(styles.Normal.Render("By Category:"))
	b.WriteString("\n")

	cats := []string{"Theory", "Chair Flying", "Garmin 430", "CFI Flights"}
	for i, cat := range cats {
		stats := pv.byCategory[cat]
		percent := stats.CalculatePercentage()
		bar := pv.renderProgressBar(percent)
		b.WriteString(fmt.Sprintf(" %d. %s ", i+1, cat))
		b.WriteString(bar)
		b.WriteString(fmt.Sprintf(" %d/%d (%.0f%%)\n", stats.Completed, stats.Total, percent))
	}

	b.WriteString("\n")

	// Today's tasks
	b.WriteString(styles.Normal.Render("Today's Tasks:"))
	b.WriteString("\n")
	todayTasks := pv.getTodayTasks()
	if len(todayTasks) == 0 {
		b.WriteString(styles.Dim.Render("  No tasks scheduled for today"))
	} else {
		for _, t := range todayTasks {
			check := "[ ]"
			if t.Completed {
				check = "[x]"
			}
			b.WriteString(fmt.Sprintf("  %s %s\n", check, t.Title))
		}
	}

	b.WriteString("\n")

	// Recent completions
	b.WriteString(styles.Normal.Render("Recent Completions:"))
	b.WriteString("\n")
	recent := pv.getRecentCompletions()
	if len(recent) == 0 {
		b.WriteString(styles.Dim.Render("  No recently completed tasks"))
	} else {
		for _, t := range recent {
			b.WriteString(fmt.Sprintf("  [x] %s\n", t.Title))
		}
	}

	b.WriteString("\n")
	b.WriteString(styles.Dim.Render("[Any key] Refresh"))

	return b.String()
}

// renderProgressBar renders a progress bar
func (pv *ProgressView) renderProgressBar(percent float64) string {
	width := 20
	filled := int(percent / 100 * float64(width))
	empty := width - filled

	bar := strings.Repeat("█", filled) + strings.Repeat("░", empty)
	return styles.ProgressBar.Render(bar[:filled]) + styles.ProgressBarEmpty.Render(bar[filled:])
}

// getTodayTasks returns tasks scheduled for today
func (pv *ProgressView) getTodayTasks() []model.DailyTask {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	var result []model.DailyTask
	for _, t := range pv.tasks {
		taskDay := time.Date(t.Date.Year(), t.Date.Month(), t.Date.Day(), 0, 0, 0, 0, time.UTC)
		if taskDay.Equal(today) {
			result = append(result, t)
		}
	}
	return result
}

// getRecentCompletions returns recently completed tasks
func (pv *ProgressView) getRecentCompletions() []model.DailyTask {
	var result []model.DailyTask
	for _, t := range pv.tasks {
		if t.Completed {
			// Check if completed recently (for now just show completed tasks)
			result = append(result, t)
			if len(result) >= 7 {
				break
			}
		}
	}
	return result
}
