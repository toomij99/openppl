package motd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime/debug"
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
	case "recall", "quiz":
		return runRecall(stdin, stdout)
	case "progress":
		return runProgress(stdout)
	case "weak":
		return runWeakAreas(stdout)
	case "install":
		return runInstall(stdout)
	case "version":
		fmt.Fprintln(stdout, "openppl motd — ACS daily quiz")
		return 0
	default:
		fmt.Fprintln(stdout, "usage: openppl motd [display|recall|quiz|progress|weak|install|version]")
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

	objective := strings.TrimSpace(entry.Text)
	if objective == "" || strings.EqualFold(objective, "[Archived]") {
		objective = "This ACS item is marked as archived in the dataset. Review current FAA ACS guidance and explain what changed from the previous standard."
	}

	fmt.Fprintf(stdout, "\n=== ACS Daily Quiz Prep ===\n")
	fmt.Fprintf(stdout, "Code:      %s\n", entry.Code)
	fmt.Fprintf(stdout, "Task:      %s\n", entry.Title)
	fmt.Fprintf(stdout, "Section:   %s\n", sectionName(entry.Section))
	fmt.Fprintf(stdout, "Category:  %s\n", entry.Category)
	fmt.Fprintf(stdout, "Objective: %s\n", objective)
	fmt.Fprintf(stdout, "\nStudy Tip: %s\n", studyTip(entry.Section))
	fmt.Fprintf(stdout, "Insight:   %s\n\n", studyInsight(entry.Section, entry.Category))

	current := currentVersionTag()
	if current != "" {
		if latest, err := services.FetchLatestReleaseTag(600 * time.Millisecond); err == nil {
			if msg, ok := services.BuildUpdateRecommendation(current, latest); ok {
				fmt.Fprintf(stdout, "%s\n", msg)
			}
		}
	}
	fmt.Fprintf(stdout, "Run: openppl motd quiz  (or wait for login prompt)\n\n")
	return 0
}

func currentVersionTag() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return ""
	}
	v := strings.TrimSpace(info.Main.Version)
	if v == "" || v == "(devel)" {
		return ""
	}
	return v
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

// runRecall prompts the user with today's multiple-choice ACS quiz and stores
// the result in the per-user motd_answers.db. If stdin is not a terminal
// (e.g. SCP, rsync, piped), it returns 0 silently. Never blocks login.
func runRecall(stdin io.Reader, stdout io.Writer) int {
	// Only run in interactive terminals — skip SCP/rsync/piped sessions.
	f, ok := stdin.(*os.File)
	if !ok || !isatty.IsTerminal(f.Fd()) {
		return 0
	}

	now := time.Now()
	quiz, err := services.BuildDailyQuiz(now)
	if err != nil {
		// Silent failure — never block login.
		return 0
	}

	fmt.Fprintf(stdout, "Daily checkride quiz — ACS %s\n", quiz.Entry.Code)
	fmt.Fprintln(stdout, quiz.Prompt)
	for _, option := range quiz.Options {
		fmt.Fprintf(stdout, "  %s) %s\n", option.Label, option.Text)
	}
	fmt.Fprint(stdout, "Choose A/B/C/D (Enter to skip): ")

	line, err := bufio.NewReader(stdin).ReadString('\n')
	if err != nil {
		// EOF or read error — treat as skip.
		return 0
	}
	choice := services.NormalizeQuizChoice(line)
	skipped := strings.TrimSpace(line) == "" || choice == ""

	db, err := services.InitMOTDDB()
	if err != nil {
		// Silent failure — never block login.
		return 0
	}

	if err := services.SaveMOTDAttempt(db, now.Format("2006-01-02"), quiz, choice, skipped); err != nil {
		// Silent failure — never block login.
		return 0
	}

	if skipped {
		fmt.Fprintf(stdout, "Skipped. Correct answer: %s\n", quiz.CorrectLabel)
		fmt.Fprintf(stdout, "%s\n\n", quiz.Explanation)
		return 0
	}

	if services.IsCorrectQuizChoice(quiz, choice) {
		fmt.Fprintln(stdout, "Correct.")
	} else {
		fmt.Fprintf(stdout, "Not quite. Correct answer: %s\n", quiz.CorrectLabel)
	}
	fmt.Fprintf(stdout, "%s\n\n", quiz.Explanation)
	return 0
}

func runProgress(stdout io.Writer) int {
	db, err := services.InitMOTDDB()
	if err != nil {
		fmt.Fprintf(stdout, "Could not open MOTD progress DB: %v\n", err)
		return 1
	}

	attempts, err := services.LoadMOTDAttempts(db)
	if err != nil {
		fmt.Fprintf(stdout, "Could not load MOTD attempts: %v\n", err)
		return 1
	}

	stats := services.ComputeMOTDReadiness(attempts, time.Now())
	fmt.Fprintln(stdout, "MOTD Readiness")
	fmt.Fprintf(stdout, "Score: %.1f/100 (%s)\n", stats.ReadinessScore, stats.ReadinessLabel)
	fmt.Fprintf(stdout, "Overall accuracy: %.1f%% (%d/%d)\n", stats.OverallAccuracy, stats.CorrectAttempts, stats.AnsweredAttempts)
	fmt.Fprintf(stdout, "Last 14 days: %.1f%%\n", stats.Last14Accuracy)
	fmt.Fprintf(stdout, "Coverage: %.1f%% (%d ACS areas tracked)\n", stats.CoverageScore, stats.TotalDistinctAreas)
	fmt.Fprintf(stdout, "Attempts: %d answered, %d skipped\n", stats.AnsweredAttempts, stats.SkippedAttempts)

	if len(stats.Areas) > 0 {
		fmt.Fprintln(stdout, "\nBy area:")
		for _, area := range stats.Areas {
			fmt.Fprintf(stdout, "  Area %s: %.1f%% (%d/%d)\n", area.Area, area.Accuracy, area.Correct, area.Attempts)
		}
	}

	return 0
}

func runWeakAreas(stdout io.Writer) int {
	db, err := services.InitMOTDDB()
	if err != nil {
		fmt.Fprintf(stdout, "Could not open MOTD progress DB: %v\n", err)
		return 1
	}

	attempts, err := services.LoadMOTDAttempts(db)
	if err != nil {
		fmt.Fprintf(stdout, "Could not load MOTD attempts: %v\n", err)
		return 1
	}

	stats := services.ComputeMOTDReadiness(attempts, time.Now())
	if len(stats.WeakAreas) == 0 {
		fmt.Fprintln(stdout, "No answered quiz data yet. Run: openppl motd quiz")
		return 0
	}

	fmt.Fprintln(stdout, "Weak ACS areas (lowest accuracy first):")
	for _, area := range stats.WeakAreas {
		fmt.Fprintf(stdout, "  Area %s: %.1f%% (%d/%d)\n", area.Area, area.Accuracy, area.Correct, area.Attempts)
	}
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
		return "Skill"
	default:
		return s
	}
}

func studyTip(section string) string {
	switch section {
	case "K":
		return "Use a 60-second teach-back: define the concept, state the FAA standard, then give one practical cockpit example."
	case "R":
		return "Use PAVE + 3P: identify one hazard, one mitigation, and one clear go/no-go trigger you would brief before flight."
	case "S":
		return "Chair-fly the maneuver aloud with callouts: setup, execution cues, tolerances, and recovery gates."
	default:
		return "State the requirement in your own words and tie it to one real preflight or in-flight decision."
	}
}

func studyInsight(section string, category string) string {
	base := "Checkride answers score higher when you connect standards to decisions, not just definitions."
	switch {
	case section == "R" && category == "Theory":
		return "Risk items are strongest when you verbalize a specific trigger point (weather, fuel, or workload) and your exact mitigation."
	case section == "S":
		return "For skill items, evaluators look for stable setup discipline first; most losses happen before the maneuver begins."
	default:
		return base
	}
}
