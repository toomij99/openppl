# Requirements: PPL Study Planner TUI

**Defined:** 2025-12-20
**Core Value:** Systematically prepare for PPL checkride by breaking down FAA Private Airplane ACS into daily actionable tasks with progress tracking and automated calendar/reminder integration.

## v1 Requirements

### Study Planning

- [ ] **PLAN-01**: User can enter forecast checkride date
- [ ] **PLAN-02**: System generates backward-scheduled study plan from checkride date
- [ ] **PLAN-03**: Daily plan includes 4 categories: Theory, Chair Flying, Garmin 430, CFI Flights
- [ ] **PLAN-04**: ACS Areas breakdown (Area I - Area XVI)
- [ ] **PLAN-05**: Task-level breakdown within each Area

### Progress Tracking

- [ ] **PROG-01**: User can mark daily tasks as complete
- [ ] **PROG-02**: Progress percentage calculated and displayed
- [ ] **PROG-03**: Checkride readiness % shown on dashboard

### Dashboard

- [ ] **DASH-01**: Dashboard view shows overall progress
- [ ] **DASH-02**: Dashboard shows days until checkride
- [ ] **DASH-03**: Dashboard shows upcoming week's study plan
- [ ] **DASH-04**: Quick stats: completed tasks, remaining, overdue

### Checkride Checklist

- [ ] **CHKL-01**: Pre-built checkride requirements checklist
- [ ] **CHKL-02**: User can check off each requirement
- [ ] **CHKL-03**: Checklist shows completion percentage
- [ ] **CHKL-04**: Categories: Documents, Aircraft, Ground, Flight

### Calendar Export

- [ ] **ICAL-01**: Export study plan to ICS file
- [ ] **ICAL-02**: ICS file compatible with Apple Calendar
- [ ] **ICAL-03**: ICS file compatible with Google Calendar

### Reminders Integration

- [ ] **AREM-01**: Create reminders in Apple Reminders via osascript
- [x] **GREM-01**: Create tasks in Google Calendar via API
- [x] **GREM-02**: Google OAuth2 authentication flow
- [x] **RBOT-01**: Export tasks in format compatible with OpenCode bot

### Budget Planning

- [ ] **BUDG-01**: User can enter flight training rates (plane $/hr, CFI $/hr)
- [ ] **BUDG-02**: User can enter living costs (travel ticket, monthly rent, monthly food, monthly car)
- [ ] **BUDG-03**: User can enter estimated total hours needed for PPL (default 60 hours)
- [ ] **BUDG-04**: System calculates total estimated flight cost (total hours × (plane rate + CFI rate))
- [ ] **BUDG-05**: System calculates total estimated budget (flight cost + travel + rent + food + car for training duration)
- [ ] **BUDG-06**: User can enter actual flight hours completed
- [ ] **BUDG-07**: System calculates current spent (actual hrs × (plane + CFI) + living costs paid to date)
- [ ] **BUDG-08**: System forecasts remaining budget to complete PPL
- [ ] **BUDG-09**: Dashboard displays budget progress: spent, remaining, %, projected total
- [ ] **BUDG-10**: System warns if projected total exceeds budget limit

## v2 Requirements

### Advanced Features

- **PLAN-06**: Time estimates per ACS element
- **PLAN-07**: Workdays only option (exclude rest days)
- [ ] **PROG-04**: Flight hour tracking (auto-sync from budget)
- [ ] **PROG-05**: Instructor endorsement tracking

### Multi-Rating Budget Planners

- [ ] **BUDG-11**: IFR Rating budget planner (separate rates, hours, costs)
- [ ] **BUDG-12**: CPL (Commercial Pilot) budget planner
- [ ] **BUDG-13**: CFI (Certified Flight Instructor) budget planner
- [ ] **BUDG-14**: ATP (Airline Transport Pilot) budget planner
- [ ] **BUDG-15**: Combined multi-rating total cost projection

### Integrations

- **GREM-03**: Two-way sync with Google Calendar
- **NOTF-01**: Push notifications for upcoming tasks

## Out of Scope

| Feature | Reason |
|---------|--------|
| IFR/CPL/CFI/ATP planning | v2 - multi-rating planners |
| Mobile app | TUI first, mobile later if needed |
| Offline-first | Requires Google API connectivity |
| Weather integration | Outside scope |
| Flight logging | Separate apps handle this |

## Traceability

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
| BUDG-01 | Phase 1 | Pending |
| BUDG-02 | Phase 1 | Pending |
| BUDG-03 | Phase 1 | Pending |
| BUDG-04 | Phase 1 | Pending |
| BUDG-05 | Phase 1 | Pending |
| BUDG-06 | Phase 1 | Pending |
| BUDG-07 | Phase 1 | Pending |
| BUDG-08 | Phase 1 | Pending |
| BUDG-09 | Phase 1 | Pending |
| BUDG-10 | Phase 1 | Pending |
| ICAL-01 | Phase 2 | Complete |
| ICAL-02 | Phase 2 | Complete |
| ICAL-03 | Phase 2 | Complete |
| AREM-01 | Phase 2 | Complete |
| GREM-01 | Phase 3 | Complete |
| GREM-02 | Phase 3 | Complete |
| RBOT-01 | Phase 3 | Complete |
| BUDG-11 | Phase 5 | Pending |
| BUDG-12 | Phase 5 | Pending |
| BUDG-13 | Phase 5 | Pending |
| BUDG-14 | Phase 5 | Pending |
| BUDG-15 | Phase 5 | Pending |

**Coverage:**
- v1 requirements: 33 total
- Mapped to phases: 33
- Unmapped: 0 ✓

---
*Requirements defined: 2025-12-20*
*Last updated: 2026-02-25 after adding budget planning feature*
