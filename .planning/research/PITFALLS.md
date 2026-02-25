# Domain Pitfalls

**Project:** PPL Study Planner TUI
**Researched:** 2026-02-25
**Confidence:** MEDIUM (TUI findings verified via docs/code, aviation app patterns inferred)

---

## Critical Pitfalls

Mistakes that cause rewrites or major issues.

### 1. Apple Reminders osascript Performance
**What goes wrong:** Creating 100+ reminders takes 18+ seconds. Each osascript call spawns Reminders app.
**Why it happens:** AppleScript is slow when called from CLI. Each reminder requires separate app launch overhead.
**Consequences:** User experience is sluggish; bulk operations feel broken.
**Prevention:** 
- Batch reminders into fewer osascript calls
- Consider EventKit alternative if performance critical
- Add loading indicators for long operations
**Detection:** Benchmark with 50, 100, 200 reminders. Expect <2 seconds for 50.
**Phase:** Integration Phase

### 2. Google Calendar OAuth2 in CLI
**What goes wrong:** OAuth2 flow expects browser redirect, but CLI app runs in terminal.
**Why it happens:** Google default OAuth flow is web-based, not terminal-native.
**Consequences:** Users can't authenticate, stuck at login step.
**Prevention:** 
- Use device flow (oauth2.DeviceConfig) for terminal auth
- Or generate tokens via web flow once, store for future runs
- Document setup process clearly
**Phase:** Integration Phase

### 3. ACS Data Structure Mismatch
**What goes wrong:** ACS parsing produces wrong task breakdown, missing objectives, or incorrect Area groupings.
**Why it happens:** 
- FAA ACS PDF structure varies between versions
- ACS has Areas → Tasks → Objectives → Knowledge/Skill/Risk elements
- Learning Statement codes (PLT codes) still used alongside ACS for written exams
- Tables in PDF don't parse cleanly to hierarchical JSON
**Consequences:** 
- Study plan missing critical objectives
- User misses checkride requirements
- Trust in app erodes
**Prevention:** 
- Lock to specific ACS version (FAA-G-ACS-2 with Changes 1 & 2)
- Build manual JSON structure first, automate later
- Validate parsed data against known ACS code patterns (e.g., "PA.I.A.K1")
- Handle both ACS and legacy PLT codes for written exam content
**Detection:** Compare parsed output against official ACS PDF manually
**Phase:** Data Ingestion Phase

### 4. ICS/Calendar Cross-Platform Incompatibility
**What goes wrong:** ICS file works in Apple Calendar but fails in Google Calendar, or shows wrong times, duplicates, or doesn't appear.
**Why it happens:** 
- RFC 5545 interpreted differently per client
- Timezone handling inconsistent (UTC vs floating)
- Microsoft 2025 Outlook update broke many ICS generators by strictly enforcing RFC
- RRULE (recurrence) especially problematic
- Apple vs Google handle METHOD:REQUEST differently
**Consequences:** 
- User misses study sessions
- Duplicated events flood calendar
- Defeats core value proposition
**Prevention:** 
- Generate simplest possible ICS (avoid complex RRULE if possible)
- Always use explicit UTC timezone (`TZID=UTC`)
- Test with multiple calendars before feature release
- Consider Google Calendar API direct integration over ICS for Google
- Use iCalendar.org validator on generated files
**Phase:** Calendar Integration Phase

### 5. TUI Event Loop Blocking
**What goes wrong:** Long-running operations (file I/O, API calls, calendar sync) freeze the entire terminal UI.
**Why it happens:** Bubble Tea and Textual both use event-driven architecture. Beginners run synchronous code in message handlers without offloading to background threads.
**Consequences:** 
- User thinks app is crashed
- Keyboard input stops responding
- Terminal may not restore properly on crash
**Prevention:** 
- Use async/await patterns (Textual) or goroutines with channel responses (Bubble Tea)
- Never block the main event loop with >10ms operations
- Show loading indicators for any operation >500ms
**Phase:** Foundation Phase

### 6. Terminal State Not Restored on Crash
**What goes wrong:** App crashes mid-execution and leaves terminal in raw mode, echo disabled, or alternate screen buffer active.
**Why it happens:** 
- Crash before `tea.Quit()` or `app.exit()`
- Missing cleanup/panic handler
- crossterm not properly restored
**Consequences:** 
- Terminal becomes unusable
- User must close and reopen terminal
**Prevention:** 
- Use framework's proper shutdown
- Install panic handler that calls restore
- Test crash scenarios during development
**Phase:** Foundation Phase

---

## Moderate Pitfalls

### 7. Backward Planning Algorithm Produces Unrealistic Schedules
**What goes wrong:** Checkride date minus required hours produces 2-hour daily study blocks that are impossible for working adults.
**Why it happens:** 
- ACS has ~90+ objectives across Areas I-XI
- Simple division doesn't account for CFI availability, weather days, fatigue from flights
**Consequences:** 
- User ignores app after week 1
- Study plan becomes disconnected from reality
**Prevention:** 
- Allow user to set available study hours/week
- Add buffer days between phases
- Suggest realistic targets (e.g., "plan for 6-9 months, not 3")
- Include rest/weather days in calculation
**Phase:** Planning Algorithm Phase

### 8. Bubble Tea Mouse Support Assumptions
**What goes wrong:** App works with mouse on developer machine but not user terminals.
**Why it happens:** Terminal emulators vary in mouse support; some disable it by default.
**Prevention:** Ensure all interactions work with keyboard alone. Test in iTerm2, Terminal.app, tmux.
**Phase:** Foundation Phase

### 9. Hardcoded Terminal Assumptions
**What goes wrong:** UI looks perfect in iTerm2 but breaks in tmux, VS Code terminal, or Windows Terminal.
**Why it happens:** 
- Different escape sequence support
- Color rendering differs (8-color vs 256 vs truecolor)
- Mouse, keyboard combinations not passed through consistently
**Prevention:** 
- Test on multiple terminals
- Test inside and outside tmux/screen
- Don't assume wide Unicode support (box-drawing characters)
- Provide fallback ASCII rendering option
**Phase:** Foundation Phase

### 10. SQLite Concurrent Access
**What goes wrong:** App crashes when accessing database from multiple goroutines.
**Why it happens:** SQLite has connection limits; GORM by default uses connection pool.
**Prevention:** Use single database connection, serialize writes. Or use `Mode: Disable`.
**Phase:** Foundation Phase

---

## Minor Pitfalls

### 11. Progress Percentage Misleading
**What goes wrong:** User marks all knowledge items complete but "readiness" shows 60% because skills/flights not tracked.
**Why it happens:** 
- ACS has knowledge, skill, AND experience requirements
- Progress = simple completion percentage misses nuance
**Prevention:** 
- Show breakdown: Written %, Oral %, Flight %
- Be explicit about what's measured vs. what's missing
**Phase:** Dashboard Phase

### 12. ANSI Color Rendering
**What goes wrong:** Colors look wrong on some terminals.
**Why it happens:** Terminal color support varies (256 colors, truecolor, ANSI).
**Prevention:** Use Lipgloss color system which handles fallback; test on target terminals.
**Phase:** UI Polish Phase

### 13. Screen Resize Handling
**What goes wrong:** UI breaks when terminal resized during use.
**Why it happens:** Bubble Tea handles this, but custom layouts may not.
**Prevention:** Use Bubble Tea's responsive components; avoid hardcoded dimensions.
**Phase:** Foundation Phase

### 14. Testing Strategy Missing
**What goes wrong:** Unit tests pass but UI is broken. No way to catch rendering bugs in CI.
**Why it happens:** 
- TUI apps harder to test than web apps
- Standard print debugging doesn't work
**Prevention:** 
- Use framework-specific testing (teatest for Bubble Tea, textual snapshot testing)
- Test at component level with mocked state
- Run visual tests in CI with consistent terminal size
**Phase:** Foundation Phase

---

## Phase-Specific Warnings

| Phase Topic | Likely Pitfall | Mitigation |
|-------------|---------------|------------|
| Foundation | Event loop blocking, terminal restoration | Use async patterns, test crash recovery |
| Data Ingestion | ACS parsing errors | Manual validation, lock to ACS version |
| Calendar Integration | ICS incompatibility | Test both platforms early, simplify RRULE |
| Planning Algorithm | Unrealistic schedules | User testing, configurable hours |
| Apple Reminders | Performance at scale | Benchmark early, batch operations |
| Google Calendar | OAuth2 terminal flow | Use device auth or token storage |
| Database | Concurrent access | Single connection or serialize writes |

---

## Sources

### TUI Development
- [Bubble Tea Tips](https://leg100.github.io/en/posts/building-bubbletea-programs/) — Event loop, message ordering, layout arithmetic
- [Ratatui Best Practices](https://github.com/ratatui/ratatui/discussions/220) — MVC pattern, error handling
- [Testing TUI Apps](https://blog.waleedkhan.name/testing-tui-apps/) — Integration vs E2E testing patterns
- [Textual Testing](https://textual.textualize.io/guide/testing/) — Snapshot testing with pytest
- [TUI Terminal Compatibility](https://hoop.dev/blog/moving-past-ncurses-modern-tui-alternatives/) — Ncurses portability issues

### Aviation Study App
- [FAA ACS Companion Guide](https://www.faa.gov/training_testing/testing/acs/acs_companion_guide_pilots.pdf) — ACS structure reference (FAA-G-ACS-2 with Changes 1 & 2)
- [PPL Exam Prep](https://www.pplexam.app/) — Existing app patterns
- [Pilot Partner - Training Dropout](https://www.pilotpartner.net/software-addresses-dismal-flight-school-failure-rate/) — Structured training importance

### Calendar/ICS
- [ICS Troubleshooting Guide 2025](https://synara.events/articles/ics-troubleshooting-guide-2025) — Cross-client compatibility
- [Add to Calendar Hidden Hell](https://add-to-calendar-pro.com/articles/hidden-hell-hand-coding-add-to-calendar-links-453d9cff) — RFC compliance issues
- [Outlook ICS RRULE Issues](https://techcommunity.microsoft.com/discussions/outlookgeneral/outlook-not-parsing-ics-invites-with-rrule-bymonthday-1-or-bysetpos-1-properly/4403885) — RRULE edge cases

### Original Sources
- [Stack Overflow: osascript performance](https://stackoverflow.com/questions/66817320)
- [Google OAuth2 Go docs](https://pkg.go.dev/golang.org/x/oauth2)
- [golang-ical timezone handling](https://github.com/arran4/golang-ical)
