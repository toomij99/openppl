---
phase: 02-calendar-export
plan: 02
subsystem: api
tags: [reminders, osascript, automation, timeout, testing]

# Dependency graph
requires:
  - phase: 01-foundation
    provides: "DailyTask persistence used as reminders source"
provides:
  - "Apple Reminders exporter adapter via osascript"
  - "Timeout and permission-aware error mapping for automation failures"
  - "Unit tests for command argument shape and error classification"
affects:
  - "Study view export wiring"
  - "Phase 2 verification"

# Tech tracking
tech-stack:
  added: []
  patterns:
    - "argv-safe osascript invocation using exec.CommandContext"
    - "Per-task export loop with bounded timeout"
    - "Typed export error categories (validation, timeout, permission, script_failure)"

key-files:
  created:
    - "internal/services/export_reminders.go"
    - "internal/services/export_reminders_test.go"
  modified: []

key-decisions:
  - "Reused existing repository script conventions (run argv + list auto-create) in Go adapter"
  - "Mapped known macOS automation denial signatures to permission-specific error kind"
  - "Used injectable command runner for deterministic tests without live Reminders dependency"

patterns-established:
  - "Exporter API pattern: `ExportAppleReminders(tasks, opts)` returns typed result + error"
  - "Tests validate argument shape and error mapping without executing osascript"

requirements-completed:
  - "AREM-01"

# Metrics
duration: 8 min
completed: 2026-02-25
---

# Phase 2 Plan 2: Apple Reminders Export Summary

**Apple Reminders export now runs through a timeout-safe osascript adapter with explicit permission and script failure reporting.**

## Performance

- **Duration:** 8 min
- **Started:** 2026-02-25T13:32:10Z
- **Completed:** 2026-02-25T13:40:36Z
- **Tasks:** 2
- **Files modified:** 2

## Accomplishments

- Added reminders exporter service that creates reminders from study tasks using argv-safe `osascript` calls.
- Added structured error categories for timeout, permission denial, and script failures.
- Added automated tests for command argument shape and error mapping behavior.

## Task Commits

Each task was committed atomically:

1. **Task 1: Create Apple Reminders exporter with timeout-safe osascript execution** - `c83cfac` (feat)
2. **Task 2: Add reminder exporter tests without requiring live Reminders access** - `35d248d` (test)

## Files Created/Modified

- `internal/services/export_reminders.go` - reminders exporter, command runner abstraction, and typed error mapping
- `internal/services/export_reminders_test.go` - tests for argv construction and timeout/permission mapping

## Decisions Made

1. Kept all command execution argv-based to avoid shell interpolation bugs.
2. Defaulted reminders list to `OpenPPL Study Tasks` when list name is not provided.
3. Added runner injection to keep tests fast and CI-safe without requiring macOS Reminders access.

## Deviations from Plan

None - plan executed exactly as written.

## Issues Encountered

None.

## User Setup Required

None - no external service configuration required.

## Next Phase Readiness

Reminders exporter is ready for Study view command wiring and user-facing export status messaging.

---
*Phase: 02-calendar-export*
*Completed: 2026-02-25*
