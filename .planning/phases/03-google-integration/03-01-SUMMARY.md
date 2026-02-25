---
phase: 03-google-integration
plan: 01
subsystem: api
tags: [google-oauth, google-calendar, oauth2, export-paths]

# Dependency graph
requires:
  - phase: 02-calendar-export
    provides: Shared exporter service boundaries for ICS and Reminders
provides:
  - Terminal-friendly Google OAuth client bootstrap with persisted token cache
  - Typed Google auth error mapping for credential, token, and exchange failures
  - Shared artifact output path enforcement under repo-level icss/
affects: [google-calendar-sync, study-view-integration, opencode-bot-export]

# Tech tracking
tech-stack:
  added: [golang.org/x/oauth2, google.golang.org/api/calendar/v3]
  patterns: [service-layer auth adapter, strict token cache permissions, artifact-path normalization]

key-files:
  created: [internal/services/google_auth.go, internal/services/google_auth_test.go]
  modified: [go.mod, go.sum, internal/services/export_paths.go, internal/services/export_ics.go, internal/services/export_reminders.go]

key-decisions:
  - "Implement terminal auth-code flow (browser URL + pasted code) as GREM-02-compatible terminal OAuth path."
  - "Persist Google tokens at data/google/token.json with 0600 file mode and fail fast on loose permissions."
  - "Keep all source files in root module and enforce icss/ as artifact-only output root."

patterns-established:
  - "Service adapters own external integration concerns and return typed errors to UI callers."
  - "Export artifact paths are normalized through a single resolver shared by all exporters."

requirements-completed: [GREM-02]

# Metrics
duration: 1 min
completed: 2026-02-25
---

# Phase 03 Plan 01: Google Auth Foundation Summary

**Google terminal OAuth bootstrap now caches reusable tokens and all export artifacts are normalized under repo-level `icss/`.**

## Performance

- **Duration:** 1 min
- **Started:** 2026-02-25T14:27:26Z
- **Completed:** 2026-02-25T14:27:58Z
- **Tasks:** 3
- **Files modified:** 9

## Accomplishments
- Added shared `ResolveArtifactOutputDir` usage across ICS and Reminders exporters so outputs stay under `<repo>/icss`.
- Implemented `EnsureGoogleAuthClient` with credentials loading, terminal consent instructions, token exchange, and durable cache reuse.
- Added lifecycle tests covering first-run auth, cached token reuse, strict `0600` enforcement, and typed error mapping.

## Task Commits

Each task was committed atomically:

1. **Task 1: Add shared artifact-path normalization enforcing output-only icss usage** - `5bf188d` (feat)
2. **Task 2: Implement Google OAuth terminal auth service with durable token cache** - `9f2fa39` (feat)
3. **Task 3: Add auth lifecycle tests for token caching and error boundaries** - `0fce2b7` (test)

**Plan metadata:** Pending

## Files Created/Modified
- `internal/services/export_paths.go` - Canonical artifact-output resolver pinned to `<repo>/icss`.
- `internal/services/export_ics.go` - Uses shared resolver for ICS output path handling.
- `internal/services/export_reminders.go` - Uses shared resolver for export artifact directory contract.
- `internal/services/google_auth.go` - Terminal OAuth service, token cache read/write, typed auth errors.
- `internal/services/google_auth_test.go` - Coverage for auth lifecycle, cache modes, and error kinds.
- `go.mod` - Adds OAuth and Google Calendar API dependencies.
- `go.sum` - Locks transitive dependencies for new Google integration stack.

## Decisions Made
- Used terminal auth-code flow (browser URL + pasted code) to satisfy terminal interaction while keeping Calendar scope compatibility.
- Enforced token cache file mode to `0600` and treat looser permissions as a typed failure for security correctness.
- Reused existing service-layer integration pattern to keep auth concerns outside Bubble Tea view logic.

## Deviations from Plan

None - plan executed exactly as written.

## Authentication Gates

None.

## Issues Encountered

None.

## User Setup Required

External services require manual configuration in Google Cloud Console:
- Set `GOOGLE_OAUTH_CREDENTIALS_PATH` to downloaded desktop OAuth credential JSON.
- Enable Google Calendar API on the target Google Cloud project.

## Next Phase Readiness

- Ready for `03-02-PLAN.md` (Google Calendar writer integration) with auth and output-path foundations in place.
- No blockers recorded for continuing phase 3.

## Self-Check: PASSED
