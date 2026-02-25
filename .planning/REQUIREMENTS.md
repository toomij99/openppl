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
- [ ] **GREM-01**: Create tasks in Google Calendar via API
- [ ] **GREM-02**: Google OAuth2 authentication flow
- [ ] **RBOT-01**: Export tasks in format compatible with OpenCode bot

## v2 Requirements

### Advanced Features

- **PLAN-06**: Time estimates per ACS element
- **PLAN-07**: Workdays only option (exclude rest days)
- **PROG-04**: Flight hour tracking
- **PROG-05**: Instructor endorsement tracking

### Integrations

- **GREM-03**: Two-way sync with Google Calendar
- **NOTF-01**: Push notifications for upcoming tasks

## Out of Scope

| Feature | Reason |
|---------|--------|
| IFR/CPL planning | PPL only for v1 |
| Mobile app | TUI first, mobile later |
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
| ICAL-01 | Phase 2 | Pending |
| ICAL-02 | Phase 2 | Pending |
| ICAL-03 | Phase 2 | Pending |
| AREM-01 | Phase 2 | Pending |
| GREM-01 | Phase 3 | Pending |
| GREM-02 | Phase 3 | Pending |
| RBOT-01 | Phase 3 | Pending |

**Coverage:**
- v1 requirements: 23 total
- Mapped to phases: 23
- Unmapped: 0 ✓

---
*Requirements defined: 2025-12-20*
*Last updated: 2025-12-20 after initial definition*
