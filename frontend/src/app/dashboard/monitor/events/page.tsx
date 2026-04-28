"use client";

import { useState } from "react";
import { ChevronDown } from "lucide-react";
import AreaChartCard from "@/components/charts/AreaChartCard";
import { cn } from "@/lib/utils";

interface EventData {
  name: string;
  stage: string;
  sources: number;
  confidence: string;
  updated: string;
  impact: number;
  summary: string;
  corroboration: string[];
}

const EVENTS: EventData[] = [
  { name: "欧盟 AI 法案实施", stage: "DEVELOPING", sources: 47, confidence: "高", updated: "5 min ago", impact: 85, summary: "欧盟委员会发布详细实施指南，涵盖高风险系统合规要求。各成员国进入两年过渡期。", corroboration: ["Reuters", "Bloomberg", "FT", "Politico EU"] },
  { name: "美联储利率政策转向", stage: "SUSTAINED", sources: 89, confidence: "高", updated: "12 min ago", impact: 92, summary: "会议纪要释放鸽派信号，市场预期年内降息2-3次。全球资本市场反应积极。", corroboration: ["Bloomberg", "WSJ", "CNBC", "Reuters", "FT"] },
  { name: "日荷半导体管制", stage: "BREAKING", sources: 23, confidence: "中", updated: "2 min ago", impact: 78, summary: "日本和荷兰相继发布新出口管制措施，限制先进制程设备出口。中国半导体产业面临新挑战。", corroboration: ["Nikkei", "TechCrunch", "Reuters"] },
  { name: "COP 气候谈判", stage: "SUSTAINED", sources: 56, confidence: "高", updated: "1 hr ago", impact: 70, summary: "各国在碳减排目标上分歧明显，发展中国家要求更多资金支持。", corroboration: ["Nature", "Guardian", "BBC", "Al Jazeera"] },
  { name: "量子计算纠错突破", stage: "BREAKING", sources: 12, confidence: "低", updated: "45 min ago", impact: 65, summary: "Google 展示百万量子比特纠错能力，但实际应用仍需时日。学界评价分歧。", corroboration: ["arXiv", "Nature", "MIT Tech Review"] },
  { name: "加密货币 ETF 审批", stage: "FADING", sources: 34, confidence: "中", updated: "3 hr ago", impact: 55, summary: "SEC 逐步放开加密 ETF 审批，市场关注度已从高峰回落。", corroboration: ["Bloomberg", "CoinDesk", "Reuters"] },
];

const STAGES = ["BREAKING", "DEVELOPING", "SUSTAINED", "FADING"];
const STAGE_STYLE: Record<string, string> = {
  BREAKING: "bg-red-100 text-red-600",
  DEVELOPING: "bg-yellow-100 text-yellow-700",
  SUSTAINED: "bg-accent-blue/10 text-accent-blue",
  FADING: "bg-muted text-muted-foreground",
};
const STAGE_DOT: Record<string, string> = {
  BREAKING: "bg-red-500",
  DEVELOPING: "bg-yellow-500",
  SUSTAINED: "bg-accent-blue",
  FADING: "bg-muted-foreground/40",
};

const EVENT_TIMELINE = Array.from({ length: 14 }, (_, i) => ({
  date: `${i + 1}日`,
  BREAKING: Math.round(1 + Math.random() * 3),
  DEVELOPING: Math.round(2 + Math.random() * 4),
  SUSTAINED: Math.round(3 + Math.sin(i / 3) * 2),
  FADING: Math.round(1 + Math.random() * 2),
}));

export default function EventsPage() {
  const [expanded, setExpanded] = useState<string | null>(null);

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold">事件追踪</h1>

      {/* Lifecycle pipeline */}
      <div className="rounded-xl border border-border bg-white p-6">
        <h3 className="text-sm font-semibold mb-4">事件生命周期分布</h3>
        <div className="grid grid-cols-4 gap-3">
          {STAGES.map((stage) => {
            const count = EVENTS.filter((e) => e.stage === stage).length;
            return (
              <div key={stage} className="text-center">
                <div className={cn("rounded-lg p-4 mb-2", STAGE_STYLE[stage])}>
                  <div className="text-2xl font-bold">{count}</div>
                </div>
                <div className="text-xs font-medium">{stage}</div>
              </div>
            );
          })}
        </div>
        <div className="flex items-center justify-center gap-0 mt-3">
          {STAGES.map((s, i) => (
            <div key={s} className="flex items-center">
              <div className={cn("h-2 w-2 rounded-full", STAGE_DOT[s])} />
              {i < 3 && <div className="w-16 h-px bg-border mx-1" />}
            </div>
          ))}
        </div>
      </div>

      {/* Stage legend */}
      <div className="flex gap-4 text-xs">
        {STAGES.map((s) => (
          <div key={s} className="flex items-center gap-1.5">
            <span className={cn("h-2 w-2 rounded-full", STAGE_DOT[s])} />
            <span className="text-muted-foreground">{s}</span>
          </div>
        ))}
      </div>

      {/* Event cards */}
      <div className="space-y-3">
        {EVENTS.map((e) => (
          <div key={e.name} className="rounded-xl border border-border bg-white transition-shadow hover:shadow-md">
            <div
              className="flex items-center justify-between p-5 cursor-pointer"
              onClick={() => setExpanded(expanded === e.name ? null : e.name)}
            >
              <div className="flex items-center gap-3">
                <span className={cn("rounded px-2 py-0.5 text-[10px] font-semibold", STAGE_STYLE[e.stage])}>{e.stage}</span>
                <h3 className="text-sm font-semibold">{e.name}</h3>
              </div>
              <div className="flex items-center gap-3">
                <span className="text-xs text-muted-foreground">{e.updated}</span>
                <ChevronDown className={cn("h-4 w-4 text-muted-foreground transition-transform", expanded === e.name && "rotate-180")} />
              </div>
            </div>
            <div className="px-5 pb-1 flex gap-4 text-xs text-muted-foreground -mt-2">
              <span>跨源佐证: {e.sources} 篇</span>
              <span>置信度: {e.confidence}</span>
              <span>影响力: {e.impact}/100</span>
            </div>

            {expanded === e.name && (
              <div className="border-t border-border px-5 py-4 space-y-3">
                <p className="text-sm text-foreground/80 leading-relaxed">{e.summary}</p>
                <div>
                  <span className="text-xs font-medium">跨源佐证来源:</span>
                  <div className="flex gap-1.5 mt-1">
                    {e.corroboration.map((src) => (
                      <span key={src} className="rounded bg-muted px-2 py-0.5 text-[10px] text-muted-foreground">{src}</span>
                    ))}
                  </div>
                </div>
                <div className="flex items-center gap-2">
                  <span className="text-xs font-medium">影响力评分:</span>
                  <div className="flex-1 h-2 rounded-full bg-muted max-w-xs">
                    <div className="h-2 rounded-full bg-accent-blue" style={{ width: `${e.impact}%` }} />
                  </div>
                  <span className="text-xs tabular-nums">{e.impact}/100</span>
                </div>
              </div>
            )}
          </div>
        ))}
      </div>

      {/* Timeline chart */}
      <AreaChartCard
        title="事件时间轴（14天）"
        subtitle="按生命周期阶段统计每日事件数量"
        data={EVENT_TIMELINE}
        xKey="date"
        areas={[
          { key: "BREAKING", color: "#EF4444" },
          { key: "DEVELOPING", color: "#F59E0B" },
          { key: "SUSTAINED", color: "#3B82F6" },
          { key: "FADING", color: "#9CA3AF" },
        ]}
        stacked
        height={240}
      />
    </div>
  );
}
