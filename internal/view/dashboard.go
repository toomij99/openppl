package view

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"ppl-study-planner/internal/styles"
)

// DashboardView displays the main dashboard with stats
type DashboardView struct {
	db            interface{}
	checkrideDate time.Time
	stats         DashboardStats
	width         int
	height        int
}

// Init implements tea.Model
func (v *DashboardView) Init() tea.Cmd {
	return nil
}

// DashboardStats holds the dashboard statistics
type DashboardStats struct {
	Completed int
	Remaining int
	Overdue   int
	Total     int
	Progress  float64
	DaysUntil int
	WeekTasks map[string]int
}

// NewDashboardView creates a new dashboard view
func NewDashboardView(db interface{}) *DashboardView {
	return &DashboardView{
		db:     db,
		stats:  DashboardStats{WeekTasks: make(map[string]int)},
		width:  80,
		height: 24,
	}
}

// Update handles dashboard updates
func (v *DashboardView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		v.width = msg.Width
		v.height = msg.Height
	}
	return v, nil
}

// View renders the dashboard
func (v *DashboardView) View() string {
	// Query stats from database if we have a DB connection
	if db, ok := v.db.(interface {
		Find(interface{}) *interface{}
	}); ok {
		v.refreshStats(db)
	}

	contentWidth := v.width - 4
	box := styles.HighlightBox.Width(contentWidth).Height(v.height - 4)

	// Calculate days until checkride
	daysUntil := v.stats.DaysUntil
	if v.stats.DaysUntil < 0 {
		daysUntil = 0
	}

	// Render days until checkride
	daysStyle := styles.Title
	if v.stats.DaysUntil > 0 && v.stats.DaysUntil < 30 {
		daysStyle = daysStyle.Foreground(styles.Error)
	}

	// Render progress bar
	progressBar := v.renderProgressBar(v.stats.Progress)

	// Render week tasks
	weekTasks := v.renderWeekTasks()

	// Render quick stats
	stats := fmt.Sprintf("  Completed: %s | Remaining: %s | Overdue: %s | Total: %s",
		styles.Success.Render(fmt.Sprintf("%d", v.stats.Completed)),
		styles.Normal.Render(fmt.Sprintf("%d", v.stats.Remaining)),
		styles.ErrorStyle.Render(fmt.Sprintf("%d", v.stats.Overdue)),
		styles.Normal.Render(fmt.Sprintf("%d", v.stats.Total)),
	)

	content := fmt.Sprintf(`%s

%s Days until checkride: %s

%s Overall Progress
%s

%s Quick Stats
%s

%s Upcoming Week
%s`,
		styles.Title.Render("Dashboard"),
		daysStyle.Render("📅"),
		daysStyle.Render(fmt.Sprintf("%d days", daysUntil)),
		styles.Normal.Render("Progress"),
		progressBar,
		styles.Normal.Render("Quick Stats"),
		stats,
		styles.Normal.Render("Upcoming Week"),
		weekTasks,
	)

	return box.Render(content)
}

func (v *DashboardView) renderProgressBar(percent float64) string {
	const barWidth = 30
	filled := int(percent / 100 * float64(barWidth))

	bar := ""
	for i := 0; i < barWidth; i++ {
		if i < filled {
			bar += styles.Success.Render("█")
		} else {
			bar += styles.Dim.Render("░")
		}
	}

	return fmt.Sprintf("  [%s] %d%%", bar, int(percent))
}

func (v *DashboardView) renderWeekTasks() string {
	if len(v.stats.WeekTasks) == 0 {
		return styles.Dim.Render("  No tasks scheduled for this week")
	}

	result := ""
	tasks := []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}
	now := time.Now()
	startOfWeek := now.AddDate(0, 0, -int(now.Weekday()))

	for i, day := range tasks {
		date := startOfWeek.AddDate(0, 0, i+1)
		dateKey := date.Format("01/02")
		count := v.stats.WeekTasks[dateKey]

		dayLabel := day
		if date.Format("01/02") == now.Format("01/02") {
			dayLabel = styles.Success.Render("● " + day)
		}

		result += fmt.Sprintf("  %s: %s\n", dayLabel, styles.Normal.Render(fmt.Sprintf("%d tasks", count)))
	}

	return result
}

func (v *DashboardView) refreshStats(db interface{}) {
	// This would query the database for real stats
	// For now, we'll set placeholder values that would come from DB queries
	// In a real implementation, this would use GORM to query

	// Placeholder - actual DB queries would go here:
	// db.Model(&model.DailyTask{}).Where("completed = ?", true).Count(&v.stats.Completed)
	// db.Model(&model.DailyTask{}).Where("completed = ?", false).Count(&v.stats.Remaining)
	// db.Model(&model.DailyTask{}).Where("date < ?", time.Now()).Where("completed = ?", false).Count(&v.stats.Overdue)
	// db.Model(&model.DailyTask{}).Count(&v.stats.Total)
}

// SetCheckrideDate sets the checkride date for the dashboard
func (v *DashboardView) SetCheckrideDate(date time.Time) {
	v.checkrideDate = date
	if !date.IsZero() {
		v.stats.DaysUntil = int(time.Until(date).Hours() / 24)
	}
}

// SetStats sets the dashboard stats directly
func (v *DashboardView) SetStats(stats DashboardStats) {
	v.stats = stats
}

// GetStats returns the current dashboard stats
func (v *DashboardView) GetStats() DashboardStats {
	return v.stats
}
