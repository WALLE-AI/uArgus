"use client";

import { Search, LayoutGrid, List, Newspaper } from "lucide-react";
import { cn } from "@/lib/utils";

const CATEGORIES = ["全部", "科技", "金融", "地缘政治", "气候", "医疗", "能源", "学术"];
const IMPORTANCE_OPTIONS = ["全部", "紧急", "高", "中", "低"];

interface FeedFilterProps {
  search: string;
  onSearchChange: (v: string) => void;
  activeCategory: string;
  onCategoryChange: (v: string) => void;
  activeImportance: string;
  onImportanceChange: (v: string) => void;
  layout: "card" | "list" | "magazine";
  onLayoutChange: (v: "card" | "list" | "magazine") => void;
}

export default function FeedFilter({
  search,
  onSearchChange,
  activeCategory,
  onCategoryChange,
  activeImportance,
  onImportanceChange,
  layout,
  onLayoutChange,
}: FeedFilterProps) {
  return (
    <div className="space-y-3">
      {/* Search + view toggle */}
      <div className="flex items-center gap-3">
        <div className="relative flex-1">
          <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
          <input
            value={search}
            onChange={(e) => onSearchChange(e.target.value)}
            placeholder="搜索文章..."
            className="w-full rounded-lg border border-border bg-white pl-9 pr-3 py-2 text-sm outline-none focus:border-accent-blue transition-colors"
          />
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

      {/* Category + Importance filters */}
      <div className="flex flex-wrap items-center gap-2">
        {CATEGORIES.map((cat) => (
          <button
            key={cat}
            onClick={() => onCategoryChange(cat)}
            className={cn(
              "rounded-full border px-3 py-1 text-xs transition-colors",
              activeCategory === cat
                ? "border-accent-blue bg-accent-blue/10 text-accent-blue"
                : "border-border bg-white text-muted-foreground hover:bg-muted"
            )}
          >
            {cat}
          </button>
        ))}
        <div className="mx-2 h-4 w-px bg-border" />
        {IMPORTANCE_OPTIONS.map((imp) => (
          <button
            key={imp}
            onClick={() => onImportanceChange(imp)}
            className={cn(
              "rounded-full border px-2.5 py-1 text-[11px] transition-colors",
              activeImportance === imp
                ? "border-accent-blue bg-accent-blue/10 text-accent-blue"
                : "border-border bg-white text-muted-foreground hover:bg-muted"
            )}
          >
            {imp}
          </button>
        ))}
      </div>
    </div>
  );
}
