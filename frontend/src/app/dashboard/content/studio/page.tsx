"use client";

import { useState } from "react";
import Markdown from "react-markdown";
import { BrainCircuit, Eye, Edit3, Save, Send, FileText } from "lucide-react";
import { cn } from "@/lib/utils";

const DEFAULT_CONTENT = `## 欧盟 AI 法案对中国科技企业的影响分析

### 背景

2024年，欧盟正式通过了全球首部全面的人工智能监管法案（AI Act）。该法案采用基于风险的分级监管方法，将 AI 系统分为四个风险等级。

### 核心影响

1. **高风险系统合规成本上升** — 中国企业在欧盟市场部署的 AI 系统需满足严格的合规要求
2. **数据治理挑战** — 训练数据的透明度和可追溯性要求
3. **市场准入门槛提高** — 特别是在生物识别、自动驾驶等领域

### 企业应对策略

- 建立内部 AI 治理框架
- 加强与欧盟本土合作伙伴的协同
- 投资合规技术和工具链建设

> 本文基于 47 篇跨源文章的 AI 综合分析生成，情感倾向偏中性。
`;

const AI_ACTIONS = [
  { label: "生成大纲", icon: "📋", desc: "基于当前主题生成结构化大纲" },
  { label: "续写段落", icon: "✍️", desc: "AI 续写下一个段落" },
  { label: "改写润色", icon: "✨", desc: "优化文风和表达" },
  { label: "生成摘要", icon: "📝", desc: "提取核心要点摘要" },
  { label: "SEO 优化", icon: "🔍", desc: "优化搜索引擎关键词" },
  { label: "多平台适配", icon: "📱", desc: "生成微信/知乎/小红书版本" },
];

const RELATED_ARTICLES = [
  { title: "欧盟委员会发布 AI 法案实施指南", source: "Reuters", relevance: 95 },
  { title: "中国企业如何应对欧盟数据合规", source: "FT", relevance: 88 },
  { title: "全球 AI 监管政策对比分析", source: "MIT Tech Review", relevance: 82 },
  { title: "AI 治理框架最佳实践白皮书", source: "IEEE", relevance: 76 },
];

export default function StudioPage() {
  const [title, setTitle] = useState("欧盟 AI 法案对中国科技企业的影响分析");
  const [content, setContent] = useState(DEFAULT_CONTENT);
  const [mode, setMode] = useState<"edit" | "preview" | "split">("split");

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold">内容中心</h1>
        <div className="flex items-center gap-2">
          <button className="inline-flex items-center gap-1.5 rounded-lg border border-border px-3 py-2 text-xs text-muted-foreground hover:bg-muted transition-colors">
            <Save className="h-3.5 w-3.5" /> 保存草稿
          </button>
          <button className="inline-flex items-center gap-1.5 rounded-lg bg-accent-blue px-4 py-2 text-xs font-medium text-white hover:bg-accent-blue/90 transition-colors">
            <Send className="h-3.5 w-3.5" /> 发布到分发
          </button>
        </div>
      </div>

      <div className="grid gap-4 lg:grid-cols-4">
        {/* Editor area */}
        <div className="lg:col-span-3 rounded-xl border border-border bg-white">
          {/* Toolbar */}
          <div className="flex items-center justify-between border-b border-border px-4 py-2">
            <input
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              className="text-lg font-bold outline-none flex-1 placeholder:text-muted-foreground/50"
              placeholder="文章标题..."
            />
            <div className="flex rounded-lg border border-border overflow-hidden">
              {([
                { key: "edit" as const, icon: Edit3, label: "编辑" },
                { key: "split" as const, icon: Eye, label: "分屏" },
                { key: "preview" as const, icon: Eye, label: "预览" },
              ]).map((m) => (
                <button
                  key={m.key}
                  onClick={() => setMode(m.key)}
                  className={cn(
                    "px-2.5 py-1 text-[11px]",
                    mode === m.key ? "bg-muted text-foreground" : "text-muted-foreground hover:bg-muted/50"
                  )}
                >
                  {m.label}
                </button>
              ))}
            </div>
          </div>

          {/* Content area */}
          <div className={cn("flex", mode === "split" ? "divide-x divide-border" : "")}>
            {(mode === "edit" || mode === "split") && (
              <textarea
                value={content}
                onChange={(e) => setContent(e.target.value)}
                className={cn(
                  "min-h-[500px] p-4 text-sm font-mono outline-none resize-none bg-transparent",
                  mode === "split" ? "w-1/2" : "w-full"
                )}
              />
            )}
            {(mode === "preview" || mode === "split") && (
              <div className={cn(
                "min-h-[500px] p-4 prose prose-sm max-w-none overflow-y-auto",
                mode === "split" ? "w-1/2" : "w-full"
              )}>
                <Markdown>{content}</Markdown>
              </div>
            )}
          </div>

          {/* Word count bar */}
          <div className="border-t border-border px-4 py-2 flex items-center justify-between text-[11px] text-muted-foreground">
            <span>{content.length} 字符 · {content.split("\n").length} 行</span>
            <span>Markdown 模式</span>
          </div>
        </div>

        {/* Sidebar */}
        <div className="space-y-4">
          {/* AI actions */}
          <div className="rounded-xl border border-border bg-white p-4">
            <div className="flex items-center gap-1.5 text-sm font-semibold mb-3">
              <BrainCircuit className="h-4 w-4 text-accent-blue" /> AI 写作助手
            </div>
            <div className="space-y-1.5">
              {AI_ACTIONS.map((a) => (
                <button key={a.label} className="w-full text-left rounded-lg border border-border px-3 py-2 hover:bg-muted/50 transition-colors group">
                  <div className="flex items-center gap-2">
                    <span className="text-xs">{a.icon}</span>
                    <span className="text-xs font-medium">{a.label}</span>
                  </div>
                  <div className="text-[10px] text-muted-foreground mt-0.5 ml-6">{a.desc}</div>
                </button>
              ))}
            </div>
          </div>

          {/* Related articles */}
          <div className="rounded-xl border border-border bg-white p-4">
            <h3 className="text-sm font-semibold mb-3">关联素材</h3>
            <div className="space-y-1.5">
              {RELATED_ARTICLES.map((a) => (
                <div key={a.title} className="flex items-start gap-2 rounded-lg p-2 hover:bg-muted/50 cursor-pointer transition-colors">
                  <FileText className="h-3.5 w-3.5 text-muted-foreground mt-0.5 shrink-0" />
                  <div className="flex-1 min-w-0">
                    <div className="text-[11px] font-medium truncate">{a.title}</div>
                    <div className="text-[10px] text-muted-foreground">{a.source} · 相关度 {a.relevance}%</div>
                  </div>
                </div>
              ))}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
