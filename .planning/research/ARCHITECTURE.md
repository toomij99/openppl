# TUI Application Architecture

**Domain:** Terminal User Interface Application  
**Project:** PPL Study Planner TUI  
**Researched:** February 2026

## Architecture Overview

TUI applications follow a **unidirectional data flow** pattern, most commonly implemented via the **Model-View-Update (MVU)** architecture (also called The Elm Architecture). This pattern provides predictable state management and clear component boundaries.

```
┌─────────────────────────────────────────────────────────────────────┐
│                        TUI Application                               │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  ┌──────────┐     ┌──────────┐     ┌──────────┐     ┌──────────┐    │
│  │  Input  │────▶│  Update  │────▶│  Model   │────▶│   View   │    │
│  │ (Keys/  │     │ (Message │     │ (State)  │     │ (Render  │    │
│  │  Mouse) │     │  Handler)│     │          │     │  String)  │    │
│  └──────────┘     └──────────┘     └──────────┘     └──────────┘    │
│       ▲                                                    │         │
│       │                                                    ▼         │
│       │                                            ┌──────────┐      │
│       └────────────────────────────────────────────│ Terminal │      │
│                   (User Action)                     │  Output  │      │
│                                                      └──────────┘      │
└─────────────────────────────────────────────────────────────────────┘
```

## Component Boundaries

### 1. Model (State Layer)

The single source of truth for application state. Contains all data needed to render the UI.

| Responsibility | Contents |
|---------------|----------|
| **Application State** | Current screen, navigation state, user session |
| **Domain Data** | ACS tasks, study plans, progress metrics |
| **UI State** | Selected items, scroll position, input values |
| **Integration State** | Calendar sync status, reminder sync status |

**For PPL Study Planner:** The Model should contain:
- ACS data (Areas, Tasks, Objectives from FAA standards)
- Study plan with daily tasks
- Checkride date and backward planning calculations
- Progress tracking (completed vs pending items)
- Integration state (Apple Reminders, Google Calendar sync status)

### 2. View (Presentation Layer)

Pure function that transforms Model state into terminal output (ANSI strings).

| Responsibility | Implementation |
|---------------|----------------|
| **Layout** | Terminal dimensions, regions, spacing |
| **Rendering** | Widget composition, styling |
| **Navigation Display** | Headers, footers, help text |

**For PPL Study Planner:**
- Daily study plan view
- Progress dashboard with completion percentage
- Checkride checklist view
- Settings/integration configuration view

### 3. Update (Logic Layer)

Handles all state transitions. Pure function: `(Model, Message) → (Model, Cmd)`

| Responsibility | Contents |
|---------------|----------|
| **Message Handling** | Keyboard input, mouse events, timer ticks |
| **State Transitions** | Navigation, form updates, selections |
| **Side Effects** | Commands for I/O operations |

**Commands** represent side effects:
- File I/O (load/save study data)
- External API calls (Google Calendar, Apple Reminders)
- Terminal operations (cursor positioning, screen clear)

### 4. Services Layer (Integration)

Separates external integrations from UI logic. Clean boundary between TUI and system integrations.

```
┌─────────────────────────────────────────────────────────────────┐
│                        Services Layer                            │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐        │
│  │   Storage    │  │   Calendar   │  │  Reminders   │        │
│  │  Service     │  │   Service    │  │   Service   │        │
│  ├──────────────┤  ├──────────────┤  ├──────────────┤        │
│  │ - Load ACS   │  │ - ICS Gen    │  │ - Apple     │        │
│  │ - Save Plan  │  │ - Google API │  │   osascript │        │
│  │ - Progress   │  │ - Parse ICS  │  │ - Google    │        │
│  └──────────────┘  └──────────────┘  │   Tasks API │        │
│                                      └──────────────┘        │
└─────────────────────────────────────────────────────────────────┘
```

## Data Flow

### Primary Flow: User Input → State Update → Render

```
1. User presses key (e.g., 'j' to move down)
      │
      ▼
2. Framework creates Message (e.g., MsgCursorDown)
      │
      ▼
3. Update function receives (Model, MsgCursorDown)
      │
      ▼
4. Update returns (newModel, optional Cmd)
      │
      ▼
5. Framework stores newModel
      │
      ▼
6. View function called with newModel
      │
      ▼
7. View returns rendered string
      │
      ▼
8. Framework writes to terminal
```

### Secondary Flow: Async Operations

```
1. User triggers action (e.g., "Export to Calendar")
      │
      ▼
2. Update returns (newModel, Cmd)
      │
      ▼
3. Cmd executes asynchronously
      │
      ├──▶ Success: Returns SuccessMsg
      │              │
      │              ▼
      │         Update processes SuccessMsg
      │              │
      ▼              ▼
      ...       Render new state
      │
      └──▶ Error: Returns ErrorMsg
                    │
                    ▼
               Update processes ErrorMsg
                    │
                    ▼
               Render error state
```

### Integration Data Flow

```
┌──────────────────────────────────────────────────────────────────┐
│                    External Integrations                         │
├──────────────────────────────────────────────────────────────────┤
│                                                                   │
│  ┌────────────┐    ┌────────────┐    ┌────────────────────┐    │
│  │    ACS     │    │  Calendar  │    │   Apple Reminders  │    │
│  │  Parser    │    │  Export    │    │                    │    │
│  └─────┬──────┘    └─────┬──────┘    └─────────┬──────────┘    │
│        │                  │                      │                │
│        ▼                  ▼                      ▼                │
│  ┌──────────────────────────────────────────────────────────┐   │
│  │                    Storage Service                        │   │
│  │                   (SQLite / JSON)                         │   │
│  └─────────────────────────┬────────────────────────────────┘   │
│                            │                                     │
│                            ▼                                     │
│  ┌──────────────────────────────────────────────────────────┐   │
│  │                      Model                                 │   │
│  │   (studyPlan, acsData, progress, integrationStatus)       │   │
│  └──────────────────────────────────────────────────────────┘   │
│                                                                   │
└──────────────────────────────────────────────────────────────────┘
```

## State Management Approaches

### Recommended: Elm Architecture (Model-View-Update)

Used by Bubble Tea (Go), Elm, and influenced Redux.

| Aspect | Approach |
|--------|----------|
| **State** | Single immutable struct (or copy-on-write) |
| **Updates** | Pure functions returning new state |
| **Side Effects** | Commands (explicit, not implicit) |
| **Testing** | Update functions are pure → easy unit test |

**Implementation Pattern:**

```go
// Model: Single state struct
type Model struct {
    CurrentScreen    Screen
    StudyPlan       *StudyPlan
    ACSData         *ACSData
    Progress        Progress
    IntegrationState IntegrationState
    // UI state
    SelectedIndex   int
    ScrollOffset    int
}

// Messages: Type-safe events
type Msg interface {
    isMsg()
}

type MsgNavigate struct{ screen Screen }
type MsgSelectTask struct{ taskID string }
type MsgExportCalendar struct{}
type MsgCalendarExported struct{ err error }
type MsgTick struct{} // For timed updates

// Update: Pure function
func (m Model) Update(msg Msg) (Model, tea.Cmd) {
    switch msg := msg.(type) {
    case MsgNavigate:
        m.CurrentScreen = msg.screen
        return m, nil
    case MsgSelectTask:
        // Update selection
        return m, nil
    case MsgExportCalendar:
        return m, calendarService.Export(m.StudyPlan)
    case MsgCalendarExported:
        if msg.err != nil {
            m.IntegrationState.LastError = msg.err
        }
        return m, nil
    }
    return m, nil
}

// View: Render from state
func (m Model) View() string {
    switch m.CurrentScreen {
    case ScreenDashboard:
        return renderDashboard(m)
    case ScreenStudyPlan:
        return renderStudyPlan(m)
    }
    return ""
}
```

### Alternative: React-like (Textual)

Textual uses reactive attributes and message passing:

```python
# Textual (Python) approach
class StudyPlanScreen(Screen):
    def compose(self) -> ComposeResult:
        yield Header()
        yield ProgressGauge()
        yield TaskList(id="tasks")
    
    def on_mount(self) -> None:
        # Load data when screen mounts
        self.load_tasks()
    
    @work
    async def load_tasks(self):
        tasks = await self.app.storage.get_tasks()
        self.task_list = tasks
```

## Component Organization

### Directory Structure (Bubble Tea / Go)

```
ppl-study-planner/
├── main.go                 # Entry point, program construction
├── model/
│   └── model.go            # Model struct, Message types
├── update/
│   └── update.go           # Update function(s)
├── view/
│   ├── view.go             # Main view dispatcher
│   ├── dashboard.go        # Dashboard screen
│   ├── study_plan.go       # Study plan screen
│   └── components/         # Reusable UI components
├── services/
│   ├── storage.go          # SQLite/JSON persistence
│   ├── acs_parser.go       # ACS PDF/data parsing
│   ├── calendar.go         # ICS generation, Google Calendar
│   └── reminders.go       # Apple Reminders, Google Tasks
├── domain/
│   ├── acs.go              # ACS domain models
│   ├── study_plan.go       # Study plan logic
│   └── progress.go         # Progress calculations
└── styles/
    └── styles.go           # Lipgloss style definitions
```

### Directory Structure (Textual / Python)

```
ppl_study_planner/
├── app.py                  # App class, message definitions
├── models/
│   └── state.py            # AppState dataclass
├── screens/
│   ├── dashboard.py        # Dashboard screen
│   └── study_plan.py       # Study plan screen
├── widgets/
│   ├── progress_gauge.py   # Custom widgets
│   └── task_list.py
├── services/
│   ├── storage.py          # SQLite persistence
│   ├── calendar.py         # Calendar integrations
│   └── reminders.py        # Reminder integrations
├── styles/
│   └──.tcss                # Textual CSS styles
└── domain/
    ├── acs.py              # ACS domain models
    └── planning.py          # Study plan logic
```

## Suggested Build Order

Based on component dependencies:

### Phase 1: Core Model & Storage
1. Define domain models (ACS, StudyPlan, Task, Progress)
2. Implement storage service (SQLite)
3. Basic Model struct with persistence

**Why:** Everything else depends on data persistence. Build this first to validate data model.

### Phase 2: Basic TUI Shell
4. Initialize Bubble Tea / Textual framework
5. Create Model, Update, View skeleton
6. Add basic navigation between screens

**Why:** Establishes the architecture contract. Verify MVU pattern works before adding complexity.

### Phase 3: Domain Logic
7. ACS parser (parse FAA ACS data)
8. Backward planning algorithm (checkride date → daily tasks)
9. Progress calculation

**Why:** Pure business logic, no UI complexity. Easy to test in isolation.

### Phase 4: UI Screens
10. Dashboard view (progress display)
11. Study plan view (daily tasks)
12. Checkride checklist view

**Why:** UI builds on stable domain layer. Can iterate on presentation without breaking logic.

### Phase 5: Integrations
13. ICS export (calendar integration)
14. Apple Reminders (osascript)
15. Google Calendar API

**Why:** External integrations have dependencies and potential failures. Build after core UX works.

### Phase 6: Polish
16. Styling (Lipgloss / CSS)
17. Error handling UI
18. Help/key binding display

## Key Architecture Decisions for PPL Study Planner

| Decision | Rationale |
|----------|-----------|
| **MVU Architecture** | Predictable state, testable updates, clear data flow |
| **Services Layer** | Separates calendar/reminder logic from TUI complexity |
| **SQLite over JSON** | Structured queries for progress tracking, future-proof for more data |
| **Command Pattern** | Async integrations don't block UI thread |
| **Single Model** | All state in one place makes serialization/save-on-exit trivial |

## Anti-Patterns to Avoid

### 1. Mutation in View
**Bad:** `view() { model.selectedIndex++ }`  
**Good:** View is pure—only reads state. Mutations happen in Update.

### 2. Blocking I/O in Update
**Bad:** `case MsgLoadData: return m, readFileSync()`  
**Good:** Use async Commands that return messages when complete.

### 3. Coupling Services to Model
**Bad:** Model directly calls `calendar.Export()`  
**Good:** Services are injected or called via Commands. Model doesn't know about implementations.

### 4. Scattered State
**Bad:** Multiple global variables tracking different state  
**Good:** Single Model struct. Everything flows through Update.

## Sources

- [Bubble Tea Framework](https://github.com/charmbracelet/bubbletea) - Go TUI framework using Elm Architecture
- [Textual Framework](https://textual.textualize.io/) - Python TUI framework with reactive patterns
- [Ratatui Elm Architecture](https://ratatui.rs/concepts/application-patterns/the-elm-architecture/) - Rust TUI library guidance
- [Elm Architecture Guide](https://guide.elm-lang.org/architecture/) - Original MVU pattern documentation
- [Bubble Tea Best Practices](https://leg100.github.io/en/posts/building-bubbletea-programs/) - Production experience writeup
