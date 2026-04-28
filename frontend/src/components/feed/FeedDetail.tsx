"use client";

import { X, Bookmark, BrainCircuit, Share2, BookOpen, ExternalLink } from "lucide-react";
import { cn } from "@/lib/utils";
import type { FeedItem } from "./FeedCard";

interface FeedDetailProps {
  item: FeedItem | null;
  onClose: () => void;
}

const IMPORTANCE_STYLE: Record<string, string> = {
  紧急: "bg-red-100 text-red-700",
  高: "bg-orange-100 text-orange-700",
  中: "bg-yellow-100 text-yellow-700",
  低: "bg-muted text-muted-foreground",
};

export default function FeedDetail({ item, onClose }: FeedDetailProps) {
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
              <span className={cn("rounded px-2 py-0.5 text-[11px] font-medium", IMPORTANCE_STYLE[item.importance])}>
                {item.importance}
              </span>
              <span className="text-xs text-muted-foreground">{item.source} · {item.time}</span>
            </div>
            <div className="flex items-center gap-1">
              <button className="rounded p-1.5 text-muted-foreground hover:bg-muted hover:text-foreground transition-colors">
                <ExternalLink className="h-4 w-4" />
              </button>
              <button onClick={onClose} className="rounded p-1.5 text-muted-foreground hover:bg-muted hover:text-foreground transition-colors">
                <X className="h-4 w-4" />
              </button>
            </div>
          </div>

          {/* Content */}
          <div className="flex-1 overflow-y-auto px-6 py-6">
            <h1 className="text-xl font-bold leading-tight">{item.title}</h1>
            <div className="mt-3 flex gap-1.5">
              {item.tags.map((tag) => (
                <span key={tag} className="rounded bg-muted px-2 py-0.5 text-[11px] text-muted-foreground">{tag}</span>
              ))}
            </div>

            {/* AI Summary */}
            <div className="mt-6 rounded-lg border border-accent-blue/20 bg-accent-blue/5 p-4">
              <div className="flex items-center gap-1.5 text-xs font-medium text-accent-blue mb-2">
                <BrainCircuit className="h-3.5 w-3.5" />
                AI 摘要
              </div>
              <p className="text-sm leading-relaxed text-foreground">{item.summary}</p>
            </div>

            {/* Article body placeholder */}
            <div className="mt-6 space-y-4 text-sm leading-relaxed text-foreground/80">
              <p>
                这是文章的详细内容区域。在实际版本中，这里会显示完整的文章正文，支持富文本渲染、
                图片展示和链接跳转。AI 分析模块会提供实体识别、关键事件提取和情感分析结果。
              </p>
              <p>
                文章来源将通过 RSS 抓取或 API 同步获取原始内容，经过清洗和结构化处理后展示。
                用户可以对文章进行标注、收藏、添加到知识库等操作。
              </p>
              <p>
                相关文章推荐基于语义相似度和主题关联性，帮助用户快速获取同一事件的多源信息，
                实现交叉验证和全面理解。
              </p>
            </div>

            {/* Related articles */}
            <div className="mt-8">
              <h3 className="text-sm font-semibold mb-3">相关文章</h3>
              <div className="space-y-2">
                {["相关分析报告：政策影响评估", "多源比较：各媒体报道角度差异", "历史回顾：类似事件发展路径"].map((t) => (
                  <div key={t} className="flex items-center gap-2 rounded-lg p-2 hover:bg-muted/50 cursor-pointer transition-colors">
                    <div className="h-1.5 w-1.5 rounded-full bg-accent-blue shrink-0" />
                    <span className="text-xs">{t}</span>
                  </div>
                ))}
              </div>
            </div>
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
