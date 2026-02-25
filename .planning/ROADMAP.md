# ROADMAP: PPL Study Planner TUI

**Project:** PPL Study Planner TUI  
**Created:** 2025-02-25  
**Depth:** Quick (3-5 phases)

---

## Phases

- [ ] **Phase 1: Foundation** - Core TUI app with study planning, progress tracking, dashboard, and checkride checklist
- [ ] **Phase 2: Calendar Export** - ICS file export and Apple Reminders integration
- [ ] **Phase 3: Google Integration** - Google Calendar API and OpenCode bot export
- [ ] **Phase 4: Polish** - UI refinements and edge case handling

---

## Phase Details

### Phase 1: Foundation

**Goal:** Users can create a study plan from checkride date, track progress, view dashboard, and manage checkride checklist in a functional TUI

**Depends on:** Nothing (first phase)

**Requirements:** PLAN-01, PLAN-02, PLAN-03, PLAN-04, PLAN-05, PROG-01, PROG-02, PROG-03, DASH-01, DASH-02, DASH-03, DASH-04, CHKL-01, CHKL-02, CHKL-03, CHKL-04

**Success Criteria** (what must be TRUE):

1. User can enter a checkride date and see a backward-scheduled study plan with daily tasks
2. User can mark any daily task as complete and see progress percentage update
3. Dashboard displays: overall progress %, days until checkride, upcoming week's tasks, quick stats (completed/remaining/overdue)
4. Checkride checklist shows all 4 categories (Documents, Aircraft, Ground, Flight) with per-category and overall completion %
5. User can navigate between Study Plan, Progress, Dashboard, and Checklist views using keyboard

**Plans:** TBD

---

### Phase 2: Calendar Export

**Goal:** Users can export their study plan to external calendars

**Depends on:** Phase 1

**Requirements:** ICAL-01, ICAL-02, ICAL-03, AREM-01

**Success Criteria** (what must be TRUE):

1. User can trigger ICS export and receive a valid .ics file compatible with Apple Calendar
2. Exported ICS file imports correctly into Google Calendar (UTC timezone)
3. User can create Apple Reminders for study tasks via osascript integration

**Plans:** TBD

---

### Phase 3: Google Integration

**Goal:** Users can sync tasks to Google Calendar and OpenCode bot

**Depends on:** Phase 2

**Requirements:** GREM-01, GREM-02, RBOT-01

**Success Criteria** (what must be TRUE):

1. User can authenticate with Google via OAuth2 device flow (works in terminal)
2. User can create tasks in Google Calendar via API
3. User can export tasks in OpenCode bot-compatible format

**Plans:** TBD

---

### Phase 4: Polish

**Goal:** Refined user experience with edge case handling

**Depends on:** Phase 3

**Requirements:** (None - refinements)

**Success Criteria** (what must be TRUE):

1. Application handles missing/invalid data gracefully with clear error messages
2. Loading indicators shown during long operations (ICS export, API calls)
3. Keyboard shortcuts documented and accessible via help command

**Plans:** TBD

---

## Progress Table

| Phase | Plans Complete | Status | Completed |
|-------|----------------|--------|-----------|
| 1. Foundation | 0/5 | Not started | - |
| 2. Calendar Export | 0/3 | Not started | - |
| 3. Google Integration | 0/3 | Not started | - |
| 4. Polish | 0/3 | Not started | - |

---

## Coverage Map

| Requirement | Phase | Status |
|-------------|-------|--------|
| PLAN-01 | Phase 1 | Pending |
| PLAN-02 | Phase 1 | Pending |
| PLAN-03 | Phase 1 | Pending |
| PLAN-04 | Phase 1 | Pending |
| PLAN-05 | Phase 1 | Pending |
| PROG-01 | Phase 1 | Pending |
| PROG-02 | Phase 1 | Pending |
| PROG-03 | Phase 1 | Pending |
| DASH-01 | Phase 1 | Pending |
| DASH-02 | Phase 1 | Pending |
| DASH-03 | Phase 1 | Pending |
| DASH-04 | Phase 1 | Pending |
| CHKL-01 | Phase 1 | Pending |
| CHKL-02 | Phase 1 | Pending |
| CHKL-03 | Phase 1 | Pending |
| CHKL-04 | Phase 1 | Pending |
| ICAL-01 | Phase 2 | Pending |
| ICAL-02 | Phase 2 | Pending |
| ICAL-03 | Phase 2 | Pending |
| AREM-01 | Phase 2 | Pending |
| GREM-01 | Phase 3 | Pending |
| GREM-02 | Phase 3 | Pending |
| RBOT-01 | Phase 3 | Pending |

**Coverage:** 23/23 requirements mapped ✓

---

*Last updated: 2025-02-25*
