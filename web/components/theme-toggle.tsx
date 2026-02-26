"use client"

import { Moon, Sun } from "lucide-react"
import { useTheme } from "next-themes"

import { Button } from "@/components/ui/button"

export function ThemeToggle() {
  const { setTheme, resolvedTheme } = useTheme()

  return (
    <div className="inline-flex items-center gap-2 rounded-md border bg-background p-1">
      <Button
        variant={resolvedTheme === "light" ? "default" : "ghost"}
        size="sm"
        onClick={() => setTheme("light")}
        aria-label="Switch to light theme"
      >
        <Sun className="mr-2 h-4 w-4" />
        Light
      </Button>
      <Button
        variant={resolvedTheme === "dark" ? "default" : "ghost"}
        size="sm"
        onClick={() => setTheme("dark")}
        aria-label="Switch to dark theme"
      >
        <Moon className="mr-2 h-4 w-4" />
        Dark
      </Button>
    </div>
  )
}
