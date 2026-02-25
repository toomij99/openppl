# Project Research Summary

**Project:** PPL Study Planner TUI
**Domain:** Terminal User Interface Application (Aviation Study Tool)
**Researched:** 2025-02-25
**Confidence:** MEDIUM-HIGH

## Executive Summary

This is a macOS-focused terminal application for PPL (Private Pilot License) students to plan and track their flight training study. The app organizes study around FAA ACS (Airman Certification Standards) structure and generates daily study tasks via backward scheduling from a user-defined checkride date.

**Key insight:** Most existing flight training apps focus on knowledge tests or checkride logistics. The differentiation opportunity is **systematic daily study planning with backward scheduling** — helping students answer "what should I study today?" based on their checkride date.

**Recommended approach:** Build with Go using Bubble Tea (Elm-architecture TUI framework), SQLite for study data, and focus on ICS export first, with Apple Reminders as the primary macOS integration. The ACS structure must be hardcoded initially to avoid PDF parsing pitfalls.

**Main risks:** 
1. Apple Reminders osascript performance at scale (batch operations)
2. Google Calendar OAuth2 terminal flow (use device auth)
3. Backward planning produces unrealistic schedules without user constraints

---

## Key Findings

### Recommended Stack

**Core technologies (from STACK.md):**
- **Go 1.21+** — Single binary deployment, native Google Calendar client
- **Bubble Tea v1.3.10** — Mature Elm-architecture TUI framework, 39k+ stars
- **Bubbles + Lipgloss** — Official UI components and declarative styling
- **SQLite + GORM** — ACID compliance for progress tracking, single-file deployment
- **arran4/golang-ical** — RFC 5545 compliant ICS generation
- **osascript** — Apple Reminders via AppleScript execution

### Expected Features

**Must have (table stakes):**
- ACS reference with hierarchical breakdown (Areas → Tasks → Objectives → Elements)
- Completion status tracking per ACS element
- ICS export (works with Apple Calendar and Google Calendar)
- Daily task view in TUI
- Checkride readiness checklist (Part 61.103 requirements)

**Should have (competitive differentiators):**
- **Backward planning from checkride date** — Primary value proposition
- Daily study task generation based on remaining time
- Progress dashboard with readiness percentage by ACS Area
- Apple Reminders integration (native macOS)

**Defer (v2+):**
- Google Calendar two-way sync (complex OAuth)
- AI-generated oral exam scenarios
- Spaced repetition scheduling
- CFI flight tracking beyond scheduled dates

### Architecture Approach

TUI applications use **Model-View-Update (MVU)** architecture. Bubble Tea implements this pattern: User Input → Message → Update Function → New Model → View Render.

**Major components:**
1. **Model** — Single state struct containing ACS data, study plan, progress metrics, integration state
2. **Services Layer** — Separates calendar, reminders, and storage logic from UI (prevents coupling)
3. **View** — Pure functions rendering terminal output via Lipgloss
4. **Update** — Message handlers with Command pattern for async operations

Directory structure: `model/`, `update/`, `view/`, `services/`, `domain/`, `styles/`

### Critical Pitfalls

1. **Apple Reminders osascript Performance** — Creating 100+ reminders takes 18+ seconds. Batch into fewer calls, add loading indicators.
2. **Google Calendar OAuth2 Terminal Flow** — Default web flow doesn't work in CLI. Use device flow or token storage.
3. **ACS Data Structure Mismatch** — PDF parsing fragile. Lock to FAA-S-ACS-6C version, build manual JSON structure first.
4. **ICS Cross-Platform Incompatibility** — Timezone handling breaks between Apple/Google. Use explicit UTC, test both platforms.
5. **TUI Event Loop Blocking** — Long operations freeze UI. Use async Commands with loading indicators.

---

## Implications for Roadmap

Based on research, suggested phase structure:

### Phase 1: Foundation
**Rationale:** Core data model and TUI shell must work before anything else. Storage and architecture patterns are foundational dependencies.

**Delivers:**
- Domain models (ACS, StudyPlan, Task, Progress)
- SQLite storage service
- Bubble Tea framework initialization
- Basic navigation between screens

**Avoids:** TUI event loop blocking (use proper async from start)

### Phase 2: Core Study Features
**Rationale:** ACS structure and backward planning are the value proposition. These must work before calendar integrations.

**Delivers:**
- ACS data with hierarchical breakdown
- Backward planning algorithm (checkride date → daily tasks)
- Progress tracking with completion percentage
- Daily task view in TUI

**Avoids:** Backward planning unrealistic schedules (add user-configurable study hours/week)

### Phase 3: Calendar Integration
**Rationale:** Core study experience works without calendar sync. Integration complexity and failure modes justify deferring.

**Delivers:**
- ICS export (RFC 5545 compliant, UTC timezone)
- Apple Reminders integration (batched osascript calls)
- Google Calendar API setup (device auth flow)

**Avoids:** ICS incompatibility (test both platforms early), Apple Reminders performance (batch operations)

### Phase 4: Polish
**Rationale:** User-facing refinements after core value works.

**Delivers:**
- Progress dashboard with Area breakdown
- Checkride checklist (documents, oral topics)
- Styling refinements
- Help/key binding display

**Avoids:** Progress percentage misleading (show Written/Oral/Flight breakdown)

### Phase Ordering Rationale

- **Foundation first:** Storage and architecture are prerequisites for everything else
- **Core study before integrations:** Backward planning is the differentiator — must work independently
- **Integrations last:** External APIs have complex auth and failure modes; don't block core UX
- **Avoids pitfalls:** Each phase addresses specific pitfalls from research (async in foundation, UTC ICS in integration, user constraints in planning)

### Research Flags

**Phases needing deeper research during planning:**
- **Phase 2 (Backward Planning):** Algorithm details, time estimation per ACS element — would benefit from `/gsd-research-phase`
- **Phase 3 (Google Calendar):** OAuth2 device flow implementation specifics — may need API research

**Phases with standard patterns (skip research-phase):**
- **Phase 1 (Foundation):** Bubble Tea + SQLite is well-documented, established Go patterns
- **Phase 4 (Polish):** UI styling with Lipgloss has clear documentation

---

## Confidence Assessment

| Area | Confidence | Notes |
|------|------------|-------|
| Stack | HIGH | Bubble Tea mature, Go standard library strong, golang-ical RFC-compliant |
| Features | MEDIUM-HIGH | ACS structure authoritative (FAA docs), backward planning gap identified but not validated |
| Architecture | HIGH | MVU pattern well-documented in Bubble Tea docs, services layer is standard pattern |
| Pitfalls | MEDIUM | TUI pitfalls verified via community posts; aviation app patterns inferred |

**Overall confidence:** MEDIUM-HIGH

### Gaps to Address

- **Backward planning algorithm validation:** No authoritative time-per-element data exists; algorithm needs user testing to tune
- **Actual PPL student usage patterns:** Could benefit from interviews to validate feature priorities
- **ACS element time estimates:** Would need to collect data or use reasonable defaults initially

---

## Sources

### Primary (HIGH confidence)
- Bubble Tea GitHub (v1.3.10, Sep 2025) — TUI framework reference
- FAA Private Pilot Airplane ACS (FAA-S-ACS-6C, November 2023) — ACS structure authoritative
- RFC 5545 — ICS standard

### Secondary (MEDIUM confidence)
- golang-ical GitHub — ICS library implementation
- FAA ACS Companion Guide (FAA-G-ACS-2) — ACS interpretation
- Community TUI development posts — Pitfall prevention

### Tertiary (LOW confidence)
- Aviation study app competitive analysis — Surface-level review, deeper validation needed
- Backward planning algorithm specifics — Inferred from general scheduling patterns

---

*Research completed: 2025-02-25*
*Ready for roadmap: yes*
