"use client"

import { CartesianGrid, Line, LineChart, ResponsiveContainer, Tooltip, XAxis, YAxis } from "recharts"

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { ChartContainer, ChartLegend } from "@/components/ui/chart"

export interface ProgressPoint {
  week: string
  completed: number
  target: number
}

export function ProgressChart({ data }: { data: ProgressPoint[] }) {
  return (
    <Card>
      <CardHeader>
        <CardTitle>Progress vs Target</CardTitle>
        <CardDescription>Weekly completion trend for current study block</CardDescription>
      </CardHeader>
      <CardContent>
        <ChartContainer>
          <ResponsiveContainer width="100%" height={260}>
            <LineChart data={data} accessibilityLayer>
              <CartesianGrid strokeDasharray="3 3" className="stroke-border" />
              <XAxis dataKey="week" stroke="hsl(var(--muted-foreground))" />
              <YAxis stroke="hsl(var(--muted-foreground))" />
              <Tooltip
                contentStyle={{
                  backgroundColor: "hsl(var(--card))",
                  borderColor: "hsl(var(--border))",
                  color: "hsl(var(--card-foreground))",
                }}
              />
              <Line type="monotone" dataKey="completed" stroke="var(--color-progress)" strokeWidth={3} dot={false} />
              <Line type="monotone" dataKey="target" stroke="var(--color-target)" strokeWidth={2} strokeDasharray="8 6" dot={false} />
            </LineChart>
          </ResponsiveContainer>
        </ChartContainer>
        <ChartLegend
          items={[
            { label: "Completed", color: "var(--color-progress)" },
            { label: "Target", color: "var(--color-target)" },
          ]}
        />
      </CardContent>
    </Card>
  )
}
