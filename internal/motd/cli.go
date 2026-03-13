package motd

import (
	"fmt"
	"io"
	"time"

	"ppl-study-planner/internal/services"
)

// Execute is the top-level dispatcher for all `openppl motd [subcommand]`
// invocations. It is called from main.go and returns a process exit code.
//
//   - args: the arguments after "motd" (e.g. []string{"display"} or nil)
//   - stdin: passed through for interactive subcommands (recall)
//   - stdout: all output is written here so callers can capture it in tests
func Execute(args []string, stdin io.Reader, stdout io.Writer) int {
	sub := ""
	if len(args) > 0 {
		sub = args[0]
	}

	switch sub {
	case "", "display":
		return runDisplay(stdout)
	case "recall":
		return runRecall(stdin, stdout)
	case "install":
		fmt.Fprintln(stdout, "install: not yet implemented")
		return 0
	case "version":
		fmt.Fprintln(stdout, "openppl motd — ACS Code of the Day")
		return 0
	default:
		fmt.Fprintln(stdout, "usage: openppl motd [display|recall|install|version]")
		return 1
	}
}

// runDisplay prints today's ACS code of the day to stdout.
// It never blocks or crashes — on any error it returns 0 silently so that
// a login MOTD script is not disrupted.
func runDisplay(stdout io.Writer) int {
	entry, err := services.TodaysACSCode(time.Now())
	if err != nil {
		// Silent failure — never block login.
		return 0
	}
	fmt.Fprintf(stdout, "\nACS Code of the Day: %s\n", entry.Code)
	fmt.Fprintf(stdout, "Task: %s — %s\n", entry.Title, sectionName(entry.Section))
	fmt.Fprintf(stdout, "%s\n\n", entry.Text)
	return 0
}

// runRecall is a stub for Plan 02. Currently it just delegates to runDisplay.
// Plan 02 will replace this with the full interactive recall implementation.
func runRecall(stdin io.Reader, stdout io.Writer) int {
	return runDisplay(stdout)
}

// sectionName maps ACS section codes to human-readable labels.
func sectionName(s string) string {
	switch s {
	case "K":
		return "Knowledge"
	case "R":
		return "Risk Management"
	case "S":
		return "Skills"
	default:
		return s
	}
}
