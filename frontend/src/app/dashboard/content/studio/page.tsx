"use client";

import { useState, useEffect } from "react";
import Markdown from "react-markdown";
import { BrainCircuit, Eye, Edit3, Save, Send, FileText, Sparkles, Lightbulb, TrendingUp } from "lucide-react";
import { cn } from "@/lib/utils";
import { fetchInsights, fetchTopicSuggestions, fetchDigest } from "@/lib/api";
import type { TopicSuggestion, SemanticFeedItem } from "@/lib/types";
import { severityToLabel } from "@/lib/types";
import { getCategoryLabel } from "@/lib/mock-data";
import { useAppStore } from "@/lib/store";

const PLACEHOLDER_CONTENT = `## 新建文章

从右侧「AI 选题推荐」中选择一个主题开始写作，或直接在此编辑。

语义层已自动分析当前热点主题并生成推荐选题。
`;

const AI_ACTIONS = [
  { label: "基于语义分析生成大纲", icon: Sparkles, desc: "根据语义层分类和关联素材自动生成结构化大纲" },
  { label: "AI 续写段落", icon: Edit3, desc: "基于上下文和跨源文章续写" },
  { label: "改写润色", icon: Sparkles, desc: "优化文风和表达" },
  { label: "生成多源摘要", icon: BrainCircuit, desc: "基于佐证文章提取核心要点" },
  { label: "SEO 优化", icon: TrendingUp, desc: "优化搜索引擎关键词" },
  { label: "多平台适配", icon: Send, desc: "生成微信/知乎/小红书版本" },
];


export default function StudioPage() {
  const [title, setTitle] = useState("");
  const [content, setContent] = useState(PLACEHOLDER_CONTENT);
  const [mode, setMode] = useState<"edit" | "preview" | "split">("split");

  // ── Load semantic data ──
  const { insights, setInsights, insightsLoading, setInsightsLoading, topicSuggestions, setTopicSuggestions, feedItems, setFeedItems } = useAppStore();
  const [relatedArticles, setRelatedArticles] = useState<SemanticFeedItem[]>([]);
  const [insightsFromApi, setInsightsFromApi] = useState(false);

  useEffect(() => {
    setInsightsLoading(true);
    Promise.all([
      fetchInsights(),
      fetchTopicSuggestions(),
      fetchDigest({ variant: "full", lang: "en" }),
    ]).then(([insRes, topicRes, digestRes]) => {
      setInsights(insRes.data);
      setInsightsFromApi(insRes.fromApi);
      setTopicSuggestions(topicRes.data);
      if (digestRes.items.length > 0 && feedItems.length === 0) {
        setFeedItems(digestRes.items);
      }
      setInsightsLoading(false);
    });
  }, [setInsights, setInsightsLoading, setTopicSuggestions, setFeedItems, feedItems.length]);

  // Find related articles when a topic is selected
  const handleSelectTopic = (topic: TopicSuggestion) => {
    setTitle(topic.title);
    // Generate initial content from insights
    const header = `## ${topic.title}\n\n`;
    const reason = `> 选题依据: ${topic.reason}\n\n`;
    const summary = insights ? `### 语义层摘要\n\n${insights.summary}\n\n` : "";
    const categories = `### 涉及领域\n\n${topic.categories.map((c) => `- ${getCategoryLabel(c)}`).join("\n")}\n\n`;
    const related = `### 关联素材\n\n基于语义分析，已匹配 ${topic.relatedCount} 篇相关文章。\n`;
    setContent(header + reason + summary + categories + related);

    // Find related feed items by matching categories
    const matched = feedItems
      .filter((item) => item.categories.some((c) => topic.categories.includes(c)))
      .sort((a, b) => b.importanceScore - a.importanceScore)
      .slice(0, 6);
    setRelatedArticles(matched);
  };

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-3">
          <h1 className="text-2xl font-bold">内容中心</h1>
          {!insightsFromApi && !insightsLoading && (
            <span className="rounded bg-amber-100 px-2 py-0.5 text-[10px] font-medium text-amber-700">
              Mock 数据
            </span>
          )}
        </div>
        <div className="flex items-center gap-2">
          <button className="inline-flex items-center gap-1.5 rounded-lg border border-border px-3 py-2 text-xs text-muted-foreground hover:bg-muted transition-colors">
            <Save className="h-3.5 w-3.5" /> 保存草稿
          </button>
          <button className="inline-flex items-center gap-1.5 rounded-lg bg-accent-blue px-4 py-2 text-xs font-medium text-white hover:bg-accent-blue/90 transition-colors">
            <Send className="h-3.5 w-3.5" /> 发布到分发
          </button>
        </div>
      </div>

      {/* Insights Summary Banner */}
      {insights && (
        <div className="rounded-xl border border-accent-blue/20 bg-accent-blue/5 p-4">
          <div className="flex items-center gap-1.5 text-xs font-medium text-accent-blue mb-2">
            <BrainCircuit className="h-3.5 w-3.5" />
            语义层智能摘要
            {insights.generatedAt && (
              <span className="ml-auto text-[10px] text-muted-foreground font-normal">
                {new Date(insights.generatedAt).toLocaleString("zh-CN")}
              </span>
            )}
          </div>
          <p className="text-sm leading-relaxed text-foreground">{insights.summary}</p>
        </div>
      )}

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
          {/* AI Topic Suggestions (from semantic layer) */}
          <div className="rounded-xl border border-border bg-white p-4">
            <div className="flex items-center gap-1.5 text-sm font-semibold mb-3">
              <Lightbulb className="h-4 w-4 text-accent-orange" /> AI 选题推荐
            </div>
            {insightsLoading ? (
              <div className="text-xs text-muted-foreground py-4 text-center">加载语义分析...</div>
            ) : (
              <div className="space-y-1.5">
                {topicSuggestions.map((t) => (
                  <button
                    key={t.title}
                    onClick={() => handleSelectTopic(t)}
                    className="w-full text-left rounded-lg border border-border px-3 py-2 hover:bg-muted/50 transition-colors"
                  >
                    <div className="flex items-center justify-between">
                      <span className="text-xs font-medium truncate flex-1">{t.title}</span>
                      <span className="text-[10px] font-bold text-accent-blue ml-2">{t.score}</span>
                    </div>
                    <div className="text-[10px] text-muted-foreground mt-0.5 line-clamp-1">{t.reason}</div>
                    <div className="flex gap-1 mt-1">
                      {t.categories.slice(0, 3).map((c) => (
                        <span key={c} className="rounded bg-muted px-1.5 py-0.5 text-[9px] text-muted-foreground">
                          {getCategoryLabel(c)}
                        </span>
                      ))}
                    </div>
                  </button>
                ))}
              </div>
            )}
          </div>

          {/* AI Writing Actions */}
          <div className="rounded-xl border border-border bg-white p-4">
            <div className="flex items-center gap-1.5 text-sm font-semibold mb-3">
              <BrainCircuit className="h-4 w-4 text-accent-blue" /> AI 写作助手
            </div>
            <div className="space-y-1.5">
              {AI_ACTIONS.map((a) => (
                <button key={a.label} className="w-full text-left rounded-lg border border-border px-3 py-2 hover:bg-muted/50 transition-colors group">
                  <div className="flex items-center gap-2">
                    <a.icon className="h-3.5 w-3.5 text-muted-foreground group-hover:text-accent-blue transition-colors" />
                    <span className="text-xs font-medium">{a.label}</span>
                  </div>
                  <div className="text-[10px] text-muted-foreground mt-0.5 ml-6">{a.desc}</div>
                </button>
              ))}
            </div>
          </div>

          {/* Related articles from semantic feed */}
          <div className="rounded-xl border border-border bg-white p-4">
            <h3 className="text-sm font-semibold mb-3">关联素材（语义匹配）</h3>
            {relatedArticles.length === 0 ? (
              <div className="text-[11px] text-muted-foreground py-2">选择选题后自动匹配关联文章</div>
            ) : (
              <div className="space-y-1.5">
                {relatedArticles.map((a) => (
                  <div key={a.hash} className="flex items-start gap-2 rounded-lg p-2 hover:bg-muted/50 cursor-pointer transition-colors">
                    <FileText className="h-3.5 w-3.5 text-muted-foreground mt-0.5 shrink-0" />
                    <div className="flex-1 min-w-0">
                      <div className="text-[11px] font-medium truncate">{a.title}</div>
                      <div className="text-[10px] text-muted-foreground">
                        {a.source} · {severityToLabel(a.severity)} · 分数 {Math.round(a.importanceScore * 100)}
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
