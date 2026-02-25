# Phase 1: Foundation - Context

**Gathered:** 2025-12-20
**Status:** Ready for planning

<domain>
## Phase Boundary

Core TUI app with study planning, progress tracking, dashboard, budget tracking, and checkride checklist. Users can:
- Enter checkride date → backward-scheduled study plan
- Track daily task completion with progress %
- View dashboard with stats and upcoming tasks
- Manage budget (flight costs, living costs)
- Complete checkride checklist (Documents, Aircraft, Ground, Flight)

</domain>

<decisions>
## Implementation Decisions

### Navigation Pattern
- Arrow keys for moving up/down within views
- Header shows current view name
- Footer shows available views and keyboard hints
- Instant switch between views (no confirmation)
- Help overlay (? or F1) to discover shortcuts

### Data Entry Flow
- Form fields for entering checkride date, flight hours, budget info
- ISO date format (YYYY-MM-DD, e.g., 2026-03-15)
- Currency displayed with formatting ($1,500)
- Auto-save changes immediately to storage

### Dashboard Display
- Main progress shown as progress bar + percentage
- Days until checkride prominently displayed
- Shows upcoming 7 days of tasks
- Today's tasks highlighted at top
- 4 quick stats below: completed, remaining, overdue, total

### Checklist Structure
- 4 categories: Documents, Aircraft, Ground, Flight (FAA ACS structure)
- Pre-populated with FAA requirements
- Space key to toggle item completion
- Progress displayed per-category and overall %

### Claude's Discretion
- Exact layout dimensions and spacing
- Color scheme (stick to terminal default or minimal colors)
- Progress bar visual style
- Help overlay exact content and styling

</decisions>

<specifics>
## Specific Ideas

- "I want it like standard TUI apps — arrow keys, space to select"
- ISO dates for unambiguous date entry
- 7-day lookahead for task planning
- FAA ACS checklist categories (Documents, Aircraft, Ground, Flight)

</specifics>

<deferred>
## Deferred Ideas

None — discussion stayed within phase scope

</deferred>

---

*Phase: 01-foundation*
*Context gathered: 2025-12-20*
