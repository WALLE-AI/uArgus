"use client";

import { useState } from "react";
import { Link2, Unlink, BrainCircuit } from "lucide-react";
import LineChartCard from "@/components/charts/LineChartCard";
import { cn } from "@/lib/utils";

const PLATFORMS = [
  { name: "微信公众号", status: "已连接", followers: "12.3K", posts: 142, views: "45.2K", color: "#07C160" },
  { name: "知乎", status: "已连接", followers: "8.7K", posts: 89, views: "28.1K", color: "#0066FF" },
  { name: "小红书", status: "未连接", followers: "—", posts: 0, views: "—", color: "#FE2C55" },
  { name: "抖音", status: "未连接", followers: "—", posts: 0, views: "—", color: "#000000" },
];

const CALENDAR_POSTS: Record<number, { title: string; platform: string }[]> = {
  2: [{ title: "AI 监管解读", platform: "微信" }],
  5: [{ title: "半导体周报", platform: "知乎" }],
  8: [{ title: "量子计算速报", platform: "微信" }, { title: "量子科普", platform: "知乎" }],
  11: [{ title: "能源转型分析", platform: "微信" }],
  14: [{ title: "周度热点综述", platform: "知乎" }],
  17: [{ title: "COP 会议跟踪", platform: "微信" }],
  20: [{ title: "AI 产业链图谱", platform: "微信" }],
  23: [{ title: "月度报告", platform: "知乎" }],
};

const PERF_DATA = Array.from({ length: 14 }, (_, i) => ({
  date: `${i + 1}日`,
  微信: Math.round(1200 + Math.sin(i / 3) * 500 + Math.random() * 300),
  知乎: Math.round(800 + Math.cos(i / 4) * 400 + Math.random() * 200),
}));

const REWRITE_TEMPLATES = [
  { platform: "微信公众号", style: "正式、专业、深度分析风格" },
  { platform: "知乎", style: "理性讨论、引经据典风格" },
  { platform: "小红书", style: "简洁、有趣、带表情包风格" },
  { platform: "抖音", style: "口语化、节奏感强、短句风格" },
];

export default function DistributionPage() {
  const [selectedDay, setSelectedDay] = useState<number | null>(null);

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold">内容分发</h1>

      {/* Platform cards */}
      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
        {PLATFORMS.map((p) => (
          <div key={p.name} className="rounded-xl border border-border bg-white p-5">
            <div className="flex items-center justify-between">
              <div className="text-sm font-semibold">{p.name}</div>
              <div className="h-3 w-3 rounded-full" style={{ backgroundColor: p.color, opacity: p.status === "已连接" ? 1 : 0.2 }} />
            </div>
            <div className={cn("text-xs mt-1 flex items-center gap-1", p.status === "已连接" ? "text-accent-green" : "text-muted-foreground")}>
              {p.status === "已连接" ? <Link2 className="h-3 w-3" /> : <Unlink className="h-3 w-3" />}
              {p.status}
            </div>
            <div className="mt-3 space-y-1 text-xs text-muted-foreground">
              <div className="flex justify-between"><span>粉丝</span><span className="font-medium text-foreground">{p.followers}</span></div>
              <div className="flex justify-between"><span>发布</span><span className="font-medium text-foreground">{p.posts}</span></div>
              <div className="flex justify-between"><span>阅读</span><span className="font-medium text-foreground">{p.views}</span></div>
            </div>
            {p.status !== "已连接" && (
              <button className="mt-3 w-full rounded-lg border border-border px-2 py-1.5 text-[11px] text-muted-foreground hover:bg-muted transition-colors">
                连接平台
              </button>
            )}
          </div>
        ))}
      </div>

      <div className="grid gap-4 lg:grid-cols-3">
        {/* Calendar */}
        <div className="lg:col-span-2 rounded-xl border border-border bg-white p-6">
          <h3 className="text-sm font-semibold mb-4">内容日历</h3>
          <div className="grid grid-cols-7 gap-2">
            {["一", "二", "三", "四", "五", "六", "日"].map((d) => (
              <div key={d} className="text-center text-xs font-medium text-muted-foreground py-1">{d}</div>
            ))}
            {Array.from({ length: 28 }, (_, i) => {
              const posts = CALENDAR_POSTS[i];
              const isSelected = selectedDay === i;
              return (
                <div
                  key={i}
                  onClick={() => setSelectedDay(isSelected ? null : i)}
                  className={cn(
                    "rounded-lg border p-1.5 text-center cursor-pointer transition-all min-h-[52px]",
                    isSelected ? "border-accent-blue bg-accent-blue/5" : "border-border hover:bg-muted/50",
                    posts && "ring-1 ring-accent-blue/20"
                  )}
                >
                  <span className="text-xs text-muted-foreground">{i + 1}</span>
                  {posts && (
                    <div className="flex justify-center gap-0.5 mt-0.5">
                      {posts.map((_, j) => (
                        <div key={j} className="h-1 w-1 rounded-full bg-accent-blue" />
                      ))}
                    </div>
                  )}
                </div>
              );
            })}
          </div>
          {selectedDay !== null && CALENDAR_POSTS[selectedDay] && (
            <div className="mt-4 rounded-lg border border-border p-3 space-y-2">
              <div className="text-xs font-medium">{selectedDay + 1}日 排期内容</div>
              {CALENDAR_POSTS[selectedDay].map((p, i) => (
                <div key={i} className="flex items-center justify-between rounded-lg bg-muted/30 px-3 py-2">
                  <span className="text-xs">{p.title}</span>
                  <span className="rounded bg-accent-blue/10 px-1.5 py-0.5 text-[10px] text-accent-blue">{p.platform}</span>
                </div>
              ))}
            </div>
          )}
        </div>

        {/* AI rewrite */}
        <div className="rounded-xl border border-border bg-white p-5">
          <div className="flex items-center gap-1.5 text-sm font-semibold mb-4">
            <BrainCircuit className="h-4 w-4 text-accent-blue" /> AI 多平台改写
          </div>
          <p className="text-xs text-muted-foreground mb-3">
            选择一篇文章，AI 自动生成适配各平台的版本
          </p>
          <div className="space-y-2">
            {REWRITE_TEMPLATES.map((t) => (
              <button key={t.platform} className="w-full text-left rounded-lg border border-border p-3 hover:bg-muted/50 transition-colors">
                <div className="text-xs font-medium">{t.platform}</div>
                <div className="text-[10px] text-muted-foreground mt-0.5">{t.style}</div>
              </button>
            ))}
          </div>
        </div>
      </div>

      {/* Performance chart */}
      <LineChartCard
        title="分发效果趋势（14天）"
        subtitle="各平台阅读量对比"
        data={PERF_DATA}
        xKey="date"
        lines={[
          { key: "微信", color: "#07C160", name: "微信公众号" },
          { key: "知乎", color: "#0066FF", name: "知乎" },
        ]}
        height={260}
      />
    </div>
  );
}
