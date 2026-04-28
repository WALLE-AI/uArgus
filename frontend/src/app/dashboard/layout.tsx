"use client";

import LeftNav from "@/components/dashboard/LeftNav";
import TopBar from "@/components/dashboard/TopBar";
import AiAssistantPanel from "@/components/dashboard/AiAssistantPanel";
import CommandPalette from "@/components/dashboard/CommandPalette";
import { useAppStore } from "@/lib/store";
import { cn } from "@/lib/utils";

export default function DashboardLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const { theme } = useAppStore();

  return (
    <div className={cn("flex h-screen flex-col overflow-hidden", theme === "dark" && "dark")}>
      <TopBar />
      <div className="flex flex-1 overflow-hidden">
        <LeftNav />
        <main className="flex-1 overflow-y-auto bg-muted/20 p-6">
          {children}
        </main>
        <AiAssistantPanel />
      </div>
      <CommandPalette />
    </div>
  );
}
