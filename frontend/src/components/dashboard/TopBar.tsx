"use client";

import { usePathname } from "next/navigation";
import { Search, BotMessageSquare, Sun, Moon } from "lucide-react";
import NotificationCenter from "@/components/dashboard/NotificationCenter";
import { Eye } from "lucide-react";
import { cn } from "@/lib/utils";
import { useAppStore } from "@/lib/store";

const BREADCRUMB_MAP: Record<string, string> = {
  "/dashboard": "看板",
  "/dashboard/monitor/sources": "监控 / 数据源管理",
  "/dashboard/monitor/categories": "监控 / 领域分类",
  "/dashboard/monitor/health": "监控 / 源健康状态",
  "/dashboard/monitor/tasks": "监控 / 采集任务",
  "/dashboard/monitor/feed": "监控 / 信息流",
  "/dashboard/monitor/analysis": "监控 / 智能分析",
  "/dashboard/monitor/sentiment": "监控 / 舆情监控",
  "/dashboard/monitor/prediction": "监控 / 预测分析",
  "/dashboard/monitor/events": "监控 / 事件追踪",
  "/dashboard/content/studio": "内容 / 内容中心",
  "/dashboard/content/knowledge": "内容 / 知识沉淀",
  "/dashboard/content/distribution": "内容 / 内容分发",
  "/dashboard/observe/channels": "观测 / 渠道运营",
  "/dashboard/observe/agents": "观测 / 智能体观测",
  "/dashboard/settings": "用户设置",
};

export default function TopBar() {
  const pathname = usePathname();
  const { setCommandPaletteOpen, toggleAssistant, assistantOpen, theme, toggleTheme } = useAppStore();

  const breadcrumb = BREADCRUMB_MAP[pathname] ?? "Dashboard";

  return (
    <header className="flex h-14 shrink-0 items-center justify-between border-b border-border bg-background px-4">
      {/* Left: Logo + Breadcrumb */}
      <div className="flex items-center gap-4">
        <div className="flex items-center gap-2 text-sm font-semibold">
          <Eye className="h-5 w-5 text-accent-blue" />
          <span className="hidden sm:inline">uArgus</span>
        </div>
        <div className="hidden items-center gap-1.5 text-xs text-muted-foreground sm:flex">
          <span>/</span>
          <span>{breadcrumb}</span>
        </div>
      </div>

      {/* Right: Search + Notifications + AI + Avatar */}
      <div className="flex items-center gap-2">
        {/* Search trigger */}
        <button
          onClick={() => setCommandPaletteOpen(true)}
          className="flex items-center gap-2 rounded-lg border border-border bg-muted/50 px-3 py-1.5 text-xs text-muted-foreground transition-colors hover:bg-muted"
        >
          <Search className="h-3.5 w-3.5" />
          <span className="hidden sm:inline">搜索...</span>
          <kbd className="hidden rounded border border-border bg-white px-1.5 py-0.5 text-[10px] font-mono sm:inline">
            ⌘K
          </kbd>
        </button>

        {/* Notifications */}
        <NotificationCenter />

        {/* Dark mode toggle */}
        <button
          onClick={toggleTheme}
          className="rounded-lg p-2 text-muted-foreground hover:bg-muted hover:text-foreground transition-colors"
          title={theme === "light" ? "切换暗色模式" : "切换亮色模式"}
        >
          {theme === "light" ? <Moon className="h-4 w-4" /> : <Sun className="h-4 w-4" />}
        </button>

        {/* AI Assistant toggle */}
        <button
          onClick={toggleAssistant}
          className={cn(
            "rounded-lg p-2 transition-colors",
            assistantOpen
              ? "bg-accent-blue/10 text-accent-blue"
              : "text-muted-foreground hover:bg-muted hover:text-foreground"
          )}
          title="AI 助手"
        >
          <BotMessageSquare className="h-4 w-4" />
        </button>

        {/* User avatar */}
        <div className="flex h-8 w-8 items-center justify-center rounded-full bg-accent-blue/20 text-xs font-medium text-accent-blue">
          U
        </div>
      </div>
    </header>
  );
}
