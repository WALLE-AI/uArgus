"use client";

import { useState, useEffect } from "react";
import {
  Loader2,
  User,
  Building2,
  MapPin,
  Tag,
  TrendingDown,
  TrendingUp,
  Minus,
  AlertTriangle,
  Clock,
  CheckSquare,
  ChevronDown,
  ChevronUp,
  Zap,
} from "lucide-react";
import { cn } from "@/lib/utils";
import type { EventAnalysis, Entity } from "@/lib/types";
import { sentimentToLabel, scopeToLabel } from "@/lib/types";
import { fetchEventAnalysis } from "@/lib/api";

interface EventAnalysisPanelProps {
  hash: string;
}

// ── Entity type styling ──────────────────────────────────────────────────

const ENTITY_STYLE: Record<Entity["type"], { bg: string; icon: React.ElementType }> = {
  PER: { bg: "bg-violet-100 text-violet-700", icon: User },
  ORG: { bg: "bg-blue-100 text-blue-700", icon: Building2 },
  LOC: { bg: "bg-emerald-100 text-emerald-700", icon: MapPin },
  MISC: { bg: "bg-gray-100 text-gray-600", icon: Tag },
};

const ENTITY_TYPE_LABEL: Record<Entity["type"], string> = {
  PER: "人物",
  ORG: "组织",
  LOC: "地点",
  MISC: "其他",
};

const SENTIMENT_STYLE: Record<string, { color: string; icon: React.ElementType }> = {
  positive: { color: "text-emerald-600", icon: TrendingUp },
  negative: { color: "text-red-600", icon: TrendingDown },
  neutral: { color: "text-gray-500", icon: Minus },
};

// ── Disruption gauge ─────────────────────────────────────────────────────

function DisruptionGauge({ value }: { value: number }) {
  const pct = Math.round(value * 100);
  const color =
    pct >= 75 ? "bg-red-500" : pct >= 50 ? "bg-orange-500" : pct >= 25 ? "bg-yellow-500" : "bg-emerald-500";
  return (
    <div className="flex items-center gap-2">
      <div className="flex-1 h-2 rounded-full bg-muted overflow-hidden">
        <div className={cn("h-full rounded-full transition-all", color)} style={{ width: `${pct}%` }} />
      </div>
      <span className="text-xs font-bold tabular-nums w-8 text-right">{pct}%</span>
    </div>
  );
}

// ── Main component ───────────────────────────────────────────────────────

export default function EventAnalysisPanel({ hash }: EventAnalysisPanelProps) {
  const [analysis, setAnalysis] = useState<EventAnalysis | null>(null);
  const [loading, setLoading] = useState(false);
  const [loaded, setLoaded] = useState(false);
  const [expanded, setExpanded] = useState(false);
  const [fromApi, setFromApi] = useState(false);

  // Load on first expand
  useEffect(() => {
    if (!expanded || loaded) return;
    setLoading(true);
    fetchEventAnalysis(hash).then(({ data, fromApi: api }) => {
      setAnalysis(data);
      setFromApi(api);
      setLoaded(true);
      setLoading(false);
    });
  }, [expanded, loaded, hash]);

  // Reset when hash changes
  useEffect(() => {
    setAnalysis(null);
    setLoaded(false);
    setExpanded(false);
  }, [hash]);

  return (
    <div className="mt-5 rounded-lg border border-border overflow-hidden">
      {/* Toggle header */}
      <button
        onClick={() => setExpanded((v) => !v)}
        className="w-full flex items-center justify-between px-4 py-3 hover:bg-muted/30 transition-colors"
      >
        <div className="flex items-center gap-2">
          <Zap className="h-4 w-4 text-accent-purple" />
          <span className="text-sm font-semibold">事件分析</span>
          {!fromApi && loaded && (
            <span className="rounded bg-amber-100 px-1.5 py-0.5 text-[9px] font-medium text-amber-700">Mock</span>
          )}
        </div>
        {expanded ? <ChevronUp className="h-4 w-4 text-muted-foreground" /> : <ChevronDown className="h-4 w-4 text-muted-foreground" />}
      </button>

      {/* Expandable content */}
      {expanded && (
        <div className="border-t border-border px-4 py-4 space-y-5">
          {loading ? (
            <div className="flex items-center justify-center py-8 text-sm text-muted-foreground gap-2">
              <Loader2 className="h-4 w-4 animate-spin" />
              正在生成事件分析...
            </div>
          ) : !analysis ? (
            <div className="text-center py-8 text-sm text-muted-foreground">
              暂无该事件的分析数据
            </div>
          ) : (
            <>
              {/* ── Entities ── */}
              <section>
                <h4 className="text-[11px] font-semibold text-muted-foreground mb-2">实体识别</h4>
                <div className="flex flex-wrap gap-1.5">
                  {analysis.entities.map((e, i) => {
                    const style = ENTITY_STYLE[e.type];
                    const Icon = style.icon;
                    return (
                      <span
                        key={`${e.text}-${i}`}
                        className={cn("inline-flex items-center gap-1 rounded-full px-2.5 py-1 text-[11px] font-medium", style.bg)}
                        title={`${ENTITY_TYPE_LABEL[e.type]} · 置信度 ${Math.round(e.confidence * 100)}%`}
                      >
                        <Icon className="h-3 w-3" />
                        {e.text}
                      </span>
                    );
                  })}
                </div>
              </section>

              {/* ── Sentiment ── */}
              <section>
                <h4 className="text-[11px] font-semibold text-muted-foreground mb-2">情感倾向</h4>
                <div className="flex items-center gap-3">
                  {(() => {
                    const s = SENTIMENT_STYLE[analysis.sentiment.label];
                    const Icon = s.icon;
                    return (
                      <>
                        <div className={cn("flex items-center gap-1.5 text-sm font-semibold", s.color)}>
                          <Icon className="h-4 w-4" />
                          {sentimentToLabel(analysis.sentiment.label)}
                        </div>
                        <div className="flex-1 h-1.5 rounded-full bg-muted overflow-hidden">
                          <div
                            className={cn(
                              "h-full rounded-full",
                              analysis.sentiment.label === "positive" ? "bg-emerald-500" :
                              analysis.sentiment.label === "negative" ? "bg-red-500" : "bg-gray-400"
                            )}
                            style={{ width: `${analysis.sentiment.score * 100}%` }}
                          />
                        </div>
                        <span className="text-xs text-muted-foreground tabular-nums">
                          {Math.round(analysis.sentiment.score * 100)}%
                        </span>
                      </>
                    );
                  })()}
                </div>
              </section>

              {/* ── Impact Assessment ── */}
              <section>
                <h4 className="text-[11px] font-semibold text-muted-foreground mb-2">影响评估</h4>
                <div className="rounded-lg border border-border p-3 space-y-3">
                  <div className="flex items-center gap-2">
                    <span className="rounded bg-accent-purple/10 px-2 py-0.5 text-[10px] font-semibold text-accent-purple">
                      {scopeToLabel(analysis.impact.scope)}影响
                    </span>
                    <div className="flex gap-1">
                      {analysis.impact.sectors.map((s) => (
                        <span key={s} className="rounded bg-muted px-1.5 py-0.5 text-[9px] text-muted-foreground">{s}</span>
                      ))}
                    </div>
                  </div>
                  <div>
                    <div className="flex items-center gap-1.5 text-[10px] text-muted-foreground mb-1">
                      <AlertTriangle className="h-3 w-3" />
                      冲击指数
                    </div>
                    <DisruptionGauge value={analysis.impact.disruption} />
                  </div>
                  <p className="text-xs leading-relaxed text-foreground/80">{analysis.impact.summary}</p>
                </div>
              </section>

              {/* ── Timeline ── */}
              {analysis.timeline.length > 0 && (
                <section>
                  <h4 className="text-[11px] font-semibold text-muted-foreground mb-2">事件时间线</h4>
                  <div className="relative pl-4 space-y-3">
                    {/* vertical line */}
                    <div className="absolute left-[7px] top-1 bottom-1 w-px bg-border" />
                    {analysis.timeline.map((m, i) => (
                      <div key={i} className="relative flex items-start gap-3">
                        <div className={cn(
                          "absolute left-[-13px] top-1.5 h-2 w-2 rounded-full border-2 border-white",
                          i === analysis.timeline.length - 1 ? "bg-accent-blue" : "bg-muted-foreground/40"
                        )} />
                        <div>
                          <div className="text-[10px] text-muted-foreground tabular-nums">{m.date}</div>
                          <div className="text-xs leading-snug">{m.label}</div>
                        </div>
                      </div>
                    ))}
                  </div>
                </section>
              )}

              {/* ── Actions ── */}
              {analysis.actions.length > 0 && (
                <section>
                  <h4 className="text-[11px] font-semibold text-muted-foreground mb-2">行动建议</h4>
                  <div className="space-y-1.5">
                    {analysis.actions.map((a, i) => (
                      <div key={i} className="flex items-start gap-2 text-xs">
                        <CheckSquare className="h-3.5 w-3.5 text-accent-blue mt-0.5 shrink-0" />
                        <span className="leading-relaxed">{a}</span>
                      </div>
                    ))}
                  </div>
                </section>
              )}

              {/* ── Generation metadata ── */}
              <div className="text-[10px] text-muted-foreground text-right">
                分析生成于 {new Date(analysis.generatedAt).toLocaleString("zh-CN")}
              </div>
            </>
          )}
        </div>
      )}
    </div>
  );
}
