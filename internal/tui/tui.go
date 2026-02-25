package tui

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"ppl-study-planner/internal/db"
	"ppl-study-planner/internal/model"
	"ppl-study-planner/internal/styles"
	"ppl-study-planner/internal/view"

	"gorm.io/gorm"
)

// Screen represents the different screens in the TUI
type Screen int

const (
	ScreenDashboard Screen = iota
	ScreenStudyPlan
	ScreenProgress
	ScreenBudget
	ScreenChecklist
)

// MainModel contains all application state
type MainModel struct {
	db            *gorm.DB
	currentScreen Screen
	helpVisible   bool
	width         int
	height        int
	studyView     *view.StudyView
	progressView  *view.ProgressView
	dashboardView *view.DashboardView
	checklistView *view.ChecklistView
	budgetView    *view.BudgetView
}

// New creates a new TUI model
func New() (*MainModel, error) {
	database, err := db.Initialize()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	return &MainModel{
		db:            database,
		currentScreen: ScreenDashboard,
		width:         80,
		height:        24,
		studyView:     view.NewStudyView(database),
		progressView:  view.NewProgressView(database),
		dashboardView: view.NewDashboardView(database),
		checklistView: view.NewChecklistView(database),
		budgetView:    view.NewBudgetView(database),
	}, nil
}

// Init implements tea.Model
func (m MainModel) Init() tea.Cmd {
	// Seed checklist items if empty
	var count int64
	m.db.Model(&model.ChecklistItem{}).Count(&count)
	if count == 0 {
		checklistItems := []model.ChecklistItem{
			{Category: model.CategoryDocuments, Title: "Pilot certificate (Airplane category)"},
			{Category: model.CategoryDocuments, Title: "Medical certificate (Class 3 or higher)"},
			{Category: model.CategoryDocuments, Title: "Logbook with required entries"},
			{Category: model.CategoryDocuments, Title: "Form AC 61-91 (if using simulator)"},
			{Category: model.CategoryAircraft, Title: "Airworthiness certificate"},
			{Category: model.CategoryAircraft, Title: "Registration certificate"},
			{Category: model.CategoryAircraft, Title: "Operating limitations"},
			{Category: model.CategoryAircraft, Title: "Weight and balance report"},
			{Category: model.CategoryAircraft, Title: "Maintenance logs"},
			{Category: model.CategoryGround, Title: "Charts and publications"},
			{Category: model.CategoryGround, Title: "Flight planner"},
			{Category: model.CategoryGround, Title: "Weather briefing documentation"},
			{Category: model.CategoryFlight, Title: "Pre-solo knowledge test passed"},
			{Category: model.CategoryFlight, Title: "Solo endorsements (3 takeoffs/landings)"},
			{Category: model.CategoryFlight, Title: "Cross-country endorsements"},
			{Category: model.CategoryFlight, Title: "Night endorsement"},
			{Category: model.CategoryFlight, Title: "Instrument proficiency"},
		}
		m.db.Create(&checklistItems)
	}

	return nil
}

// Update implements tea.Model
func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "?" || msg.String() == "f1" {
			m.helpVisible = !m.helpVisible
			return m, nil
		}

		if m.helpVisible {
			if msg.String() == "esc" {
				m.helpVisible = false
				return m, nil
			}
			if msg.String() != "ctrl+c" && msg.String() != "q" {
				return m, nil
			}
		}

		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "1":
			m.currentScreen = ScreenDashboard
		case "2":
			m.currentScreen = ScreenStudyPlan
		case "3":
			m.currentScreen = ScreenProgress
		case "4":
			m.currentScreen = ScreenBudget
		case "5":
			m.currentScreen = ScreenChecklist
		}

		// Route messages to dashboard view when on dashboard screen
		if m.currentScreen == ScreenDashboard && m.dashboardView != nil {
			updated, cmd := m.dashboardView.Update(msg)
			m.dashboardView = updated.(*view.DashboardView)
			return m, cmd
		}

		// Route messages to study view when on study plan screen
		if m.currentScreen == ScreenStudyPlan && m.studyView != nil {
			updated, cmd := m.studyView.Update(msg)
			m.studyView = updated.(*view.StudyView)
			return m, cmd
		}

		// Route messages to progress view when on progress screen
		if m.currentScreen == ScreenProgress && m.progressView != nil {
			updated, _ := m.progressView.Update(msg)
			m.progressView = updated.(*view.ProgressView)
		}

		// Route messages to budget view when on budget screen
		if m.currentScreen == ScreenBudget && m.budgetView != nil {
			updated, cmd := m.budgetView.Update(msg)
			m.budgetView = updated.(*view.BudgetView)
			return m, cmd
		}

		// Route messages to checklist view when on checklist screen
		if m.currentScreen == ScreenChecklist && m.checklistView != nil {
			updated, cmd := m.checklistView.Update(msg)
			m.checklistView = updated.(*view.ChecklistView)
			return m, cmd
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, nil
}

// View implements tea.Model
func (m MainModel) View() string {
	header := renderHeader(m.currentScreen)
	footer := renderFooter()
	content := renderContent(m.currentScreen, m.width, m.height, m.dashboardView, m.studyView, m.progressView, m.budgetView, m.checklistView)
	if m.helpVisible {
		content = renderHelpOverlay(m.currentScreen)
	}

	return header + "\n" + content + "\n" + footer
}

func renderHeader(screen Screen) string {
	title := styles.ScreenTitle[screen]
	return styles.Header.Width(styles.Footer.GetWidth()).Render(
		fmt.Sprintf(" openppl | %s ", title),
	)
}

func renderFooter() string {
	parts := make([]string, 0, len(FooterShortcuts()))
	for _, shortcut := range FooterShortcuts() {
		parts = append(parts, fmt.Sprintf("[%s] %s", shortcut.Keys, shortcut.Action))
	}

	return styles.Footer.Width(80).Render(
		" " + strings.Join(parts, " | ") + " ",
	)
}

func renderHelpOverlay(screen Screen) string {
	sections := HelpSections(screen)
	var b strings.Builder
	b.WriteString(styles.Title.Render("Keyboard Shortcuts"))
	b.WriteString("\n\n")

	for _, section := range sections {
		b.WriteString(styles.Subtitle.Render(section.Title))
		b.WriteString("\n")
		for _, shortcut := range section.Shortcuts {
			b.WriteString(fmt.Sprintf("  %-14s %s\n", shortcut.Keys, shortcut.Action))
		}
		b.WriteString("\n")
	}

	b.WriteString(styles.Dim.Render("Press ? or F1 to close help"))
	return styles.HighlightBox.Width(76).Render(b.String())
}

func renderContent(screen Screen, width, height int, dashboardView *view.DashboardView, studyView *view.StudyView, progressView *view.ProgressView, budgetView *view.BudgetView, checklistView *view.ChecklistView) string {
	contentWidth := width - 4
	contentHeight := height - 4

	switch screen {
	case ScreenDashboard:
		if dashboardView != nil {
			return dashboardView.View()
		}
		return renderDashboard(contentWidth, contentHeight)
	case ScreenStudyPlan:
		if studyView != nil {
			return studyView.View()
		}
		return renderStudyPlan(contentWidth, contentHeight)
	case ScreenProgress:
		if progressView != nil {
			return progressView.View()
		}
		return renderProgress(contentWidth, contentHeight)
	case ScreenBudget:
		if budgetView != nil {
			return budgetView.View()
		}
		return renderBudget(contentWidth, contentHeight)
	case ScreenChecklist:
		if checklistView != nil {
			return checklistView.View()
		}
		return renderChecklist(contentWidth, contentHeight)
	default:
		return ""
	}
}

func renderDashboard(width, height int) string {
	box := styles.HighlightBox.Width(width).Height(height)
	return box.Render(
		styles.Title.Render("Dashboard") + "\n\n" +
			styles.Normal.Render("Days until checkride: --") + "\n\n" +
			styles.Normal.Render("Today's Tasks:") + "\n" +
			styles.Dim.Render("  No tasks scheduled") + "\n\n" +
			styles.Normal.Render("Quick Stats:") + "\n" +
			styles.Dim.Render("  Completed: 0") + "\n" +
			styles.Dim.Render("  Remaining: 0") + "\n" +
			styles.Dim.Render("  Overdue: 0") + "\n" +
			styles.Dim.Render("  Total: 0"),
	)
}

func renderStudyPlan(width, height int) string {
	box := styles.HighlightBox.Width(width).Height(height)
	return box.Render(
		styles.Title.Render("Study Plan") + "\n\n" +
			styles.Normal.Render("Checkride Date: Not set") + "\n\n" +
			styles.Dim.Render("Use the Study Plan screen to manage your daily tasks.") + "\n" +
			styles.Dim.Render("Set a checkride date to begin backward planning."),
	)
}

func renderProgress(width, height int) string {
	box := styles.HighlightBox.Width(width).Height(height)
	return box.Render(
		styles.Title.Render("Progress") + "\n\n" +
			styles.Normal.Render("Overall Progress:") + "\n" +
			styles.Dim.Render("  [                    ] 0%") + "\n\n" +
			styles.Normal.Render("By Category:") + "\n" +
			styles.Dim.Render("  Documents: 0/4") + "\n" +
			styles.Dim.Render("  Aircraft: 0/5") + "\n" +
			styles.Dim.Render("  Ground: 0/3") + "\n" +
			styles.Dim.Render("  Flight: 0/5"),
	)
}

func renderBudget(width, height int) string {
	box := styles.HighlightBox.Width(width).Height(height)
	return box.Render(
		styles.Title.Render("Budget") + "\n\n" +
			styles.Normal.Render("Estimated Costs:") + "\n" +
			styles.Dim.Render("  Flight Hours: $0.00") + "\n" +
			styles.Dim.Render("  Plane Rental: $0.00") + "\n" +
			styles.Dim.Render("  CFI Rate: $0.00") + "\n" +
			styles.Dim.Render("  Living Expenses: $0.00") + "\n\n" +
			styles.Normal.Render("Total: $0.00"),
	)
}

func renderChecklist(width, height int) string {
	box := styles.HighlightBox.Width(width).Height(height)
	return box.Render(
		styles.Title.Render("Pre-Checkride Checklist") + "\n\n" +
			styles.Normal.Render("Documents (0/4)") + "\n" +
			styles.Dim.Render("  [ ] Pilot certificate") + "\n" +
			styles.Dim.Render("  [ ] Medical certificate") + "\n" +
			styles.Dim.Render("  [ ] Logbook") + "\n" +
			styles.Dim.Render("  [ ] Form AC 61-91") + "\n\n" +
			styles.Normal.Render("Press SPACE to toggle items"),
	)
}

// Run starts the TUI application
func Run() error {
	model, err := New()
	if err != nil {
		return fmt.Errorf("failed to create model: %w", err)
	}
	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running TUI: %v\n", err)
		os.Exit(1)
	}

	return nil
}
