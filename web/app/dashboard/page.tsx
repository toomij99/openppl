import { AppSidebar } from "@/components/app-sidebar"
import { OverviewCards, type OverviewMetrics } from "@/components/dashboard/overview-cards"
import { ProgressChart, type ProgressPoint } from "@/components/dashboard/progress-chart"
import { TasksTable, type TaskRow } from "@/components/dashboard/tasks-table"
import { ThemeToggle } from "@/components/theme-toggle"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { SidebarInset, SidebarProvider, SidebarTrigger } from "@/components/ui/sidebar"

type ViewState = "idle" | "loading" | "success" | "warning" | "error" | "empty"

type DashboardData = {
  state: ViewState
  metrics: OverviewMetrics
  progress: ProgressPoint[]
  tasks: TaskRow[]
}

const baseMetrics: OverviewMetrics = {
  completionPercent: 64,
  completedTasks: 58,
  totalTasks: 91,
  daysToCheckride: 42,
  weeklyHours: 11.5,
  budgetSpent: 9800,
}

const baseProgress: ProgressPoint[] = [
  { week: "W1", completed: 12, target: 10 },
  { week: "W2", completed: 18, target: 16 },
  { week: "W3", completed: 29, target: 26 },
  { week: "W4", completed: 39, target: 36 },
  { week: "W5", completed: 50, target: 46 },
  { week: "W6", completed: 58, target: 56 },
]

const baseTasks: TaskRow[] = [
  { id: "1", task: "Review cross-country weather minimums", category: "Theory", status: "success", dueDate: "Mar 02" },
  { id: "2", task: "Chair-fly soft-field pattern", category: "Chair Flying", status: "loading", dueDate: "Mar 03" },
  { id: "3", task: "Program VOR direct-to on GNS 430", category: "Garmin 430", status: "warning", dueDate: "Mar 04" },
  { id: "4", task: "Preflight briefing with CFI", category: "CFI Flights", status: "idle", dueDate: "Mar 05" },
]

function buildDashboardData(state: ViewState): DashboardData {
  if (state === "empty") {
    return { state, metrics: { ...baseMetrics, completionPercent: 0, completedTasks: 0 }, progress: [], tasks: [] }
  }

  if (state === "error") {
    return { state, metrics: baseMetrics, progress: [], tasks: [] }
  }

  return { state, metrics: baseMetrics, progress: baseProgress, tasks: baseTasks }
}

function resolveState(value?: string): ViewState {
  if (value === "loading" || value === "empty" || value === "error" || value === "warning" || value === "idle") {
    return value
  }
  return "success"
}

export default function DashboardPage({ searchParams }: { searchParams?: { state?: string } }) {
  const state = resolveState(searchParams?.state)
  const dashboard = buildDashboardData(state)

  const stateBanner =
    state === "warning"
      ? "Warning: some study exports need attention."
      : state === "idle"
        ? "Idle: dashboard is ready for the next action."
        : null

  return (
    <SidebarProvider>
      <div className="flex min-h-screen w-full">
        <AppSidebar />
        <SidebarInset>
          <header className="flex h-16 items-center justify-between border-b bg-background px-6">
            <div>
              <h1 className="text-xl font-semibold tracking-tight">PPL Dashboard</h1>
              <p className="text-sm text-muted-foreground">Terminal and web progress at a glance</p>
            </div>
            <div className="flex items-center gap-2">
              <ThemeToggle />
              <SidebarTrigger />
            </div>
          </header>

          <div className="space-y-6 p-6">
            {stateBanner ? (
              <Card>
                <CardHeader>
                  <CardTitle className="text-base">Status</CardTitle>
                  <CardDescription>{stateBanner}</CardDescription>
                </CardHeader>
              </Card>
            ) : null}

            {dashboard.state === "loading" ? (
              <Card>
                <CardHeader>
                  <CardTitle>Loading dashboard data</CardTitle>
                  <CardDescription>Pulling study plan metrics and recent task activity.</CardDescription>
                </CardHeader>
                <CardContent className="space-y-3">
                  <div className="h-4 w-1/2 animate-pulse rounded bg-muted" />
                  <div className="h-4 w-2/3 animate-pulse rounded bg-muted" />
                </CardContent>
              </Card>
            ) : null}

            {dashboard.state === "error" ? (
              <Card>
                <CardHeader>
                  <CardTitle>Dashboard unavailable</CardTitle>
                  <CardDescription>Could not load study data. Retry or inspect service health.</CardDescription>
                </CardHeader>
              </Card>
            ) : null}

            {dashboard.state === "empty" ? (
              <Card>
                <CardHeader>
                  <CardTitle>No study data yet</CardTitle>
                  <CardDescription>Set a checkride date in TUI or seed data to populate this dashboard.</CardDescription>
                </CardHeader>
              </Card>
            ) : null}

            {dashboard.state !== "loading" && dashboard.state !== "error" && dashboard.state !== "empty" ? (
              <>
                <OverviewCards metrics={dashboard.metrics} />
                <div className="grid gap-6 xl:grid-cols-5">
                  <div className="xl:col-span-3">
                    <ProgressChart data={dashboard.progress} />
                  </div>
                  <div className="xl:col-span-2">
                    <TasksTable data={dashboard.tasks} />
                  </div>
                </div>
              </>
            ) : null}
          </div>
        </SidebarInset>
      </div>
    </SidebarProvider>
  )
}
