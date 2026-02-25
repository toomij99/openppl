---
phase: 04-polish
plan: 03
subsystem: ui
tags: [tui, keyboard-shortcuts, help-overlay, bubbletea]

# Dependency graph
requires:
  - phase: 03-google-integration
    provides: Study view export/sync key actions (`e`, `r`, `g`, `o`) used by help docs
provides:
  - Centralized shortcut registry for footer and full help content
  - Global `?`/`F1` help overlay toggle from any screen
  - Regression tests for shortcut registry and help rendering contracts
affects: [internal/tui, keyboard discoverability, phase-04 polish UX]

# Tech tracking
tech-stack:
  added: []
  patterns:
    - Single-source shortcut metadata powering both compact and full UI docs
    - App-level modal help toggle handled before view-level key routing

key-files:
  created:
    - internal/tui/shortcuts.go
    - internal/tui/help_overlay_test.go
  modified:
    - internal/tui/tui.go

key-decisions:
  - "Model shortcuts as typed data in `internal/tui/shortcuts.go` so footer/help text cannot drift from each other."
  - "Handle help visibility globally in `MainModel.Update` so `?` and `F1` work consistently across all screens."

patterns-established:
  - "Shortcut Registry: global and view actions live in one `Shortcut` source, then renderers derive footer/help variants."
  - "Help Overlay Contract: overlay toggles with `?`/`F1` and blocks non-exit navigation while open."

requirements-completed: []

# Metrics
duration: 4min
completed: 2026-02-25
---

# Phase 4 Plan 3: Global help command + centralized shortcut documentation Summary

**Global keyboard help now ships as a `?`/`F1` overlay backed by a centralized shortcut registry, with tests locking footer and help content alignment.**

## Performance

- **Duration:** 4 min
- **Started:** 2026-02-25T15:12:48Z
- **Completed:** 2026-02-25T15:16:40Z
- **Tasks:** 3
- **Files modified:** 3

## Accomplishments
- Added `Shortcut` registry helpers (`AllShortcuts`, `FooterShortcuts`, `HelpSections`) in one place.
- Added global help state and `?`/`F1` toggle handling in `MainModel`, with overlay rendering in `View()`.
- Added regression tests for shortcut coverage, help toggle behavior, and Study action presence in help output.

## Task Commits

Each task was committed atomically:

1. **Task 1: Create centralized shortcut registry for global and screen-specific keys** - `f6396dc` (feat)
2. **Task 2: Add global help toggle and overlay rendering in main TUI model** - `ce696cc` (feat)
3. **Task 3: Add regression tests for help toggle and shortcut coverage contract** - `badd804` (test)

**Plan metadata:** pending

## Files Created/Modified
- `internal/tui/shortcuts.go` - Defines centralized shortcut data and grouped helper accessors.
- `internal/tui/tui.go` - Adds global help visibility, key toggle handling, help overlay render, and footer hint rendering from registry.
- `internal/tui/help_overlay_test.go` - Adds regression tests for shortcut contracts and help overlay output.

## Decisions Made
- Use a shared shortcut data structure instead of duplicating strings between footer and help overlays.
- Keep help toggle at the TUI shell layer so every screen gets identical `?`/`F1` behavior.

## Deviations from Plan

None - plan executed exactly as written.

## Issues Encountered
None.

## User Setup Required

None - no external service configuration required.

## Next Phase Readiness
- Keyboard shortcut discoverability criterion is implemented and test-covered.
- Remaining Phase 4 plans (`04-01`, `04-02`) are still pending for full phase completion.

---
*Phase: 04-polish*
*Completed: 2026-02-25*

## Self-Check: PASSED

- Verified summary file exists on disk.
- Verified all task commit hashes (`f6396dc`, `ce696cc`, `badd804`) exist in git history.
