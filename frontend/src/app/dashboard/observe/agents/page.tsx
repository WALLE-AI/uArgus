"use client";

import { useState } from "react";
import AreaChartCard from "@/components/charts/AreaChartCard";
import BarChartCard from "@/components/charts/BarChartCard";
import PieChartCard from "@/components/charts/PieChartCard";
import LineChartCard from "@/components/charts/LineChartCard";
import { cn } from "@/lib/utils";

const TOKEN_TREND = Array.from({ length: 7 }, (_, i) => ({
  day: ["一", "二", "三", "四", "五", "六", "日"][i],
  "GPT-4o": Math.round(300 + Math.random() * 200),
  "Claude 3.5": Math.round(200 + Math.random() * 150),
  "BGE-M3": Math.round(100 + Math.random() * 80),
  "Groq Llama": Math.round(80 + Math.random() * 60),
}));

const TASK_PIE = [
  { name: "Summarize", value: 35, color: "#3B82F6" },
  { name: "Classify", value: 25, color: "#8B5CF6" },
  { name: "Embed", value: 20, color: "#10B981" },
  { name: "Chat", value: 12, color: "#F59E0B" },
  { name: "Translate", value: 8, color: "#EC4899" },
];

const TOOL_DATA = [
  { tool: "RSS Fetcher", calls: 2341 },
  { tool: "Summarizer", calls: 1024 },
  { tool: "Classifier", calls: 856 },
  { tool: "Embedder", calls: 612 },
  { tool: "Translator", calls: 188 },
  { tool: "Web Search", calls: 145 },
];

const LATENCY_TREND = Array.from({ length: 24 }, (_, i) => ({
  hour: `${i}:00`,
  "GPT-4o": Math.round(280 + Math.sin(i / 3) * 60 + Math.random() * 30),
  "Claude 3.5": Math.round(240 + Math.cos(i / 4) * 50 + Math.random() * 25),
  "Groq Llama": Math.round(120 + Math.sin(i / 2) * 100 + Math.random() * 50),
}));

const COST_TABLE = [
  { model: "GPT-4o", tokens: "520K", price: "¥0.15/1K", cost: "¥78.00" },
  { model: "Claude 3.5", tokens: "310K", price: "¥0.12/1K", cost: "¥37.20" },
  { model: "BGE-M3", tokens: "180K", price: "¥0.01/1K", cost: "¥1.80" },
  { model: "Groq Llama", tokens: "140K", price: "¥0.02/1K", cost: "¥2.80" },
];

const FAIL_LOGS = [
  { time: "10:23", tool: "Web Search", error: "Timeout after 30s", retry: true },
  { time: "09:47", tool: "Summarizer", error: "Rate limit exceeded (429)", retry: true },
  { time: "08:15", tool: "RSS Fetcher", error: "Connection refused: nasa.gov", retry: false },
  { time: "昨日 23:50", tool: "Embedder", error: "OOM: batch size too large", retry: false },
];

const TASK_LOG = [
  { id: "T-4521", name: "RSS 全量抓取", status: "completed", duration: "2m 15s" },
  { id: "T-4520", name: "新文章摘要生成", status: "running", duration: "45s" },
  { id: "T-4519", name: "情感分析批处理", status: "completed", duration: "1m 32s" },
  { id: "T-4518", name: "知识图谱更新", status: "pending", duration: "—" },
  { id: "T-4517", name: "舆情预警检查", status: "failed", duration: "12s" },
];

const TASK_STATUS_STYLE: Record<string, string> = {
  completed: "bg-accent-green/10 text-accent-green",
  running: "bg-accent-blue/10 text-accent-blue",
  pending: "bg-muted text-muted-foreground",
  failed: "bg-red-50 text-red-500",
};

const TABS = ["Token 消耗", "工具调用", "Agent 任务", "模型健康"];

export default function AgentsPage() {
  const [tab, setTab] = useState("Token 消耗");

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold">智能体观测</h1>

      {/* Overview cards */}
      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
        {[
          { label: "今日 Token 消耗", value: "1.2M", sub: "≈ ¥119.80" },
          { label: "工具调用次数", value: "4,521", sub: "成功率 99.2%" },
          { label: "活跃 Agent", value: "8", sub: "12 个已注册" },
          { label: "错误率", value: "0.8%", sub: "较昨日 -0.2%" },
        ].map((m) => (
          <div key={m.label} className="rounded-xl border border-border bg-white p-5">
            <div className="text-xs text-muted-foreground">{m.label}</div>
            <div className="mt-1 text-2xl font-bold">{m.value}</div>
            <div className="text-xs text-muted-foreground">{m.sub}</div>
          </div>
        ))}
      </div>

      {/* Tab switcher */}
      <div className="flex gap-1 rounded-lg border border-border bg-white p-1">
        {TABS.map((t) => (
          <button
            key={t}
            onClick={() => setTab(t)}
            className={cn(
              "flex-1 rounded-md px-3 py-2 text-xs font-medium transition-colors",
              tab === t ? "bg-accent-blue text-white" : "text-muted-foreground hover:bg-muted"
            )}
          >
            {t}
          </button>
        ))}
      </div>

      {/* Token panel */}
      {tab === "Token 消耗" && (
        <div className="space-y-4">
          <div className="grid gap-4 lg:grid-cols-2">
            <BarChartCard
              title="按模型维度 Token 消耗（周）"
              data={TOKEN_TREND}
              xKey="day"
              bars={[
                { key: "GPT-4o", color: "#3B82F6" },
                { key: "Claude 3.5", color: "#8B5CF6" },
                { key: "BGE-M3", color: "#10B981" },
                { key: "Groq Llama", color: "#F59E0B" },
              ]}
              stacked
              height={260}
            />
            <PieChartCard
              title="按任务类型维度 Token 分布"
              data={TASK_PIE}
              height={260}
            />
          </div>
          <div className="rounded-xl border border-border bg-white p-6">
            <h3 className="text-sm font-semibold mb-4">费用估算表</h3>
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-border text-left text-xs text-muted-foreground">
                  <th className="pb-2 font-medium">模型</th>
                  <th className="pb-2 font-medium">用量</th>
                  <th className="pb-2 font-medium">单价</th>
                  <th className="pb-2 font-medium text-right">费用</th>
                </tr>
              </thead>
              <tbody>
                {COST_TABLE.map((r) => (
                  <tr key={r.model} className="border-b border-border last:border-0">
                    <td className="py-2 font-medium">{r.model}</td>
                    <td className="py-2 text-muted-foreground tabular-nums">{r.tokens}</td>
                    <td className="py-2 text-muted-foreground">{r.price}</td>
                    <td className="py-2 text-right font-medium">{r.cost}</td>
                  </tr>
                ))}
                <tr className="font-semibold">
                  <td className="pt-2">合计</td><td></td><td></td>
                  <td className="pt-2 text-right text-accent-blue">¥119.80</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      )}

      {/* Tool calls panel */}
      {tab === "工具调用" && (
        <div className="space-y-4">
          <div className="grid gap-4 lg:grid-cols-2">
            <BarChartCard
              title="工具调用排行"
              data={TOOL_DATA}
              xKey="tool"
              bars={[{ key: "calls", color: "#3B82F6", name: "调用次数" }]}
              layout="vertical"
              height={260}
            />
            <div className="rounded-xl border border-border bg-white p-6">
              <h3 className="text-sm font-semibold mb-2">成功率指标</h3>
              <div className="grid grid-cols-2 gap-4 mt-4">
                {[
                  { label: "总成功率", value: "99.2%", color: "text-accent-green" },
                  { label: "平均延迟", value: "142ms", color: "text-accent-blue" },
                  { label: "失败次数", value: "36", color: "text-red-500" },
                  { label: "重试成功", value: "28", color: "text-accent-orange" },
                ].map((g) => (
                  <div key={g.label} className="text-center p-3 rounded-lg bg-muted/30">
                    <div className={`text-xl font-bold ${g.color}`}>{g.value}</div>
                    <div className="text-[11px] text-muted-foreground">{g.label}</div>
                  </div>
                ))}
              </div>
            </div>
          </div>
          <div className="rounded-xl border border-border bg-white p-6">
            <h3 className="text-sm font-semibold mb-4">失败调用日志</h3>
            <div className="space-y-2">
              {FAIL_LOGS.map((l, i) => (
                <div key={i} className="flex items-center gap-3 rounded-lg p-3 hover:bg-muted/30 transition-colors">
                  <span className="text-xs text-muted-foreground w-20">{l.time}</span>
                  <span className="text-xs font-medium w-24">{l.tool}</span>
                  <span className="flex-1 text-xs text-red-500">{l.error}</span>
                  <span className={cn("rounded px-2 py-0.5 text-[10px]", l.retry ? "bg-accent-green/10 text-accent-green" : "bg-muted text-muted-foreground")}>
                    {l.retry ? "可重试" : "不可重试"}
                  </span>
                </div>
              ))}
            </div>
          </div>
        </div>
      )}

      {/* Agent tasks panel */}
      {tab === "Agent 任务" && (
        <div className="space-y-4">
          <div className="grid gap-4 sm:grid-cols-3">
            <div className="rounded-xl border border-border bg-white p-5">
              <div className="text-2xl font-bold">3</div>
              <div className="text-xs text-muted-foreground">队列深度</div>
            </div>
            <div className="rounded-xl border border-border bg-white p-5">
              <div className="text-2xl font-bold">1m 28s</div>
              <div className="text-xs text-muted-foreground">平均执行时长</div>
            </div>
            <div className="rounded-xl border border-border bg-white p-5">
              <div className="text-2xl font-bold text-accent-green">94.7%</div>
              <div className="text-xs text-muted-foreground">任务成功率</div>
            </div>
          </div>
          <div className="rounded-xl border border-border bg-white p-6">
            <h3 className="text-sm font-semibold mb-4">任务状态表</h3>
            <div className="space-y-2">
              {TASK_LOG.map((t) => (
                <div key={t.id} className="flex items-center gap-3 rounded-lg p-3 hover:bg-muted/30 transition-colors">
                  <span className="text-xs text-muted-foreground font-mono w-16">{t.id}</span>
                  <span className="flex-1 text-xs font-medium">{t.name}</span>
                  <span className={cn("rounded px-2 py-0.5 text-[10px] font-medium", TASK_STATUS_STYLE[t.status])}>{t.status}</span>
                  <span className="text-xs text-muted-foreground tabular-nums w-16 text-right">{t.duration}</span>
                </div>
              ))}
            </div>
          </div>
        </div>
      )}

      {/* Model health panel */}
      {tab === "模型健康" && (
        <div className="space-y-4">
          <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
            {[
              { name: "GPT-4o", status: "正常", p50: "320ms", p95: "580ms", p99: "1.2s", fallback: 0, rateLimit: 0 },
              { name: "Claude 3.5", status: "正常", p50: "280ms", p95: "520ms", p99: "980ms", fallback: 2, rateLimit: 1 },
              { name: "BGE-M3", status: "正常", p50: "45ms", p95: "82ms", p99: "150ms", fallback: 0, rateLimit: 0 },
              { name: "Groq Llama", status: "降级", p50: "120ms", p95: "890ms", p99: "2.1s", fallback: 12, rateLimit: 5 },
            ].map((m) => (
              <div key={m.name} className="rounded-xl border border-border bg-white p-5">
                <div className="flex items-center gap-2 mb-3">
                  <span className={`h-3 w-3 rounded-full ${m.status === "正常" ? "bg-accent-green" : "bg-accent-orange"}`} />
                  <span className="text-sm font-semibold">{m.name}</span>
                </div>
                <div className="space-y-1.5 text-xs">
                  <div className="flex justify-between"><span className="text-muted-foreground">P50</span><span className="font-medium">{m.p50}</span></div>
                  <div className="flex justify-between"><span className="text-muted-foreground">P95</span><span className="font-medium">{m.p95}</span></div>
                  <div className="flex justify-between"><span className="text-muted-foreground">P99</span><span className="font-medium">{m.p99}</span></div>
                  <div className="flex justify-between"><span className="text-muted-foreground">Fallback</span><span className={m.fallback > 0 ? "text-accent-orange font-medium" : "text-muted-foreground"}>{m.fallback}</span></div>
                  <div className="flex justify-between"><span className="text-muted-foreground">Rate Limit</span><span className={m.rateLimit > 0 ? "text-red-500 font-medium" : "text-muted-foreground"}>{m.rateLimit}</span></div>
                </div>
              </div>
            ))}
          </div>
          <LineChartCard
            title="模型响应延迟趋势（24h P50）"
            data={LATENCY_TREND}
            xKey="hour"
            lines={[
              { key: "GPT-4o", color: "#3B82F6" },
              { key: "Claude 3.5", color: "#8B5CF6" },
              { key: "Groq Llama", color: "#F59E0B" },
            ]}
            height={260}
          />
        </div>
      )}
    </div>
  );
}
