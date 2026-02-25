package styles

import "github.com/charmbracelet/lipgloss"

// Color palette from Charm
var (
	Primary   = lipgloss.Color("#BD93F9") // Purple
	Secondary = lipgloss.Color("#6272A4") // Muted blue
	Accent    = lipgloss.Color("#50FA7B") // Green
	Error     = lipgloss.Color("#FF5555") // Red
	Warning   = lipgloss.Color("#F1FA8C") // Yellow
	Info      = lipgloss.Color("#8BE9FD") // Cyan
)

// Common styles
var (
	Title = lipgloss.NewStyle().
		Foreground(Primary).
		Bold(true).
		Align(lipgloss.Center)

	Subtitle = lipgloss.NewStyle().
			Foreground(Secondary).
			Align(lipgloss.Center)

	Header = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#F8F8F2")).
		Background(Primary).
		Bold(true).
		Padding(0, 1)

	Footer = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#F8F8F2")).
		Background(Secondary).
		Padding(0, 1)

	Selected = lipgloss.NewStyle().
			Foreground(Accent).
			Bold(true)

	Normal = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#F8F8F2"))

	Dim = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6272A4"))

	Success = lipgloss.NewStyle().
		Foreground(Accent)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(Error)

	WarningStyle = lipgloss.NewStyle().
			Foreground(Warning)

	// Container styles
	BorderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Primary)

	HighlightBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Accent).
			Padding(1)

	// List styles
	ListItem = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F8F8F2")).
			Padding(0, 2)

	ListItemSelected = lipgloss.NewStyle().
				Foreground(Accent).
				Background(lipgloss.Color("#44475A")).
				Padding(0, 2)
)

// Screen titles
var ScreenTitle = [5]string{
	"Dashboard",
	"Study Plan",
	"Progress",
	"Budget",
	"Checklist",
}

// Category colors for study tasks
var (
	CategoryTheory      = lipgloss.NewStyle().Foreground(lipgloss.Color("99"))
	CategoryChairFlying = lipgloss.NewStyle().Foreground(lipgloss.Color("141"))
	CategoryGarmin      = lipgloss.NewStyle().Foreground(lipgloss.Color("86"))
	CategoryCFI         = lipgloss.NewStyle().Foreground(lipgloss.Color("208"))
	SelectedTask        = lipgloss.NewStyle().Background(lipgloss.Color("57")).Foreground(lipgloss.Color("255"))
	SelectedFilter      = lipgloss.NewStyle().Background(lipgloss.Color("57")).Foreground(lipgloss.Color("255"))
	ProgressBar         = lipgloss.NewStyle().Foreground(Accent)
	ProgressBarEmpty    = lipgloss.NewStyle().Foreground(Secondary)
)
