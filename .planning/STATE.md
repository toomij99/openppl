# STATE: PPL Study Planner TUI

**Last updated:** 2026-02-25

---

## Current Position

| Field | Value |
|-------|-------|
| **Phase** | 04-polish |
| **Plan** | 03 |
| **Total Plans** | 3 |
| **Status** | Complete |
| **Last Completed** | 04-polish-02 (Loading indicators + async operation lifecycle UX) |

---

## Roadmap Summary

| Phase | Goal | Requirements |
|-------|------|--------------|
| 1 - Foundation | Core TUI with study planning, tracking, dashboard, checklist | 16 |
| 2 - Calendar Export | ICS export + Apple Reminders | 4 |
| 3 - Google Integration | Google Calendar + OpenCode bot | 3 |
| 4 - Polish | UI refinements | 0 |

**Total:** 4 phases, 23 requirements

---

## Accumulated Context

### Phase 1 Decisions (from discuss-phase)

**Navigation:**
- Arrow keys for movement, header+footer for view indication
- Instant view switching, help overlay (? or F1)

**Data Entry:**
- Form fields, ISO dates (YYYY-MM-DD), formatted currency ($1,500)
- Auto-save to storage

**Dashboard:**
- Progress bar + % with days until checkride
- 7-day lookahead, today's tasks highlighted
- 4 quick stats (completed, remaining, overdue, total)

**Checklist:**
- 4 FAA categories (Documents, Aircraft, Ground, Flight)
- Pre-populated with FAA requirements
- Space to toggle, per-category + overall %

### Decisions Made

- **Phase structure:** Derived from requirements and research recommendations
- **Depth setting:** Quick (3-5 phases) applied to compress into 4 natural phases
- **Phase 1 scope:** All study planning, tracking, dashboard, and checklist features (16 reqs) - largest phase but cohesive deliverable
- [Phase 03-google-integration]: Implemented terminal auth-code flow (browser URL + pasted code) for Google Calendar-compatible terminal OAuth.
- [Phase 03-google-integration]: Persist Google OAuth tokens at data/google/token.json with strict 0600 file permissions and typed auth errors.
- [Phase 03-google-integration]: Standardized exporter artifact paths through ResolveArtifactOutputDir so outputs remain under repo icss/.
- [Phase 04-polish]: Translate service-layer typed errors at the Study view boundary to avoid leaking raw internal errors
- [Phase 04-polish]: Treat invalid or incomplete date input as warning state and keep input mode active until corrected
- [Phase 04-polish]: Centralized keyboard shortcuts in internal/tui/shortcuts.go now power both footer hints and help overlay sections.
- [Phase 04-polish]: MainModel handles ?/F1 globally with modal help visibility so discoverability works consistently across screens.
- [Phase 04-polish]: Use explicit Study operation state (label + loading flag) while keeping severity/message in shared studyStatus.
- [Phase 04-polish]: Block e/r/g/o while async operation is active and show warning status instead of launching duplicate commands.

### Research Context

Key research findings incorporated:
- Bubble Tea (Go) selected as TUI framework
- SQLite for data storage
- ICS-first export with UTC timezone
- Apple Reminders via osascript (batched)
- Google OAuth2 device flow for terminal

### Known Gaps

- Backward planning algorithm needs user testing to tune time estimates
- No actual PPL student usage validation yet

---

## Session Continuity

**Last session:** 2026-02-25T20:41:06.045Z

**Next action:** Transition to Phase 5 planning/execution

---

## Decisions

- Added dedicated service-layer exporters for both ICS and Apple Reminders in Phase 2.
- Kept Phase 2 scoped to local ICS + Apple Reminders and deferred Google API integration to Phase 3.
- Added Google Calendar sync adapter with deterministic identity mapping and bounded retry handling.
- Added OpenCode bot `v1` exporter and Study view key actions for Google sync (`g`) and bot export (`o`).

---

## Performance Metrics

| Phase | Plan | Duration | Tasks | Files |
|-------|------|----------|-------|-------|
| 01-foundation | 01 | 35min | 3 | 5 |
| 01-foundation | 02 | 58min | 3 | 6 |
| 01-foundation | 03 | 74min | 4 | 7 |
| 01-foundation | 04 | 4min | 2 | 2 |
| 01-foundation | 05 | 1min | 3 | 3 |
| 03-google-integration | 02 | 8min | 2 | 3 |
| 03-google-integration | 03 | 6min | 2 | 3 |
| 04-polish | 01 | 3 min | 3 | 3 |
| 04-polish | 02 | 1 min | 3 | 2 |

---

*State managed by GSD roadmapper*
