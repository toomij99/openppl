package tui

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"ppl-study-planner/internal/db"
	"ppl-study-planner/internal/model"
	"ppl-study-planner/internal/styles"

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
	width         int
	height        int
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
	content := renderContent(m.currentScreen, m.width, m.height)

	return header + "\n" + content + "\n" + footer
}

func renderHeader(screen Screen) string {
	title := styles.ScreenTitle[screen]
	return styles.Header.Width(styles.Footer.GetWidth()).Render(
		fmt.Sprintf(" openppl | %s ", title),
	)
}

func renderFooter() string {
	return styles.Footer.Width(80).Render(
		" [1] Dashboard | [2] Study Plan | [3] Progress | [4] Budget | [5] Checklist | [q] Quit ",
	)
}

func renderContent(screen Screen, width, height int) string {
	contentWidth := width - 4
	contentHeight := height - 4

	switch screen {
	case ScreenDashboard:
		return renderDashboard(contentWidth, contentHeight)
	case ScreenStudyPlan:
		return renderStudyPlan(contentWidth, contentHeight)
	case ScreenProgress:
		return renderProgress(contentWidth, contentHeight)
	case ScreenBudget:
		return renderBudget(contentWidth, contentHeight)
	case ScreenChecklist:
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
