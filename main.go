package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"ppl-study-planner/internal/automation"
	"ppl-study-planner/internal/db"
	"ppl-study-planner/internal/motd"
	"ppl-study-planner/internal/onboarding"
	"ppl-study-planner/internal/tui"
	"ppl-study-planner/internal/web"
)

var (
	needsSetupCheck = needsSetup
	runOnboardingFn = runOnboarding
	initDatabaseFn  = db.Initialize
	runWebServerFn  = web.Run
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "automation" {
		os.Exit(runAutomationCommand(os.Args[2:], os.Stdout, os.Stderr))
	}

	defer func() {
		if recovered := recover(); recovered != nil {
			err := fmt.Errorf("unexpected panic: %v", recovered)
			handleFatalError(err, debug.Stack())
		}
	}()

	if err := run(); err != nil {
		handleFatalError(err, nil)
	}
}

func handleFatalError(err error, stack []byte) {
	logPath, logErr := writeErrorLog(err, stack)

	fmt.Fprintln(os.Stderr, "openppl encountered an error.")
	fmt.Fprintf(os.Stderr, "Reason: %v\n", err)

	if logErr == nil {
		fmt.Fprintf(os.Stderr, "Details saved to: %s\n", logPath)
	} else {
		fmt.Fprintf(os.Stderr, "Could not write error log: %v\n", logErr)
	}

	os.Exit(1)
}

func writeErrorLog(err error, stack []byte) (string, error) {
	if mkErr := os.MkdirAll("data", 0755); mkErr != nil {
		return "", mkErr
	}

	path := filepath.Join("data", "errors.log")
	f, openErr := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if openErr != nil {
		return "", openErr
	}
	defer f.Close()

	stamp := time.Now().Format(time.RFC3339)
	if _, wErr := fmt.Fprintf(f, "[%s] %v\n", stamp, err); wErr != nil {
		return "", wErr
	}
	if len(stack) > 0 {
		if _, wErr := fmt.Fprintf(f, "%s\n", string(stack)); wErr != nil {
			return "", wErr
		}
	}
	if _, wErr := fmt.Fprintln(f, "---"); wErr != nil {
		return "", wErr
	}

	return path, nil
}

func run() error {
	args := os.Args[1:]

	if len(args) > 0 {
		command, remaining := resolveCommand(args)
		switch command {
		case "onboard":
			return runOnboarding(false)
		case "web":
			return runWeb(remaining)
		case "motd":
			os.Exit(runMotdCommand(remaining))
			return nil
		case "configure":
			return runOnboarding(true)
		case "logs":
			return showLogs()
		case "automation":
			return fmt.Errorf("automation commands must be invoked as `openppl automation ...`")
		case "help":
			printUsage()
			return nil
		case "highlights":
			printHighlights()
			return nil
		case "quickstart":
			printQuickStart()
			return nil
		case "examples":
			printExamples()
			return nil
		case "guide":
			printGuide()
			return nil
		default:
			if suggestion := suggestCommand(args); suggestion != "" {
				return fmt.Errorf("unknown command %q\nDid you mean `openppl %s`?\n\n%s", strings.Join(args, " "), suggestion, usageText())
			}
			return fmt.Errorf("unknown command %q\n\n%s", strings.Join(args, " "), usageText())
		}
	}

	needs, err := needsSetup()
	if err != nil {
		return err
	}
	if needs {
		fmt.Println("First run detected. Starting onboarding setup...")
		if err := runOnboarding(true); err != nil {
			return err
		}
		fmt.Println("Launching openppl...")
	}

	if err := tui.Run(); err != nil {
		return err
	}

	return nil
}

func needsSetup() (bool, error) {
	database, err := initDatabaseFn()
	if err != nil {
		return false, fmt.Errorf("failed to initialize database: %w", err)
	}

	needed, err := onboarding.NeedsOnboarding(database)
	if err != nil {
		return false, fmt.Errorf("failed to check onboarding status: %w", err)
	}
	return needed, nil
}

func runOnboarding(force bool) error {
	database, err := initDatabaseFn()
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}

	if err := onboarding.RunInteractiveSetup(database, os.Stdin, os.Stdout, force); err != nil {
		return fmt.Errorf("onboarding failed: %w", err)
	}

	return nil
}

func printUsage() {
	fmt.Print(usageText())
}

func printHighlights() {
	fmt.Println("openppl highlights:")
	fmt.Println("- Interactive terminal (TUI) experience for study planning")
	fmt.Println("- Web dashboard mode")
	fmt.Println("- Automation commands for integrations")
	fmt.Println("- Onboarding and configuration flows")
}

func printQuickStart() {
	fmt.Println("openppl quick start:")
	fmt.Println("1) openppl help")
	fmt.Println("2) openppl")
	fmt.Println("3) openppl onboard")
	fmt.Println("4) openppl web --hostname 0.0.0.0 --port 5016")
	fmt.Println("5) openppl automation status")
}

func printExamples() {
	fmt.Println("openppl command examples:")
	fmt.Println("- openppl")
	fmt.Println("- openppl help")
	fmt.Println("- openppl onboard")
	fmt.Println("- openppl --configure")
	fmt.Println("- openppl web")
	fmt.Println("- openppl web --hostname 0.0.0.0 --port 5016")
	fmt.Println("- openppl automation status")
	fmt.Println("- openppl automation action --name remind --request-id req-001")
}

func printGuide() {
	fmt.Println("openppl guide:")
	fmt.Println("- Start app: openppl")
	fmt.Println("- Setup wizard: openppl onboard")
	fmt.Println("- Reconfigure: openppl --configure")
	fmt.Println("- Web UI: openppl web --hostname 0.0.0.0 --port 5016")
	fmt.Println("- Automation: openppl automation status")
	fmt.Println("- More examples: openppl examples")
}

func resolveCommand(args []string) (string, []string) {
	if len(args) == 0 {
		return "", nil
	}

	first := normalizeCommandToken(args[0])
	if first == "" {
		return "", args[1:]
	}

	if len(args) > 1 {
		combined := normalizeCommandToken(args[0] + args[1])
		if combined == "quickstart" {
			return "quickstart", args[2:]
		}
	}

	switch first {
	case "h", "help":
		return "help", args[1:]
	case "logs", "log":
		return "logs", args[1:]
	case "configure", "config":
		return "configure", args[1:]
	case "web", "dashboard":
		return "web", args[1:]
	case "motd":
		return "motd", args[1:]
	case "onboard", "onboarding":
		return "onboard", args[1:]
	case "automation", "auto":
		return "automation", args[1:]
	case "highlights", "highlight":
		return "highlights", args[1:]
	case "quickstart", "quick":
		return "quickstart", args[1:]
	case "examples", "example":
		return "examples", args[1:]
	case "guide":
		return "guide", args[1:]
	default:
		return "", args[1:]
	}
}

func suggestCommand(args []string) string {
	if len(args) == 0 {
		return ""
	}

	if len(args) > 1 {
		combined := normalizeCommandToken(args[0] + args[1])
		if combined == "quickstart" {
			return "quickstart"
		}
	}

	normalized := normalizeCommandToken(args[0])
	if normalized == "" {
		return ""
	}

	suggestions := map[string]string{
		"high":       "highlights",
		"highlight":  "highlights",
		"highlights": "highlights",
		"quick":      "quickstart",
		"quickstart": "quickstart",
		"start":      "guide",
		"guide":      "guide",
		"example":    "examples",
		"examples":   "examples",
		"web":        "web",
		"dashboard":  "web",
		"motd":       "motd",
		"onboarding": "onboard",
		"onboard":    "onboard",
		"auto":       "automation",
		"automation": "automation",
		"config":     "configure",
		"configure":  "configure",
		"help":       "help",
		"log":        "logs",
		"logs":       "logs",
	}

	if suggestion, ok := suggestions[normalized]; ok {
		return suggestion
	}

	return ""
}

func normalizeCommandToken(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	value = strings.TrimPrefix(value, "--")
	value = strings.TrimPrefix(value, "-")
	value = strings.ReplaceAll(value, "-", "")
	value = strings.ReplaceAll(value, "_", "")
	value = strings.ReplaceAll(value, " ", "")
	return value
}

func showLogs() error {
	path := filepath.Join("data", "errors.log")
	content, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("No error logs found yet.")
			return nil
		}
		return fmt.Errorf("failed to read logs: %w", err)
	}

	text := strings.TrimSpace(string(content))
	if text == "" {
		fmt.Println("Error log is empty.")
		return nil
	}

	entries := strings.Split(text, "\n---\n")
	start := 0
	if len(entries) > 10 {
		start = len(entries) - 10
	}

	fmt.Printf("Showing %d of %d error log entries from %s\n\n", len(entries[start:]), len(entries), path)
	for i := start; i < len(entries); i++ {
		fmt.Println(entries[i])
		if i != len(entries)-1 {
			fmt.Println("---")
		}
	}

	return nil
}

func usageText() string {
	return `openppl usage:
  openppl               Launch TUI (auto-runs onboarding on first run)
  openppl automation     Run non-interactive automation commands
  openppl automation status
  openppl automation action --name remind --request-id <id>
  openppl web           Launch web UI and open browser
  openppl web --hostname 0.0.0.0 --port 5016
  openppl onboard       Run onboarding setup wizard
  openppl --configure   Reconfigure core planning settings
  openppl highlights    Show product highlights in terminal
  openppl quickstart    Show copy-paste quick start commands
  openppl examples      Show common command examples
  openppl guide         Show command guide by workflow
  openppl logs          Show recent error logs (last 10 entries)
  openppl help          Show this help
`
}

func runWeb(args []string) error {
	fs := flag.NewFlagSet("web", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	hostname := fs.String("hostname", "127.0.0.1", "host interface to bind web server")
	port := fs.Int("port", 5016, "port to bind web server")

	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("invalid web command arguments: %w", err)
	}

	if *hostname == "" {
		return fmt.Errorf("invalid --hostname: must not be empty")
	}
	if *port < 1 || *port > 65535 {
		return fmt.Errorf("invalid --port %s: must be 1-65535", strconv.Itoa(*port))
	}

	needs, err := needsSetupCheck()
	if err != nil {
		return err
	}
	if needs {
		fmt.Println("First run detected. Starting onboarding setup...")
		if err := runOnboardingFn(true); err != nil {
			return err
		}
	}

	database, err := initDatabaseFn()
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}

	return runWebServerFn(database, *hostname, *port)
}

func runMotdCommand(args []string) int {
	// IMPORTANT: do NOT call initDatabaseFn() here.
	// The display subcommand (default) is database-free and must be fast.
	// The recall subcommand initializes its own DB via services.InitMOTDDB().
	return motd.Execute(args, os.Stdin, os.Stdout)
}

func runAutomationCommand(args []string, stdout io.Writer, stderr io.Writer) int {
	database, err := initDatabaseFn()
	if err != nil {
		fmt.Fprintf(stderr, "failed to initialize database: %v\n", err)
		return 1
	}

	return automation.Execute(database, args, stdout, stderr)
}
