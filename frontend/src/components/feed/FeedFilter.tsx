"use client";

import { Search, LayoutGrid, List, Newspaper, ArrowDownUp } from "lucide-react";
import { cn } from "@/lib/utils";
import { CATEGORY_LABELS } from "@/lib/mock-data";

/** Semantic categories mapped from backend classify output */
const CATEGORY_FILTERS = [
  { key: "all", label: "全部" },
  ...Object.entries(CATEGORY_LABELS).map(([key, label]) => ({ key, label })),
];

/** Severity-based importance filter aligned with semantic scoring */
const SEVERITY_FILTERS = [
  { key: "all", label: "全部" },
  { key: "critical", label: "紧急", min: 6 },
  { key: "high", label: "高", min: 4 },
  { key: "medium", label: "中", min: 2 },
  { key: "low", label: "低", min: 0 },
];

/** Sort options using semantic fields */
const SORT_OPTIONS = [
  { key: "score", label: "重要度" },
  { key: "time", label: "时间" },
  { key: "corroboration", label: "佐证数" },
  { key: "severity", label: "严重度" },
];

export type SortKey = "score" | "time" | "corroboration" | "severity";

interface FeedFilterProps {
  search: string;
  onSearchChange: (v: string) => void;
  activeCategory: string;
  onCategoryChange: (v: string) => void;
  activeSeverity: string;
  onSeverityChange: (v: string) => void;
  layout: "card" | "list" | "magazine";
  onLayoutChange: (v: "card" | "list" | "magazine") => void;
  sortBy: SortKey;
  onSortChange: (v: SortKey) => void;
}

export default function FeedFilter({
  search,
  onSearchChange,
  activeCategory,
  onCategoryChange,
  activeSeverity,
  onSeverityChange,
  layout,
  onLayoutChange,
  sortBy,
  onSortChange,
}: FeedFilterProps) {
  return (
    <div className="space-y-3">
      {/* Search + sort + view toggle */}
      <div className="flex items-center gap-3">
        <div className="relative flex-1">
          <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
          <input
            value={search}
            onChange={(e) => onSearchChange(e.target.value)}
            placeholder="搜索文章（标题 / 来源 / 分类）..."
            className="w-full rounded-lg border border-border bg-white pl-9 pr-3 py-2 text-sm outline-none focus:border-accent-blue transition-colors"
          />
        </div>
        {/* Sort selector */}
        <div className="flex items-center gap-1 rounded-lg border border-border bg-white px-2 py-1">
          <ArrowDownUp className="h-3.5 w-3.5 text-muted-foreground" />
          <select
            value={sortBy}
            onChange={(e) => onSortChange(e.target.value as SortKey)}
            className="text-xs bg-transparent outline-none cursor-pointer text-muted-foreground"
          >
            {SORT_OPTIONS.map((opt) => (
              <option key={opt.key} value={opt.key}>{opt.label}</option>
            ))}
          </select>
        </div>
        <div className="flex rounded-lg border border-border bg-white">
          {([
            { key: "card" as const, icon: LayoutGrid },
            { key: "list" as const, icon: List },
            { key: "magazine" as const, icon: Newspaper },
          ]).map((opt) => (
            <button
              key={opt.key}
              onClick={() => onLayoutChange(opt.key)}
              className={cn(
                "p-2 transition-colors",
                layout === opt.key ? "bg-muted text-foreground" : "text-muted-foreground hover:text-foreground"
              )}
            >
              <opt.icon className="h-4 w-4" />
            </button>
          ))}
        </div>
      </div>

      {/* Category + Severity filters */}
      <div className="flex flex-wrap items-center gap-2">
        {CATEGORY_FILTERS.slice(0, 9).map((cat) => (
          <button
            key={cat.key}
            onClick={() => onCategoryChange(cat.key)}
            className={cn(
              "rounded-full border px-3 py-1 text-xs transition-colors",
              activeCategory === cat.key
                ? "border-accent-blue bg-accent-blue/10 text-accent-blue"
                : "border-border bg-white text-muted-foreground hover:bg-muted"
            )}
          >
            {cat.label}
          </button>
        ))}
        <div className="mx-2 h-4 w-px bg-border" />
        {SEVERITY_FILTERS.map((sev) => (
          <button
            key={sev.key}
            onClick={() => onSeverityChange(sev.key)}
            className={cn(
              "rounded-full border px-2.5 py-1 text-[11px] transition-colors",
              activeSeverity === sev.key
                ? "border-accent-blue bg-accent-blue/10 text-accent-blue"
                : "border-border bg-white text-muted-foreground hover:bg-muted"
            )}
          >
            {sev.label}
          </button>
        ))}
      </div>
    </div>
  );
}
