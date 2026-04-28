"use client";

import { useState } from "react";
import Link from "next/link";
import { usePathname } from "next/navigation";
import {
  LayoutDashboard,
  Radar,
  FileText,
  Eye,
  ChevronRight,
  PanelLeftClose,
  PanelLeft,
  Database,
  FolderTree,
  HeartPulse,
  Clock,
  Rss,
  BrainCircuit,
  TrendingUp,
  BarChart3,
  Crosshair,
  PenTool,
  BookOpen,
  Share2,
  Megaphone,
  Bot,
  Settings,
  User,
  LogOut,
  Moon,
  Sun,
  Key,
} from "lucide-react";
import { cn } from "@/lib/utils";
import { useAppStore } from "@/lib/store";

interface NavItem {
  label: string;
  href: string;
  icon: React.ElementType;
}

interface NavGroup {
  label: string;
  icon: React.ElementType;
  items: NavItem[];
}

const NAV_GROUPS: NavGroup[] = [
  {
    label: "看板",
    icon: LayoutDashboard,
    items: [
      { label: "Dashboard 总览", href: "/dashboard", icon: LayoutDashboard },
    ],
  },
  {
    label: "监控",
    icon: Radar,
    items: [
      { label: "数据源管理", href: "/dashboard/monitor/sources", icon: Database },
      { label: "领域分类", href: "/dashboard/monitor/categories", icon: FolderTree },
      { label: "源健康状态", href: "/dashboard/monitor/health", icon: HeartPulse },
      { label: "采集任务", href: "/dashboard/monitor/tasks", icon: Clock },
      { label: "智能分析", href: "/dashboard/monitor/analysis", icon: BrainCircuit },
      { label: "舆情监控", href: "/dashboard/monitor/sentiment", icon: TrendingUp },
      { label: "预测分析", href: "/dashboard/monitor/prediction", icon: BarChart3 },
      { label: "事件追踪", href: "/dashboard/monitor/events", icon: Crosshair },
    ],
  },
  {
    label: "内容",
    icon: FileText,
    items: [
      { label: "AI 内容流", href: "/dashboard/content/feed", icon: Rss },
      { label: "内容中心", href: "/dashboard/content/studio", icon: PenTool },
      { label: "知识沉淀", href: "/dashboard/content/knowledge", icon: BookOpen },
      { label: "内容分发", href: "/dashboard/content/distribution", icon: Share2 },
    ],
  },
  {
    label: "观测",
    icon: Eye,
    items: [
      { label: "渠道运营", href: "/dashboard/observe/channels", icon: Megaphone },
      { label: "智能体观测", href: "/dashboard/observe/agents", icon: Bot },
    ],
  },
];

export default function LeftNav() {
  const pathname = usePathname();
  const { sidebarCollapsed, toggleSidebar, theme, toggleTheme } = useAppStore();
  const [expandedGroup, setExpandedGroup] = useState<string | null>("监控");
  const [userMenuOpen, setUserMenuOpen] = useState(false);

  const toggleGroup = (label: string) => {
    if (sidebarCollapsed) return;
    setExpandedGroup((prev) => (prev === label ? null : label));
  };

  return (
    <aside
      className={cn(
        "flex h-full flex-col border-r border-border bg-muted/30 transition-all duration-300",
        sidebarCollapsed ? "w-[60px]" : "w-[240px]"
      )}
    >
      {/* Toggle button */}
      <div className="flex items-center justify-end px-3 py-2">
        <button
          onClick={toggleSidebar}
          className="rounded-md p-1.5 text-muted-foreground hover:bg-muted hover:text-foreground transition-colors"
          title={sidebarCollapsed ? "展开侧边栏" : "收起侧边栏"}
        >
          {sidebarCollapsed ? (
            <PanelLeft className="h-4 w-4" />
          ) : (
            <PanelLeftClose className="h-4 w-4" />
          )}
        </button>
      </div>

      {/* Navigation groups */}
      <nav className="flex-1 overflow-y-auto px-2 pb-2">
        {NAV_GROUPS.map((group) => {
          const isExpanded = expandedGroup === group.label;
          const isGroupActive = group.items.some((item) => pathname === item.href);

          return (
            <div key={group.label} className="mb-1">
              {/* Group header */}
              <button
                onClick={() => toggleGroup(group.label)}
                className={cn(
                  "flex w-full items-center gap-2 rounded-lg px-3 py-2 text-sm font-medium transition-colors",
                  isGroupActive
                    ? "text-accent-blue"
                    : "text-muted-foreground hover:bg-muted hover:text-foreground"
                )}
                title={sidebarCollapsed ? group.label : undefined}
              >
                <group.icon className="h-4 w-4 shrink-0" />
                {!sidebarCollapsed && (
                  <>
                    <span className="flex-1 text-left">{group.label}</span>
                    <ChevronRight
                      className={cn(
                        "h-3 w-3 transition-transform",
                        isExpanded && "rotate-90"
                      )}
                    />
                  </>
                )}
              </button>

              {/* Group items */}
              {!sidebarCollapsed && isExpanded && (
                <div className="ml-3 mt-0.5 space-y-0.5 border-l border-border pl-3">
                  {group.items.map((item) => {
                    const isActive = pathname === item.href;
                    return (
                      <Link
                        key={item.href}
                        href={item.href}
                        className={cn(
                          "flex items-center gap-2 rounded-md px-2 py-1.5 text-xs transition-colors",
                          isActive
                            ? "bg-accent-blue/10 text-accent-blue font-medium"
                            : "text-muted-foreground hover:bg-muted hover:text-foreground"
                        )}
                      >
                        <item.icon className="h-3.5 w-3.5 shrink-0" />
                        <span className="truncate">{item.label}</span>
                      </Link>
                    );
                  })}
                </div>
              )}
            </div>
          );
        })}
      </nav>

      {/* Bottom user section */}
      <div className="relative border-t border-border p-2">
        <button
          onClick={() => setUserMenuOpen(!userMenuOpen)}
          className="flex w-full items-center gap-2 rounded-lg px-3 py-2 text-sm text-muted-foreground hover:bg-muted hover:text-foreground transition-colors"
          title={sidebarCollapsed ? "用户菜单" : undefined}
        >
          <div className="flex h-6 w-6 shrink-0 items-center justify-center rounded-full bg-accent-blue/20 text-xs font-medium text-accent-blue">
            U
          </div>
          {!sidebarCollapsed && (
            <span className="flex-1 truncate text-left text-xs">User</span>
          )}
        </button>

        {/* User popup menu */}
        {userMenuOpen && !sidebarCollapsed && (
          <div className="absolute bottom-full left-2 right-2 mb-1 rounded-lg border border-border bg-background p-1 shadow-lg">
            {[
              { label: "个人设置", icon: Settings, href: "/dashboard/settings" },
              { label: "API Key 管理", icon: Key, href: "#" },
            ].map((item) => (
              <Link
                key={item.label}
                href={item.href}
                className="flex items-center gap-2 rounded-md px-3 py-2 text-xs text-muted-foreground hover:bg-muted hover:text-foreground"
                onClick={() => setUserMenuOpen(false)}
              >
                <item.icon className="h-3.5 w-3.5" />
                {item.label}
              </Link>
            ))}
            <div className="my-1 border-t border-border" />
            <button
              onClick={() => { toggleTheme(); }}
              className="flex w-full items-center gap-2 rounded-md px-3 py-2 text-xs text-muted-foreground hover:bg-muted hover:text-foreground"
            >
              {theme === "light" ? <Moon className="h-3.5 w-3.5" /> : <Sun className="h-3.5 w-3.5" />}
              {theme === "light" ? "深色模式" : "浅色模式"}
            </button>
            <button className="flex w-full items-center gap-2 rounded-md px-3 py-2 text-xs text-red-500 hover:bg-muted">
              <LogOut className="h-3.5 w-3.5" />
              登出
            </button>
          </div>
        )}
      </div>
    </aside>
  );
}
