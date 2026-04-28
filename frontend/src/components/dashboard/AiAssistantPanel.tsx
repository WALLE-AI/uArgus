"use client";

import { useState } from "react";
import { X, Send, Sparkles, FileText, Search, Share2 } from "lucide-react";
import { cn } from "@/lib/utils";
import { useAppStore } from "@/lib/store";

interface Message {
  id: string;
  role: "user" | "assistant";
  content: string;
}

const QUICK_ACTIONS = [
  { label: "摘要本文", icon: FileText },
  { label: "查找关联事件", icon: Search },
  { label: "生成社媒帖子", icon: Share2 },
  { label: "分析情感趋势", icon: Sparkles },
];

const INITIAL_MESSAGES: Message[] = [
  {
    id: "1",
    role: "assistant",
    content: "你好！我是 uArgus AI 助手。我可以帮你摘要文章、分析趋势、生成内容，或回答你关于当前数据的问题。有什么需要帮助的吗？",
  },
];

export default function AiAssistantPanel() {
  const { assistantOpen, setAssistantOpen } = useAppStore();
  const [messages, setMessages] = useState<Message[]>(INITIAL_MESSAGES);
  const [input, setInput] = useState("");

  if (!assistantOpen) return null;

  const handleSend = () => {
    if (!input.trim()) return;
    const userMsg: Message = {
      id: Date.now().toString(),
      role: "user",
      content: input,
    };
    const aiMsg: Message = {
      id: (Date.now() + 1).toString(),
      role: "assistant",
      content: "我正在分析你的问题...（这是一个静态原型，实际版本会通过 AI 流式回复）",
    };
    setMessages((prev) => [...prev, userMsg, aiMsg]);
    setInput("");
  };

  return (
    <aside className="flex h-full w-[380px] shrink-0 flex-col border-l border-border bg-background">
      {/* Header */}
      <div className="flex items-center justify-between border-b border-border px-4 py-3">
        <div className="flex items-center gap-2">
          <Sparkles className="h-4 w-4 text-accent-blue" />
          <span className="text-sm font-semibold">AI 助手</span>
        </div>
        <button
          onClick={() => setAssistantOpen(false)}
          className="rounded-md p-1 text-muted-foreground hover:bg-muted hover:text-foreground transition-colors"
        >
          <X className="h-4 w-4" />
        </button>
      </div>

      {/* Quick actions */}
      <div className="flex flex-wrap gap-1.5 border-b border-border px-4 py-3">
        {QUICK_ACTIONS.map((action) => (
          <button
            key={action.label}
            className="inline-flex items-center gap-1 rounded-full border border-border px-2.5 py-1 text-[11px] text-muted-foreground hover:bg-muted hover:text-foreground transition-colors"
          >
            <action.icon className="h-3 w-3" />
            {action.label}
          </button>
        ))}
      </div>

      {/* Messages */}
      <div className="flex-1 overflow-y-auto px-4 py-4 space-y-4">
        {messages.map((msg) => (
          <div
            key={msg.id}
            className={cn(
              "flex",
              msg.role === "user" ? "justify-end" : "justify-start"
            )}
          >
            <div
              className={cn(
                "max-w-[85%] rounded-2xl px-4 py-2.5 text-sm leading-relaxed",
                msg.role === "user"
                  ? "rounded-br-md bg-accent-blue text-white"
                  : "rounded-bl-md bg-muted text-foreground"
              )}
            >
              {msg.content}
            </div>
          </div>
        ))}
      </div>

      {/* Input */}
      <div className="border-t border-border p-3">
        <div className="flex items-center gap-2 rounded-xl border border-border bg-muted/30 px-3 py-2">
          <input
            type="text"
            value={input}
            onChange={(e) => setInput(e.target.value)}
            onKeyDown={(e) => e.key === "Enter" && handleSend()}
            placeholder="输入消息..."
            className="flex-1 bg-transparent text-sm outline-none placeholder:text-muted-foreground"
          />
          <button
            onClick={handleSend}
            disabled={!input.trim()}
            className={cn(
              "rounded-lg p-1.5 transition-colors",
              input.trim()
                ? "bg-accent-blue text-white hover:bg-accent-blue/90"
                : "text-muted-foreground"
            )}
          >
            <Send className="h-4 w-4" />
          </button>
        </div>
        <p className="mt-1.5 text-center text-[10px] text-muted-foreground">
          AI 回复仅供参考，请核实关键信息
        </p>
      </div>
    </aside>
  );
}
