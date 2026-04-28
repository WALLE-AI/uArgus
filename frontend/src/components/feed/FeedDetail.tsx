"use client";

import { X, Bookmark, Share2, BookOpen, ExternalLink, Shield, Users, Radio, Layers, BrainCircuit } from "lucide-react";
import EventAnalysisPanel from "@/components/feed/EventAnalysisPanel";
import { cn } from "@/lib/utils";
import type { SemanticFeedItem, SeverityLevel } from "@/lib/types";
import { severityToLevel, severityToLabel, stageToLabel, tierToLabel } from "@/lib/types";
import { getCategoryLabel } from "@/lib/mock-data";

interface FeedDetailProps {
  item: SemanticFeedItem | null;
  relatedItems?: SemanticFeedItem[];
  onClose: () => void;
  onSelectRelated?: (item: SemanticFeedItem) => void;
}

const SEVERITY_STYLE: Record<SeverityLevel, string> = {
  critical: "bg-red-100 text-red-700",
  high: "bg-orange-100 text-orange-700",
  medium: "bg-yellow-100 text-yellow-700",
  low: "bg-muted text-muted-foreground",
  info: "bg-blue-50 text-blue-600",
};

const STAGE_STYLE: Record<string, string> = {
  BREAKING: "bg-red-500 text-white",
  DEVELOPING: "bg-orange-500 text-white",
  SUSTAINED: "bg-blue-500 text-white",
  FADING: "bg-gray-400 text-white",
};

const TIER_STYLE: Record<number, string> = {
  1: "text-emerald-600",
  2: "text-blue-600",
  3: "text-amber-600",
  4: "text-gray-500",
};

export default function FeedDetail({ item, relatedItems = [], onClose, onSelectRelated }: FeedDetailProps) {
  return (
    <div
      className={cn(
        "fixed inset-y-0 right-0 z-50 w-full max-w-xl transform border-l border-border bg-white shadow-2xl transition-transform duration-300",
        item ? "translate-x-0" : "translate-x-full"
      )}
    >
      {item && (
        <div className="flex h-full flex-col">
          {/* Header */}
          <div className="flex items-center justify-between border-b border-border px-6 py-4">
            <div className="flex items-center gap-2">
              <span className={cn("rounded px-2 py-0.5 text-[11px] font-medium", SEVERITY_STYLE[severityToLevel(item.severity)])}>
                {severityToLabel(item.severity)}
              </span>
              {item.stage && (
                <span className={cn("rounded px-1.5 py-0.5 text-[10px] font-medium", STAGE_STYLE[item.stage])}>
                  {stageToLabel(item.stage)}
                </span>
              )}
              <span className="text-xs text-muted-foreground">{item.source}</span>
            </div>
            <div className="flex items-center gap-1">
              <a href={item.link} target="_blank" rel="noopener noreferrer" className="rounded p-1.5 text-muted-foreground hover:bg-muted hover:text-foreground transition-colors">
                <ExternalLink className="h-4 w-4" />
              </a>
              <button onClick={onClose} className="rounded p-1.5 text-muted-foreground hover:bg-muted hover:text-foreground transition-colors">
                <X className="h-4 w-4" />
              </button>
            </div>
          </div>

          {/* Content */}
          <div className="flex-1 overflow-y-auto px-6 py-6">
            <h1 className="text-xl font-bold leading-tight">{item.title}</h1>
            <div className="mt-3 flex gap-1.5">
              {item.categories.map((cat) => (
                <span key={cat} className="rounded bg-muted px-2 py-0.5 text-[11px] text-muted-foreground">{getCategoryLabel(cat)}</span>
              ))}
            </div>

            {/* Semantic Analysis Panel */}
            <div className="mt-5 grid grid-cols-2 gap-3">
              <div className="rounded-lg border border-border p-3">
                <div className="flex items-center gap-1.5 text-[10px] text-muted-foreground mb-1.5">
                  <Radio className="h-3 w-3" />
                  重要度评分
                </div>
                <div className="flex items-center gap-2">
                  <div className="text-lg font-bold text-accent-blue">{Math.round(item.importanceScore * 100)}</div>
                  <div className="flex-1">
                    <div className="h-1.5 rounded-full bg-muted overflow-hidden">
                      <div className="h-full rounded-full bg-accent-blue" style={{ width: `${item.importanceScore * 100}%` }} />
                    </div>
                  </div>
                </div>
              </div>
              <div className="rounded-lg border border-border p-3">
                <div className="flex items-center gap-1.5 text-[10px] text-muted-foreground mb-1.5">
                  <Users className="h-3 w-3" />
                  跨源佐证
                </div>
                <div className="text-lg font-bold">
                  {item.corroboration} <span className="text-xs font-normal text-muted-foreground">个独立来源</span>
                </div>
              </div>
              <div className="rounded-lg border border-border p-3">
                <div className="flex items-center gap-1.5 text-[10px] text-muted-foreground mb-1.5">
                  <Shield className="h-3 w-3" />
                  来源可信度
                </div>
                <div className={cn("text-sm font-semibold", TIER_STYLE[item.tier])}>
                  Tier {item.tier} — {tierToLabel(item.tier)}
                </div>
              </div>
              <div className="rounded-lg border border-border p-3">
                <div className="flex items-center gap-1.5 text-[10px] text-muted-foreground mb-1.5">
                  <Layers className="h-3 w-3" />
                  严重度分级
                </div>
                <div className="text-sm font-semibold">
                  Level {item.severity}/7 — {severityToLabel(item.severity)}
                </div>
              </div>
            </div>

            {/* AI Summary */}
            {item.description && (
              <div className="mt-5 rounded-lg border border-accent-blue/20 bg-accent-blue/5 p-4">
                <div className="flex items-center gap-1.5 text-xs font-medium text-accent-blue mb-2">
                  <BrainCircuit className="h-3.5 w-3.5" />
                  AI 摘要
                </div>
                <p className="text-sm leading-relaxed text-foreground">{item.description}</p>
              </div>
            )}

            {/* Event Stage Lifecycle */}
            {item.stage && (
              <div className="mt-5">
                <h3 className="text-xs font-semibold text-muted-foreground mb-2">事件生命周期</h3>
                <div className="flex gap-1">
                  {(["BREAKING", "DEVELOPING", "SUSTAINED", "FADING"] as const).map((s) => (
                    <div
                      key={s}
                      className={cn(
                        "flex-1 rounded py-1 text-center text-[10px] font-medium transition-all",
                        item.stage === s ? STAGE_STYLE[s] : "bg-muted/50 text-muted-foreground"
                      )}
                    >
                      {stageToLabel(s)}
                    </div>
                  ))}
                </div>
              </div>
            )}

            {/* Event Analysis Panel */}
            <EventAnalysisPanel hash={item.hash} />

            {/* Related articles (semantic) */}
            {relatedItems.length > 0 && (
              <div className="mt-8">
                <h3 className="text-sm font-semibold mb-3">相关文章（语义关联）</h3>
                <div className="space-y-2">
                  {relatedItems.map((r) => (
                    <div
                      key={r.hash}
                      onClick={() => onSelectRelated?.(r)}
                      className="flex items-center gap-2 rounded-lg p-2 hover:bg-muted/50 cursor-pointer transition-colors"
                    >
                      <div className="h-1.5 w-1.5 rounded-full bg-accent-blue shrink-0" />
                      <div className="flex-1 min-w-0">
                        <span className="text-xs truncate block">{r.title}</span>
                        <span className="text-[10px] text-muted-foreground">{r.source}</span>
                      </div>
                      <span className="text-[10px] text-muted-foreground">{Math.round(r.importanceScore * 100)}</span>
                    </div>
                  ))}
                </div>
              </div>
            )}
          </div>

          {/* Actions bar */}
          <div className="flex items-center justify-between border-t border-border px-6 py-3">
            <div className="flex gap-2">
              {[
                { icon: Bookmark, label: "收藏" },
                { icon: BrainCircuit, label: "深度分析" },
                { icon: Share2, label: "分享" },
                { icon: BookOpen, label: "加入知识库" },
              ].map((action) => (
                <button
                  key={action.label}
                  className="inline-flex items-center gap-1.5 rounded-lg border border-border px-3 py-1.5 text-xs text-muted-foreground hover:bg-muted hover:text-foreground transition-colors"
                >
                  <action.icon className="h-3.5 w-3.5" />
                  {action.label}
                </button>
              ))}
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
