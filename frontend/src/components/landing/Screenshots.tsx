"use client";

import { useState } from "react";
import { motion, AnimatePresence } from "framer-motion";
import { LayoutDashboard, BotMessageSquare, BarChart3, Share2 } from "lucide-react";
import { cn } from "@/lib/utils";
import { useI18n } from "@/lib/i18n";

const TABS = [
  {
    id: "dashboard",
    labelKey: "screenshots.dashboard",
    descKey: "screenshots.dashboard.desc",
    icon: LayoutDashboard,
    mockContent: (
      <div className="grid grid-cols-4 gap-3 p-6">
        {/* Metric cards */}
        {["Sources Active", "Articles Today", "Trending Topics", "Alerts"].map((label, i) => (
          <div key={label} className="rounded-xl bg-white/10 p-4">
            <div className="text-2xl font-bold text-white">{[3247, 1842, 23, 7][i]}</div>
            <div className="text-xs text-white/50">{label}</div>
          </div>
        ))}
        {/* Chart placeholder */}
        <div className="col-span-3 h-40 rounded-xl bg-white/5 border border-white/10 flex items-center justify-center">
          <div className="flex items-end gap-1">
            {[40, 65, 45, 80, 55, 70, 90, 60, 75, 85, 50, 95].map((h, i) => (
              <div key={i} className="w-4 rounded-t bg-accent-blue/60" style={{ height: `${h}%` }} />
            ))}
          </div>
        </div>
        {/* Side widget */}
        <div className="h-40 rounded-xl bg-white/5 border border-white/10 p-3">
          <div className="text-xs text-white/40 mb-2">Source Health</div>
          {[92, 87, 95].map((v, i) => (
            <div key={i} className="flex items-center gap-2 mb-1.5">
              <div className="h-1.5 flex-1 rounded-full bg-white/10">
                <div className="h-1.5 rounded-full bg-accent-green/70" style={{ width: `${v}%` }} />
              </div>
              <span className="text-[10px] text-white/40">{v}%</span>
            </div>
          ))}
        </div>
      </div>
    ),
  },
  {
    id: "assistant",
    labelKey: "screenshots.assistant",
    descKey: "screenshots.assistant.desc",
    icon: BotMessageSquare,
    mockContent: (
      <div className="flex h-full p-6 gap-4">
        <div className="flex-1 flex flex-col gap-3">
          <div className="self-end max-w-[70%] rounded-2xl rounded-br-md bg-accent-blue/80 px-4 py-2.5 text-sm text-white">
            Summarize the latest geopolitics trends from today&apos;s feed
          </div>
          <div className="self-start max-w-[80%] rounded-2xl rounded-bl-md bg-white/10 px-4 py-2.5 text-sm text-white/80">
            Based on 47 articles from today, here are the key geopolitical trends:
            <br /><br />
            <strong className="text-white">1. US-China Trade</strong> — New tariff negotiations...
            <br />
            <strong className="text-white">2. EU Energy Policy</strong> — Green transition...
            <br />
            <strong className="text-white">3. Middle East</strong> — Diplomatic shifts...
          </div>
          <div className="mt-auto flex gap-2">
            {["Analyze deeper", "Generate article", "Push to WeChat"].map((a) => (
              <button key={a} className="rounded-full border border-white/20 px-3 py-1 text-xs text-white/60">{a}</button>
            ))}
          </div>
        </div>
      </div>
    ),
  },
  {
    id: "analytics",
    labelKey: "screenshots.analytics",
    descKey: "screenshots.analytics.desc",
    icon: BarChart3,
    mockContent: (
      <div className="grid grid-cols-2 gap-3 p-6">
        <div className="rounded-xl bg-white/5 border border-white/10 p-4">
          <div className="text-xs text-white/40 mb-3">Topic Clusters</div>
          <div className="flex flex-wrap gap-2">
            {[
              { name: "AI Regulation", size: "text-lg" },
              { name: "Energy", size: "text-base" },
              { name: "Climate", size: "text-sm" },
              { name: "Crypto", size: "text-xs" },
              { name: "Trade", size: "text-base" },
              { name: "Defense", size: "text-sm" },
            ].map((t) => (
              <span key={t.name} className={`${t.size} text-accent-blue/80 font-medium`}>{t.name}</span>
            ))}
          </div>
        </div>
        <div className="rounded-xl bg-white/5 border border-white/10 p-4">
          <div className="text-xs text-white/40 mb-3">Sentiment Trend</div>
          <div className="flex items-end gap-1 h-20">
            {[60, 55, 65, 40, 35, 50, 45, 55, 60, 70, 65, 75].map((v, i) => (
              <div key={i} className="flex-1 rounded-t" style={{ height: `${v}%`, background: v > 50 ? "rgba(16,185,129,0.5)" : "rgba(239,68,68,0.5)" }} />
            ))}
          </div>
        </div>
        <div className="col-span-2 rounded-xl bg-white/5 border border-white/10 p-4">
          <div className="text-xs text-white/40 mb-3">Event Tracking</div>
          <div className="flex items-center gap-3 text-xs">
            {["BREAKING", "DEVELOPING", "SUSTAINED", "FADING"].map((s, i) => (
              <div key={s} className="flex items-center gap-1.5">
                <div className={cn("h-2 w-2 rounded-full", i === 0 ? "bg-red-400" : i === 1 ? "bg-yellow-400" : i === 2 ? "bg-accent-blue" : "bg-white/30")} />
                <span className="text-white/50">{s}</span>
              </div>
            ))}
          </div>
        </div>
      </div>
    ),
  },
  {
    id: "distribution",
    labelKey: "screenshots.distribution",
    descKey: "screenshots.distribution.desc",
    icon: Share2,
    mockContent: (
      <div className="p-6 space-y-3">
        <div className="grid grid-cols-4 gap-3">
          {["WeChat", "Zhihu", "Xiaohongshu", "Douyin"].map((p, i) => (
            <div key={p} className="rounded-xl bg-white/5 border border-white/10 p-3 text-center">
              <div className={cn("mx-auto mb-1 h-8 w-8 rounded-lg", i % 2 === 0 ? "bg-accent-green/20" : "bg-accent-blue/20")} />
              <div className="text-xs text-white/70">{p}</div>
              <div className="text-[10px] text-accent-green">Connected</div>
            </div>
          ))}
        </div>
        <div className="rounded-xl bg-white/5 border border-white/10 p-4">
          <div className="text-xs text-white/40 mb-2">Content Calendar — This Week</div>
          <div className="grid grid-cols-7 gap-1">
            {["Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"].map((d) => (
              <div key={d} className="text-center text-[10px] text-white/30">{d}</div>
            ))}
            {[2, 0, 3, 1, 2, 0, 0].map((count, i) => (
              <div key={i} className="flex justify-center">
                {count > 0 && <div className="h-1.5 w-1.5 rounded-full bg-accent-blue" />}
              </div>
            ))}
          </div>
        </div>
      </div>
    ),
  },
];

export default function Screenshots() {
  const [activeTab, setActiveTab] = useState("dashboard");
  const { t } = useI18n();
  const active = TABS.find((tab) => tab.id === activeTab)!;

  return (
    <section className="py-24 bg-white">
      <div className="mx-auto max-w-7xl px-6">
        <div className="mx-auto max-w-2xl text-center mb-12">
          <h2 className="text-3xl font-bold tracking-tight sm:text-4xl">
            {t("screenshots.title")}
          </h2>
          <p className="mt-4 text-lg text-muted-foreground">
            {t("screenshots.subtitle")}
          </p>
        </div>

        {/* Tab buttons */}
        <div className="flex justify-center gap-2 mb-8 flex-wrap">
          {TABS.map((tab) => (
            <button
              key={tab.id}
              onClick={() => setActiveTab(tab.id)}
              className={cn(
                "inline-flex items-center gap-2 rounded-lg px-4 py-2.5 text-sm font-medium transition-all",
                activeTab === tab.id
                  ? "bg-hero-from text-white shadow-lg"
                  : "bg-muted text-muted-foreground hover:bg-muted/80"
              )}
            >
              <tab.icon className="h-4 w-4" />
              {t(tab.labelKey)}
            </button>
          ))}
        </div>

        {/* Mock screen */}
        <div className="mx-auto max-w-5xl">
          <div className="rounded-2xl bg-hero-from border border-white/10 overflow-hidden shadow-2xl">
            {/* Window chrome */}
            <div className="flex items-center gap-2 border-b border-white/10 px-4 py-3">
              <div className="h-3 w-3 rounded-full bg-red-400/60" />
              <div className="h-3 w-3 rounded-full bg-yellow-400/60" />
              <div className="h-3 w-3 rounded-full bg-green-400/60" />
              <div className="ml-4 flex-1 rounded-md bg-white/10 px-3 py-1 text-xs text-white/40">
                app.uargus.io/{active.id}
              </div>
            </div>
            {/* Content */}
            <AnimatePresence mode="wait">
              <motion.div
                key={activeTab}
                initial={{ opacity: 0, y: 10 }}
                animate={{ opacity: 1, y: 0 }}
                exit={{ opacity: 0, y: -10 }}
                transition={{ duration: 0.3 }}
                className="min-h-[320px]"
              >
                {active.mockContent}
              </motion.div>
            </AnimatePresence>
          </div>
          <p className="mt-4 text-center text-sm text-muted-foreground">{t(active.descKey)}</p>
        </div>
      </div>
    </section>
  );
}
