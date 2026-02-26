package automation

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"strings"
	"time"

	"gorm.io/gorm"

	"ppl-study-planner/internal/services"
)

type CommandRunner func(req services.AutomationActionRequest) (services.AutomationActionResponse, error)

func Execute(database *gorm.DB, args []string, stdout io.Writer, stderr io.Writer) int {
	if stdout == nil {
		stdout = io.Discard
	}
	if stderr == nil {
		stderr = io.Discard
	}
	if database == nil {
		writeError(stderr, services.AutomationStatusResponse{
			Version:     services.AutomationVersionV1,
			ResultState: services.AutomationResultStateError,
			Timestamp:   time.Now().UTC().Format(time.RFC3339),
			Error:       &services.AutomationError{Code: "automation.db_required", Message: "database is required"},
		})
		return 2
	}

	if len(args) == 0 {
		writeError(stderr, services.AutomationStatusResponse{
			Version:     services.AutomationVersionV1,
			ResultState: services.AutomationResultStateError,
			Timestamp:   time.Now().UTC().Format(time.RFC3339),
			Error:       &services.AutomationError{Code: "automation.subcommand_required", Message: "expected `status` or `action`"},
		})
		return 2
	}

	subcommand := strings.ToLower(strings.TrimSpace(args[0]))
	switch subcommand {
	case "status":
		return runStatus(database, args[1:], stdout, stderr)
	case "action":
		return runAction(database, args[1:], stdout, stderr)
	default:
		writeError(stderr, services.AutomationStatusResponse{
			Version:     services.AutomationVersionV1,
			ResultState: services.AutomationResultStateError,
			Timestamp:   time.Now().UTC().Format(time.RFC3339),
			Error:       &services.AutomationError{Code: "automation.unknown_subcommand", Message: fmt.Sprintf("unknown subcommand %q", subcommand)},
		})
		return 2
	}
}

func runStatus(database *gorm.DB, args []string, stdout io.Writer, stderr io.Writer) int {
	if len(args) > 0 {
		writeError(stderr, services.AutomationStatusResponse{
			Version:     services.AutomationVersionV1,
			ResultState: services.AutomationResultStateError,
			Timestamp:   time.Now().UTC().Format(time.RFC3339),
			Error:       &services.AutomationError{Code: "status.invalid_arguments", Message: "status command does not accept extra arguments"},
		})
		return 2
	}

	response, err := services.BuildAutomationStatus(database, time.Now())
	if err != nil {
		statusErr := mapCommandError(err, services.AutomationResultStateError)
		writeError(stderr, statusErr)
		return 1
	}
	writeJSON(stdout, response)
	return 0
}

func runAction(database *gorm.DB, args []string, stdout io.Writer, stderr io.Writer) int {
	fs := flag.NewFlagSet("action", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	name := fs.String("name", "", "allowlisted action name")
	requestID := fs.String("request-id", "", "idempotency request identifier")
	actorScope := fs.String("actor-scope", "default", "actor scope for idempotency")

	if err := fs.Parse(args); err != nil {
		writeError(stderr, services.AutomationActionResponse{
			Version:     services.AutomationVersionV1,
			ResultState: services.AutomationResultStateRejected,
			Timestamp:   time.Now().UTC().Format(time.RFC3339),
			Error:       &services.AutomationError{Code: "action.invalid_flags", Message: err.Error()},
		})
		return 2
	}
	if fs.NArg() > 0 {
		writeError(stderr, services.AutomationActionResponse{
			Version:     services.AutomationVersionV1,
			ResultState: services.AutomationResultStateRejected,
			Timestamp:   time.Now().UTC().Format(time.RFC3339),
			Error:       &services.AutomationError{Code: "action.unexpected_arguments", Message: "unexpected positional arguments"},
		})
		return 2
	}

	response, err := services.RunAutomationAction(database, services.AutomationActionRequest{
		Name:       strings.TrimSpace(*name),
		RequestID:  strings.TrimSpace(*requestID),
		ActorScope: strings.TrimSpace(*actorScope),
		Args:       map[string]string{},
	})
	if err != nil {
		actionErr := mapCommandError(err, services.AutomationResultStateRejected)
		writeError(stderr, actionErr)
		return 1
	}

	writeJSON(stdout, response)
	return 0
}

func mapCommandError(err error, fallbackResult string) services.AutomationActionResponse {
	var commandErr *services.AutomationCommandError
	if errors.As(err, &commandErr) {
		resultState := fallbackResult
		if commandErr.Kind == "validation" {
			resultState = services.AutomationResultStateRejected
		}
		return services.AutomationActionResponse{
			Version:     services.AutomationVersionV1,
			ResultState: resultState,
			Timestamp:   time.Now().UTC().Format(time.RFC3339),
			Error:       &services.AutomationError{Code: commandErr.Code, Message: commandErr.Error()},
		}
	}

	return services.AutomationActionResponse{
		Version:     services.AutomationVersionV1,
		ResultState: fallbackResult,
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
		Error:       &services.AutomationError{Code: "automation.unknown_error", Message: err.Error()},
	}
}

func writeJSON(writer io.Writer, payload any) {
	encoder := json.NewEncoder(writer)
	encoder.SetEscapeHTML(false)
	_ = encoder.Encode(payload)
}

func writeError(stderr io.Writer, payload any) {
	writeJSON(stderr, payload)
}
