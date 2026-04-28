"use client";

import { useState } from "react";
import { Bell, X, AlertTriangle, Rss, CheckCircle2, Info } from "lucide-react";
import { cn } from "@/lib/utils";

interface Notification {
  id: string;
  type: "alert" | "source" | "task" | "system";
  title: string;
  detail: string;
  time: string;
  read: boolean;
}

const MOCK_NOTIFICATIONS: Notification[] = [
  { id: "1", type: "alert", title: "舆情预警：AI 监管负面情感突增", detail: "负面占比从 18% 升至 42%", time: "5 min ago", read: false },
  { id: "2", type: "source", title: "数据源失败：NASA Earth RSS", detail: "Connection refused，已重试 3 次", time: "32 min ago", read: false },
  { id: "3", type: "task", title: "RSS 全量抓取完成", detail: "采集 284 篇新文章", time: "1 hr ago", read: false },
  { id: "4", type: "alert", title: "红海航运事件升级为 BREAKING", detail: "跨 12 源佐证，影响力 92/100", time: "2 hr ago", read: true },
  { id: "5", type: "system", title: "系统更新完成", detail: "v1.2.3 已部署，新增预测分析模块", time: "5 hr ago", read: true },
  { id: "6", type: "source", title: "数据源降级：arXiv CS.AI", detail: "响应延迟 >5s，自动降低抓取频率", time: "6 hr ago", read: true },
];

const TYPE_ICON = {
  alert: AlertTriangle,
  source: Rss,
  task: CheckCircle2,
  system: Info,
};
const TYPE_COLOR = {
  alert: "text-red-500 bg-red-50",
  source: "text-accent-orange bg-orange-50",
  task: "text-accent-green bg-green-50",
  system: "text-accent-blue bg-blue-50",
};

export default function NotificationCenter() {
  const [open, setOpen] = useState(false);
  const [notifications, setNotifications] = useState(MOCK_NOTIFICATIONS);

  const unreadCount = notifications.filter((n) => !n.read).length;

  const markAllRead = () => {
    setNotifications((prev) => prev.map((n) => ({ ...n, read: true })));
  };

  return (
    <div className="relative">
      <button
        onClick={() => setOpen(!open)}
        className="relative rounded-lg p-2 text-muted-foreground hover:bg-muted transition-colors"
      >
        <Bell className="h-5 w-5" />
        {unreadCount > 0 && (
          <span className="absolute -top-0.5 -right-0.5 flex h-4 w-4 items-center justify-center rounded-full bg-red-500 text-[9px] font-bold text-white">
            {unreadCount}
          </span>
        )}
      </button>

      {open && (
        <>
          <div className="fixed inset-0 z-40" onClick={() => setOpen(false)} />
          <div className="absolute right-0 top-full z-50 mt-2 w-96 rounded-xl border border-border bg-background shadow-2xl">
            {/* Header */}
            <div className="flex items-center justify-between border-b border-border px-4 py-3">
              <div className="flex items-center gap-2">
                <span className="text-sm font-semibold">通知</span>
                {unreadCount > 0 && (
                  <span className="rounded-full bg-red-100 px-1.5 py-0.5 text-[10px] font-medium text-red-600">
                    {unreadCount} 条未读
                  </span>
                )}
              </div>
              <div className="flex items-center gap-2">
                <button
                  onClick={markAllRead}
                  className="text-[11px] text-accent-blue hover:underline"
                >
                  全部已读
                </button>
                <button onClick={() => setOpen(false)} className="text-muted-foreground hover:text-foreground">
                  <X className="h-4 w-4" />
                </button>
              </div>
            </div>

            {/* Notification list */}
            <div className="max-h-96 overflow-y-auto">
              {notifications.map((n) => {
                const Icon = TYPE_ICON[n.type];
                return (
                  <div
                    key={n.id}
                    className={cn(
                      "flex items-start gap-3 px-4 py-3 border-b border-border last:border-0 hover:bg-muted/30 cursor-pointer transition-colors",
                      !n.read && "bg-accent-blue/5"
                    )}
                  >
                    <div className={cn("mt-0.5 rounded-lg p-1.5 shrink-0", TYPE_COLOR[n.type])}>
                      <Icon className="h-3.5 w-3.5" />
                    </div>
                    <div className="flex-1 min-w-0">
                      <div className="flex items-center gap-1.5">
                        <span className={cn("text-xs font-medium", !n.read && "font-semibold")}>{n.title}</span>
                        {!n.read && <span className="h-1.5 w-1.5 rounded-full bg-accent-blue shrink-0" />}
                      </div>
                      <div className="text-[11px] text-muted-foreground mt-0.5">{n.detail}</div>
                      <div className="text-[10px] text-muted-foreground mt-1">{n.time}</div>
                    </div>
                  </div>
                );
              })}
            </div>

            {/* Footer */}
            <div className="border-t border-border px-4 py-2 text-center">
              <button className="text-[11px] text-accent-blue hover:underline">
                查看全部通知
              </button>
            </div>
          </div>
        </>
      )}
    </div>
  );
}
