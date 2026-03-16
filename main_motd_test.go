package main

import (
	"bytes"
	"os"
	"regexp"
	"testing"

	"ppl-study-planner/internal/motd"
)

// TestMotdCommand_DisplayNoDB verifies that the motd display subcommand exits 0,
// prints the "ACS Daily Quiz Prep" header, and includes a recognizable ACS code
// pattern — without requiring the main app database to be initialized.
func TestMotdCommand_DisplayNoDB(t *testing.T) {
	var buf bytes.Buffer
	// Call motd.Execute directly so we can capture output;
	// runMotdCommand is hard-wired to os.Stdout.
	code := motd.Execute([]string{}, os.Stdin, &buf)
	if code != 0 {
		t.Errorf("motd display exit code = %d; want 0", code)
	}

	out := buf.String()
	if !regexp.MustCompile(`ACS Daily Quiz Prep`).MatchString(out) {
		t.Errorf("output missing 'ACS Daily Quiz Prep'; got:\n%s", out)
	}

	// ACS code pattern: two uppercase letters, Roman numeral section,
	// uppercase letter, section code letter, digit(s) — e.g. "PA.I.A.K1"
	acsPattern := regexp.MustCompile(`[A-Z]{2}\.[IVX]+\.[A-Z]\.[KRS]\d+`)
	if !acsPattern.MatchString(out) {
		t.Errorf("output does not contain ACS code pattern %s; got:\n%s", acsPattern, out)
	}
}

// TestMotdCommand_UnknownSubcommand verifies that an unknown subcommand exits 1.
func TestMotdCommand_UnknownSubcommand(t *testing.T) {
	code := runMotdCommand([]string{"nonexistent"})
	if code != 1 {
		t.Errorf("runMotdCommand(nonexistent) = %d; want 1", code)
	}
}
