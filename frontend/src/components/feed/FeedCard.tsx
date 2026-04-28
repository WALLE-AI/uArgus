"use client";

import { Bookmark, BrainCircuit, Share2, BookOpen } from "lucide-react";
import { cn } from "@/lib/utils";
import type { SemanticFeedItem, SeverityLevel } from "@/lib/types";
import { severityToLevel, severityToLabel, stageToLabel } from "@/lib/types";
import { getCategoryLabel } from "@/lib/mock-data";

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

interface FeedCardProps {
  item: SemanticFeedItem;
  selected?: boolean;
  onClick?: () => void;
  layout?: "card" | "list";
}

function formatTimeAgo(pubDate: string): string {
  const diff = Date.now() - new Date(pubDate).getTime();
  const mins = Math.floor(diff / 60000);
  if (mins < 1) return "刚刚";
  if (mins < 60) return `${mins} min ago`;
  const hrs = Math.floor(mins / 60);
  if (hrs < 24) return `${hrs} hr ago`;
  const days = Math.floor(hrs / 24);
  return `${days} day ago`;
}

/** Score bar: 0-1 mapped to a small visual indicator */
function ScoreBar({ score }: { score: number }) {
  const pct = Math.round(score * 100);
  return (
    <div className="flex items-center gap-1.5" title={`重要度 ${pct}`}>
      <div className="h-1 w-12 rounded-full bg-muted overflow-hidden">
        <div
          className="h-full rounded-full bg-accent-blue transition-all"
          style={{ width: `${pct}%` }}
        />
      </div>
      <span className="text-[9px] text-muted-foreground tabular-nums">{pct}</span>
    </div>
  );
}

export default function FeedCard({ item, selected, onClick, layout = "card" }: FeedCardProps) {
  const level = severityToLevel(item.severity);
  const timeAgo = formatTimeAgo(item.pubDate);

  if (layout === "list") {
    return (
      <div
        onClick={onClick}
        className={cn(
          "flex items-center gap-4 rounded-lg border border-border bg-white px-4 py-3 cursor-pointer transition-all hover:shadow-sm",
          selected && "border-accent-blue ring-1 ring-accent-blue/30"
        )}
      >
        <span className={cn("shrink-0 rounded px-1.5 py-0.5 text-[10px] font-medium", SEVERITY_STYLE[level])}>
          {severityToLabel(item.severity)}
        </span>
        {item.stage && (
          <span className={cn("shrink-0 rounded px-1.5 py-0.5 text-[9px] font-medium", STAGE_STYLE[item.stage])}>
            {stageToLabel(item.stage)}
          </span>
        )}
        <div className="flex-1 min-w-0">
          <div className="text-sm font-medium truncate">{item.title}</div>
        </div>
        <ScoreBar score={item.importanceScore} />
        <span className="shrink-0 text-xs text-muted-foreground">{item.source}</span>
        <span className="shrink-0 text-xs text-muted-foreground">{timeAgo}</span>
      </div>
    );
  }

  return (
    <div
      onClick={onClick}
      className={cn(
        "rounded-xl border border-border bg-white p-5 cursor-pointer transition-all hover:shadow-md",
        selected && "border-accent-blue ring-1 ring-accent-blue/30"
      )}
    >
      {/* Top row: severity + stage + source + time */}
      <div className="flex items-center gap-2 mb-1.5">
        <span className={cn("rounded px-1.5 py-0.5 text-[10px] font-medium", SEVERITY_STYLE[level])}>
          {severityToLabel(item.severity)}
        </span>
        {item.stage && (
          <span className={cn("rounded px-1 py-0.5 text-[9px] font-medium", STAGE_STYLE[item.stage])}>
            {stageToLabel(item.stage)}
          </span>
        )}
        <span className="text-xs text-muted-foreground">{item.source}</span>
        <span className="text-xs text-muted-foreground">{timeAgo}</span>
        <div className="ml-auto">
          <ScoreBar score={item.importanceScore} />
        </div>
      </div>

      <h3 className="text-sm font-semibold leading-snug">{item.title}</h3>
      {item.description && (
        <p className="mt-1 text-xs text-muted-foreground line-clamp-2">{item.description}</p>
      )}

      {/* Bottom row: categories + corroboration + tier + actions */}
      <div className="mt-3 flex items-center justify-between">
        <div className="flex items-center gap-1.5">
          {item.categories.slice(0, 3).map((cat) => (
            <span key={cat} className="rounded bg-muted px-2 py-0.5 text-[10px] text-muted-foreground">
              {getCategoryLabel(cat)}
            </span>
          ))}
          {item.corroboration > 1 && (
            <span className="rounded bg-accent-blue/10 px-1.5 py-0.5 text-[9px] text-accent-blue font-medium" title="跨源佐证数">
              {item.corroboration} 源
            </span>
          )}
        </div>
        <div className="flex gap-1">
          {[Bookmark, BrainCircuit, Share2, BookOpen].map((Icon, i) => (
            <button
              key={i}
              onClick={(e) => e.stopPropagation()}
              className="rounded p-1 text-muted-foreground/50 hover:text-foreground hover:bg-muted transition-colors"
            >
              <Icon className="h-3.5 w-3.5" />
            </button>
          ))}
        </div>
      </div>
    </div>
  );
}
