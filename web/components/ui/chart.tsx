import * as React from "react"

import { cn } from "@/lib/utils"

export function ChartContainer({ className, children }: { className?: string; children: React.ReactNode }) {
  return <div className={cn("min-h-[240px] w-full", className)}>{children}</div>
}

export function ChartLegend({ items }: { items: Array<{ label: string; color: string }> }) {
  return (
    <div className="mt-4 flex flex-wrap gap-4 text-xs text-muted-foreground">
      {items.map((item) => (
        <div key={item.label} className="inline-flex items-center gap-2">
          <span className="h-2.5 w-2.5 rounded-full" style={{ backgroundColor: item.color }} />
          <span>{item.label}</span>
        </div>
      ))}
    </div>
  )
}
