---
phase: 04-polish
plan: 04
subsystem: ui
tags: [nextjs, shadcn, next-themes, recharts, tanstack-table]
requires:
  - phase: 04-polish
    provides: TUI polish semantics and global shortcut/help baseline from Plans 01-03
provides:
  - Next.js App Router web dashboard workspace under web/
  - shadcn-style dashboard shell with sidebar/inset composition
  - light/dark theme architecture with persisted explicit toggle
  - chart and table dashboard modules with loading/empty/error fallbacks
affects: [web-mode, dashboard, ui-polish, future-phase-04.1]
tech-stack:
  added: [nextjs, tailwindcss, shadcn-patterns, next-themes, recharts, @tanstack/react-table]
  patterns: [class-based theme tokens, typed dashboard modules, explicit async state branching]
key-files:
  created:
    - web/package.json
    - web/components/theme-provider.tsx
    - web/components/theme-toggle.tsx
    - web/components/dashboard/overview-cards.tsx
    - web/components/dashboard/progress-chart.tsx
    - web/components/dashboard/tasks-table.tsx
  modified:
    - web/app/layout.tsx
    - web/app/globals.css
    - web/app/dashboard/page.tsx
key-decisions:
  - "Bootstrap a dedicated web/ Next.js workspace so Go/TUI build paths remain isolated."
  - "Implement light and dark tokens together (including chart variables) to prevent late-stage contrast drift."
  - "Encode dashboard UX states with deterministic page-level branches for loading, empty, and error."
patterns-established:
  - "Dashboard uses SidebarProvider/AppSidebar/SidebarInset shell composition with typed module boundaries."
  - "Charts consume CSS variables (var(--color-*), var(--chart-*)) for theme-safe rendering."
requirements-completed: []
duration: 5 min
completed: 2026-02-26
---

# Phase 4 Plan 04: shadcn web dashboard baseline + required dark/light themes Summary

**A dedicated Next.js dashboard workspace now ships with shadcn-style shell composition, persisted light/dark switching, and explicit loading/empty/error module states.**

## Performance

- **Duration:** 5 min
- **Started:** 2026-02-26T08:14:23Z
- **Completed:** 2026-02-26T08:20:01Z
- **Tasks:** 3
- **Files modified:** 26

## Accomplishments
- Bootstrapped `web/` workspace with Next.js App Router, shadcn config conventions, and dashboard shell primitives.
- Added hydration-safe `next-themes` provider wiring, explicit light/dark toggle controls, and paired light/dark design + chart tokens.
- Implemented typed dashboard modules (overview cards, progress chart, tasks table) and explicit async-state UI branches for loading, empty, and error.

## Task Commits

Each task was committed atomically:

1. **Task 1: Bootstrap web dashboard workspace with shadcn shell primitives** - `5f0874d` (feat)
2. **Task 2: Implement required dark + light theme architecture and switcher** - `78f3578` (feat)
3. **Task 3: Add dashboard data modules with loading/empty/error polish states** - `8ec14c4` (feat)

**Plan metadata:** `(pending docs commit)`

## Files Created/Modified
- `web/package.json` - Next.js workspace scripts/dependencies for dashboard stack.
- `web/components.json` - shadcn registry configuration and aliases.
- `web/app/layout.tsx` - root theme provider wiring with `suppressHydrationWarning`.
- `web/app/globals.css` - paired light/dark design tokens and chart color variables.
- `web/app/dashboard/page.tsx` - shell composition and explicit async state branches.
- `web/components/theme-toggle.tsx` - user-facing light/dark switch control.
- `web/components/dashboard/progress-chart.tsx` - Recharts module using tokenized chart colors.
- `web/components/dashboard/tasks-table.tsx` - TanStack table module for task queue rendering.

## Decisions Made
- Used a standalone `web/` workspace and kept Go root untouched for minimal coupling between TUI and web build systems.
- Chose class-based theme tokens with `next-themes` at the root layout so chart/table readability stays consistent across themes.
- Drove async UX states from route-level state resolution (`loading`, `empty`, `error`, plus shared semantics) to avoid blank/frozen sections.

## Deviations from Plan

### Auto-fixed Issues

**1. [Rule 3 - Blocking] Replaced unsupported Next config format**
- **Found during:** Task 2 (theme architecture verification)
- **Issue:** `next lint` failed because Next.js 14 does not accept `next.config.ts`.
- **Fix:** Switched to `next.config.mjs` with equivalent config content.
- **Files modified:** `web/next.config.mjs`
- **Verification:** `npm --prefix web run lint && npm --prefix web run build` passed after change.
- **Committed in:** `5f0874d` (part of task commit)

---

**Total deviations:** 1 auto-fixed (1 blocking)
**Impact on plan:** Fix was required for workspace validity; no scope creep.

## Authentication Gates

None.

## Issues Encountered
None.

## User Setup Required

None - no external service configuration required.

## Next Phase Readiness
- Phase 4 plan set is now complete from a dashboard UX perspective.
- Ready for Phase 04.1 planning and web command/runtime integration.

---
*Phase: 04-polish*
*Completed: 2026-02-26*

## Self-Check: PASSED
- Verified summary file exists on disk.
- Verified task commits `5f0874d`, `78f3578`, and `8ec14c4` exist in git history.
