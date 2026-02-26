"use client"

import * as React from "react"

import { cn } from "@/lib/utils"

type SidebarContextValue = {
  open: boolean
  toggle: () => void
}

const SidebarContext = React.createContext<SidebarContextValue | null>(null)

export function SidebarProvider({ children }: { children: React.ReactNode }) {
  const [open, setOpen] = React.useState(true)
  const value = React.useMemo(
    () => ({
      open,
      toggle: () => setOpen((prev) => !prev),
    }),
    [open],
  )

  return <SidebarContext.Provider value={value}>{children}</SidebarContext.Provider>
}

export function useSidebar() {
  const context = React.useContext(SidebarContext)
  if (!context) {
    throw new Error("useSidebar must be used within SidebarProvider")
  }
  return context
}

export function SidebarInset({ className, ...props }: React.HTMLAttributes<HTMLElement>) {
  return <main className={cn("min-h-screen flex-1 bg-muted/30", className)} {...props} />
}

export function SidebarTrigger({ className, ...props }: React.ButtonHTMLAttributes<HTMLButtonElement>) {
  const { toggle } = useSidebar()
  return (
    <button
      type="button"
      onClick={toggle}
      className={cn("rounded-md border bg-background px-2 py-1 text-sm", className)}
      {...props}
    >
      Menu
    </button>
  )
}
