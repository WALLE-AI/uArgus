"use client";

import { useState } from "react";
import { Search, FolderOpen, FileText, Plus, X } from "lucide-react";
import { cn } from "@/lib/utils";

interface GraphNode {
  id: string;
  label: string;
  type: "person" | "org" | "event" | "topic";
  x: number;
  y: number;
  articles: number;
}

interface GraphEdge {
  from: string;
  to: string;
}

const NODES: GraphNode[] = [
  { id: "1", label: "AI 监管", type: "topic", x: 50, y: 40, articles: 142 },
  { id: "2", label: "欧盟委员会", type: "org", x: 25, y: 25, articles: 87 },
  { id: "3", label: "OpenAI", type: "org", x: 70, y: 20, articles: 65 },
  { id: "4", label: "半导体管制", type: "topic", x: 80, y: 55, articles: 98 },
  { id: "5", label: "台积电", type: "org", x: 90, y: 35, articles: 43 },
  { id: "6", label: "地缘政治", type: "topic", x: 55, y: 70, articles: 87 },
  { id: "7", label: "能源转型", type: "topic", x: 20, y: 65, articles: 72 },
  { id: "8", label: "气候变化", type: "event", x: 30, y: 80, articles: 48 },
  { id: "9", label: "Sam Altman", type: "person", x: 65, y: 30, articles: 34 },
  { id: "10", label: "量子计算", type: "topic", x: 40, y: 55, articles: 54 },
  { id: "11", label: "Google", type: "org", x: 45, y: 25, articles: 29 },
  { id: "12", label: "COP 大会", type: "event", x: 15, y: 50, articles: 56 },
];

const EDGES: GraphEdge[] = [
  { from: "1", to: "2" }, { from: "1", to: "3" }, { from: "1", to: "9" },
  { from: "3", to: "9" }, { from: "3", to: "11" }, { from: "4", to: "5" },
  { from: "4", to: "6" }, { from: "6", to: "1" }, { from: "7", to: "8" },
  { from: "7", to: "12" }, { from: "8", to: "12" }, { from: "10", to: "11" },
  { from: "6", to: "4" }, { from: "1", to: "10" },
];

const TYPE_COLOR: Record<string, string> = {
  person: "#EC4899",
  org: "#3B82F6",
  event: "#F59E0B",
  topic: "#10B981",
};
const TYPE_LABEL: Record<string, string> = {
  person: "人物",
  org: "组织",
  event: "事件",
  topic: "主题",
};

const COLLECTIONS = [
  { name: "AI 政策合集", count: 23, icon: FolderOpen },
  { name: "半导体产业链", count: 18, icon: FolderOpen },
  { name: "气候变化追踪", count: 14, icon: FolderOpen },
  { name: "量子计算前沿", count: 9, icon: FolderOpen },
];

const NOTES = [
  { title: "欧盟 AI 法案要点整理", updated: "2 hr ago", words: 1240 },
  { title: "美联储政策时间线", updated: "5 hr ago", words: 860 },
  { title: "能源转型路线图", updated: "1 day ago", words: 2100 },
  { title: "半导体出口管制对比表", updated: "2 days ago", words: 650 },
];

export default function KnowledgePage() {
  const [selectedNode, setSelectedNode] = useState<GraphNode | null>(null);
  const [search, setSearch] = useState("");

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold">知识沉淀</h1>
        <button className="inline-flex items-center gap-1.5 rounded-lg bg-accent-blue px-4 py-2 text-sm font-medium text-white hover:bg-accent-blue/90 transition-colors">
          <Plus className="h-4 w-4" /> 新建笔记
        </button>
      </div>

      {/* Search */}
      <div className="relative">
        <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
        <input
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          placeholder="搜索知识库（全文 + 语义搜索）..."
          className="w-full rounded-lg border border-border bg-white pl-9 pr-3 py-2 text-sm outline-none focus:border-accent-blue"
        />
      </div>

      <div className="grid gap-4 lg:grid-cols-3">
        {/* Knowledge Graph SVG */}
        <div className="lg:col-span-2 rounded-xl border border-border bg-white p-6">
          <div className="flex items-center justify-between mb-4">
            <h3 className="text-sm font-semibold">知识图谱</h3>
            <div className="flex gap-3">
              {Object.entries(TYPE_LABEL).map(([type, label]) => (
                <div key={type} className="flex items-center gap-1">
                  <span className="h-2 w-2 rounded-full" style={{ backgroundColor: TYPE_COLOR[type] }} />
                  <span className="text-[10px] text-muted-foreground">{label}</span>
                </div>
              ))}
            </div>
          </div>
          <div className="relative">
            <svg viewBox="0 0 100 100" className="w-full h-80">
              {/* Edges */}
              {EDGES.map((e, i) => {
                const from = NODES.find((n) => n.id === e.from)!;
                const to = NODES.find((n) => n.id === e.to)!;
                return (
                  <line key={i} x1={from.x} y1={from.y} x2={to.x} y2={to.y}
                    stroke="hsl(var(--border))" strokeWidth="0.3" />
                );
              })}
              {/* Nodes */}
              {NODES.map((node) => {
                const r = 2 + (node.articles / 142) * 3;
                return (
                  <g key={node.id} className="cursor-pointer" onClick={() => setSelectedNode(node)}>
                    <circle cx={node.x} cy={node.y} r={r}
                      fill={TYPE_COLOR[node.type]}
                      opacity={selectedNode?.id === node.id ? 1 : 0.7}
                      stroke={selectedNode?.id === node.id ? TYPE_COLOR[node.type] : "none"}
                      strokeWidth="0.5"
                    />
                    <text x={node.x} y={node.y + r + 2.5}
                      textAnchor="middle" fontSize="2.2"
                      fill="hsl(var(--muted-foreground))"
                    >
                      {node.label}
                    </text>
                  </g>
                );
              })}
            </svg>
          </div>

          {/* Detail panel */}
          {selectedNode && (
            <div className="mt-4 rounded-lg border border-border p-4">
              <div className="flex items-center justify-between">
                <div className="flex items-center gap-2">
                  <span className="h-3 w-3 rounded-full" style={{ backgroundColor: TYPE_COLOR[selectedNode.type] }} />
                  <span className="text-sm font-semibold">{selectedNode.label}</span>
                  <span className="rounded bg-muted px-1.5 py-0.5 text-[10px] text-muted-foreground">{TYPE_LABEL[selectedNode.type]}</span>
                </div>
                <button onClick={() => setSelectedNode(null)} className="text-muted-foreground hover:text-foreground">
                  <X className="h-4 w-4" />
                </button>
              </div>
              <div className="mt-2 text-xs text-muted-foreground">
                关联文章: {selectedNode.articles} 篇 · 关联实体: {EDGES.filter((e) => e.from === selectedNode.id || e.to === selectedNode.id).length} 个
              </div>
              <div className="mt-3 flex gap-2">
                <button className="rounded-lg border border-border px-3 py-1.5 text-xs hover:bg-muted transition-colors">查看文章</button>
                <button className="rounded-lg border border-border px-3 py-1.5 text-xs hover:bg-muted transition-colors">查看关联</button>
                <button className="rounded-lg border border-border px-3 py-1.5 text-xs hover:bg-muted transition-colors">添加笔记</button>
              </div>
            </div>
          )}
        </div>

        {/* Sidebar */}
        <div className="space-y-4">
          {/* Collections */}
          <div className="rounded-xl border border-border bg-white p-5">
            <div className="flex items-center justify-between mb-3">
              <h3 className="text-sm font-semibold">收藏夹</h3>
              <button className="text-muted-foreground hover:text-foreground"><Plus className="h-3.5 w-3.5" /></button>
            </div>
            <div className="space-y-1">
              {COLLECTIONS.map((c) => (
                <div key={c.name} className="flex items-center gap-2 rounded-lg p-2 hover:bg-muted/50 cursor-pointer transition-colors">
                  <c.icon className="h-3.5 w-3.5 text-muted-foreground" />
                  <span className="flex-1 text-sm">{c.name}</span>
                  <span className="text-[10px] text-muted-foreground">{c.count}</span>
                </div>
              ))}
            </div>
          </div>

          {/* Recent notes */}
          <div className="rounded-xl border border-border bg-white p-5">
            <h3 className="text-sm font-semibold mb-3">最近笔记</h3>
            <div className="space-y-1">
              {NOTES.map((n) => (
                <div key={n.title} className="flex items-start gap-2 rounded-lg p-2 hover:bg-muted/50 cursor-pointer transition-colors">
                  <FileText className="h-3.5 w-3.5 text-muted-foreground mt-0.5" />
                  <div className="flex-1 min-w-0">
                    <div className="text-sm truncate">{n.title}</div>
                    <div className="text-[10px] text-muted-foreground">{n.updated} · {n.words} 字</div>
                  </div>
                </div>
              ))}
            </div>
          </div>

          {/* Stats */}
          <div className="rounded-xl border border-border bg-white p-5">
            <h3 className="text-sm font-semibold mb-3">知识库统计</h3>
            <div className="grid grid-cols-2 gap-3">
              {[
                { label: "实体数", value: "1,247" },
                { label: "关系数", value: "3,891" },
                { label: "笔记数", value: "86" },
                { label: "收藏文章", value: "342" },
              ].map((s) => (
                <div key={s.label} className="text-center p-2 rounded-lg bg-muted/30">
                  <div className="text-sm font-bold">{s.value}</div>
                  <div className="text-[10px] text-muted-foreground">{s.label}</div>
                </div>
              ))}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
