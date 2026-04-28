"use client";

import { useState, useRef, useCallback } from "react";
import { useVirtualizer } from "@tanstack/react-virtual";
import FeedCard from "@/components/feed/FeedCard";
import FeedFilter from "@/components/feed/FeedFilter";
import FeedDetail from "@/components/feed/FeedDetail";
import type { FeedItem } from "@/components/feed/FeedCard";

const MOCK_FEED: FeedItem[] = [
  { id: "1", title: "AI 监管新政：欧盟发布最新人工智能法案实施细则", source: "Reuters", time: "3 min ago", summary: "欧盟委员会今日发布了 AI 法案的详细实施指南，涵盖高风险 AI 系统的合规要求、透明度义务和处罚标准。", tags: ["科技政策", "AI"], importance: "紧急" },
  { id: "2", title: "美联储会议纪要释放鸽派信号，市场反应积极", source: "Bloomberg", time: "12 min ago", summary: "最新公布的联邦公开市场委员会会议纪要显示，多数委员倾向于在未来几个月内开始降息周期。", tags: ["金融", "宏观"], importance: "高" },
  { id: "3", title: "全球碳排放达到新高，COP 会议各方博弈加剧", source: "Nature", time: "28 min ago", summary: "最新研究数据显示全球碳排放量再创历史新高，各国在即将召开的 COP 大会上立场分歧明显加大。", tags: ["气候", "环境"], importance: "中" },
  { id: "4", title: "半导体出口管制升级：日荷最新限制措施分析", source: "TechCrunch", time: "45 min ago", summary: "日本和荷兰政府相继发布新的半导体设备出口管制措施，进一步限制先进制程技术出口。", tags: ["地缘政治", "半导体"], importance: "高" },
  { id: "5", title: "量子计算突破：Google 实现百万量子比特纠错", source: "arXiv", time: "1 hr ago", summary: "Google Quantum AI 团队在最新论文中展示了百万量子比特级别的纠错能力，推动量子优越性实现。", tags: ["前沿科技", "量子"], importance: "中" },
  { id: "6", title: "全球粮食安全预警：厄尔尼诺影响评估报告", source: "World Bank", time: "2 hr ago", summary: "世界银行发布最新报告，评估厄尔尼诺气候现象对全球粮食供应链的潜在影响和政策建议。", tags: ["经济", "农业"], importance: "低" },
  { id: "7", title: "SpaceX 星舰第六次试飞成功回收助推器", source: "SpaceNews", time: "2.5 hr ago", summary: "SpaceX 完成星舰第六次试飞，成功实现助推器空中回收，标志着全面可重复使用火箭的重要里程碑。", tags: ["航天", "科技"], importance: "中" },
  { id: "8", title: "OpenAI 推出新一代推理模型 o3-mini", source: "The Verge", time: "3 hr ago", summary: "OpenAI 发布 o3-mini 推理模型，在数学和编程任务上表现大幅提升，同时推理成本降低 60%。", tags: ["AI", "科技"], importance: "高" },
  { id: "9", title: "中东局势：红海航运危机持续升级", source: "Al Jazeera", time: "4 hr ago", summary: "胡塞武装持续攻击红海商船，多家航运公司绕道好望角，全球供应链成本显著上升。", tags: ["地缘政治", "能源"], importance: "紧急" },
  { id: "10", title: "新型固态电池技术突破量产瓶颈", source: "MIT Tech Review", time: "5 hr ago", summary: "日本丰田与中国宁德时代分别宣布固态电池技术突破，预计2026年实现量产。", tags: ["能源", "科技"], importance: "中" },
  { id: "11", title: "全球央行数字货币进展追踪", source: "BIS", time: "6 hr ago", summary: "国际清算银行最新报告显示全球超过90%的央行正在研究CBDC，其中11个已正式启动试点。", tags: ["金融", "数字货币"], importance: "低" },
  { id: "12", title: "DeepMind AlphaFold3 解析蛋白质-药物互作", source: "Science", time: "7 hr ago", summary: "DeepMind 发布 AlphaFold3 升级版，能够精确预测蛋白质与小分子药物的结合模式，加速药物研发。", tags: ["医疗", "AI"], importance: "中" },
];

export default function FeedPage() {
  const [search, setSearch] = useState("");
  const [category, setCategory] = useState("全部");
  const [importance, setImportance] = useState("全部");
  const [layout, setLayout] = useState<"card" | "list" | "magazine">("card");
  const [selectedId, setSelectedId] = useState<string | null>(null);

  const parentRef = useRef<HTMLDivElement>(null);

  const filtered = MOCK_FEED.filter((item) => {
    if (search && !item.title.toLowerCase().includes(search.toLowerCase())) return false;
    if (category !== "全部" && !item.tags.some((t) => t.includes(category))) return false;
    if (importance !== "全部" && item.importance !== importance) return false;
    return true;
  });

  const virtualizer = useVirtualizer({
    count: filtered.length,
    getScrollElement: () => parentRef.current,
    estimateSize: useCallback(() => (layout === "list" ? 56 : 160), [layout]),
    overscan: 5,
  });

  const selectedItem = filtered.find((i) => i.id === selectedId) ?? null;

  return (
    <div className="space-y-4 h-full flex flex-col">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold">信息流</h1>
        <span className="text-xs text-muted-foreground">{filtered.length} 篇文章</span>
      </div>

      <FeedFilter
        search={search}
        onSearchChange={setSearch}
        activeCategory={category}
        onCategoryChange={setCategory}
        activeImportance={importance}
        onImportanceChange={setImportance}
        layout={layout}
        onLayoutChange={setLayout}
      />

      {/* Virtualized list */}
      <div ref={parentRef} className="flex-1 overflow-y-auto min-h-0">
        <div
          className="relative"
          style={{ height: `${virtualizer.getTotalSize()}px` }}
        >
          {virtualizer.getVirtualItems().map((vRow) => {
            const item = filtered[vRow.index];
            return (
              <div
                key={item.id}
                className="absolute left-0 right-0 px-0.5"
                style={{
                  top: `${vRow.start}px`,
                  height: `${vRow.size}px`,
                  paddingBottom: 8,
                }}
              >
                <FeedCard
                  item={item}
                  selected={selectedId === item.id}
                  onClick={() => setSelectedId(item.id)}
                  layout={layout === "list" ? "list" : "card"}
                />
              </div>
            );
          })}
        </div>
      </div>

      {/* Detail slide-in panel */}
      <FeedDetail item={selectedItem} onClose={() => setSelectedId(null)} />
      {selectedItem && (
        <div
          className="fixed inset-0 z-40 bg-black/20"
          onClick={() => setSelectedId(null)}
        />
      )}
    </div>
  );
}
