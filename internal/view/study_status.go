package view

import (
	"errors"
	"fmt"
	"strings"

	"ppl-study-planner/internal/services"
)

const (
	studyStatusSeverityInfo    = "info"
	studyStatusSeveritySuccess = "success"
	studyStatusSeverityWarning = "warning"
	studyStatusSeverityError   = "error"
)

type studyStatus struct {
	severity string
	message  string
}

func newStudyStatusInfo(message string) studyStatus {
	return studyStatus{severity: studyStatusSeverityInfo, message: strings.TrimSpace(message)}
}

func newStudyStatusSuccess(message string) studyStatus {
	return studyStatus{severity: studyStatusSeveritySuccess, message: strings.TrimSpace(message)}
}

func newStudyStatusWarning(message string) studyStatus {
	return studyStatus{severity: studyStatusSeverityWarning, message: strings.TrimSpace(message)}
}

func newStudyStatusError(message string) studyStatus {
	return studyStatus{severity: studyStatusSeverityError, message: strings.TrimSpace(message)}
}

func newStudyStatusFromError(operation string, err error) studyStatus {
	if err == nil {
		return newStudyStatusSuccess(fmt.Sprintf("%s complete", strings.TrimSpace(operation)))
	}

	var googleAuthErr *services.GoogleAuthError
	if errors.As(err, &googleAuthErr) {
		return newStudyStatusError(translateGoogleAuthError(googleAuthErr))
	}

	var googleCalendarErr *services.GoogleCalendarError
	if errors.As(err, &googleCalendarErr) {
		if googleCalendarErr.Kind == services.GoogleCalendarErrorAuth {
			if unwrapped := newStudyStatusFromError(operation, errors.Unwrap(googleCalendarErr)); unwrapped.message != "" {
				return unwrapped
			}
		}

		switch googleCalendarErr.Kind {
		case services.GoogleCalendarErrorValidation:
			return newStudyStatusError("Google sync request is invalid. Check task data and try again.")
		case services.GoogleCalendarErrorServiceInit:
			return newStudyStatusError("Google Calendar service setup failed. Verify configuration and retry.")
		default:
			return newStudyStatusError("Google sync failed. Check credentials and try again.")
		}
	}

	var remindersErr *services.RemindersExportError
	if errors.As(err, &remindersErr) {
		switch remindersErr.Kind {
		case "validation":
			return newStudyStatusError("Reminders export request is invalid. Refresh tasks and try again.")
		case "timeout":
			return newStudyStatusError("Reminders export timed out. Open Reminders and try again.")
		case "permission":
			return newStudyStatusError("Reminders access denied. Allow automation permissions and retry.")
		default:
			return newStudyStatusError("Reminders export failed. Check Reminders setup and try again.")
		}
	}

	return newStudyStatusError(fmt.Sprintf("%s failed. Check configuration and try again.", strings.TrimSpace(operation)))
}

func translateGoogleAuthError(err *services.GoogleAuthError) string {
	if err == nil {
		return "Google authentication failed. Re-authenticate and try again."
	}

	switch err.Kind {
	case services.GoogleAuthErrorMissingCredentials:
		return "Google credentials missing. Set GOOGLE_OAUTH_CREDENTIALS_PATH and try again."
	case services.GoogleAuthErrorCredentialsRead:
		return "Google credentials could not be read. Check file path and permissions."
	case services.GoogleAuthErrorInvalidCredentials:
		return "Google credentials are invalid. Download a valid OAuth client file and retry."
	case services.GoogleAuthErrorAuthExchange:
		return "Google authorization failed. Re-authenticate and try again."
	case services.GoogleAuthErrorTokenRead:
		return "Google token could not be read. Re-authenticate and retry."
	case services.GoogleAuthErrorInvalidToken:
		return "Google token is invalid. Remove it and authenticate again."
	case services.GoogleAuthErrorTokenPermissions:
		return "Google token permissions are invalid. Set token file permissions to 0600."
	case services.GoogleAuthErrorTokenSave:
		return "Google token could not be saved. Check write permissions and retry."
	default:
		return "Google authentication failed. Verify credentials and try again."
	}
}
