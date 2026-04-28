"use client";

import { Bookmark, BrainCircuit, Share2, BookOpen } from "lucide-react";
import { cn } from "@/lib/utils";

export interface FeedItem {
  id: string;
  title: string;
  source: string;
  time: string;
  summary: string;
  tags: string[];
  importance: "紧急" | "高" | "中" | "低";
}

const IMPORTANCE_STYLE: Record<string, string> = {
  紧急: "bg-red-100 text-red-700",
  高: "bg-orange-100 text-orange-700",
  中: "bg-yellow-100 text-yellow-700",
  低: "bg-muted text-muted-foreground",
};

interface FeedCardProps {
  item: FeedItem;
  selected?: boolean;
  onClick?: () => void;
  layout?: "card" | "list";
}

export default function FeedCard({ item, selected, onClick, layout = "card" }: FeedCardProps) {
  if (layout === "list") {
    return (
      <div
        onClick={onClick}
        className={cn(
          "flex items-center gap-4 rounded-lg border border-border bg-white px-4 py-3 cursor-pointer transition-all hover:shadow-sm",
          selected && "border-accent-blue ring-1 ring-accent-blue/30"
        )}
      >
        <span className={cn("shrink-0 rounded px-1.5 py-0.5 text-[10px] font-medium", IMPORTANCE_STYLE[item.importance])}>
          {item.importance}
        </span>
        <div className="flex-1 min-w-0">
          <div className="text-sm font-medium truncate">{item.title}</div>
        </div>
        <span className="shrink-0 text-xs text-muted-foreground">{item.source}</span>
        <span className="shrink-0 text-xs text-muted-foreground">{item.time}</span>
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
      <div className="flex items-center gap-2 mb-1.5">
        <span className={cn("rounded px-1.5 py-0.5 text-[10px] font-medium", IMPORTANCE_STYLE[item.importance])}>
          {item.importance}
        </span>
        <span className="text-xs text-muted-foreground">{item.source}</span>
        <span className="text-xs text-muted-foreground">· {item.time}</span>
      </div>
      <h3 className="text-sm font-semibold leading-snug">{item.title}</h3>
      <p className="mt-1 text-xs text-muted-foreground line-clamp-2">{item.summary}</p>
      <div className="mt-3 flex items-center justify-between">
        <div className="flex gap-1.5">
          {item.tags.map((tag) => (
            <span key={tag} className="rounded bg-muted px-2 py-0.5 text-[10px] text-muted-foreground">{tag}</span>
          ))}
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
