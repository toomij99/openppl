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
		switch args[0] {
		case "onboard":
			return runOnboarding(false)
		case "web":
			return runWeb(args[1:])
		case "--configure", "configure":
			return runOnboarding(true)
		case "logs", "--logs":
			return showLogs()
		case "automation":
			return fmt.Errorf("automation commands must be invoked as `openppl automation ...`")
		case "--help", "-h", "help":
			printUsage()
			return nil
		default:
			return fmt.Errorf("unknown command %q\n\n%s", args[0], usageText())
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

func runAutomationCommand(args []string, stdout io.Writer, stderr io.Writer) int {
	database, err := initDatabaseFn()
	if err != nil {
		fmt.Fprintf(stderr, "failed to initialize database: %v\n", err)
		return 1
	}

	return automation.Execute(database, args, stdout, stderr)
}
