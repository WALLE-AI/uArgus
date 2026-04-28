"use client";

import { useState, useRef, useCallback, useEffect, useMemo } from "react";
import { useVirtualizer } from "@tanstack/react-virtual";
import FeedCard from "@/components/feed/FeedCard";
import FeedFilter from "@/components/feed/FeedFilter";
import FeedDetail from "@/components/feed/FeedDetail";
import type { SortKey } from "@/components/feed/FeedFilter";
import type { SemanticFeedItem } from "@/lib/types";
import { fetchDigest } from "@/lib/api";
import { getCategoryLabel } from "@/lib/mock-data";
import { useAppStore } from "@/lib/store";

/** Severity threshold map for filtering */
const SEVERITY_MIN: Record<string, number> = {
  all: -1,
  critical: 6,
  high: 4,
  medium: 2,
  low: 0,
};

/** Sort comparators using semantic fields */
function sortItems(items: SemanticFeedItem[], sortBy: SortKey): SemanticFeedItem[] {
  const sorted = [...items];
  switch (sortBy) {
    case "score":
      return sorted.sort((a, b) => b.importanceScore - a.importanceScore);
    case "time":
      return sorted.sort((a, b) => new Date(b.pubDate).getTime() - new Date(a.pubDate).getTime());
    case "corroboration":
      return sorted.sort((a, b) => b.corroboration - a.corroboration);
    case "severity":
      return sorted.sort((a, b) => b.severity - a.severity);
    default:
      return sorted;
  }
}

export default function FeedPage() {
  const [search, setSearch] = useState("");
  const [category, setCategory] = useState("all");
  const [severity, setSeverity] = useState("all");
  const [sortBy, setSortBy] = useState<SortKey>("score");
  const [layout, setLayout] = useState<"card" | "list" | "magazine">("card");
  const [selectedHash, setSelectedHash] = useState<string | null>(null);

  const parentRef = useRef<HTMLDivElement>(null);

  // ── Load feed from API (with mock fallback) ──
  const { feedItems, setFeedItems, feedLoading, setFeedLoading, feedFromApi, setFeedFromApi, feedError, setFeedError } = useAppStore();

  useEffect(() => {
    let cancelled = false;
    setFeedLoading(true);
    setFeedError(null);
    fetchDigest({ variant: "full", lang: "en" }).then(({ items, fromApi }) => {
      if (!cancelled) {
        setFeedItems(items);
        setFeedFromApi(fromApi);
        setFeedLoading(false);
      }
    }).catch((err) => {
      if (!cancelled) {
        setFeedError(String(err));
        setFeedLoading(false);
      }
    });
    return () => { cancelled = true; };
  }, [setFeedItems, setFeedLoading, setFeedFromApi, setFeedError]);

  // ── Filter + sort using semantic fields ──
  const filtered = useMemo(() => {
    let items = feedItems.filter((item) => {
      // text search across title, source, categories
      if (search) {
        const q = search.toLowerCase();
        const matchTitle = item.title.toLowerCase().includes(q);
        const matchSource = item.source.toLowerCase().includes(q);
        const matchCat = item.categories.some((c) => getCategoryLabel(c).toLowerCase().includes(q));
        if (!matchTitle && !matchSource && !matchCat) return false;
      }
      // category filter
      if (category !== "all" && !item.categories.includes(category)) return false;
      // severity filter
      if (severity !== "all" && item.severity < (SEVERITY_MIN[severity] ?? 0)) return false;
      return true;
    });
    return sortItems(items, sortBy);
  }, [feedItems, search, category, severity, sortBy]);

  const virtualizer = useVirtualizer({
    count: filtered.length,
    getScrollElement: () => parentRef.current,
    estimateSize: useCallback(() => (layout === "list" ? 56 : 160), [layout]),
    overscan: 5,
  });

  const selectedItem = filtered.find((i) => i.hash === selectedHash) ?? null;

  // Related items: same categories, excluding current
  const relatedItems = useMemo(() => {
    if (!selectedItem) return [];
    return feedItems
      .filter((i) => i.hash !== selectedItem.hash && i.categories.some((c) => selectedItem.categories.includes(c)))
      .sort((a, b) => b.importanceScore - a.importanceScore)
      .slice(0, 5);
  }, [selectedItem, feedItems]);

  return (
    <div className="space-y-4 h-full flex flex-col">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-3">
          <h1 className="text-2xl font-bold">AI 内容流</h1>
          {!feedFromApi && !feedLoading && (
            <span className="rounded bg-amber-100 px-2 py-0.5 text-[10px] font-medium text-amber-700">
              Mock 数据
            </span>
          )}
        </div>
        <span className="text-xs text-muted-foreground">
          {feedLoading ? "加载中..." : `${filtered.length} 篇文章`}
        </span>
      </div>

      {feedError && (
        <div className="rounded-lg border border-red-200 bg-red-50 p-3 text-xs text-red-700">
          数据加载失败: {feedError}
        </div>
      )}

      <FeedFilter
        search={search}
        onSearchChange={setSearch}
        activeCategory={category}
        onCategoryChange={setCategory}
        activeSeverity={severity}
        onSeverityChange={setSeverity}
        layout={layout}
        onLayoutChange={setLayout}
        sortBy={sortBy}
        onSortChange={setSortBy}
      />

      {/* Virtualized list */}
      <div ref={parentRef} className="flex-1 overflow-y-auto min-h-0">
        {feedLoading ? (
          <div className="flex items-center justify-center py-20 text-sm text-muted-foreground">
            正在从语义层加载数据...
          </div>
        ) : filtered.length === 0 ? (
          <div className="flex items-center justify-center py-20 text-sm text-muted-foreground">
            无匹配结果
          </div>
        ) : (
          <div
            className="relative"
            style={{ height: `${virtualizer.getTotalSize()}px` }}
          >
            {virtualizer.getVirtualItems().map((vRow) => {
              const item = filtered[vRow.index];
              return (
                <div
                  key={item.hash}
                  className="absolute left-0 right-0 px-0.5"
                  style={{
                    top: `${vRow.start}px`,
                    height: `${vRow.size}px`,
                    paddingBottom: 8,
                  }}
                >
                  <FeedCard
                    item={item}
                    selected={selectedHash === item.hash}
                    onClick={() => setSelectedHash(item.hash)}
                    layout={layout === "list" ? "list" : "card"}
                  />
                </div>
              );
            })}
          </div>
        )}
      </div>

      {/* Detail slide-in panel */}
      <FeedDetail
        item={selectedItem}
        relatedItems={relatedItems}
        onClose={() => setSelectedHash(null)}
        onSelectRelated={(r) => setSelectedHash(r.hash)}
      />
      {selectedItem && (
        <div
          className="fixed inset-0 z-40 bg-black/20"
          onClick={() => setSelectedHash(null)}
        />
      )}
    </div>
  );
}
