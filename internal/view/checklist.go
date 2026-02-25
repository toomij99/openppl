package view

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"ppl-study-planner/internal/model"
	"ppl-study-planner/internal/styles"
)

// ChecklistView displays the FAA pre-checkride checklist with toggleable items
type ChecklistView struct {
	items         []model.ChecklistItem
	selectedIndex int
	categoryIndex int // 0=All, 1=Documents, 2=Aircraft, 3=Ground, 4=Flight
	width         int
	height        int
	db            interface{}
}

// categories represents the filter categories
var categories = []string{"All", "Documents", "Aircraft", "Ground", "Flight"}

// NewChecklistView creates a new checklist view
func NewChecklistView(db interface{}) *ChecklistView {
	return &ChecklistView{
		items:         []model.ChecklistItem{},
		selectedIndex: 0,
		categoryIndex: 0,
		width:         80,
		height:        24,
		db:            db,
	}
}

// Init loads the checklist items from database
func (v *ChecklistView) Init() tea.Cmd {
	// Load items from database
	// This would be: v.db.Find(&v.items)
	return nil
}

// Update handles checklist updates
func (v *ChecklistView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if v.selectedIndex > 0 {
				v.selectedIndex--
			}
		case "down", "j":
			filtered := v.getFilteredItems()
			if v.selectedIndex < len(filtered)-1 {
				v.selectedIndex++
			}
		case "enter", " ":
			v.toggleItem()
		case "tab":
			v.categoryIndex = (v.categoryIndex + 1) % len(categories)
			v.selectedIndex = 0
		}
	case tea.WindowSizeMsg:
		v.width = msg.Width
		v.height = msg.Height
	}
	return v, nil
}

// View renders the checklist
func (v *ChecklistView) View() string {
	contentWidth := v.width - 4
	box := styles.HighlightBox.Width(contentWidth).Height(v.height - 4)

	filtered := v.getFilteredItems()
	categoryStats := v.calculateCategoryStats()

	// Header with category filter
	header := styles.Title.Render("Pre-Checkride Checklist")

	// Category tabs
	categoryTabs := v.renderCategoryTabs()

	// Stats line
	stats := fmt.Sprintf("  Overall: %s | %s",
		styles.Success.Render(fmt.Sprintf("%d/%d completed", categoryStats.Completed, categoryStats.Total)),
		styles.Normal.Render(fmt.Sprintf("%d%%", categoryStats.Percentage)),
	)

	// Category stats
	catStats := fmt.Sprintf("  Documents: %d/%d | Aircraft: %d/%d | Ground: %d/%d | Flight: %d/%d",
		catStatsFor(filtered, model.CategoryDocuments).Completed,
		catStatsFor(filtered, model.CategoryDocuments).Total,
		catStatsFor(filtered, model.CategoryAircraft).Completed,
		catStatsFor(filtered, model.CategoryAircraft).Total,
		catStatsFor(filtered, model.CategoryGround).Completed,
		catStatsFor(filtered, model.CategoryGround).Total,
		catStatsFor(filtered, model.CategoryFlight).Completed,
		catStatsFor(filtered, model.CategoryFlight).Total,
	)

	// Checklist items
	items := v.renderItems(filtered)

	// Footer help
	footer := styles.Dim.Render(" [↑↓] Navigate | [Enter/Space] Toggle | [Tab] Filter by category")

	content := fmt.Sprintf(`%s

%s
%s

%s

%s

%s`,
		header,
		categoryTabs,
		stats,
		catStats,
		items,
		footer,
	)

	return box.Render(content)
}

func (v *ChecklistView) renderCategoryTabs() string {
	var tabs []string
	for i, cat := range categories {
		if i == v.categoryIndex {
			tabs = append(tabs, styles.Selected.Render("[ "+cat+" ]"))
		} else {
			tabs = append(tabs, styles.Dim.Render(" "+cat+" "))
		}
	}
	return strings.Join(tabs, "")
}

func (v *ChecklistView) renderItems(items []model.ChecklistItem) string {
	if len(items) == 0 {
		return styles.Dim.Render("  No items in this category")
	}

	var lines []string
	currentCategory := ""

	for i, item := range items {
		// Add category header if category changed
		if string(item.Category) != currentCategory && v.categoryIndex == 0 {
			currentCategory = string(item.Category)
			lines = append(lines, "")
			lines = append(lines, styles.Subtitle.Render("■ "+currentCategory))
		}

		// Render item
		checkbox := "[ ]"
		if item.Completed {
			checkbox = styles.Success.Render("[✓]")
		}

		line := fmt.Sprintf("  %s %s", checkbox, item.Title)

		if i == v.selectedIndex && v.categoryIndex > 0 ||
			(v.categoryIndex == 0 && isSelectedCategory(item.Category, items, v.selectedIndex)) {
			line = styles.Selected.Render("▶ " + strings.TrimPrefix(line, "  "))
		}

		if item.Completed {
			line = styles.Dim.Render(line)
		}

		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}

func isSelectedCategory(cat model.ChecklistCategory, items []model.ChecklistItem, selectedIdx int) bool {
	// Find the actual index considering category headers
	idx := 0
	for _, item := range items {
		if idx == selectedIdx {
			return item.Category == cat
		}
		if item.Category == cat {
			idx++
		}
	}
	return false
}

func (v *ChecklistView) getFilteredItems() []model.ChecklistItem {
	if v.categoryIndex == 0 {
		return v.items
	}

	category := categories[v.categoryIndex]
	var filtered []model.ChecklistItem
	for _, item := range v.items {
		if string(item.Category) == category {
			filtered = append(filtered, item)
		}
	}
	return filtered
}

func (v *ChecklistView) toggleItem() {
	filtered := v.getFilteredItems()
	if v.selectedIndex >= 0 && v.selectedIndex < len(filtered) {
		filtered[v.selectedIndex].Completed = !filtered[v.selectedIndex].Completed
		// Save to database: v.db.Save(&filtered[v.selectedIndex])
	}
}

type categoryStats struct {
	Completed  int
	Total      int
	Percentage int
}

func (v *ChecklistView) calculateCategoryStats() categoryStats {
	var completed, total int
	for _, item := range v.items {
		if item.Completed {
			completed++
		}
		total++
	}
	percentage := 0
	if total > 0 {
		percentage = (completed * 100) / total
	}
	return categoryStats{Completed: completed, Total: total, Percentage: percentage}
}

func catStatsFor(items []model.ChecklistItem, category model.ChecklistCategory) categoryStats {
	var completed, total int
	for _, item := range items {
		if item.Category == category {
			total++
			if item.Completed {
				completed++
			}
		}
	}
	percentage := 0
	if total > 0 {
		percentage = (completed * 100) / total
	}
	return categoryStats{Completed: completed, Total: total, Percentage: percentage}
}

// SetItems sets the checklist items
func (v *ChecklistView) SetItems(items []model.ChecklistItem) {
	v.items = items
}

// GetItems returns the checklist items
func (v *ChecklistView) GetItems() []model.ChecklistItem {
	return v.items
}
