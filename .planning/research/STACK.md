# Technology Stack

**Project:** PPL Study Planner TUI
**Researched:** 2025-02-25
**Confidence:** HIGH

## Recommended Stack

### Core Framework
| Technology | Version | Purpose | Why |
|------------|---------|---------|-----|
| **Bubble Tea** | v1.3.10 (Sep 2025) | TUI framework | Mature, Elm-architecture based, 39k+ stars, production-tested. Best Go TUI option. |
| **Bubbles** | Latest | UI components | Official component library (list, textinput, spinner, progress). Works seamlessly with Bubble Tea. |
| **Lipgloss** | Latest | Styling | Declarative styling from Charm ecosystem. CSS-like approach for terminal rendering. |

### Language
| Technology | Version | Purpose | Why |
|------------|---------|---------|-----|
| **Go** | 1.21+ | Primary language | Excellent for CLI tools, single binary deployment, native Google Calendar API client. |

### Calendar/ICS Generation
| Library | Purpose | Why |
|---------|---------|-----|
| **arran4/golang-ical** | ICS file generation | Most comprehensive Go ICS library. Supports parsing and creation, RRULE recurrence, timezones. Apache 2.0 licensed. |
| **emersion/go-ical** | Alternative minimal | Lightweight, MIT licensed. Less feature-rich but simpler API. |

**Recommendation:** Use `arran4/golang-ical` for robust ICS export that works with Apple Calendar and Google Calendar.

### Calendar Integrations
| Integration | Technology | Why |
|-------------|------------|-----|
| **Apple Reminders** | osascript (AppleScript) | Native macOS integration. Project explicitly requires macOS. Run via `exec.Command` in Go. |
| **Google Calendar** | google.golang.org/api/calendar/v3 | Official Google client library. Supports full Calendar API (create events, sync, reminders). |

### Data Storage
| Technology | Purpose | Why |
|------------|---------|-----|
| **SQLite** + **GORM** | Structured data | ACID compliance, complex queries for progress tracking, single-file deployment. |
| **JSON** | Configuration | Simple settings storage (checkride date, preferences). |

**Recommendation:** SQLite for study data (tasks, progress, ACS breakdown). JSON for app configuration.

## Alternative Options

### TUI Frameworks Considered

| Framework | Language | Verdict | Why Not |
|-----------|----------|---------|---------|
| **Bubble Tea** | Go | **Recommended** | Best ecosystem, Google Calendar official Go client matches stack |
| **Textual** | Python | Strong alternative | Python has rich calendar libraries but requires separate stack |
| **Ratatui** | Rust | Good alternative | If Rust preferred; less friendly for Google Calendar integration |
| **ncurses** | C | Avoid | Brittle, poor cross-platform support, outdated |

### ICS Libraries Compared

| Library | Go/Python | Verdict | Notes |
|---------|-----------|---------|-------|
| **arran4/golang-ical** | Go | **Recommended** | Full RFC 5545 support |
| **emersion/go-ical** | Go | Good alternative | Minimal, MIT licensed |
| **ics.py** | Python | Recommended | If Python stack chosen |
| **icalendar** | Python | Good alternative | Mature, widely used |

### Storage Compared

| Storage | Use Case | Verdict | Notes |
|---------|----------|---------|-------|
| **SQLite** | Study data | **Recommended** | Relational queries, progress tracking, ACID |
| **JSON** | Config | **Recommended** | Simple key-value settings |
| **BoltDB** | Key-value | Alternative | Pure Go, but SQLite more versatile |

## Installation

### Go Dependencies
```bash
# Core TUI
go get github.com/charmbracelet/bubbletea
go get github.com/charmbracelet/bubbles
go get github.com/charmbracelet/lipgloss

# Database
go get gorm.io/gorm
go get gorm.io/driver/sqlite

# Calendar
go get github.com/arran4/golang-ical
go get google.golang.org/api/calendar/v3
go get golang.org/x/oauth2/google

# Utilities
go get github.com/spf13/cobra   # CLI (if needed beyond TUI)
```

## Platform Considerations

### macOS Focus
- **Apple Reminders**: Uses `osascript` to invoke AppleScript. Requires:
  - macOS with Reminders app
  - Automation permission granted by user
  - Works well with CLI invocation: `osascript -e 'tell app "Reminders"'`
  
- **EventKit**: Alternative native API (Swift/Objective-C only), not available from Go directly

### Google Calendar
- OAuth2 authentication required
- Works on macOS with browser-based auth flow
- API quota: 1M requests/day (generous for personal use)

## Confidence Assessment

| Area | Level | Reason |
|------|-------|--------|
| TUI Framework | HIGH | Bubble Tea is well-documented, actively maintained (2025 releases), extensive community |
| ICS Libraries | HIGH | golang-ical is stable, well-tested, RFC 5545 compliant |
| Google Calendar API | HIGH | Official Google client, comprehensive docs |
| Apple Reminders | MEDIUM | osascript approach works but can be slow for bulk operations |
| Storage | HIGH | SQLite + GORM is standard Go pattern |

## Sources

- [Bubble Tea GitHub](https://github.com/charmbracelet/bubbletea) - v1.3.10 (Sep 2025)
- [Bubbles Components](https://github.com/charmbracelet/bubbles) - Official component library
- [Lipgloss Styling](https://github.com/charmbracelet/lipgloss) - Declarative terminal styling
- [golang-ical](https://github.com/arran4/golang-ical) - ICS parser/serializer
- [Google Calendar API Go](https://pkg.go.dev/google.golang.org/api/calendar/v3) - Official client v0.233.0
- [Google Calendar Quickstart Go](https://developers.google.com/workspace/calendar/api/quickstart/go)
- [Apple Reminders osascript](https://apple.stackexchange.com/questions/416280) - Stack Exchange guidance
