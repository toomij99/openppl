# PPL Study Planner TUI

## What This Is

A terminal user interface (TUI) application that helps plan and track Private Pilot License (PPL) study progress according to FAA ACS standards. Generates ICS calendars and Apple/Google Reminders for daily study tasks, flight training, and checkride preparation.

## Core Value

Systematically prepare for PPL checkride by breaking down FAA Private Airplane ACS into daily actionable tasks with progress tracking and automated calendar/reminder integration.

## Requirements

### Validated

(None yet — ship to validate)

### Active

- [ ] TUI app with daily study plan view
- [ ] ACS-based task breakdown (Areas, Tasks, Objectives)
- [ ] Checkride date forecast input with backward planning
- [ ] Daily plan categories: Theory, Chair Flying, Garmin 430, CFI Flights
- [ ] ICS export for Apple Calendar and Google Calendar
- [ ] Apple Reminders integration via osascript
- [ ] Google Tasks/Calendar API integration
- [ ] Progress dashboard with checkride readiness %
- [ ] Checkride checklist with requirement tracking

### Out of Scope

- [Mobile app] — TUI first, mobile later if needed
- [IFR/CPL] — PPL only for v1
- [Offline-first] — Requires internet for Google sync

## Context

**Technical Environment:**
- TUI framework: Bubble Tea (Go) or Textual (Python)
- Calendar formats: ICS standard
- Integrations: Apple Reminders (osascript), Google Calendar API
- Data storage: SQLite or JSON files

**User Profile:**
- Aviation student preparing for PPL checkride
- Uses Mac, Apple Reminders, Google Calendar
- OpenCode bot for task reminders

**Prior Work:**
- Existing import-reminders scripts (bash/Python)
- CSV-based study plans in csvs/ directory

## Constraints

- **[Platform]**: macOS primary (Apple Reminders, osascript)
- **[Calendar]**: Must support both Apple and Google Calendar export
- **[Input]**: ACS PDF reference (will parse/ingest)
- **[Reminders]**: OpenCode bot integration for automated checks

## Key Decisions

| Decision | Rationale | Outcome |
|----------|-----------|---------|
| TUI over web | User preference, faster to build, terminal-native | — Pending |
| Bubble Tea (Go) vs Textual (Python) | Need to decide based on library availability | — Pending |
| ICS-first export | Universal format, works with Apple + Google | — Pending |

---
*Last updated: 2025-12-20 after initialization*
