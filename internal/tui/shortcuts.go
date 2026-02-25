package tui

// Shortcut defines a keyboard shortcut and its display context.
type Shortcut struct {
	Keys    string
	Action  string
	Section string
	Footer  bool
}

// HelpSection groups shortcuts for help overlay rendering.
type HelpSection struct {
	Title     string
	Shortcuts []Shortcut
}

var shortcutRegistry = []Shortcut{
	{Keys: "1-5", Action: "Switch screens (Dashboard, Study, Progress, Budget, Checklist)", Section: "Global Navigation", Footer: true},
	{Keys: "up/down", Action: "Move selection", Section: "Global Navigation", Footer: false},
	{Keys: "enter", Action: "Select or toggle focused item", Section: "Global Navigation", Footer: false},
	{Keys: "q", Action: "Quit app", Section: "App Controls", Footer: true},
	{Keys: "ctrl+c", Action: "Force quit", Section: "App Controls", Footer: false},
	{Keys: "? / F1", Action: "Toggle help", Section: "App Controls", Footer: true},
	{Keys: "/", Action: "Set or edit checkride date", Section: "Study Actions", Footer: false},
	{Keys: "tab", Action: "Cycle study category filter", Section: "Study Actions", Footer: false},
	{Keys: "1-5", Action: "Filter study categories", Section: "Study Actions", Footer: false},
	{Keys: "e", Action: "Export ICS", Section: "Study Actions", Footer: false},
	{Keys: "r", Action: "Export Apple Reminders", Section: "Study Actions", Footer: false},
	{Keys: "g", Action: "Sync Google Calendar", Section: "Study Actions", Footer: false},
	{Keys: "o", Action: "Export OpenCode bot tasks", Section: "Study Actions", Footer: false},
	{Keys: "up/down + enter", Action: "Toggle study task completion", Section: "Study Actions", Footer: false},
}

// AllShortcuts returns all shortcut definitions.
func AllShortcuts() []Shortcut {
	shortcuts := make([]Shortcut, len(shortcutRegistry))
	copy(shortcuts, shortcutRegistry)
	return shortcuts
}

// FooterShortcuts returns shortcuts displayed in the footer hint.
func FooterShortcuts() []Shortcut {
	all := AllShortcuts()
	footer := make([]Shortcut, 0, len(all))
	for _, shortcut := range all {
		if shortcut.Footer {
			footer = append(footer, shortcut)
		}
	}
	return footer
}

// HelpSections returns grouped help content for the current screen.
func HelpSections(_ Screen) []HelpSection {
	globalNavigation := make([]Shortcut, 0)
	appControls := make([]Shortcut, 0)
	studyActions := make([]Shortcut, 0)

	for _, shortcut := range shortcutRegistry {
		switch shortcut.Section {
		case "Global Navigation":
			globalNavigation = append(globalNavigation, shortcut)
		case "App Controls":
			appControls = append(appControls, shortcut)
		case "Study Actions":
			studyActions = append(studyActions, shortcut)
		}
	}

	return []HelpSection{
		{Title: "Global Navigation", Shortcuts: globalNavigation},
		{Title: "App Controls", Shortcuts: appControls},
		{Title: "Study Actions", Shortcuts: studyActions},
	}
}
