import { BookCheck, CalendarClock, Target, Wallet } from "lucide-react"

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"

export interface OverviewMetrics {
  completionPercent: number
  completedTasks: number
  totalTasks: number
  daysToCheckride: number
  weeklyHours: number
  budgetSpent: number
}

export function OverviewCards({ metrics }: { metrics: OverviewMetrics }) {
  const cards = [
    {
      title: "Study Completion",
      value: `${metrics.completionPercent}%`,
      subtext: `${metrics.completedTasks} of ${metrics.totalTasks} tasks complete`,
      icon: BookCheck,
    },
    {
      title: "Days to Checkride",
      value: `${metrics.daysToCheckride}`,
      subtext: "Target date still achievable",
      icon: CalendarClock,
    },
    {
      title: "Weekly Hours",
      value: `${metrics.weeklyHours.toFixed(1)}h`,
      subtext: "Including chair-flying sessions",
      icon: Target,
    },
    {
      title: "Budget Spent",
      value: `$${metrics.budgetSpent.toLocaleString()}`,
      subtext: "Fuel + CFI + living costs",
      icon: Wallet,
    },
  ]

  return (
    <section className="grid gap-4 md:grid-cols-2 xl:grid-cols-4">
      {cards.map((card) => (
        <Card key={card.title}>
          <CardHeader className="pb-2">
            <CardDescription className="flex items-center justify-between gap-2">
              <span>{card.title}</span>
              <card.icon className="h-4 w-4 text-muted-foreground" />
            </CardDescription>
            <CardTitle className="text-3xl">{card.value}</CardTitle>
          </CardHeader>
          <CardContent className="text-sm text-muted-foreground">{card.subtext}</CardContent>
        </Card>
      ))}
    </section>
  )
}
