# ROADMAP: PPL Study Planner TUI

**Project:** PPL Study Planner TUI  
**Created:** 2025-02-25  
**Depth:** Quick (3-5 phases)

---

## Phases

- [x] **Phase 1: Foundation** - Core TUI app with study planning, progress tracking, budget planner, dashboard, and checkride checklist
- [x] **Phase 2: Calendar Export** - ICS file export and Apple Reminders integration
- [ ] **Phase 3: Google Integration** - Google Calendar API and OpenCode bot export
- [ ] **Phase 4: Polish** - UI refinements and edge case handling
- [ ] **Phase 5: Multi-Rating Planners** - IFR, CPL, CFI, ATP budget planners

---

## Phase Details

### Phase 1: Foundation

**Goal:** Users can create a study plan from checkride date, track progress, manage budget, view dashboard, and manage checkride checklist in a functional TUI

**Depends on:** Nothing (first phase)

**Requirements:** PLAN-01, PLAN-02, PLAN-03, PLAN-04, PLAN-05, PROG-01, PROG-02, PROG-03, DASH-01, DASH-02, DASH-03, DASH-04, CHKL-01, CHKL-02, CHKL-03, CHKL-04, BUDG-01, BUDG-02, BUDG-03, BUDG-04, BUDG-05, BUDG-06, BUDG-07, BUDG-08, BUDG-09, BUDG-10

**Success Criteria** (what must be TRUE):

1. User can enter a checkride date and see a backward-scheduled study plan with daily tasks
2. User can mark any daily task as complete and see progress percentage update
3. Dashboard displays: overall progress %, days until checkride, upcoming week's tasks, quick stats (completed/remaining/overdue)
4. Checkride checklist shows all 4 categories (Documents, Aircraft, Ground, Flight) with per-category and overall completion %
5. User can navigate between Study Plan, Progress, Dashboard, Budget, and Checklist views using keyboard
6. User can enter/edit flight rates (plane $/hr, CFI $/hr) and living costs (travel, rent, food, car)
7. User can enter estimated PPL hours and actual hours flown
8. Dashboard shows budget: current spent, estimated total, remaining, % used, and projected total
9. System warns if projected budget exceeds entered budget limit

**Plans:** 5 plans
- 01-01-PLAN.md — Core TUI Foundation + Database Models
- 01-02-PLAN.md — Study Planning + Progress Tracking
- 01-03-PLAN.md — Dashboard + Checklist + Budget
- 01-04-PLAN.md — Gap Closure - View Integration and Database Queries
- 01-05-PLAN.md — Gap Closure - Critical Bug Fixes (Dashboard, Checklist, Budget)

---

### Phase 2: Calendar Export

**Goal:** Users can export their study plan to external calendars

**Depends on:** Phase 1

**Requirements:** ICAL-01, ICAL-02, ICAL-03, AREM-01

**Success Criteria** (what must be TRUE):

1. User can trigger ICS export and receive a valid .ics file compatible with Apple Calendar
2. Exported ICS file imports correctly into Google Calendar (UTC timezone)
3. User can create Apple Reminders for study tasks via osascript integration

**Plans:** 3 plans
- [x] 02-01-PLAN.md — ICS Export Service + Compatibility Tests
- [x] 02-02-PLAN.md — Apple Reminders Export Adapter
- [x] 02-03-PLAN.md — Study View Export Wiring + User Feedback

---

### Phase 3: Google Integration

**Goal:** Users can sync tasks to Google Calendar and OpenCode bot

**Depends on:** Phase 2

**Requirements:** GREM-01, GREM-02, RBOT-01

**Success Criteria** (what must be TRUE):

1. User can authenticate with Google via OAuth2 device flow (works in terminal)
2. User can create tasks in Google Calendar via API
3. User can export tasks in OpenCode bot-compatible format

**Plans:** 3 plans
- [x] 03-01-PLAN.md — Google OAuth foundation + export path normalization
- [ ] 03-02-PLAN.md — Google Calendar sync adapter + retry-safe event writes
- [ ] 03-03-PLAN.md — OpenCode bot exporter + Study view Google wiring

---

### Phase 4: Polish

**Goal:** Refined user experience with edge case handling

**Depends on:** Phase 3

**Requirements:** (None - refinements)

**Success Criteria** (what must be TRUE):

1. Application handles missing/invalid data gracefully with clear error messages
2. Loading indicators shown during long operations (ICS export, API calls)
3. Keyboard shortcuts documented and accessible via help command

**Plans:** 3 plans
- 01-01-PLAN.md — Core TUI Foundation + Database Models
- 01-02-PLAN.md — Study Planning + Progress Tracking
- 01-03-PLAN.md — Dashboard + Checklist + Budget

---

### Phase 5: Multi-Rating Planners

**Goal:** Extended budget planners for IFR, CPL, CFI, and ATP ratings

**Depends on:** Phase 4

**Requirements:** BUDG-11, BUDG-12, BUDG-13, BUDG-14, BUDG-15

**Success Criteria** (what must be TRUE):

1. User can create separate budget plans for IFR, CPL, CFI, and ATP ratings
2. Each rating planner has its own rates (plane, CFI), hours estimate, and living costs
3. System calculates combined multi-rating total cost projection
4. User can view all ratings progress and budget in unified view

**Plans:** 3 plans
- 01-01-PLAN.md — Core TUI Foundation + Database Models
- 01-02-PLAN.md — Study Planning + Progress Tracking
- 01-03-PLAN.md — Dashboard + Checklist + Budget

---

## Progress Table

| Phase | Plans Complete | Status | Completed |
|-------|----------------|--------|-----------|
| 1. Foundation | 5/5 | Complete | 2026-02-25 |
| 2. Calendar Export | 3/3 | Complete | 2026-02-25 |
| 3. Google Integration | 1/3 | In Progress | - |
| 4. Polish | 0/3 | Not started | - |
| 5. Multi-Rating Planners | 0/4 | Not started | - |

---

## Coverage Map

| Requirement | Phase | Status |
|-------------|-------|--------|
| PLAN-01 | Phase 1 | Complete |
| PLAN-02 | Phase 1 | Complete |
| PLAN-03 | Phase 1 | Complete |
| PLAN-04 | Phase 1 | Complete |
| PLAN-05 | Phase 1 | Complete |
| PROG-01 | Phase 1 | Complete |
| PROG-02 | Phase 1 | Complete |
| PROG-03 | Phase 1 | Complete |
| DASH-01 | Phase 1 | Complete |
| DASH-02 | Phase 1 | Complete |
| DASH-03 | Phase 1 | Complete |
| DASH-04 | Phase 1 | Complete |
| CHKL-01 | Phase 1 | Complete |
| CHKL-02 | Phase 1 | Complete |
| CHKL-03 | Phase 1 | Complete |
| CHKL-04 | Phase 1 | Complete |
| BUDG-01 | Phase 1 | Complete |
| BUDG-02 | Phase 1 | Complete |
| BUDG-03 | Phase 1 | Complete |
| BUDG-04 | Phase 1 | Complete |
| BUDG-05 | Phase 1 | Complete |
| BUDG-06 | Phase 1 | Complete |
| BUDG-07 | Phase 1 | Complete |
| BUDG-08 | Phase 1 | Complete |
| BUDG-09 | Phase 1 | Complete |
| BUDG-10 | Phase 1 | Complete |
| ICAL-01 | Phase 2 | Complete |
| ICAL-02 | Phase 2 | Complete |
| ICAL-03 | Phase 2 | Complete |
| AREM-01 | Phase 2 | Complete |
| GREM-01 | Phase 3 | Pending |
| GREM-02 | Phase 3 | Complete |
| RBOT-01 | Phase 3 | Pending |
| BUDG-11 | Phase 5 | Pending |
| BUDG-12 | Phase 5 | Pending |
| BUDG-13 | Phase 5 | Pending |
| BUDG-14 | Phase 5 | Pending |
| BUDG-15 | Phase 5 | Pending |

**Coverage:** 37/37 requirements mapped ✓

---

*Last updated: 2026-02-25*
