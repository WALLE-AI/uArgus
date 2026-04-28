"use client";

import { useState, useMemo } from "react";
import {
  useReactTable,
  getCoreRowModel,
  getSortedRowModel,
  getFilteredRowModel,
  flexRender,
  createColumnHelper,
  type SortingState,
} from "@tanstack/react-table";
import { ArrowUpDown, Plus, X } from "lucide-react";
import { cn } from "@/lib/utils";

interface Source {
  id: string;
  name: string;
  url: string;
  category: string;
  status: "正常" | "降级" | "失败";
  lastFetch: string;
  articles: number;
  interval: string;
}

const MOCK_SOURCES: Source[] = [
  { id: "1", name: "Reuters Top News", url: "reuters.com/rss/topNews", category: "新闻", status: "正常", lastFetch: "2 min ago", articles: 12483, interval: "15min" },
  { id: "2", name: "Hacker News", url: "news.ycombinator.com/rss", category: "科技", status: "正常", lastFetch: "5 min ago", articles: 8921, interval: "15min" },
  { id: "3", name: "arXiv CS.AI", url: "arxiv.org/rss/cs.AI", category: "学术", status: "降级", lastFetch: "1 hr ago", articles: 4523, interval: "1hr" },
  { id: "4", name: "Bloomberg Markets", url: "bloomberg.com/feed", category: "金融", status: "正常", lastFetch: "3 min ago", articles: 15672, interval: "15min" },
  { id: "5", name: "NASA Earth", url: "nasa.gov/rss/earth", category: "科学", status: "失败", lastFetch: "6 hr ago", articles: 892, interval: "6hr" },
  { id: "6", name: "TechCrunch", url: "techcrunch.com/feed", category: "科技", status: "正常", lastFetch: "8 min ago", articles: 6734, interval: "15min" },
  { id: "7", name: "Nature News", url: "nature.com/nature.rss", category: "学术", status: "正常", lastFetch: "30 min ago", articles: 3291, interval: "1hr" },
  { id: "8", name: "Al Jazeera", url: "aljazeera.com/xml/rss/all.xml", category: "新闻", status: "正常", lastFetch: "4 min ago", articles: 9845, interval: "15min" },
  { id: "9", name: "Reddit r/worldnews", url: "reddit.com/r/worldnews/.rss", category: "社区", status: "正常", lastFetch: "1 min ago", articles: 21342, interval: "5min" },
  { id: "10", name: "CoinDesk", url: "coindesk.com/arc/outboundfeeds/rss", category: "金融", status: "降级", lastFetch: "45 min ago", articles: 2143, interval: "30min" },
];

const STATUS_STYLE: Record<string, { dot: string; text: string }> = {
  正常: { dot: "bg-accent-green", text: "text-accent-green" },
  降级: { dot: "bg-accent-orange", text: "text-accent-orange" },
  失败: { dot: "bg-red-500", text: "text-red-500" },
};

const col = createColumnHelper<Source>();

export default function SourcesPage() {
  const [sorting, setSorting] = useState<SortingState>([]);
  const [globalFilter, setGlobalFilter] = useState("");
  const [selected, setSelected] = useState<Set<string>>(new Set());
  const [showAddModal, setShowAddModal] = useState(false);

  const columns = useMemo(
    () => [
      col.display({
        id: "select",
        header: () => (
          <input
            type="checkbox"
            checked={selected.size === MOCK_SOURCES.length}
            onChange={(e) =>
              setSelected(e.target.checked ? new Set(MOCK_SOURCES.map((s) => s.id)) : new Set())
            }
            className="accent-accent-blue"
          />
        ),
        cell: ({ row }) => (
          <input
            type="checkbox"
            checked={selected.has(row.original.id)}
            onChange={(e) => {
              const next = new Set(selected);
              e.target.checked ? next.add(row.original.id) : next.delete(row.original.id);
              setSelected(next);
            }}
            className="accent-accent-blue"
          />
        ),
      }),
      col.accessor("name", {
        header: "名称",
        cell: (info) => <span className="font-medium">{info.getValue()}</span>,
      }),
      col.accessor("url", {
        header: "URL",
        cell: (info) => <span className="text-muted-foreground">{info.getValue()}</span>,
      }),
      col.accessor("category", {
        header: "分类",
        cell: (info) => <span className="rounded bg-muted px-2 py-0.5 text-xs">{info.getValue()}</span>,
      }),
      col.accessor("status", {
        header: "状态",
        cell: (info) => {
          const s = info.getValue();
          return (
            <span className={cn("inline-flex items-center gap-1.5 text-xs", STATUS_STYLE[s].text)}>
              <span className={cn("h-1.5 w-1.5 rounded-full", STATUS_STYLE[s].dot)} />
              {s}
            </span>
          );
        },
      }),
      col.accessor("articles", {
        header: "文章数",
        cell: (info) => <span className="tabular-nums">{info.getValue().toLocaleString()}</span>,
      }),
      col.accessor("interval", { header: "频率" }),
      col.accessor("lastFetch", { header: "最后抓取" }),
    ],
    [selected]
  );

  const table = useReactTable({
    data: MOCK_SOURCES,
    columns,
    state: { sorting, globalFilter },
    onSortingChange: setSorting,
    onGlobalFilterChange: setGlobalFilter,
    getCoreRowModel: getCoreRowModel(),
    getSortedRowModel: getSortedRowModel(),
    getFilteredRowModel: getFilteredRowModel(),
  });

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold">数据源管理</h1>
        <button
          onClick={() => setShowAddModal(true)}
          className="inline-flex items-center gap-1.5 rounded-lg bg-accent-blue px-4 py-2 text-sm font-medium text-white hover:bg-accent-blue/90 transition-colors"
        >
          <Plus className="h-4 w-4" /> 添加数据源
        </button>
      </div>

      {/* Toolbar */}
      <div className="flex items-center gap-3">
        <input
          value={globalFilter}
          onChange={(e) => setGlobalFilter(e.target.value)}
          placeholder="搜索数据源..."
          className="flex-1 rounded-lg border border-border bg-white px-3 py-2 text-sm outline-none focus:border-accent-blue"
        />
        <button className="rounded-lg border border-border bg-white px-3 py-2 text-xs text-muted-foreground hover:bg-muted transition-colors">导入 OPML</button>
        <button className="rounded-lg border border-border bg-white px-3 py-2 text-xs text-muted-foreground hover:bg-muted transition-colors">导出 OPML</button>
        {selected.size > 0 && (
          <div className="flex items-center gap-2 ml-2">
            <span className="text-xs text-muted-foreground">已选 {selected.size}</span>
            <button className="rounded-lg border border-border px-2.5 py-1.5 text-[11px] text-muted-foreground hover:bg-muted">批量启用</button>
            <button className="rounded-lg border border-border px-2.5 py-1.5 text-[11px] text-red-500 hover:bg-red-50">批量删除</button>
          </div>
        )}
      </div>

      {/* Table */}
      <div className="rounded-xl border border-border bg-white overflow-hidden">
        <table className="w-full text-sm">
          <thead>
            {table.getHeaderGroups().map((hg) => (
              <tr key={hg.id} className="border-b border-border text-left text-xs text-muted-foreground">
                {hg.headers.map((header) => (
                  <th
                    key={header.id}
                    className="px-4 py-3 font-medium cursor-pointer select-none"
                    onClick={header.column.getToggleSortingHandler()}
                  >
                    <div className="flex items-center gap-1">
                      {flexRender(header.column.columnDef.header, header.getContext())}
                      {header.column.getCanSort() && <ArrowUpDown className="h-3 w-3 opacity-40" />}
                    </div>
                  </th>
                ))}
              </tr>
            ))}
          </thead>
          <tbody>
            {table.getRowModel().rows.map((row) => (
              <tr key={row.id} className="border-b border-border last:border-0 hover:bg-muted/30 transition-colors">
                {row.getVisibleCells().map((cell) => (
                  <td key={cell.id} className="px-4 py-3">
                    {flexRender(cell.column.columnDef.cell, cell.getContext())}
                  </td>
                ))}
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {/* Add Source Modal */}
      {showAddModal && (
        <div className="fixed inset-0 z-50 flex items-center justify-center">
          <div className="absolute inset-0 bg-black/40" onClick={() => setShowAddModal(false)} />
          <div className="relative z-10 w-full max-w-md rounded-xl border border-border bg-white p-6 shadow-2xl">
            <div className="flex items-center justify-between mb-4">
              <h3 className="text-sm font-semibold">添加数据源</h3>
              <button onClick={() => setShowAddModal(false)} className="text-muted-foreground hover:text-foreground">
                <X className="h-4 w-4" />
              </button>
            </div>
            <div className="space-y-3">
              <div>
                <label className="text-xs text-muted-foreground">RSS URL</label>
                <input placeholder="https://example.com/rss" className="mt-1 w-full rounded-lg border border-border px-3 py-2 text-sm outline-none focus:border-accent-blue" />
              </div>
              <div>
                <label className="text-xs text-muted-foreground">名称</label>
                <input placeholder="数据源名称" className="mt-1 w-full rounded-lg border border-border px-3 py-2 text-sm outline-none focus:border-accent-blue" />
              </div>
              <div className="grid grid-cols-2 gap-3">
                <div>
                  <label className="text-xs text-muted-foreground">分类</label>
                  <select className="mt-1 w-full rounded-lg border border-border px-3 py-2 text-sm outline-none">
                    <option>科技</option><option>金融</option><option>新闻</option><option>学术</option><option>社区</option>
                  </select>
                </div>
                <div>
                  <label className="text-xs text-muted-foreground">抓取频率</label>
                  <select className="mt-1 w-full rounded-lg border border-border px-3 py-2 text-sm outline-none">
                    <option>5 分钟</option><option>15 分钟</option><option>30 分钟</option><option>1 小时</option><option>6 小时</option>
                  </select>
                </div>
              </div>
              <button className="w-full mt-2 rounded-lg bg-accent-blue px-4 py-2 text-sm font-medium text-white hover:bg-accent-blue/90 transition-colors">
                添加
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
