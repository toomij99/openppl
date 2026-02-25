package view

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"ppl-study-planner/internal/model"
	"ppl-study-planner/internal/services"
	"ppl-study-planner/internal/styles"

	"gorm.io/gorm"
)

// StudyView represents the study plan screen
type StudyView struct {
	db            *gorm.DB
	checkrideDate time.Time
	hasCheckride  bool
	tasks         []model.DailyTask
	filteredTasks []model.DailyTask
	selectedIdx   int
	inputMode     bool
	dateInput     string
	category      string
	statusMessage string
}

type icsExportDoneMsg struct {
	result services.ICSExportResult
	err    error
}

type remindersExportDoneMsg struct {
	result services.RemindersExportResult
	err    error
}

// NewStudyView creates a new study view
func NewStudyView(db *gorm.DB) *StudyView {
	sv := &StudyView{
		db:       db,
		category: "",
	}
	sv.loadData()
	return sv
}

// loadData loads study plan data from database
func (sv *StudyView) loadData() {
	var plan model.StudyPlan
	if err := sv.db.Preload("DailyTasks").Last(&plan).Error; err == nil {
		sv.checkrideDate = plan.CheckrideDate
		sv.hasCheckride = true
		sv.tasks = plan.DailyTasks
	}
	sv.applyFilter()
}

// applyFilter filters tasks by selected category
func (sv *StudyView) applyFilter() {
	if sv.category == "" {
		sv.filteredTasks = sv.tasks
	} else {
		sv.filteredTasks = nil
		for _, t := range sv.tasks {
			if t.Category == sv.category {
				sv.filteredTasks = append(sv.filteredTasks, t)
			}
		}
	}
	if sv.selectedIdx >= len(sv.filteredTasks) {
		sv.selectedIdx = 0
	}
}

// Init implements tea.Model
func (sv *StudyView) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (sv *StudyView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case icsExportDoneMsg:
		if msg.err != nil {
			sv.statusMessage = fmt.Sprintf("ICS export failed: %v", msg.err)
		} else {
			sv.statusMessage = fmt.Sprintf("ICS export complete: %s (%d events)", msg.result.Path, msg.result.EventCount)
		}
		return sv, nil
	case remindersExportDoneMsg:
		if msg.err != nil {
			sv.statusMessage = fmt.Sprintf("Reminders export failed: %v", msg.err)
		} else {
			sv.statusMessage = fmt.Sprintf("Reminders export complete: %d created in '%s'", msg.result.Created, msg.result.ListName)
		}
		return sv, nil
	case tea.KeyMsg:
		if sv.inputMode {
			return sv.handleInput(msg)
		}
		return sv.handleNav(msg)
	}
	return sv, nil
}

// handleInput handles date input mode
func (sv *StudyView) handleInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		if len(sv.dateInput) >= 6 {
			// Try parsing
			layout := "01/02/2006"
			if len(sv.dateInput) == 6 {
				layout = "01/02/06"
			}
			if d, err := time.Parse(layout, sv.dateInput); err == nil {
				sv.saveStudyPlan(d)
			}
		}
		sv.inputMode = false
		sv.dateInput = ""
	case "esc":
		sv.inputMode = false
		sv.dateInput = ""
	case "backspace":
		if len(sv.dateInput) > 0 {
			sv.dateInput = sv.dateInput[:len(sv.dateInput)-1]
		}
	default:
		if len(msg.String()) == 1 && len(sv.dateInput) < 10 {
			sv.dateInput += msg.String()
		}
	}
	return sv, nil
}

// handleNav handles navigation keys
func (sv *StudyView) handleNav(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up":
		if sv.selectedIdx > 0 {
			sv.selectedIdx--
		}
	case "down":
		if sv.selectedIdx < len(sv.filteredTasks)-1 {
			sv.selectedIdx++
		}
	case "enter":
		if len(sv.filteredTasks) > 0 {
			return sv, sv.toggleTask(sv.filteredTasks[sv.selectedIdx].ID)
		}
	case "/":
		sv.inputMode = true
		sv.dateInput = ""
	case "tab":
		sv.cycleCategory()
	case "1":
		sv.category = ""
		sv.applyFilter()
	case "2":
		sv.category = "Theory"
		sv.applyFilter()
	case "3":
		sv.category = "Chair Flying"
		sv.applyFilter()
	case "4":
		sv.category = "Garmin 430"
		sv.applyFilter()
	case "5":
		sv.category = "CFI Flights"
		sv.applyFilter()
	case "e":
		return sv, sv.exportICS()
	case "r":
		return sv, sv.exportReminders()
	}
	return sv, nil
}

func (sv *StudyView) exportICS() tea.Cmd {
	if len(sv.tasks) == 0 {
		sv.statusMessage = "ICS export skipped: no tasks to export"
		return nil
	}

	tasks := append([]model.DailyTask(nil), sv.tasks...)
	return func() tea.Msg {
		result, err := services.ExportICS(tasks, "exports")
		return icsExportDoneMsg{result: result, err: err}
	}
}

func (sv *StudyView) exportReminders() tea.Cmd {
	if len(sv.tasks) == 0 {
		sv.statusMessage = "Reminders export skipped: no tasks to export"
		return nil
	}

	tasks := append([]model.DailyTask(nil), sv.tasks...)
	return func() tea.Msg {
		result, err := services.ExportAppleReminders(tasks, services.RemindersExportOptions{})
		return remindersExportDoneMsg{result: result, err: err}
	}
}

// cycleCategory cycles through category filters
func (sv *StudyView) cycleCategory() {
	cats := []string{"", "Theory", "Chair Flying", "Garmin 430", "CFI Flights"}
	idx := 0
	for i, c := range cats {
		if sv.category == c {
			idx = i
			break
		}
	}
	sv.category = cats[(idx+1)%len(cats)]
	sv.applyFilter()
}

// toggleTask toggles task completion in database
func (sv *StudyView) toggleTask(id uint) tea.Cmd {
	return func() tea.Msg {
		var task model.DailyTask
		if err := sv.db.First(&task, id).Error; err == nil {
			task.Completed = !task.Completed
			sv.db.Save(&task)
			sv.loadData()
		}
		return nil
	}
}

// saveStudyPlan saves the checkride date and generates tasks
func (sv *StudyView) saveStudyPlan(date time.Time) {
	var plan model.StudyPlan
	if err := sv.db.Last(&plan).Error; err != nil {
		plan = model.StudyPlan{CheckrideDate: date}
		sv.db.Create(&plan)
	} else {
		plan.CheckrideDate = date
		sv.db.Save(&plan)
	}

	// Clear old tasks
	sv.db.Where("study_plan_id = ?", plan.ID).Delete(&model.DailyTask{})

	// Generate new tasks
	tasks := services.GenerateStudyPlan(date, 90)
	for i := range tasks {
		tasks[i].StudyPlanID = plan.ID
	}
	sv.db.Create(&tasks)

	sv.loadData()
}

// View implements tea.Model
func (sv *StudyView) View() string {
	var b strings.Builder

	b.WriteString(styles.Title.Render("Study Plan"))
	b.WriteString("\n\n")

	// Date display
	if sv.hasCheckride {
		b.WriteString(fmt.Sprintf("Checkride: %s ", sv.checkrideDate.Format("01/02/2006")))
		b.WriteString(styles.Dim.Render("(press / to change)"))
	} else if sv.inputMode {
		b.WriteString("Enter date (MM/DD/YYYY): ")
		b.WriteString(styles.HighlightBox.Width(12).Render(sv.dateInput + "_"))
	} else {
		b.WriteString("[Not set] ")
		b.WriteString(styles.Dim.Render("Press / to set checkride date"))
	}
	b.WriteString("\n\n")

	// Category filters
	b.WriteString("Filter: ")
	catMap := map[string]string{"": "All", "Theory": "Theory", "Chair Flying": "Chair", "Garmin 430": "GPS", "CFI Flights": "Flights"}
	for i, cat := range []string{"", "Theory", "Chair Flying", "Garmin 430", "CFI Flights"} {
		if cat == sv.category {
			b.WriteString(styles.SelectedFilter.Render(fmt.Sprintf("[%d] %s", i+1, catMap[cat])))
		} else {
			b.WriteString(fmt.Sprintf("[%d] %s", i+1, catMap[cat]))
		}
		b.WriteString(" ")
	}
	b.WriteString("\n\n")

	// Tasks
	if len(sv.filteredTasks) == 0 {
		b.WriteString(styles.Dim.Render("No study plan. Set a checkride date to begin."))
		b.WriteString("\n")
	} else {
		// Group by week
		weeks := sv.groupByWeek()
		for week, weekTasks := range weeks {
			b.WriteString(styles.Subtitle.Render("Week of " + week.Format("Jan 2")))
			b.WriteString("\n")
			for _, t := range weekTasks {
				idx := sv.findIndex(t.ID)
				sel := idx == sv.selectedIdx
				b.WriteString(sv.renderTask(t, sel))
			}
			b.WriteString("\n")
		}
	}

	// Help
	b.WriteString(styles.Dim.Render("\n[↑↓] Navigate  [Enter] Toggle  [/] Date  [Tab/1-5] Filter"))

	return b.String()
}

// findIndex finds task index in filtered list
func (sv *StudyView) findIndex(id uint) int {
	for i, t := range sv.filteredTasks {
		if t.ID == id {
			return i
		}
	}
	return -1
}

// groupByWeek groups tasks by week
func (sv *StudyView) groupByWeek() map[time.Time][]model.DailyTask {
	m := make(map[time.Time][]model.DailyTask)
	for _, t := range sv.filteredTasks {
		d := t.Date
		wd := int(d.Weekday())
		if wd == 0 {
			wd = 7
		}
		weekStart := time.Date(d.Year(), d.Month(), d.Day()-wd+1, 0, 0, 0, 0, time.UTC)
		m[weekStart] = append(m[weekStart], t)
	}
	return m
}

// renderTask renders a single task
func (sv *StudyView) renderTask(t model.DailyTask, selected bool) string {
	check := "[ ]"
	if t.Completed {
		check = "[x]"
	}
	catStyle := sv.categoryStyle(t.Category)
	if selected {
		return styles.SelectedTask.Render(fmt.Sprintf(" > %s %s - %s", check, catStyle.Render(t.Category), t.Title)) + "\n"
	}
	return fmt.Sprintf("   %s %s - %s\n", check, catStyle.Render(t.Category), t.Title)
}

// categoryStyle returns the style for a category
func (sv *StudyView) categoryStyle(cat string) lipgloss.Style {
	switch cat {
	case "Theory":
		return styles.CategoryTheory
	case "Chair Flying":
		return styles.CategoryChairFlying
	case "Garmin 430":
		return styles.CategoryGarmin
	case "CFI Flights":
		return styles.CategoryCFI
	}
	return styles.Dim
}
