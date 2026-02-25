package tui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestShortcutRegistry(t *testing.T) {
	all := AllShortcuts()
	if len(all) == 0 {
		t.Fatal("expected shortcut registry to be populated")
	}

	required := map[string]bool{
		"q":      false,
		"ctrl+c": false,
		"? / F1": false,
		"/":      false,
		"tab":    false,
		"1-5":    false,
		"e":      false,
		"r":      false,
		"g":      false,
		"o":      false,
	}

	for _, shortcut := range all {
		if _, ok := required[shortcut.Keys]; ok {
			required[shortcut.Keys] = true
		}
	}

	for key, seen := range required {
		if !seen {
			t.Fatalf("expected key %q in centralized shortcut registry", key)
		}
	}

	footer := FooterShortcuts()
	if len(footer) == 0 {
		t.Fatal("expected footer shortcuts to include help and navigation hints")
	}

	helpFound := false
	for _, shortcut := range footer {
		if shortcut.Keys == "? / F1" {
			helpFound = true
			break
		}
	}
	if !helpFound {
		t.Fatal("expected footer shortcuts to include ? / F1 help hint")
	}

	sections := HelpSections(ScreenDashboard)
	if len(sections) < 3 {
		t.Fatalf("expected at least 3 help sections, got %d", len(sections))
	}
}

func TestHelpOverlayToggle(t *testing.T) {
	model := MainModel{currentScreen: ScreenDashboard, width: 80, height: 24}

	updatedModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
	toggled := updatedModel.(MainModel)
	if !toggled.helpVisible {
		t.Fatal("expected ? to open help overlay")
	}

	updatedModel, _ = toggled.Update(tea.KeyMsg{Type: tea.KeyF1})
	toggled = updatedModel.(MainModel)
	if toggled.helpVisible {
		t.Fatal("expected f1 to close help overlay")
	}

	footer := renderFooter()
	if !strings.Contains(footer, "? / F1") {
		t.Fatalf("expected footer to advertise ? / F1 help shortcut, got: %s", footer)
	}
}

func TestHelpOutputIncludesStudyActions(t *testing.T) {
	model := MainModel{currentScreen: ScreenStudyPlan, helpVisible: true, width: 80, height: 24}
	view := model.View()

	checks := []string{
		"Keyboard Shortcuts",
		"Global Navigation",
		"Study Actions",
		"Export ICS",
		"Export Apple Reminders",
		"Sync Google Calendar",
		"Export OpenCode bot tasks",
		"? / F1",
	}

	for _, check := range checks {
		if !strings.Contains(view, check) {
			t.Fatalf("expected help output to include %q", check)
		}
	}
}
