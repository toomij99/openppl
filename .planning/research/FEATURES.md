# Feature Landscape

**Domain:** PPL (Private Pilot License) Study Planning Applications
**Researched:** 2025-02-25
**Confidence:** MEDIUM-HIGH

## Executive Summary

PPL study planning apps occupy a niche between general task managers and flight training management systems. The ACS (Airman Certification Standards) framework provides a natural structure that most apps build around, but there's significant variation in how deeply they integrate ACS task breakdown versus treating it as a reference document.

**Key insight from research:** Most existing apps focus on **knowledge test preparation** (written exam) or **checkride day logistics** (scheduling, checklists). Few address **systematic daily study planning with backward scheduling** from a checkride date — this is the primary differentiation opportunity.

---

## Category 1: ACS-Based Study Planning Features

Features that organize study according to FAA ACS structure. This is the foundation — the ACS defines what students must know for the checkride.

### Table Stakes (Must Have)

| Feature | Why Expected | Complexity | Notes |
|---------|--------------|------------|-------|
| **ACS reference access** | The ACS is the authoritative source for checkride requirements. Students must demonstrate competency in every element. | Low | Must reference FAA-S-ACS-6C (Private Pilot Airplane, Nov 2023) |
| **Area of Operation breakdown** | ACS organizes content into 7 Areas (I-VII for PPL). Basic structure students expect. | Low | Standard ACS structure: Preflight Prep, Preflight Procedures, Airport and Seaplane Base Operations, etc. |
| **Task/Objective listing** | Below Areas are Tasks (A-F per Area), then Objectives. | Medium | Hierarchical structure: Area → Task → Objective → Element |
| **Completion status per element** | Track which ACS elements have been studied/practiced. | Medium | Core progress tracking mechanism |

### Differentiators

| Feature | Value Proposition | Complexity | Notes |
|---------|-------------------|------------|-------|
| **Backward planning from checkride date** | "I have a checkride in 60 days, what should I study today?" — **Primary differentiator** | High | Requires estimating task duration, accounting for CFI availability, calculating daily workload |
| **Daily study task generation** | Automatically generate daily tasks based on remaining time and ACS scope | High | Algorithm needs to balance: ground school vs flight prep, knowledge areas vs skills |
| **ACS element time estimation** | Estimate time per ACS element based on complexity | Medium | Could use historical data or fixed multipliers per element type |
| **Priority ordering based on weakness** | Highlight ACS areas where student shows knowledge gaps | Medium | Requires integration with practice test results |
| **Spaced repetition scheduling** | Schedule review of previously-learned material at optimal intervals | High | Research-backed learning technique; significant implementation effort |

### Anti-Features

| Anti-Feature | Why Avoid | What to Do Instead |
|--------------|-----------|-------------------|
| **CFI scheduling integration** | Complex, requires flight school partnerships | Keep as manual input: user enters "CFI flight scheduled for Tuesday" |
| **Flight hour tracking/logbook** | Separate apps exist (ForeFlight, LogTen Pro); scope creep | Allow manual entry of flight completion dates only |

---

## Category 2: Progress Tracking Features

Features that show student advancement toward checkride readiness.

### Table Stakes (Must Have)

| Feature | Why Expected | Complexity | Notes |
|---------|--------------|------------|-------|
| **Completion percentage** | Basic progress indicator; "how ready am I?" | Low | Simple: (completed elements / total elements) × 100 |
| **Checkride readiness checklist** | Regulatory requirements: minimum flight hours, endorsements, documents | Low | Standard FAA Part 61.103 requirements |
| **Area-by-area progress view** | Show which ACS Areas are most complete vs least | Medium | Visual breakdown by Area of Operation |

### Differentiators

| Feature | Value Proposition | Complexity | Notes |
|---------|-------------------|------------|-------|
| **Readiness % by ACS Area** | See which areas need more attention | Medium | Weighted by ACS importance or question frequency |
| **Time-based forecast** | "At current pace, you'll be ready by [date]" | Medium | Track study time per element, extrapolate |
| **Study streak tracking** | Gamify daily study habit | Low | Simple counter; powerful motivation |
| **Flight vs ground balance** | Track ratio of ground study to flight preparation | Low | Manual input of flight dates; auto-categorize study tasks |
| **Checkride date countdown** | Prominent display of days until checkride | Low | Simple but effective motivation |

### Anti-Features

| Anti-Feature | Why Avoid | What to Do Instead |
|--------------|-----------|-------------------|
| **Detailed flight hour logging** | Logbook apps do this better | Link out to user's preferred logbook |
| **Multi-user/CFI dashboard** | Flight school management feature | Keep individual student focus |

---

## Category 3: Calendar & Export Features

Features for integrating study tasks into calendar systems.

### Table Stakes (Must Have)

| Feature | Why Expected | Complexity | Notes |
|---------|--------------|------------|-------|
| **ICS export** | Universal calendar format; works with Apple Calendar, Google Calendar, Outlook | Medium | RFC 5545 standard; well-documented |
| **Calendar event creation** | Create events in primary calendar | Low | ICS import covers this |

### Differentiators

| Feature | Value Proposition | Complexity | Notes |
|---------|-------------------|------------|-------|
| **Apple Reminders integration** | Native reminders in macOS ecosystem | Medium | Use osascript; unique value for macOS users |
| **Google Calendar API sync** | Two-way sync with Google Calendar | High | OAuth required; ongoing sync complexity |
| **Daily task scheduling** | Place tasks at optimal times (morning vs evening) | Medium | Could add preferred time preference |
| **Recurring study blocks** | Set up recurring "study Tuesday/Thursday" blocks | Low | Standard ICS RRULE support |
| **OpenCode bot integration** | Automated task checking via bot | Low | Per PROJECT.md requirement |

### Anti-Features

| Anti-Feature | Why Avoid | What to Do Instead |
|--------------|-----------|-------------------|
| **Real-time calendar sync (polling)** | Complex, rate-limited | Export ICS; let user import manually or use calendar URL |
| **Multi-calendar support** | Overcomplication | Single default calendar export |
| **Outlook/Exchange integration** | Niche; macOS user primarily uses Apple/Google | Focus on Apple + Google |

---

## Category 4: Checkride Preparation Features

Features specifically for the checkride (practical test) itself.

### Table Stakes (Must Have)

| Feature | Why Expected | Complexity | Notes |
|---------|--------------|------------|-------|
| **Oral exam topic checklist** | DPE will ask about these topics | Low | ACS knowledge elements map to oral topics |
| **Required documents list** | What to bring: ID, logbook, medical, aircraft docs | Low | Standard FAA requirements |
| **Aircraft airworthiness checklist** | Verify aircraft inspections current | Low | Annual, VOR check, ELT, etc. |
| **Weight & balance form** | Required for practical test | Low | Can be simplified placeholder |
| **Navigation log template** | Cross-country planning requirement | Medium | ACS specifies XC planning task |

### Differentiators

| Feature | Value Proposition | Complexity | Notes |
|---------|-------------------|------------|-------|
| **Scenario-based practice** | Practice realistic oral scenarios | High | Generate scenarios from ACS risk management elements |
| **DPE-specific tips** | Local DPE knowledge (valuable but hard to maintain) | Low | Could allow user-submitted notes |
| **Mock oral examination** | AI or preset Q&A for practice | High | Significant development effort |
| **Pre-checkride day checklist** | What to do the night before / morning of | Low | Simple but valuable |
| **Flight deck organization guide** | What to have ready in the aircraft | Low | Quick reference |

### Anti-Features

| Anti-Feature | Why Avoid | What to Do Instead |
|--------------|-----------|-------------------|
| **DPE booking integration** | Complex; requires DPE network | User schedules separately |
| **Real-time weather briefing** | Already covered by aviation weather apps | Link out to ForeFlight, aviationweather.gov |
| **Flight plan filing** | Not needed until checkride day; handled separately | Not in scope |

---

## Feature Dependencies

```
Checkride Date Input
    ↓
Backward Planning Algorithm
    ↓
Daily Task Generation
    ↓
ICS Export / Reminder Creation
    ↓
Progress Tracking ← Completion Status per ACS Element
```

**Key dependency:** ACS task breakdown must exist before progress tracking can work. This is foundational.

---

## MVP Recommendation

Based on research findings, prioritize in this order:

### Phase 1: Foundation (Must Build)

1. **ACS reference with hierarchical breakdown** — Foundation of entire app
   - Areas, Tasks, Objectives, Elements structure
   - Completion status tracking per element

2. **Checkride readiness checklist** — Basic regulatory requirements
   - Part 61.103 requirements (night, XC, etc.)
   - Document checklist

3. **ICS export** — Core calendar integration
   - RFC 5545 compliant
   - Works with Apple Calendar and Google Calendar

4. **Daily task view** — Primary TUI interface
   - Shows tasks for today
   - Category filtering (Theory, Chair Flying, Garmin 430, CFI Flights)

### Phase 2: Differentiate

5. **Backward planning from checkride date** — Primary value proposition
   - Input checkride date
   - Calculate required daily study load

6. **Daily study task generation** — Works with backward planning
   - Generate tasks from ACS elements
   - Distribute across available days

7. **Progress dashboard with readiness %** — Visual motivation
   - Overall completion percentage
   - Area-by-area breakdown

### Phase 3: Integrations

8. **Apple Reminders integration** — macOS ecosystem fit
   - osascript-based creation
   - Per-project lists

9. **OpenCode bot integration** — Automated task checking

### Defer to v2+

- Spaced repetition scheduling
- AI-generated oral scenarios
- Google Calendar API two-way sync
- CFI flight tracking beyond scheduled dates

---

## Sources

- FAA Private Pilot Airplane ACS (FAA-S-ACS-6C, November 2023)
- FAA ACS Companion Guide (FAA-G-ACS-2)
- ASA Student Flight Record (ACS-aligned)
- Answers to the ACS app (ACS-based checkride prep)
- Sporty's Checkride Insights (Annotated ACS)
- Flight Schedule Pro (student progress tracking)
- PilotPractice.app (aviation mental math training)
- Checkride.io (checkride scheduling platform)
- "5 Must-Have Apps for Student Pilots" (Jared Ailstock, Feb 2026)

---

## Confidence Assessment

| Area | Confidence | Notes |
|------|------------|-------|
| ACS Features | HIGH | FAA documentation is authoritative; ACS structure is fixed |
| Progress Tracking | MEDIUM | Common patterns from flight school software |
| Calendar/Export | HIGH | Standard ICS format; osascript well-documented |
| Checkride Prep | MEDIUM | Based on general flight training practices |
| Differentiators | MEDIUM | Backward planning is identified gap but not exhaustively validated |

## Research Gaps

- **Actual usage patterns**: Could benefit from interviews with active PPL students
- **Time-to-complete per ACS element**: No authoritative data; would need to estimate or collect
- **Competitive analysis depth**: Only surface-level review of existing apps; deeper dive could reveal more differentiators
