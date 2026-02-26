"use client"

import { BookOpen, CalendarClock, ClipboardList, Gauge, GraduationCap, Wallet } from "lucide-react"

import { cn } from "@/lib/utils"
import { useSidebar } from "@/components/ui/sidebar"

const navItems = [
  { label: "Dashboard", icon: Gauge },
  { label: "Study Plan", icon: BookOpen },
  { label: "Progress", icon: GraduationCap },
  { label: "Checklist", icon: ClipboardList },
  { label: "Budget", icon: Wallet },
  { label: "Schedule", icon: CalendarClock },
]

export function AppSidebar({ className }: { className?: string }) {
  const { open } = useSidebar()

  return (
    <aside
      className={cn(
        "sticky top-0 h-screen border-r bg-background transition-all duration-200",
        open ? "w-64" : "w-20",
        className,
      )}
    >
      <div className="flex h-16 items-center border-b px-4">
        <span className="font-semibold tracking-tight">openppl</span>
      </div>

      <nav className="space-y-1 p-3">
        {navItems.map(({ label, icon: Icon }) => (
          <button
            type="button"
            key={label}
            className="flex w-full items-center gap-3 rounded-md px-3 py-2 text-sm text-foreground transition-colors hover:bg-accent"
          >
            <Icon className="h-4 w-4" />
            {open ? <span>{label}</span> : null}
          </button>
        ))}
      </nav>
    </aside>
  )
}
