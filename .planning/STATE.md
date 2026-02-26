# STATE: PPL Study Planner TUI

**Last updated:** 2026-02-26

---

## Current Position

| Field | Value |
|-------|-------|
| **Phase** | 04-polish |
| **Plan** | 04 |
| **Total Plans** | 4 |
| **Status** | Complete |
| **Last Completed** | 04-polish-04 (shadcn web dashboard baseline + required dark/light themes) |

---

## Roadmap Summary

| Phase | Goal | Requirements |
|-------|------|--------------|
| 1 - Foundation | Core TUI with study planning, tracking, dashboard, checklist | 16 |
| 2 - Calendar Export | ICS export + Apple Reminders | 4 |
| 3 - Google Integration | Google Calendar + OpenCode bot | 3 |
| 4 - Polish | UI refinements + web dashboard baseline | 0 |

**Total:** 4 phases, 23 requirements

---

## Accumulated Context

### Roadmap Evolution

- Phase 04.1 inserted after Phase 4: Add web mode command (openppl web) + browser launch + host/port flags (URGENT)
- Phase 04.2 inserted after Phase 4: Add support of Openclaw. As user I would like to interact with openppl with telegram: ask Openclaw to show current status, send reminder, and other actions. (URGENT)

### Decisions Made

- **Phase structure:** Derived from requirements and research recommendations
- **Depth setting:** Quick (3-5 phases) applied to compress into 4 natural phases
- **Phase 1 scope:** All study planning, tracking, dashboard, and checklist features (16 reqs)
- [Phase 03-google-integration]: Implemented terminal auth-code flow for Google Calendar-compatible terminal OAuth.
- [Phase 03-google-integration]: Persist Google OAuth tokens at `data/google/token.json` with strict 0600 file permissions.
- [Phase 03-google-integration]: Standardized exporter artifact paths through `ResolveArtifactOutputDir` under repo `icss/`.
- [Phase 04-polish]: Translate typed service errors at Study view boundary to avoid leaking raw internals.
- [Phase 04-polish]: Keep invalid/incomplete date input in warning state until corrected.
- [Phase 04-polish]: Centralized shortcuts in `internal/tui/shortcuts.go` for footer hints + help overlay.
- [Phase 04-polish]: MainModel handles `?`/`F1` globally with modal help visibility.
- [Phase 04-polish]: Use explicit Study operation state while preserving shared `studyStatus` severity/message.
- [Phase 04-polish]: Block `e/r/g/o` during async operations and return warning feedback for duplicate triggers.
- [Phase 04-polish]: Bootstrap a dedicated `web/` Next.js workspace to isolate Go/TUI and web toolchains.
- [Phase 04-polish]: Ship light/dark token sets together, including chart tokens, via `next-themes` class mode.
- [Phase 04-polish]: Drive dashboard UX with explicit loading/empty/error page branches and typed module contracts.
- [Phase 04]: Bootstrap web dashboard in isolated web/ Next.js workspace to avoid coupling with Go/TUI build paths.

### Research Context

Key research findings incorporated:
- Bubble Tea (Go) selected as TUI framework
- SQLite for data storage
- ICS-first export with UTC timezone
- Apple Reminders via osascript (batched)
- Google OAuth2 device flow for terminal
- shadcn + Next.js App Router for web dashboard baseline with dual-theme support

### Known Gaps

- Backward planning algorithm needs user testing to tune time estimates
- No actual PPL student usage validation yet

---

## Session Continuity

**Last session:** 2026-02-26T08:21:39.666Z

**Next action:** Run `/gsd-plan-phase 04.1` for openppl web-mode command integration

---

## Decisions

- Added dedicated service-layer exporters for both ICS and Apple Reminders in Phase 2.
- Kept Phase 2 scoped to local ICS + Apple Reminders and deferred Google API integration to Phase 3.
- Added Google Calendar sync adapter with deterministic identity mapping and bounded retry handling.
- Added OpenCode bot `v1` exporter and Study view key actions for Google sync (`g`) and bot export (`o`).
- Added web dashboard shell under `web/` with shadcn-style composition and dual-theme controls.

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
| 04-polish | 04 | 5 min | 3 | 26 |

---

*State managed by GSD roadmapper*
