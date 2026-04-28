"use client";

import { Activity, Rss, TrendingUp, AlertTriangle } from "lucide-react";
import AreaChartCard from "@/components/charts/AreaChartCard";
import BarChartCard from "@/components/charts/BarChartCard";
import PieChartCard from "@/components/charts/PieChartCard";

const METRICS = [
  { label: "活跃数据源", value: "3,247", change: "+12", icon: Rss, color: "text-accent-blue", bg: "bg-accent-blue/10" },
  { label: "今日文章", value: "1,842", change: "+284", icon: Activity, color: "text-accent-green", bg: "bg-accent-green/10" },
  { label: "热点话题", value: "23", change: "+5", icon: TrendingUp, color: "text-accent-purple", bg: "bg-accent-purple/10" },
  { label: "预警通知", value: "7", change: "+2", icon: AlertTriangle, color: "text-accent-orange", bg: "bg-accent-orange/10" },
];

const ARTICLE_TREND = Array.from({ length: 24 }, (_, i) => ({
  hour: `${i}:00`,
  articles: Math.round(40 + Math.sin(i / 3) * 30 + Math.random() * 20),
  alerts: Math.round(2 + Math.random() * 5),
}));

const TOPIC_DATA = [
  { topic: "AI 监管", count: 142 },
  { topic: "半导体", count: 98 },
  { topic: "地缘政治", count: 87 },
  { topic: "能源转型", count: 72 },
  { topic: "量子计算", count: 54 },
  { topic: "气候变化", count: 48 },
  { topic: "数字货币", count: 41 },
  { topic: "医疗 AI", count: 35 },
];

const HEALTH_DATA = [
  { name: "正常", value: 3102, color: "#10B981" },
  { name: "降级", value: 108, color: "#F59E0B" },
  { name: "失败", value: 37, color: "#EF4444" },
];

const SENTIMENT_HEATMAP = Array.from({ length: 7 }, (_, d) =>
  Array.from({ length: 24 }, (_, h) => ({
    d, h,
    v: Math.random(),
  }))
).flat();

const DAYS = ["Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"];

export default function DashboardPage() {
  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold">Dashboard 总览</h1>

      {/* Metric cards */}
      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
        {METRICS.map((m) => (
          <div key={m.label} className="rounded-xl border border-border bg-white p-5">
            <div className="flex items-center justify-between">
              <span className="text-sm text-muted-foreground">{m.label}</span>
              <div className={`rounded-lg ${m.bg} p-2`}>
                <m.icon className={`h-4 w-4 ${m.color}`} />
              </div>
            </div>
            <div className="mt-2 text-3xl font-bold">{m.value}</div>
            <div className="mt-1 text-xs text-accent-green">{m.change} 较昨日</div>
          </div>
        ))}
      </div>

      {/* Charts row 1 */}
      <div className="grid gap-4 lg:grid-cols-3">
        <div className="lg:col-span-2">
          <AreaChartCard
            title="文章采集趋势（24h）"
            subtitle="每小时文章数量与预警数量"
            data={ARTICLE_TREND}
            xKey="hour"
            areas={[
              { key: "articles", color: "#3B82F6", name: "文章" },
              { key: "alerts", color: "#EF4444", name: "预警" },
            ]}
            height={260}
          />
        </div>
        <PieChartCard
          title="数据源健康"
          data={HEALTH_DATA}
          height={260}
        />
      </div>

      {/* Charts row 2 */}
      <div className="grid gap-4 lg:grid-cols-2">
        <BarChartCard
          title="热点话题 TOP 8"
          subtitle="按文章提及频次排序"
          data={TOPIC_DATA}
          xKey="topic"
          bars={[{ key: "count", color: "#3B82F6", name: "提及次数" }]}
          layout="vertical"
          height={280}
        />

        {/* Sentiment heatmap */}
        <div className="rounded-xl border border-border bg-white p-6">
          <h3 className="text-sm font-semibold">情感热力图（7×24）</h3>
          <p className="text-xs text-muted-foreground mb-4">按日/时段的正面情感占比</p>
          <div className="space-y-1">
            {DAYS.map((day, d) => (
              <div key={day} className="flex items-center gap-1">
                <span className="w-8 text-[10px] text-muted-foreground">{day}</span>
                <div className="flex flex-1 gap-0.5">
                  {Array.from({ length: 24 }, (_, h) => {
                    const cell = SENTIMENT_HEATMAP.find((c) => c.d === d && c.h === h);
                    const v = cell?.v ?? 0.5;
                    const r = Math.round(239 - v * 190);
                    const g = Math.round(68 + v * 117);
                    const b = Math.round(68 + v * 63);
                    return (
                      <div
                        key={h}
                        className="flex-1 aspect-square rounded-[2px]"
                        style={{ backgroundColor: `rgb(${r},${g},${b})`, opacity: 0.7 + v * 0.3 }}
                        title={`${day} ${h}:00 — ${(v * 100).toFixed(0)}% 正面`}
                      />
                    );
                  })}
                </div>
              </div>
            ))}
            <div className="flex items-center gap-1 mt-1 ml-9">
              {Array.from({ length: 24 }, (_, h) => (
                <div key={h} className="flex-1 text-center text-[8px] text-muted-foreground">
                  {h % 6 === 0 ? h : ""}
                </div>
              ))}
            </div>
          </div>
        </div>
      </div>

      {/* Activity feed + live list */}
      <div className="grid gap-4 lg:grid-cols-3">
        <div className="lg:col-span-2 rounded-xl border border-border bg-white p-6">
          <h3 className="text-sm font-semibold mb-4">实时信息流</h3>
          <div className="space-y-3">
            {[
              { title: "AI 监管新政：欧盟发布最新人工智能法案实施细则", source: "Reuters", time: "3 min ago", tag: "科技政策" },
              { title: "美联储会议纪要释放鸽派信号，市场反应积极", source: "Bloomberg", time: "12 min ago", tag: "金融" },
              { title: "全球碳排放达到新高，COP 会议各方博弈加剧", source: "Nature", time: "28 min ago", tag: "气候" },
              { title: "半导体出口管制升级：日荷最新限制措施分析", source: "TechCrunch", time: "45 min ago", tag: "地缘政治" },
              { title: "量子计算突破：Google 实现百万量子比特纠错", source: "arXiv", time: "1 hr ago", tag: "前沿科技" },
            ].map((item) => (
              <div key={item.title} className="flex items-start gap-3 rounded-lg p-3 hover:bg-muted/50 transition-colors cursor-pointer">
                <div className="mt-1 h-2 w-2 shrink-0 rounded-full bg-accent-blue" />
                <div className="flex-1 min-w-0">
                  <div className="text-sm font-medium truncate">{item.title}</div>
                  <div className="mt-0.5 flex items-center gap-2 text-xs text-muted-foreground">
                    <span>{item.source}</span>
                    <span>·</span>
                    <span>{item.time}</span>
                    <span className="rounded bg-muted px-1.5 py-0.5 text-[10px]">{item.tag}</span>
                  </div>
                </div>
              </div>
            ))}
          </div>
        </div>

        {/* Activity timeline */}
        <div className="rounded-xl border border-border bg-white p-6">
          <h3 className="text-sm font-semibold mb-4">活动时间线</h3>
          <div className="space-y-4">
            {[
              { time: "10:32", event: "RSS 全量抓取完成", detail: "采集 284 篇新文章", color: "bg-accent-green" },
              { time: "10:15", event: "舆情预警触发", detail: "「AI 监管」负面情感突增", color: "bg-red-500" },
              { time: "09:48", event: "AI 分析完成", detail: "生成 12 篇摘要, 5 个选题", color: "bg-accent-blue" },
              { time: "09:30", event: "数据源降级", detail: "arXiv RSS 响应超时", color: "bg-accent-orange" },
              { time: "09:00", event: "定时任务启动", detail: "开始每日数据源健康检查", color: "bg-muted-foreground" },
              { time: "08:00", event: "系统启动", detail: "所有 Agent 已就绪", color: "bg-accent-green" },
            ].map((a, i) => (
              <div key={i} className="flex gap-3">
                <div className="flex flex-col items-center">
                  <div className={`h-2 w-2 rounded-full ${a.color} mt-1.5`} />
                  {i < 5 && <div className="flex-1 w-px bg-border mt-1" />}
                </div>
                <div>
                  <div className="text-xs font-medium">{a.event}</div>
                  <div className="text-[11px] text-muted-foreground">{a.detail}</div>
                  <div className="text-[10px] text-muted-foreground mt-0.5">{a.time}</div>
                </div>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
}
