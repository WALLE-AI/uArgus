"use client";

import { useEffect, useState, useCallback } from "react";
import { Search, Rss, BarChart3, BookOpen, FileText, Settings } from "lucide-react";
import { cn } from "@/lib/utils";
import { useAppStore } from "@/lib/store";

interface CommandItem {
  id: string;
  label: string;
  group: string;
  icon: React.ElementType;
  href: string;
}

const COMMANDS: CommandItem[] = [
  { id: "dashboard", label: "Dashboard 总览", group: "导航", icon: BarChart3, href: "/dashboard" },
  { id: "feed", label: "AI 内容流", group: "导航", icon: Rss, href: "/dashboard/content/feed" },
  { id: "sources", label: "数据源管理", group: "导航", icon: Rss, href: "/dashboard/monitor/sources" },
  { id: "analysis", label: "智能分析", group: "导航", icon: BarChart3, href: "/dashboard/monitor/analysis" },
  { id: "sentiment", label: "舆情监控", group: "导航", icon: BarChart3, href: "/dashboard/monitor/sentiment" },
  { id: "events", label: "事件追踪", group: "导航", icon: BarChart3, href: "/dashboard/monitor/events" },
  { id: "studio", label: "内容中心", group: "导航", icon: FileText, href: "/dashboard/content/studio" },
  { id: "knowledge", label: "知识沉淀", group: "导航", icon: BookOpen, href: "/dashboard/content/knowledge" },
  { id: "distribution", label: "内容分发", group: "导航", icon: FileText, href: "/dashboard/content/distribution" },
  { id: "agents", label: "智能体观测", group: "导航", icon: BarChart3, href: "/dashboard/observe/agents" },
  { id: "channels", label: "渠道运营", group: "导航", icon: BarChart3, href: "/dashboard/observe/channels" },
  { id: "settings", label: "用户设置", group: "系统", icon: Settings, href: "/dashboard/settings" },
];

export default function CommandPalette() {
  const { commandPaletteOpen, setCommandPaletteOpen } = useAppStore();
  const [query, setQuery] = useState("");

  const close = useCallback(() => {
    setCommandPaletteOpen(false);
    setQuery("");
  }, [setCommandPaletteOpen]);

  useEffect(() => {
    const handleKey = (e: KeyboardEvent) => {
      if ((e.metaKey || e.ctrlKey) && e.key === "k") {
        e.preventDefault();
        setCommandPaletteOpen(!commandPaletteOpen);
      }
      if (e.key === "Escape" && commandPaletteOpen) {
        close();
      }
    };
    window.addEventListener("keydown", handleKey);
    return () => window.removeEventListener("keydown", handleKey);
  }, [commandPaletteOpen, setCommandPaletteOpen, close]);

  if (!commandPaletteOpen) return null;

  const filtered = COMMANDS.filter(
    (cmd) =>
      cmd.label.toLowerCase().includes(query.toLowerCase()) ||
      cmd.group.toLowerCase().includes(query.toLowerCase())
  );

  const groups = filtered.reduce<Record<string, CommandItem[]>>((acc, cmd) => {
    if (!acc[cmd.group]) acc[cmd.group] = [];
    acc[cmd.group].push(cmd);
    return acc;
  }, {});

  return (
    <div className="fixed inset-0 z-[100] flex items-start justify-center pt-[20vh]">
      {/* Backdrop */}
      <div className="absolute inset-0 bg-black/50 backdrop-blur-sm" onClick={close} />

      {/* Palette */}
      <div className="relative z-10 w-full max-w-lg rounded-xl border border-border bg-background shadow-2xl">
        {/* Search input */}
        <div className="flex items-center gap-3 border-b border-border px-4 py-3">
          <Search className="h-4 w-4 text-muted-foreground" />
          <input
            type="text"
            autoFocus
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            placeholder="搜索页面、功能..."
            className="flex-1 bg-transparent text-sm outline-none placeholder:text-muted-foreground"
          />
          <kbd className="rounded border border-border px-1.5 py-0.5 text-[10px] font-mono text-muted-foreground">
            ESC
          </kbd>
        </div>

        {/* Results */}
        <div className="max-h-80 overflow-y-auto p-2">
          {Object.entries(groups).map(([group, items]) => (
            <div key={group}>
              <div className="px-2 py-1.5 text-[10px] font-medium uppercase tracking-wider text-muted-foreground">
                {group}
              </div>
              {items.map((item) => (
                <a
                  key={item.id}
                  href={item.href}
                  onClick={close}
                  className="flex items-center gap-3 rounded-lg px-3 py-2 text-sm text-foreground hover:bg-muted transition-colors"
                >
                  <item.icon className="h-4 w-4 text-muted-foreground" />
                  {item.label}
                </a>
              ))}
            </div>
          ))}
          {filtered.length === 0 && (
            <div className="px-4 py-8 text-center text-sm text-muted-foreground">
              没有找到匹配结果
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
