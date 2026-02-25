package view

import (
	"fmt"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"ppl-study-planner/internal/model"
	"ppl-study-planner/internal/styles"
)

// BudgetView displays the flight training budget planner
type BudgetView struct {
	// Flight rates
	PlaneRate float64
	CfiRate   float64

	// Estimated hours
	DualGivenHours float64
	SoloHours      float64
	XcHours        float64
	SimulatorHours float64

	// Living costs
	TravelCost float64
	RentCost   float64
	FoodCost   float64
	CarCost    float64

	// Budget
	BudgetLimit float64

	// UI state
	selectedField int
	width         int
	height        int
	db            interface{}
}

// BudgetField represents a field in the budget form
type BudgetField struct {
	Name   string
	Value  *float64
	Label  string
	IsCost bool // true for costs (one-time), false for rates/hours
}

var budgetFields = []BudgetField{
	{"plane_rate", nil, "Plane $/hr", false},
	{"cfi_rate", nil, "CFI $/hr", false},
	{"dual_given", nil, "Dual Given hrs", false},
	{"solo_hours", nil, "Solo hrs", false},
	{"xc_hours", nil, "XC hrs", false},
	{"simulator_hours", nil, "Simulator hrs", false},
	{"travel", nil, "Travel $", true},
	{"rent", nil, "Rent/Month $", true},
	{"food", nil, "Food/Month $", true},
	{"car", nil, "Car one-time $", true},
	{"budget_limit", nil, "Total Budget $", true},
}

// NewBudgetView creates a new budget view
func NewBudgetView(db interface{}) *BudgetView {
	return &BudgetView{
		PlaneRate:      150,   // Default: $150/hr
		CfiRate:        60,    // Default: $60/hr
		DualGivenHours: 40,    // Default: 40 hours
		SoloHours:      10,    // Default: 10 hours
		XcHours:        5,     // Default: 5 hours
		SimulatorHours: 10,    // Default: 10 hours
		TravelCost:     500,   // Default: $500
		RentCost:       0,     // Default: $0
		FoodCost:       0,     // Default: $0
		CarCost:        0,     // Default: $0
		BudgetLimit:    10000, // Default: $10,000
		selectedField:  0,
		width:          80,
		height:         24,
		db:             db,
	}
}

// Init implements tea.Model
func (v *BudgetView) Init() tea.Cmd {
	return nil
}

// Update handles budget updates
func (v *BudgetView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if v.selectedField > 0 {
				v.selectedField--
			}
		case "down", "j":
			if v.selectedField < len(budgetFields)-1 {
				v.selectedField++
			}
		case "left", "h":
			v.adjustValue(-1)
		case "right", "l":
			v.adjustValue(1)
		case "tab":
			v.selectedField = (v.selectedField + 1) % len(budgetFields)
		}
	case tea.WindowSizeMsg:
		v.width = msg.Width
		v.height = msg.Height
	}
	return v, nil
}

// View renders the budget view
func (v *BudgetView) View() string {
	contentWidth := v.width - 4
	box := styles.HighlightBox.Width(contentWidth).Height(v.height - 4)

	// Calculate costs
	costs := v.calculateCosts()

	// Header
	header := styles.Title.Render("Flight Training Budget")

	// Input fields
	inputs := v.renderInputs()

	// Calculations
	calculations := v.renderCalculations(costs)

	// Warning if over budget
	warning := ""
	if costs.Remaining < 0 {
		warning = styles.ErrorStyle.Render(fmt.Sprintf(`
  ⚠️ OVER BUDGET by $%.2f

  Reduce flight hours or living costs to stay within budget.`,
			-costs.Remaining))
	}

	// Help
	help := styles.Dim.Render(" [↑↓] Navigate fields | [←→] Adjust value | [Tab] Next field")

	content := fmt.Sprintf(`%s

%s

%s

%s`,
		header,
		inputs,
		calculations,
		warning+"\n\n"+help,
	)

	return box.Render(content)
}

func (v *BudgetView) renderInputs() string {
	lines := []string{
		styles.Normal.Render("┌─ Flight Rates ──────────────────────────────────────────────┐"),
	}

	fields := []struct {
		label string
		value float64
		field int
	}{
		{"Plane $/hr:", v.PlaneRate, 0},
		{"CFI $/hr:", v.CfiRate, 1},
		{"Dual Given hrs:", v.DualGivenHours, 2},
		{"Solo hrs:", v.SoloHours, 3},
		{"XC hrs:", v.XcHours, 4},
		{"Simulator hrs:", v.SimulatorHours, 5},
	}

	for _, f := range fields {
		value := fmt.Sprintf("%.2f", f.value)
		if f.field == v.selectedField {
			lines = append(lines, fmt.Sprintf("│ ● %-14s %-15s                           │", f.label, styles.Selected.Render(value)))
		} else {
			lines = append(lines, fmt.Sprintf("│   %-14s %-15s                           │", f.label, value))
		}
	}

	lines = append(lines, styles.Normal.Render("├─ Living Costs (one-time or monthly) ───────────────────────┤"))

	lifeFields := []struct {
		label string
		value float64
		field int
	}{
		{"Travel:", v.TravelCost, 6},
		{"Rent/Month:", v.RentCost, 7},
		{"Food/Month:", v.FoodCost, 8},
		{"Car (one-time):", v.CarCost, 9},
	}

	for _, f := range lifeFields {
		value := fmt.Sprintf("%.2f", f.value)
		if f.field == v.selectedField {
			lines = append(lines, fmt.Sprintf("│ ● %-14s %-15s                           │", f.label, styles.Selected.Render(value)))
		} else {
			lines = append(lines, fmt.Sprintf("│   %-14s %-15s                           │", f.label, value))
		}
	}

	lines = append(lines, styles.Normal.Render("├─ Budget ────────────────────────────────────────────────────┤"))

	budgetLabel := "Total Budget $:"
	budgetValue := fmt.Sprintf("%.2f", v.BudgetLimit)
	if v.selectedField == 10 {
		lines = append(lines, fmt.Sprintf("│ ● %-14s %-15s                           │", budgetLabel, styles.Selected.Render(budgetValue)))
	} else {
		lines = append(lines, fmt.Sprintf("│   %-14s %-15s                           │", budgetLabel, budgetValue))
	}

	lines = append(lines, styles.Normal.Render("└──────────────────────────────────────────────────────────────┘"))

	return strings.Join(lines, "\n")
}

func (v *BudgetView) renderCalculations(costs BudgetCosts) string {
	lines := []string{
		"",
		styles.Normal.Render("┌─ Cost Breakdown ────────────────────────────────────────────┐"),
	}

	flightLine := fmt.Sprintf("│ Flight Costs:                                           │")
	lines = append(lines, flightLine)

	lines = append(lines, fmt.Sprintf("│   Plane rental:    $%9.2f × %5.2f hrs = $%9.2f    │",
		v.PlaneRate, v.DualGivenHours+v.SoloHours+v.XcHours+v.SimulatorHours, costs.FlightCost))
	lines = append(lines, fmt.Sprintf("│   CFI instruction: $%9.2f × %5.2f hrs = $%9.2f    │",
		v.CfiRate, v.DualGivenHours, costs.CfiCost))

	lines = append(lines, styles.Dim.Render("├───────────────────────────────────────────────────────────────┤"))
	lines = append(lines, fmt.Sprintf("│   %-30s = $%9.2f          │", "Total Flight Cost:", costs.FlightCost+costs.CfiCost))

	lines = append(lines, "")
	lines = append(lines, styles.Normal.Render("│ Living Costs:                                           │"))
	lines = append(lines, fmt.Sprintf("│   Travel:        $%9.2f                           │", v.TravelCost))
	lines = append(lines, fmt.Sprintf("│   Rent:          $%9.2f                           │", v.RentCost))
	lines = append(lines, fmt.Sprintf("│   Food:          $%9.2f                           │", v.FoodCost))
	lines = append(lines, fmt.Sprintf("│   Car:           $%9.2f                           │", v.CarCost))

	lines = append(lines, styles.Dim.Render("├───────────────────────────────────────────────────────────────┤"))
	lines = append(lines, fmt.Sprintf("│   %-30s = $%9.2f          │", "Total Living Cost:", costs.LivingCost))

	lines = append(lines, "")
	lines = append(lines, styles.Normal.Render("├═════════════════════════════════════════════════════════════│"))
	lines = append(lines, fmt.Sprintf("│ %-30s = $%9.2f          │", "TOTAL PROJECTED:", costs.Total))
	lines = append(lines, fmt.Sprintf("│ %-30s = $%9.2f          │", "Budget Limit:", v.BudgetLimit))

	remainingStyle := styles.Success
	if costs.Remaining < 0 {
		remainingStyle = styles.ErrorStyle
	}
	lines = append(lines, fmt.Sprintf("│ %-30s = %s│", "Remaining:", remainingStyle.Render(fmt.Sprintf("$%.2f", costs.Remaining))))
	lines = append(lines, fmt.Sprintf("│ %-30s = %s                         │", "Budget Used:", styles.WarningStyle.Render(fmt.Sprintf("%.1f%%", costs.PercentUsed))))

	lines = append(lines, styles.Normal.Render("└───────────────────────────────────────────────────────────────┘"))

	return strings.Join(lines, "\n")
}

// BudgetCosts holds the calculated budget breakdown
type BudgetCosts struct {
	FlightCost  float64
	CfiCost     float64
	LivingCost  float64
	Total       float64
	Remaining   float64
	PercentUsed float64
}

func (v *BudgetView) calculateCosts() BudgetCosts {
	// Flight costs
	totalHours := v.DualGivenHours + v.SoloHours + v.XcHours + v.SimulatorHours
	flightCost := v.PlaneRate * totalHours
	cfiCost := v.CfiRate * v.DualGivenHours

	// Living costs (one-time)
	livingCost := v.TravelCost + v.RentCost + v.FoodCost + v.CarCost

	// Totals
	total := flightCost + cfiCost + livingCost
	remaining := v.BudgetLimit - total
	percentUsed := 0.0
	if v.BudgetLimit > 0 {
		percentUsed = (total / v.BudgetLimit) * 100
	}

	return BudgetCosts{
		FlightCost:  flightCost,
		CfiCost:     cfiCost,
		LivingCost:  livingCost,
		Total:       total,
		Remaining:   remaining,
		PercentUsed: percentUsed,
	}
}

func (v *BudgetView) adjustValue(delta int) {
	// Adjust value based on selected field
	// Different step sizes for different fields
	step := 1.0
	switch v.selectedField {
	case 0: // plane_rate
		step = 5 // $5/hr steps
		v.PlaneRate = adjustFloat(v.PlaneRate, delta, step, 0, 500)
	case 1: // cfi_rate
		step = 5 // $5/hr steps
		v.CfiRate = adjustFloat(v.CfiRate, delta, step, 0, 200)
	case 2: // dual_given
		step = 1 // 1 hour steps
		v.DualGivenHours = adjustFloat(v.DualGivenHours, delta, step, 0, 100)
	case 3: // solo_hours
		step = 1
		v.SoloHours = adjustFloat(v.SoloHours, delta, step, 0, 50)
	case 4: // xc_hours
		step = 1
		v.XcHours = adjustFloat(v.XcHours, delta, step, 0, 50)
	case 5: // simulator_hours
		step = 1
		v.SimulatorHours = adjustFloat(v.SimulatorHours, delta, step, 0, 50)
	case 6: // travel
		step = 50
		v.TravelCost = adjustFloat(v.TravelCost, delta, step, 0, 10000)
	case 7: // rent
		step = 100
		v.RentCost = adjustFloat(v.RentCost, delta, step, 0, 5000)
	case 8: // food
		step = 50
		v.FoodCost = adjustFloat(v.FoodCost, delta, step, 0, 3000)
	case 9: // car
		step = 500
		v.CarCost = adjustFloat(v.CarCost, delta, step, 0, 50000)
	case 10: // budget_limit
		step = 500
		v.BudgetLimit = adjustFloat(v.BudgetLimit, delta, step, 1000, 100000)
	}
}

func adjustFloat(current float64, delta int, step float64, min float64, max float64) float64 {
	newVal := current + (float64(delta) * step)
	if newVal < min {
		return min
	}
	if newVal > max {
		return max
	}
	// Round to 2 decimal places
	return float64(int64(newVal*100)) / 100
}

// SetRates sets the flight rates
func (v *BudgetView) SetRates(planeRate, cfiRate float64) {
	v.PlaneRate = planeRate
	v.CfiRate = cfiRate
}

// SetHours sets the flight hours
func (v *BudgetView) SetHours(dual, solo, xc, sim float64) {
	v.DualGivenHours = dual
	v.SoloHours = solo
	v.XcHours = xc
	v.SimulatorHours = sim
}

// SetLivingCosts sets the living costs
func (v *BudgetView) SetLivingCosts(travel, rent, food, car float64) {
	v.TravelCost = travel
	v.RentCost = rent
	v.FoodCost = food
	v.CarCost = car
}

// SetBudgetLimit sets the budget limit
func (v *BudgetView) SetBudgetLimit(limit float64) {
	v.BudgetLimit = limit
}

// Budget returns the current budget settings
func (v *BudgetView) Budget() (planeRate, cfiRate, budget float64, costs model.Budget) {
	return v.PlaneRate, v.CfiRate, v.BudgetLimit, model.Budget{}
}

// ParseCurrency parses a currency string to float64
func ParseCurrency(s string) (float64, error) {
	s = strings.ReplaceAll(s, "$", "")
	s = strings.ReplaceAll(s, ",", "")
	return strconv.ParseFloat(s, 64)
}

// FormatCurrency formats a float64 as currency
func FormatCurrency(f float64) string {
	return fmt.Sprintf("$%.2f", f)
}
