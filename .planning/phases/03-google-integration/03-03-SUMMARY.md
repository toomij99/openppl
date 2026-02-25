---
phase: 03-google-integration
plan: 03
subsystem: ui
tags: [opencode-bot, study-view, bubbletea, google-sync]

# Dependency graph
requires:
  - phase: 03-google-integration-01
    provides: "OAuth auth client and output path normalization"
  - phase: 03-google-integration-02
    provides: "Google Calendar sync adapter and structured sync results"
provides:
  - "Versioned OpenCode bot exporter with deterministic JSON schema"
  - "Study view keybindings and async commands for Google sync and bot export"
  - "In-view status messaging for sync/export outcomes"
affects: [phase-3-user-workflow, opencode-bot-integration]

# Tech tracking
tech-stack:
  added: []
  patterns:
    - "Golden schema contract tests for integration payload stability"
    - "Bubble Tea async command wiring for external sync/export actions"

key-files:
  created:
    - "internal/services/export_opencode_bot.go"
    - "internal/services/export_opencode_bot_test.go"
  modified:
    - "internal/view/study.go"

key-decisions:
  - "Defined explicit OpenCode bot payload version (`v1`) and locked structure with a golden contract test"
  - "Kept Study view integration boundary at service calls only, with no inline OAuth/API logic"
  - "Preserved existing phase-2 export keys while adding dedicated Google/OpenCode actions"

patterns-established:
  - "User-triggered Study view operations follow key -> async cmd -> done msg -> status banner pattern"

requirements-completed: [GREM-01, RBOT-01]

# Metrics
duration: 6 min
completed: 2026-02-25
---

# Phase 3 Plan 3: Study View Google + OpenCode Integration Summary

**Study view now supports non-blocking Google Calendar sync and deterministic OpenCode bot export, including clear operation status and versioned bot payload output.**

## Performance

- **Duration:** 6 min
- **Started:** 2026-02-25T14:27:00Z
- **Completed:** 2026-02-25T14:33:29Z
- **Tasks:** 2
- **Files modified:** 3

## Accomplishments

- Implemented `ExportOpenCodeBotTasks` with explicit `v1` payload shape, deterministic ordering, and normalized artifact output path handling.
- Added contract-focused tests including a golden-schema assertion to catch payload drift.
- Extended Study view key handling (`g`, `o`) and async status updates for Google sync and OpenCode export actions.

## Task Commits

Each task was committed atomically:

1. **Task 1: Implement versioned OpenCode bot exporter with explicit schema contract** - `243e7a4` (feat)
2. **Task 2: Wire Study view commands for Google sync and OpenCode bot export** - `4066cd8` (feat)

## Files Created/Modified

- `internal/services/export_opencode_bot.go` - deterministic versioned OpenCode bot payload writer
- `internal/services/export_opencode_bot_test.go` - schema and output contract verification
- `internal/view/study.go` - Google sync and OpenCode export commands with in-view status messaging

## Decisions Made

- Chose JSON `v1` payload with required task metadata and deterministic ordering to support stable bot parsing.
- Kept exporter output under `icss/` via shared resolver to preserve source/output boundary.
- Exposed separate keybindings and status messages so users can distinguish sync and export outcomes quickly.

## Deviations from Plan

None - plan executed exactly as written.

## Issues Encountered

None.

## User Setup Required

None - no additional manual setup required beyond Google credentials from Plan 03-01.

## Next Phase Readiness

Phase 3 end-user integration is complete: Study view can sync to Google and export OpenCode bot artifacts with deterministic contracts and observable status.

## Self-Check: PASSED

---
*Phase: 03-google-integration*
*Completed: 2026-02-25*
