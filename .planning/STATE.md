# STATE: PPL Study Planner TUI

**Last updated:** 2026-02-25

---

## Current Position

| Field | Value |
|-------|-------|
| **Phase** | 04-polish |
| **Plan** | 01 |
| **Total Plans** | 3 |
| **Status** | Ready |
| **Last Completed** | 03-google-integration-03 (OpenCode bot exporter + Study view Google wiring) |

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

**Last session:** 2026-02-25T14:33:29Z

**Next action:** Plan and execute 04-polish-01 (UI refinements and edge-case handling)

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

---

*State managed by GSD roadmapper*
| Phase 01-foundation Complete | 246min | 15 tasks | 23 files |
| Phase 03-google-integration P01 | 1 min | 3 tasks | 9 files |
| Phase 03-google-integration Complete | 35min | 7 tasks | 16 files |
