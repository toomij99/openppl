import { AppSidebar } from "@/components/app-sidebar"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { SidebarInset, SidebarProvider, SidebarTrigger } from "@/components/ui/sidebar"

export default function DashboardPage() {
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
            <SidebarTrigger />
          </header>

          <div className="grid gap-4 p-6 md:grid-cols-2 xl:grid-cols-4">
            <Card>
              <CardHeader className="pb-2">
                <CardDescription>Study Completion</CardDescription>
                <CardTitle className="text-3xl">64%</CardTitle>
              </CardHeader>
              <CardContent className="text-sm text-muted-foreground">58 of 91 tasks complete</CardContent>
            </Card>
            <Card>
              <CardHeader className="pb-2">
                <CardDescription>Days to Checkride</CardDescription>
                <CardTitle className="text-3xl">42</CardTitle>
              </CardHeader>
              <CardContent className="text-sm text-muted-foreground">Target date on track</CardContent>
            </Card>
            <Card>
              <CardHeader className="pb-2">
                <CardDescription>Weekly Hours</CardDescription>
                <CardTitle className="text-3xl">11.5</CardTitle>
              </CardHeader>
              <CardContent className="text-sm text-muted-foreground">1.3h above target pace</CardContent>
            </Card>
            <Card>
              <CardHeader className="pb-2">
                <CardDescription>Budget Burn</CardDescription>
                <CardTitle className="text-3xl">$9.8k</CardTitle>
              </CardHeader>
              <CardContent className="text-sm text-muted-foreground">Under projected spend by 6%</CardContent>
            </Card>
          </div>
        </SidebarInset>
      </div>
    </SidebarProvider>
  )
}
