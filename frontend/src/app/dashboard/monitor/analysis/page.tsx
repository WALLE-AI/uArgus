"use client";

import { useState } from "react";
import AreaChartCard from "@/components/charts/AreaChartCard";
import BarChartCard from "@/components/charts/BarChartCard";
import { cn } from "@/lib/utils";

const TOPICS = [
  { name: "AI 监管", articles: 142, sentiment: 0.35, trend: "up", color: "#3B82F6" },
  { name: "半导体管制", articles: 98, sentiment: 0.28, trend: "up", color: "#8B5CF6" },
  { name: "地缘政治", articles: 87, sentiment: 0.18, trend: "stable", color: "#EF4444" },
  { name: "能源转型", articles: 72, sentiment: 0.62, trend: "up", color: "#10B981" },
  { name: "量子计算", articles: 54, sentiment: 0.78, trend: "up", color: "#F59E0B" },
  { name: "气候变化", articles: 48, sentiment: 0.41, trend: "down", color: "#06B6D4" },
  { name: "数字货币", articles: 41, sentiment: 0.55, trend: "stable", color: "#EC4899" },
  { name: "医疗 AI", articles: 35, sentiment: 0.82, trend: "up", color: "#14B8A6" },
];

const TOPIC_TIMELINE = Array.from({ length: 14 }, (_, i) => ({
  date: `${i + 1}日`,
  "AI 监管": Math.round(8 + Math.sin(i / 2) * 5 + Math.random() * 3),
  "半导体管制": Math.round(5 + Math.cos(i / 3) * 4 + Math.random() * 2),
  "地缘政治": Math.round(6 + Math.random() * 4),
  "能源转型": Math.round(4 + Math.sin(i / 4) * 3 + Math.random() * 2),
}));

const SUGGESTIONS = [
  { title: "欧盟 AI 法案对中国科技企业的影响分析", reason: "热度上升 +42%，多源佐证充分", score: 92 },
  { title: "半导体管制升级下的国产替代进展", reason: "BREAKING 事件，跨源 23 篇报道", score: 88 },
  { title: "碳交易市场最新动态与投资机会", reason: "COP 会议热度持续，金融关联度高", score: 85 },
  { title: "固态电池量产突破对新能源汽车的影响", reason: "技术突破类选题，正面情感 82%", score: 79 },
  { title: "全球央行数字货币试点进展对比", reason: "长线话题，知识积累充足", score: 74 },
];

export default function AnalysisPage() {
  const [selectedTopic, setSelectedTopic] = useState<string | null>(null);

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold">智能分析</h1>

      {/* Topic bubbles */}
      <div className="rounded-xl border border-border bg-white p-6">
        <h3 className="text-sm font-semibold mb-4">主题聚类气泡</h3>
        <div className="flex flex-wrap items-center justify-center gap-4 py-4">
          {TOPICS.map((t) => {
            const size = 48 + (t.articles / 142) * 72;
            return (
              <button
                key={t.name}
                onClick={() => setSelectedTopic(selectedTopic === t.name ? null : t.name)}
                className={cn(
                  "rounded-full flex items-center justify-center transition-all hover:scale-105",
                  selectedTopic === t.name && "ring-2 ring-offset-2"
                )}
                style={{
                  width: size,
                  height: size,
                  backgroundColor: t.color + "20",
                  color: t.color,
                  outlineColor: selectedTopic === t.name ? t.color : undefined,
                }}
              >
                <div className="text-center">
                  <div className="text-[10px] font-semibold leading-tight">{t.name}</div>
                  <div className="text-[9px] opacity-70">{t.articles}</div>
                </div>
              </button>
            );
          })}
        </div>
        {selectedTopic && (
          <div className="mt-4 rounded-lg border border-border p-4">
            <div className="flex items-center justify-between">
              <span className="text-sm font-semibold">{selectedTopic}</span>
              <span className="text-xs text-muted-foreground">
                {TOPICS.find((t) => t.name === selectedTopic)?.articles} 篇文章
              </span>
            </div>
            <div className="mt-2 flex gap-4 text-xs text-muted-foreground">
              <span>情感指数: {((TOPICS.find((t) => t.name === selectedTopic)?.sentiment ?? 0) * 100).toFixed(0)}% 正面</span>
              <span>趋势: {TOPICS.find((t) => t.name === selectedTopic)?.trend === "up" ? "📈 上升" : TOPICS.find((t) => t.name === selectedTopic)?.trend === "down" ? "📉 下降" : "➡️ 稳定"}</span>
            </div>
          </div>
        )}
      </div>

      {/* Charts row */}
      <div className="grid gap-4 lg:grid-cols-2">
        <AreaChartCard
          title="主题时间线（14天）"
          subtitle="各主题每日文章数量变化"
          data={TOPIC_TIMELINE}
          xKey="date"
          areas={[
            { key: "AI 监管", color: "#3B82F6" },
            { key: "半导体管制", color: "#8B5CF6" },
            { key: "地缘政治", color: "#EF4444" },
            { key: "能源转型", color: "#10B981" },
          ]}
          stacked
          height={260}
        />
        <BarChartCard
          title="主题热度排行"
          subtitle="按文章提及频次排序"
          data={TOPICS.map((t) => ({ name: t.name, articles: t.articles }))}
          xKey="name"
          bars={[{ key: "articles", color: "#3B82F6", name: "文章数" }]}
          layout="vertical"
          height={260}
        />
      </div>

      {/* AI topic suggestions */}
      <div className="rounded-xl border border-border bg-white p-6">
        <h3 className="text-sm font-semibold mb-4">AI 选题推荐</h3>
        <div className="space-y-3">
          {SUGGESTIONS.map((s, i) => (
            <div key={s.title} className="flex items-center gap-4 rounded-lg p-3 hover:bg-muted/50 cursor-pointer transition-colors">
              <span className="flex h-8 w-8 shrink-0 items-center justify-center rounded-full bg-accent-blue/10 text-sm font-bold text-accent-blue">
                {i + 1}
              </span>
              <div className="flex-1 min-w-0">
                <div className="text-sm font-medium">{s.title}</div>
                <div className="text-xs text-muted-foreground mt-0.5">{s.reason}</div>
              </div>
              <div className="shrink-0 text-right">
                <div className="text-sm font-bold text-accent-blue">{s.score}</div>
                <div className="text-[10px] text-muted-foreground">推荐分</div>
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}
