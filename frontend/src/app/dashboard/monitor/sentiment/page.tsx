"use client";

import LineChartCard from "@/components/charts/LineChartCard";
import BarChartCard from "@/components/charts/BarChartCard";

const SENTIMENT_TREND = Array.from({ length: 30 }, (_, i) => ({
  date: `${i + 1}日`,
  正面: Math.round(35 + Math.sin(i / 4) * 12 + Math.random() * 5),
  中性: Math.round(30 + Math.cos(i / 5) * 8 + Math.random() * 5),
  负面: Math.round(20 + Math.sin(i / 3) * 10 + Math.random() * 5),
}));

const SOURCE_SENTIMENT = [
  { source: "Reuters", 正面: 45, 中性: 35, 负面: 20 },
  { source: "Bloomberg", 正面: 38, 中性: 40, 负面: 22 },
  { source: "arXiv", 正面: 62, 中性: 30, 负面: 8 },
  { source: "TechCrunch", 正面: 52, 中性: 28, 负面: 20 },
  { source: "Al Jazeera", 正面: 22, 中性: 33, 负面: 45 },
  { source: "Nature", 正面: 58, 中性: 32, 负面: 10 },
];

const ALERTS = [
  { time: "10:15", topic: "AI 监管", type: "负面突增", detail: "负面情感占比从 18% 升至 42%，触发阈值", severity: "高" },
  { time: "09:32", topic: "红海航运", type: "负面突增", detail: "负面情感占比从 35% 升至 68%，跨 12 源", severity: "紧急" },
  { time: "昨日 22:10", topic: "半导体管制", type: "情感反转", detail: "正面情感从 55% 下降至 22%，跨 8 源", severity: "中" },
  { time: "昨日 14:50", topic: "量子计算", type: "正面激增", detail: "正面情感从 60% 升至 92%，突破性进展", severity: "低" },
];

const SEV_STYLE: Record<string, string> = {
  紧急: "bg-red-100 text-red-700",
  高: "bg-orange-100 text-orange-700",
  中: "bg-yellow-100 text-yellow-700",
  低: "bg-muted text-muted-foreground",
};

export default function SentimentPage() {
  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold">舆情监控</h1>

      {/* Metric cards */}
      <div className="grid gap-4 sm:grid-cols-3">
        {[
          { label: "正面", value: "42%", color: "text-accent-green", bg: "bg-accent-green/10" },
          { label: "中性", value: "35%", color: "text-accent-blue", bg: "bg-accent-blue/10" },
          { label: "负面", value: "23%", color: "text-red-500", bg: "bg-red-50" },
        ].map((s) => (
          <div key={s.label} className={`rounded-xl border border-border bg-white p-5`}>
            <div className={`text-3xl font-bold ${s.color}`}>{s.value}</div>
            <div className="text-sm text-muted-foreground">{s.label}情感</div>
          </div>
        ))}
      </div>

      {/* Trend + Source charts */}
      <div className="grid gap-4 lg:grid-cols-2">
        <LineChartCard
          title="情感趋势（30天）"
          subtitle="正面/中性/负面随时间变化"
          data={SENTIMENT_TREND}
          xKey="date"
          lines={[
            { key: "正面", color: "#10B981", name: "正面" },
            { key: "中性", color: "#3B82F6", name: "中性" },
            { key: "负面", color: "#EF4444", name: "负面" },
          ]}
          height={280}
        />
        <BarChartCard
          title="各源情感分布"
          subtitle="按数据源分类的正面/中性/负面占比"
          data={SOURCE_SENTIMENT}
          xKey="source"
          bars={[
            { key: "正面", color: "#10B981", name: "正面" },
            { key: "中性", color: "#3B82F6", name: "中性" },
            { key: "负面", color: "#EF4444", name: "负面" },
          ]}
          stacked
          height={280}
        />
      </div>

      {/* Alert config + Recent alerts */}
      <div className="grid gap-4 lg:grid-cols-3">
        {/* Alert thresholds */}
        <div className="rounded-xl border border-border bg-white p-6">
          <h3 className="text-sm font-semibold mb-4">预警阈值配置</h3>
          <div className="space-y-4">
            {[
              { label: "负面情感占比", value: 40, unit: "%" },
              { label: "情感变化幅度", value: 25, unit: "%" },
              { label: "最低源数量", value: 3, unit: "源" },
            ].map((cfg) => (
              <div key={cfg.label}>
                <div className="flex items-center justify-between mb-1.5">
                  <span className="text-xs text-muted-foreground">{cfg.label}</span>
                  <span className="text-xs font-medium">{cfg.value}{cfg.unit}</span>
                </div>
                <input
                  type="range"
                  min={0}
                  max={100}
                  defaultValue={cfg.value}
                  className="w-full h-1.5 rounded-full appearance-none bg-muted accent-accent-blue"
                />
              </div>
            ))}
            <button className="w-full rounded-lg bg-accent-blue px-3 py-2 text-xs font-medium text-white hover:bg-accent-blue/90 transition-colors">
              保存配置
            </button>
          </div>
        </div>

        {/* Recent alerts */}
        <div className="lg:col-span-2 rounded-xl border border-border bg-white p-6">
          <h3 className="text-sm font-semibold mb-4">近期预警</h3>
          <div className="space-y-3">
            {ALERTS.map((a, i) => (
              <div key={i} className="flex items-start gap-3 rounded-lg p-3 hover:bg-muted/30 transition-colors cursor-pointer">
                <span className={`shrink-0 mt-0.5 rounded px-1.5 py-0.5 text-[10px] font-medium ${SEV_STYLE[a.severity]}`}>
                  {a.severity}
                </span>
                <div className="flex-1 min-w-0">
                  <div className="flex items-center gap-2">
                    <span className="text-sm font-medium">{a.topic}</span>
                    <span className="rounded bg-muted px-1.5 py-0.5 text-[10px] text-muted-foreground">{a.type}</span>
                  </div>
                  <div className="text-xs text-muted-foreground mt-0.5">{a.detail}</div>
                </div>
                <span className="shrink-0 text-[11px] text-muted-foreground">{a.time}</span>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
}
