package motd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mattn/go-isatty"

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
		return runInstall(stdout)
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

// runInstall writes the two system integration files that make login-time MOTD
// display and recall work. It requires root — returns 1 with an error message
// if not running as root.
func runInstall(stdout io.Writer) int {
	if os.Getuid() != 0 {
		fmt.Fprint(stdout, "motd install requires root: run with sudo openppl motd install\n")
		return 1
	}

	exePath, err := os.Executable()
	if err != nil {
		fmt.Fprintf(stdout, "motd install: failed to determine binary path: %v\n", err)
		return 1
	}
	// Resolve any symlinks so the script contains the real binary path.
	exePath, err = filepath.EvalSymlinks(exePath)
	if err != nil {
		fmt.Fprintf(stdout, "motd install: failed to resolve binary symlinks: %v\n", err)
		return 1
	}

	// Write /etc/update-motd.d/99-openppl-acs (no extension — run-parts skips files with dots)
	motdScript := fmt.Sprintf(`#!/bin/sh
# ACS Code of the Day - openppl PPL study tool
OPENPPL_BIN="%s"
if [ -x "$OPENPPL_BIN" ]; then
    "$OPENPPL_BIN" motd
fi
`, exePath)
	if err := os.WriteFile("/etc/update-motd.d/99-openppl-acs", []byte(motdScript), 0755); err != nil {
		fmt.Fprintf(stdout, "motd install: failed to write MOTD script: %v\n", err)
		return 1
	}

	// Write /etc/profile.d/openppl-recall.sh (mode 0644)
	recallScript := fmt.Sprintf(`#!/bin/sh
# ACS recall prompt for openppl PPL study tool
# Only runs in interactive login shells - skips SCP/rsync sessions
case "$-" in
    *i*) ;;
    *)   return ;;
esac
# Allow users to opt out: export OPENPPL_MOTD_RECALL=0
[ "${OPENPPL_MOTD_RECALL:-1}" = "0" ] && return
OPENPPL_BIN="%s"
if [ -x "$OPENPPL_BIN" ]; then
    "$OPENPPL_BIN" motd recall
fi
`, exePath)
	if err := os.WriteFile("/etc/profile.d/openppl-recall.sh", []byte(recallScript), 0644); err != nil {
		fmt.Fprintf(stdout, "motd install: failed to write recall script: %v\n", err)
		return 1
	}

	// Check for PrintMotd yes in sshd_config — could cause double display.
	checkSSHDConfig(stdout)

	fmt.Fprintf(stdout, `Installed:
  /etc/update-motd.d/99-openppl-acs  (MOTD display, runs as root)
  /etc/profile.d/openppl-recall.sh   (recall prompt, interactive shells only)
Binary path: %s
Run 'sudo run-parts --test /etc/update-motd.d' to verify MOTD script.
`, exePath)
	return 0
}

// checkSSHDConfig reads /etc/ssh/sshd_config and warns if PrintMotd yes is set,
// since that would cause the ACS code to appear twice during SSH login.
func checkSSHDConfig(stdout io.Writer) {
	content, err := os.ReadFile("/etc/ssh/sshd_config")
	if err != nil {
		return
	}
	for _, line := range strings.Split(string(content), "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "#") {
			continue
		}
		lower := strings.ToLower(trimmed)
		if strings.Contains(lower, "printmotd") && strings.Contains(lower, "yes") {
			fmt.Fprint(stdout, "Warning: sshd PrintMotd yes detected in /etc/ssh/sshd_config — ACS code may appear twice. Set PrintMotd no to fix.\n")
			return
		}
	}
}

// runRecall prompts the user to describe today's ACS requirement and saves
// their answer to the per-user motd_answers.db. If stdin is not a terminal
// (e.g. SCP, rsync, piped), it returns 0 silently. Never blocks login.
func runRecall(stdin io.Reader, stdout io.Writer) int {
	// Only run in interactive terminals — skip SCP/rsync/piped sessions.
	f, ok := stdin.(*os.File)
	if !ok || !isatty.IsTerminal(f.Fd()) {
		return 0
	}

	entry, err := services.TodaysACSCode(time.Now())
	if err != nil {
		// Silent failure — never block login.
		return 0
	}

	fmt.Fprintf(stdout, "Recall check — ACS %s: %s\n", entry.Code, entry.Text)
	fmt.Fprint(stdout, "Describe this requirement (Enter to skip): ")

	line, err := bufio.NewReader(stdin).ReadString('\n')
	if err != nil {
		// EOF or read error — treat as skip.
		return 0
	}
	answer := strings.TrimSpace(line)
	if answer == "" {
		return 0
	}

	db, err := services.InitMOTDDB()
	if err != nil {
		// Silent failure — never block login.
		return 0
	}

	if err := services.SaveMOTDAnswer(db, entry.Code, answer); err != nil {
		// Silent failure — never block login.
		return 0
	}

	fmt.Fprint(stdout, "Answer saved.\n\n")
	return 0
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
