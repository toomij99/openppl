package main

import (
	"strings"
	"testing"
)

func TestResolveCommand_RecognizesAliases(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		wantCmd   string
		wantAfter int
	}{
		{name: "help short", args: []string{"-h"}, wantCmd: "help", wantAfter: 0},
		{name: "configure flag", args: []string{"--configure"}, wantCmd: "configure", wantAfter: 0},
		{name: "onboarding alias", args: []string{"onboarding"}, wantCmd: "onboard", wantAfter: 0},
		{name: "version alias", args: []string{"ver"}, wantCmd: "version", wantAfter: 0},
		{name: "quick start phrase", args: []string{"Quick", "start"}, wantCmd: "quickstart", wantAfter: 0},
		{name: "web keeps flags", args: []string{"web", "--hostname", "0.0.0.0"}, wantCmd: "web", wantAfter: 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCmd, gotAfter := resolveCommand(tt.args)
			if gotCmd != tt.wantCmd {
				t.Fatalf("expected command %q, got %q", tt.wantCmd, gotCmd)
			}
			if len(gotAfter) != tt.wantAfter {
				t.Fatalf("expected %d remaining args, got %d", tt.wantAfter, len(gotAfter))
			}
		})
	}
}

func TestSuggestCommand(t *testing.T) {
	if got := suggestCommand([]string{"Highlights"}); got != "highlights" {
		t.Fatalf("expected highlights suggestion, got %q", got)
	}
	if got := suggestCommand([]string{"Quick", "start"}); got != "quickstart" {
		t.Fatalf("expected quickstart suggestion, got %q", got)
	}
	if got := suggestCommand([]string{"version"}); got != "version" {
		t.Fatalf("expected version suggestion, got %q", got)
	}
	if got := suggestCommand([]string{"nonsense"}); got != "" {
		t.Fatalf("expected empty suggestion, got %q", got)
	}
}

func TestUsageText_IncludesNewCommands(t *testing.T) {
	text := usageText()
	expected := []string{
		"openppl version",
		"openppl highlights",
		"openppl quickstart",
		"openppl examples",
		"openppl guide",
	}

	for _, item := range expected {
		if !strings.Contains(text, item) {
			t.Fatalf("usage text missing %q", item)
		}
	}
}

func TestVersionText(t *testing.T) {
	original := appVersion
	t.Cleanup(func() { appVersion = original })

	appVersion = "v0.1.11"
	if got := versionText(); got != "openppl version v0.1.11" {
		t.Fatalf("unexpected version text: %q", got)
	}
}
