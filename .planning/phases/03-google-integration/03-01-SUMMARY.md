---
phase: 03-google-integration
plan: 01
subsystem: auth
tags: [google-oauth, oauth2, token-cache, exports, icss]

# Dependency graph
requires:
  - phase: 02-calendar-export
    provides: "ICS and Apple Reminders exporters with Study view wiring"
provides:
  - "Shared artifact path resolver that keeps outputs under icss"
  - "Terminal-friendly Google OAuth client bootstrap with durable token cache"
  - "Google auth lifecycle tests for credentials, cache, and permissions"
affects: [google-calendar-sync, study-view-google-actions, opencode-export]

# Tech tracking
tech-stack:
  added: [golang.org/x/oauth2, google.golang.org/api]
  patterns:
    - "Service-level output path normalization via ResolveArtifactOutputDir"
    - "Terminal OAuth URL + pasted code flow with typed error mapping"

key-files:
  created:
    - "internal/services/export_paths.go"
    - "internal/services/export_paths_test.go"
    - "internal/services/google_auth.go"
    - "internal/services/google_auth_test.go"
  modified:
    - "internal/services/export_ics.go"
    - "internal/services/export_ics_test.go"
    - "internal/services/export_reminders.go"
    - "go.mod"
    - "go.sum"

key-decisions:
  - "Normalized all artifact writes through one resolver to enforce icss as output-only"
  - "Used terminal auth-code flow with persistent token cache to satisfy terminal OAuth usability"
  - "Enforced token cache mode 0600 and explicit typed auth errors for UI-safe messaging"

patterns-established:
  - "Google auth service pattern: resolve credentials -> load cache -> fallback exchange -> persist"
  - "Exporter path policy pattern: ResolveArtifactOutputDir for all artifact outputs"

requirements-completed: [GREM-02]

# Metrics
duration: 21 min
completed: 2026-02-25
---

# Phase 3 Plan 1: Google OAuth Foundation + Export Path Normalization Summary

**Terminal Google OAuth auth with durable token reuse now works while all export artifacts are consistently constrained to the `icss/` output boundary.**

## Performance

- **Duration:** 21 min
- **Started:** 2026-02-25T14:12:24Z
- **Completed:** 2026-02-25T14:33:29Z
- **Tasks:** 3
- **Files modified:** 10

## Accomplishments

- Added `ResolveArtifactOutputDir` and refactored ICS/Reminders exporters to honor a single artifact-output contract.
- Implemented `EnsureGoogleAuthClient` with browser URL + pasted code terminal flow and durable token cache under `data/google/token.json`.
- Added focused auth lifecycle tests for first-run auth, cached token reuse, strict permission enforcement, and typed credential/token error mapping.

## Task Commits

Each task was committed atomically:

1. **Task 1: Add shared artifact-path normalization enforcing output-only icss usage** - `5bf188d` (feat)
2. **Task 2: Implement Google OAuth terminal auth service with durable token cache** - `9f2fa39` (feat)
3. **Task 3: Add auth lifecycle tests for token caching and error boundaries** - `0fce2b7` (test)

## Files Created/Modified

- `internal/services/export_paths.go` - shared artifact output resolver with traversal guards
- `internal/services/export_paths_test.go` - resolver behavior coverage for defaults, normalization, and rejection cases
- `internal/services/google_auth.go` - terminal OAuth config/token cache bootstrap with typed errors
- `internal/services/google_auth_test.go` - auth flow and token-cache regression coverage
- `internal/services/export_ics.go` - switched to shared resolver
- `internal/services/export_reminders.go` - aligned output contract handling

## Decisions Made

- Kept Go source ownership in root module while using `icss/` only for runtime artifacts.
- Used default Google Calendar scope and cache-backed client reuse to avoid repeated browser login prompts.
- Required strict token file mode (`0600`) before token reuse to protect local credentials.

## Deviations from Plan

### Auto-fixed Issues

**1. [Rule 1 - Bug] Updated ICS tests to use normalized artifact output paths**
- **Found during:** Task 1 (artifact path normalization)
- **Issue:** Existing tests wrote to temp directories outside `icss/` and failed under new resolver contract.
- **Fix:** Routed test output through `ResolveArtifactOutputDir` and cleaned test artifacts post-run.
- **Files modified:** `internal/services/export_ics_test.go`
- **Verification:** `go test ./internal/services -run 'Test(ICS|Reminders|ExportPaths)' -v`
- **Committed in:** `5bf188d`

---

**Total deviations:** 1 auto-fixed (1 bug)
**Impact on plan:** Fix was required to align test behavior with the new output policy; no scope creep.

## Issues Encountered

None.

## User Setup Required

External services require manual configuration. Set `GOOGLE_OAUTH_CREDENTIALS_PATH` to your downloaded desktop OAuth client JSON and ensure Google Calendar API is enabled in your Google Cloud project.

## Next Phase Readiness

Google auth and output path policy are in place, enabling resilient Google Calendar write adapter implementation in Plan 03-02.

## Self-Check: PASSED

---
*Phase: 03-google-integration*
*Completed: 2026-02-25*
