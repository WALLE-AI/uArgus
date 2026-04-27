"use client";

import { createContext, useContext, useState, useCallback, type ReactNode } from "react";

export type Locale = "en" | "zh";

interface I18nContextValue {
  locale: Locale;
  t: (key: string) => string;
  toggle: () => void;
}

const I18nContext = createContext<I18nContextValue | null>(null);

const dict: Record<Locale, Record<string, string>> = {
  en: {
    // Navbar
    "nav.features": "Features",
    "nav.howItWorks": "How It Works",
    "nav.cases": "Cases",
    "nav.pricing": "Pricing",
    "nav.startFree": "Start Free",
    "nav.lang": "中文",

    // Hero
    "hero.badge": "Monitoring 3000+ sources worldwide",
    "hero.title.1": "Turn Global ",
    "hero.title.highlight": "Intelligence",
    "hero.title.2": " Into Actionable Insight",
    "hero.subtitle": "Monitor thousands of RSS feeds, analyze trends with AI, accumulate structured knowledge, and distribute content across all major platforms — automatically.",
    "hero.cta.trial": "Start Free Trial",
    "hero.cta.demo": "Watch Demo",
    "hero.stat.sources": "RSS Sources",
    "hero.stat.dimensions": "Data Dimensions",
    "hero.stat.analysis": "Real-time Analysis",
    "hero.stat.cost": "Cost Reduction",

    // Trust Bar
    "trust.title": "Aggregating intelligence from trusted sources worldwide",

    // Features
    "features.title.1": "Everything You Need to ",
    "features.title.highlight": "Own the Information Edge",
    "features.subtitle": "From data collection to content distribution — a complete intelligence operation platform.",
    "features.monitoring.title": "Global Intelligence Monitoring",
    "features.monitoring.desc": "Monitor 3000+ RSS sources, API feeds, crawlers, and RSSHub channels. Categorize by domain — tech, finance, geopolitics, energy, and more.",
    "features.analysis.title": "AI-Powered Content Analysis",
    "features.analysis.desc": "Automated topic clustering, event tracking, sentiment analysis, and trend prediction. Transform raw data into structured intelligence.",
    "features.knowledge.title": "Knowledge Graph",
    "features.knowledge.desc": "Defeat information fragmentation. Build structured, searchable knowledge with entity graphs, semantic links, and curated collections.",
    "features.assistant.title": "Personal AI Assistant",
    "features.assistant.desc": "Persistent conversations with long-term memory. Summarize articles, find related events, generate content — all context-aware.",
    "features.distribution.title": "Multi-Platform Distribution",
    "features.distribution.desc": "Auto-publish to WeChat, Zhihu, Xiaohongshu, Douyin, and more. AI-powered rewriting for platform-specific formatting and tone.",

    // How It Works
    "how.title": "How It Works",
    "how.subtitle": "Four steps from raw information to audience engagement.",
    "how.collect.title": "Collect",
    "how.collect.desc": "3000+ RSS sources, APIs, crawlers, and RSSHub channels aggregated in real-time.",
    "how.analyze.title": "Analyze",
    "how.analyze.desc": "AI semantic layer processes topics, events, sentiment, and predictions automatically.",
    "how.generate.title": "Generate",
    "how.generate.desc": "AI-assisted content creation with human review. From raw intelligence to publishable articles.",
    "how.distribute.title": "Distribute",
    "how.distribute.desc": "Auto-publish across WeChat, Zhihu, Xiaohongshu, Douyin with platform-specific optimization.",

    // Screenshots
    "screenshots.title": "See It in Action",
    "screenshots.subtitle": "Explore the key product views that power your intelligence workflow.",
    "screenshots.dashboard": "Dashboard",
    "screenshots.assistant": "AI Assistant",
    "screenshots.analytics": "Analytics",
    "screenshots.distribution": "Distribution",
    "screenshots.dashboard.desc": "Real-time overview with configurable widgets, source health, trending topics, and activity timeline.",
    "screenshots.assistant.desc": "Context-aware AI chat with persistent memory, quick actions, and multi-channel push.",
    "screenshots.analytics.desc": "Topic clustering, event lifecycle tracking, sentiment trends, and prediction analysis.",
    "screenshots.distribution.desc": "Multi-platform content calendar, AI rewriting, and cross-platform performance metrics.",

    // ROI
    "roi.title": "Real Results, Real Impact",
    "roi.subtitle": "Our users report transformative improvements in their content operations.",
    "roi.cost.label": "Lower Operation Cost",
    "roi.cost.desc": "Automate content collection, analysis, and distribution. Replace manual work with AI-powered workflows.",
    "roi.exposure.label": "More Exposure",
    "roi.exposure.desc": "New accounts gain visibility through data-driven content strategies and multi-platform optimization.",
    "roi.retention.label": "Data Asset Retention",
    "roi.retention.desc": "Every piece of intelligence, analysis, and performance metric is structured and preserved as company assets.",

    // Use Cases
    "cases.title": "Who Uses uArgus",
    "cases.subtitle": "From solo creators to enterprise teams — built for anyone who needs the information edge.",
    "cases.media.title": "Self-Media Operator",
    "cases.media.before": "8 hours/day on content research",
    "cases.media.after": "1 hour/day with AI-assisted workflow",
    "cases.media.quote": "uArgus replaced 3 tools and 2 interns. My content quality went up while my workload dropped dramatically.",
    "cases.live.title": "Live-Streaming Company",
    "cases.live.before": "Product selection based on gut feeling",
    "cases.live.after": "Data-driven selection with trend prediction",
    "cases.live.quote": "We now pick trending products 3 days before competitors. Our data assets compound with every stream.",
    "cases.research.title": "Research Team",
    "cases.research.before": "Fragmented notes across 5 tools",
    "cases.research.after": "Unified knowledge graph with semantic search",
    "cases.research.quote": "The knowledge graph turned our scattered notes into a searchable intelligence database.",
    "cases.news.title": "News Agency",
    "cases.news.before": "Missing breaking news, slow reaction",
    "cases.news.after": "Real-time event tracking with cross-source verification",
    "cases.news.quote": "Event lifecycle tracking catches stories at the BREAKING stage. We're always first to report.",
    "cases.before": "Before",
    "cases.after": "After",

    // Pricing
    "pricing.title": "Simple, Transparent Pricing",
    "pricing.subtitle": "Start free, scale as you grow. No hidden fees.",
    "pricing.free": "Free",
    "pricing.free.price": "$0",
    "pricing.free.period": "forever",
    "pricing.free.desc": "Get started with core monitoring features.",
    "pricing.free.f1": "50 RSS source monitoring",
    "pricing.free.f2": "Basic feed reader",
    "pricing.free.f3": "Daily AI summary (5/day)",
    "pricing.free.f4": "Community support",
    "pricing.free.cta": "Get Started",
    "pricing.pro": "Pro",
    "pricing.pro.price": "$29",
    "pricing.pro.period": "/month",
    "pricing.pro.desc": "Full intelligence suite for serious creators.",
    "pricing.pro.f1": "Unlimited RSS sources",
    "pricing.pro.f2": "All analysis modules",
    "pricing.pro.f3": "AI Assistant (unlimited)",
    "pricing.pro.f4": "Knowledge graph",
    "pricing.pro.f5": "3 platform distributions",
    "pricing.pro.f6": "Priority support",
    "pricing.pro.cta": "Start Free Trial",
    "pricing.enterprise": "Enterprise",
    "pricing.enterprise.price": "Custom",
    "pricing.enterprise.period": "",
    "pricing.enterprise.desc": "For teams and organizations at scale.",
    "pricing.enterprise.f1": "Everything in Pro",
    "pricing.enterprise.f2": "Custom data sources & crawlers",
    "pricing.enterprise.f3": "Team collaboration",
    "pricing.enterprise.f4": "Unlimited distributions",
    "pricing.enterprise.f5": "API access",
    "pricing.enterprise.f6": "Dedicated account manager",
    "pricing.enterprise.f7": "SLA guarantee",
    "pricing.enterprise.cta": "Contact Sales",
    "pricing.popular": "Most Popular",

    // Footer
    "footer.desc": "Turn global intelligence into actionable insight. Monitor, analyze, create, distribute — all in one platform.",
    "footer.product": "Product",
    "footer.company": "Company",
    "footer.legal": "Legal",
    "footer.connect": "Connect",
    "footer.copyright": "uArgus. All rights reserved.",
  },
  zh: {
    // Navbar
    "nav.features": "功能特性",
    "nav.howItWorks": "工作原理",
    "nav.cases": "案例",
    "nav.pricing": "定价",
    "nav.startFree": "免费开始",
    "nav.lang": "EN",

    // Hero
    "hero.badge": "全球 3000+ 数据源实时监控",
    "hero.title.1": "将全球",
    "hero.title.highlight": "情报",
    "hero.title.2": "转化为可执行洞察",
    "hero.subtitle": "监控数千个 RSS 源，AI 智能分析趋势，构建结构化知识体系，多平台自动化分发 —— 一站式智能运营。",
    "hero.cta.trial": "免费试用",
    "hero.cta.demo": "观看演示",
    "hero.stat.sources": "RSS 数据源",
    "hero.stat.dimensions": "数据维度",
    "hero.stat.analysis": "实时分析",
    "hero.stat.cost": "成本降低",

    // Trust Bar
    "trust.title": "聚合全球可信数据源情报",

    // Features
    "features.title.1": "你需要的一切，",
    "features.title.highlight": "掌握信息优势",
    "features.subtitle": "从数据采集到内容分发 —— 完整的智能运营平台。",
    "features.monitoring.title": "全球情报监控",
    "features.monitoring.desc": "监控 3000+ RSS 源、API 接口、爬虫和 RSSHub 频道。按科技、金融、地缘政治、能源等领域分类管理。",
    "features.analysis.title": "AI 智能内容分析",
    "features.analysis.desc": "自动主题聚类、事件追踪、舆情分析和趋势预测。将原始数据转化为结构化情报。",
    "features.knowledge.title": "知识图谱",
    "features.knowledge.desc": "告别信息碎片化。通过实体图谱、语义关联和策展收藏构建可搜索的结构化知识库。",
    "features.assistant.title": "个人 AI 助手",
    "features.assistant.desc": "支持长期记忆的持久对话。摘要文章、发现关联事件、生成内容 —— 全程上下文感知。",
    "features.distribution.title": "多平台内容分发",
    "features.distribution.desc": "自动发布到微信公众号、知乎、小红书、抖音等平台。AI 智能改写，适配各平台格式与风格。",

    // How It Works
    "how.title": "工作原理",
    "how.subtitle": "四步实现从原始信息到受众触达。",
    "how.collect.title": "采集",
    "how.collect.desc": "3000+ RSS 源、API、爬虫和 RSSHub 频道实时聚合。",
    "how.analyze.title": "分析",
    "how.analyze.desc": "AI 语义层自动处理主题、事件、舆情和趋势预测。",
    "how.generate.title": "生成",
    "how.generate.desc": "AI 辅助内容创作，人工审核把关。从原始情报到可发布文章。",
    "how.distribute.title": "分发",
    "how.distribute.desc": "自动发布到微信、知乎、小红书、抖音，针对各平台优化内容。",

    // Screenshots
    "screenshots.title": "产品展示",
    "screenshots.subtitle": "探索驱动你情报工作流的核心产品视图。",
    "screenshots.dashboard": "仪表盘",
    "screenshots.assistant": "AI 助手",
    "screenshots.analytics": "智能分析",
    "screenshots.distribution": "内容分发",
    "screenshots.dashboard.desc": "实时概览，可配置组件，数据源健康状态，热点话题和活动时间线。",
    "screenshots.assistant.desc": "上下文感知的 AI 对话，持久记忆，快捷操作，多渠道推送。",
    "screenshots.analytics.desc": "主题聚类，事件生命周期追踪，舆情趋势，预测分析。",
    "screenshots.distribution.desc": "多平台内容日历，AI 改写，跨平台效果指标。",

    // ROI
    "roi.title": "真实效果，真实影响",
    "roi.subtitle": "我们的用户在内容运营中实现了显著提升。",
    "roi.cost.label": "运营成本降低",
    "roi.cost.desc": "自动化内容采集、分析和分发，用 AI 工作流替代人工操作。",
    "roi.exposure.label": "曝光量提升",
    "roi.exposure.desc": "新账号通过数据驱动的内容策略和多平台优化快速获得流量。",
    "roi.retention.label": "数据资产沉淀",
    "roi.retention.desc": "每一条情报、分析和效果指标都被结构化保存，成为企业数据资产。",

    // Use Cases
    "cases.title": "谁在使用 uArgus",
    "cases.subtitle": "从个人创作者到企业团队 —— 为每一个需要信息优势的人而建。",
    "cases.media.title": "自媒体运营者",
    "cases.media.before": "每天 8 小时内容调研",
    "cases.media.after": "AI 辅助每天仅需 1 小时",
    "cases.media.quote": "uArgus 替代了 3 个工具和 2 个实习生。内容质量提升的同时，工作量大幅下降。",
    "cases.live.title": "直播电商公司",
    "cases.live.before": "凭经验选品",
    "cases.live.after": "数据驱动 + 趋势预测选品",
    "cases.live.quote": "我们现在比竞品提前 3 天发现热门产品。每场直播的数据资产持续积累。",
    "cases.research.title": "研究团队",
    "cases.research.before": "笔记散落在 5 个工具中",
    "cases.research.after": "统一知识图谱 + 语义搜索",
    "cases.research.quote": "知识图谱将我们散乱的笔记变成了可搜索的情报数据库。",
    "cases.news.title": "新闻机构",
    "cases.news.before": "错过突发新闻，反应迟缓",
    "cases.news.after": "实时事件追踪 + 跨源验证",
    "cases.news.quote": "事件生命周期追踪在 BREAKING 阶段就能捕获新闻。我们总是第一个报道。",
    "cases.before": "使用前",
    "cases.after": "使用后",

    // Pricing
    "pricing.title": "简单透明的定价",
    "pricing.subtitle": "免费起步，按需升级。没有隐藏费用。",
    "pricing.free": "免费版",
    "pricing.free.price": "¥0",
    "pricing.free.period": "永久免费",
    "pricing.free.desc": "体验核心监控功能。",
    "pricing.free.f1": "50 个 RSS 源监控",
    "pricing.free.f2": "基础信息流阅读",
    "pricing.free.f3": "每日 AI 摘要（5次/天）",
    "pricing.free.f4": "社区支持",
    "pricing.free.cta": "免费开始",
    "pricing.pro": "专业版",
    "pricing.pro.price": "¥199",
    "pricing.pro.period": "/月",
    "pricing.pro.desc": "为专业创作者打造的完整情报套件。",
    "pricing.pro.f1": "无限 RSS 源",
    "pricing.pro.f2": "全部分析模块",
    "pricing.pro.f3": "AI 助手（无限使用）",
    "pricing.pro.f4": "知识图谱",
    "pricing.pro.f5": "3 个平台分发",
    "pricing.pro.f6": "优先支持",
    "pricing.pro.cta": "免费试用",
    "pricing.enterprise": "企业版",
    "pricing.enterprise.price": "定制",
    "pricing.enterprise.period": "",
    "pricing.enterprise.desc": "为团队和组织量身定制。",
    "pricing.enterprise.f1": "专业版全部功能",
    "pricing.enterprise.f2": "自定义数据源和爬虫",
    "pricing.enterprise.f3": "团队协作",
    "pricing.enterprise.f4": "无限平台分发",
    "pricing.enterprise.f5": "API 接口",
    "pricing.enterprise.f6": "专属客户经理",
    "pricing.enterprise.f7": "SLA 保障",
    "pricing.enterprise.cta": "联系销售",
    "pricing.popular": "最受欢迎",

    // Footer
    "footer.desc": "将全球情报转化为可执行洞察。监控、分析、创作、分发 —— 一站式平台。",
    "footer.product": "产品",
    "footer.company": "公司",
    "footer.legal": "法律",
    "footer.connect": "联系我们",
    "footer.copyright": "uArgus 版权所有",
  },
};

export function I18nProvider({ children }: { children: ReactNode }) {
  const [locale, setLocale] = useState<Locale>("en");

  const toggle = useCallback(() => {
    setLocale((prev) => (prev === "en" ? "zh" : "en"));
  }, []);

  const t = useCallback(
    (key: string) => dict[locale][key] ?? key,
    [locale]
  );

  return (
    <I18nContext.Provider value={{ locale, t, toggle }}>
      {children}
    </I18nContext.Provider>
  );
}

export function useI18n() {
  const ctx = useContext(I18nContext);
  if (!ctx) throw new Error("useI18n must be used within I18nProvider");
  return ctx;
}
