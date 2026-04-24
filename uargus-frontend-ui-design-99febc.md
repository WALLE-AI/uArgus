# uArgus Frontend UI/UX Design Plan

Next.js + TailwindCSS fullstack frontend for uArgus — a world intelligence monitoring platform with Landing Page and Main Dashboard.

---

## Tech Stack

- **Framework**: Next.js 14 (App Router, SSR/SSG)
- **Styling**: TailwindCSS 4 + shadcn/ui components
- **Icons**: Lucide React
- **i18n**: next-intl (Chinese/English bilingual)
- **Charts**: Recharts or ECharts (data visualization)
- **Animations**: Framer Motion (landing page hero + transitions)
- **State**: Zustand (lightweight global state)
- **AI Chat**: Vercel AI SDK (streaming chat UI)
- **Package Manager**: pnpm

---

## Part 1: Landing Page

### 1.1 Page Sections (Top to Bottom)

1. **Navbar** — Logo + Nav links (Features/Pricing/Cases/Docs) + Language toggle (CN/EN) + CTA "Start Free"
2. **Hero Section** (Data Intelligence Style)
   - Dark gradient background with animated particle/globe visualization
   - Main headline: "Monitor 3000+ global sources, turn noise into insight"
   - Sub-headline: value proposition in one sentence
   - Two CTAs: "Start Free Trial" (primary) + "Watch Demo" (ghost)
   - Live stats ticker: "3000+ RSS Sources | 50+ Data Dimensions | Real-time Analysis"
   - Background: animated world map with data flow lines or rotating 3D globe
3. **Trust Bar** — Logo strip of data sources (RSS, arXiv, HackerNews, Reddit, etc.)
4. **Core Value Proposition** (SaaS Card Style)
   - 5 feature cards in grid layout, each with icon + title + description:
     - **Global Intelligence Monitoring** — 3000+ RSS, API, crawlers, RSSHub
     - **AI-Powered Content Analysis** — Topic/event/sentiment/prediction analysis
     - **Knowledge Graph** — Anti-fragmentation, structured knowledge accumulation
     - **Personal AI Assistant** — Persistent chat, scheduled tasks, long-term memory
     - **Multi-Platform Distribution** — WeChat, Zhihu, Xiaohongshu, Douyin auto-publish
5. **How It Works** — 4-step horizontal flow diagram:
   - Collect (3000+ sources) -> Analyze (AI semantic layer) -> Generate (human-reviewed) -> Distribute (multi-platform)
   - Each step has an animated icon and brief description
6. **Product Screenshots / Interactive Demo**
   - Tab-switchable screenshots: Dashboard / AI Assistant / Analytics / Distribution
   - Optional: embedded interactive mini-demo
7. **ROI / Numbers Section** (Media Style)
   - "10x lower operation cost" / "10x more exposure for new accounts" / "Data asset accumulation"
   - Animated counter numbers on scroll
8. **Use Cases / Testimonials**
   - 3-4 case cards: Self-media operator / Live-streaming company / Research team / News agency
   - Each with before/after metrics
9. **Pricing Section** — Free / Pro / Enterprise tiers
10. **Footer** — Links, social media, legal, language toggle

### 1.2 Landing Page Design Tokens

- **Colors**: Dark hero (#0A0F1C -> #1A1F36), Light sections (#FAFBFC), Accent blue (#3B82F6), Accent green (#10B981)
- **Typography**: Inter (EN) + Noto Sans SC (CN), hero 64px bold, section title 36px
- **Spacing**: Section padding 96px vertical, card gap 24px
- **Animations**: Scroll-triggered fade-in, number counters, particle background

---

## Part 2: Main Dashboard (Product App)

### 2.1 Layout Architecture (Mixed Mode)

```
+----------------------------------------------------------+
|  Top Bar: Logo | Breadcrumb | Search (Cmd+K) | Notifications | User Avatar  |
+------+-----------------------------------+---------------+
| Left |        Center Content             |    Right      |
| Nav  |   (Feed / Dashboard / Analytics   |  AI Assistant |
| Bar  |    switchable views)              |    Panel      |
|      |                                   |  (collapsible)|
| 60px |         flex-1                    |    380px      |
+------+-----------------------------------+---------------+
```

### 2.2 Left Navigation (60px collapsed / 240px expanded)

Icon-based navigation with tooltip on hover, expandable groups. **4 groups + bottom user**:

**看板 (Dashboard)** — LayoutDashboard icon
- Dashboard 总览 — 核心指标卡片 + 可配置 widget 网格

**监控 (Monitor)** — Radar icon
- 数据源管理 — 3000+ RSS 源配置、添加、OPML 导入/导出
- 领域分类 — 数据源按领域/行业分类管理（科技/金融/地缘政治/能源/医疗等），标签体系、分类规则配置
- 源健康状态 — 健康/降级/失败状态看板
- 采集任务 — RSS 抓取频率配置、定时任务
- 信息流 — Feed 流式阅读、筛选、搜索（原内容模块）
- 智能分析 — 主题/事件/选题分析（原内容模块）
- 舆情监控 — 情感趋势、舆情预警、热点追踪
- 预测分析 — 趋势预测、置信区间、历史准确率
- 事件追踪 — 事件生命周期（BREAKING→DEVELOPING→SUSTAINED→FADING）、跨源佐证

**内容 (Content)** — FileText icon
- 内容中心 — AI 辅助内容生成/编辑/校核工作台（从分析到成稿的创作中心）
- 知识沉淀 — 知识图谱、结构化笔记、收藏夹
- 内容分发 — 多平台发布管理、内容日历、AI 改写编辑器

**观测 (Observe)** — Eye icon
- 渠道运营 — 各平台（微信公众号/知乎/小红书/抖音）发布后流量监控、粉丝增长曲线、互动数据、ROI 分析
- 智能体观测 — Agent 运行状态总览：
  - Token 消耗：日/周/月用量趋势图、按模型/Agent 分维度统计、费用估算
  - 工具调用：各工具调用次数、成功/失败率、平均延迟
  - Agent 任务：任务队列深度、执行时长分布、错误日志
  - 模型健康：各 LLM provider 可用性、响应延迟、fallback 触发次数

**底部区域**:
- 用户头像 + 用户名 — 点击展开：个人设置、API Key 管理、语言切换、深色/浅色模式、登出

### 2.3 Center Content Views

#### A. Dashboard Overview (Home)
- **Top row**: Key metric cards (Sources Active, Articles Today, Trending Topics, Alerts)
- **Middle**: Configurable widget grid (drag-and-drop):
  - Real-time feed mini-widget
  - Trending topics word cloud
  - Sentiment heatmap
  - Source health status
  - Recent AI analysis summaries
- **Bottom**: Activity timeline

#### B. Feed View (Core Information Stream)
- **Filter bar**: Category tags, source filter, time range, importance level, search
- **Feed layout**: Card list (switchable: list/card/magazine view)
- Each feed item card:
  - Source icon + name + publish time
  - Title (with importance badge: Critical/High/Medium/Low)
  - AI-generated 2-line summary
  - Tags (topic classification)
  - Actions: Save / Analyze / Share / Add to Knowledge
- **Infinite scroll** with virtual list for performance
- **Detail panel**: Click card -> slide-in detail view with full article + AI analysis

#### C. Analytics View
- **Sub-tabs**: Topic Analysis | Event Tracking | Sentiment Analysis | Prediction
- **Topic Analysis**:
  - Topic clustering visualization (bubble chart)
  - Topic timeline (when topics emerge/peak/fade)
  - Related articles drill-down
- **Event Tracking**:
  - Event lifecycle: BREAKING -> DEVELOPING -> SUSTAINED -> FADING
  - Event impact scoring
  - Cross-source corroboration display
- **Sentiment Analysis**:
  - Sentiment trend line chart (positive/negative/neutral over time)
  - Sentiment by source/region heatmap
  - Alert thresholds configuration
- **Prediction Analysis**:
  - Trend prediction charts
  - Confidence intervals
  - Historical accuracy tracking

#### D. Knowledge View
- **Knowledge Graph**: Interactive node-link visualization (D3.js or vis.js)
  - Nodes = entities (people, orgs, events, topics)
  - Edges = relationships
  - Click node -> detail panel
- **Structured Notes**: Markdown editor with AI assist
- **Collections**: User-curated topic folders
- **Search**: Full-text + semantic search across accumulated knowledge

#### E. Distribution View
- **Platform connections**: WeChat Official Account, Zhihu, Xiaohongshu, Douyin status cards
- **Content queue**: Drag-and-drop content scheduling calendar
- **Content editor**: Rich text editor with AI rewrite for platform-specific formatting
- **Performance dashboard**: Cross-platform metrics (views, likes, shares, followers)
- **Auto-publish rules**: Trigger-based auto-distribution configuration

#### F. Sources Management
- **Source list**: Table with status, health, last fetch, record count
- **Add source**: RSS URL input, category assignment, fetch interval config
- **Health monitoring**: Source health dashboard (healthy/degraded/failing)
- **Bulk operations**: Import/export OPML, batch enable/disable

#### G. Agent Observability (智能体观测)
- **Overview cards**: Total tokens today / Total tool calls / Active agents / Error rate
- **Token 消耗面板**:
  - 日/周/月 token 用量折线图（按时间粒度切换）
  - 按模型维度堆叠柱状图（GPT-4 / Claude / BGE-M3 / Groq 等）
  - 按 Agent/任务类型维度饼图（Summarize / Classify / Embed / Chat）
  - 费用估算表：模型 × 用量 × 单价 = 费用，支持导出
- **工具调用面板**:
  - 各工具调用次数排行榜（bar chart）
  - 成功率/失败率 gauge 指标
  - 平均响应延迟趋势图
  - 失败调用日志表（时间、工具名、错误信息、可重试标记）
- **Agent 任务面板**:
  - 任务队列深度实时指标
  - 执行时长分布直方图
  - 任务状态表：pending / running / completed / failed
  - 错误日志详情（可展开查看完整 stack trace）
- **模型健康面板**:
  - 各 LLM provider 可用性状态灯（green/yellow/red）
  - 响应延迟 P50/P95/P99 趋势
  - Fallback 触发次数及链路展示
  - Rate limit 触发记录

### 2.4 Right Panel — AI Assistant (380px, collapsible)

- **Chat interface**: Streaming conversation with markdown rendering
- **Context awareness**: Can reference current feed item / analysis / knowledge
- **Quick actions**:
  - "Summarize this article"
  - "Find related events"
  - "Generate social media post"
  - "Analyze sentiment trend for [topic]"
- **Persistent conversation**: Long-term memory across sessions
- **Channel push**: Configure push to WeChat/Telegram/Slack from chat
- **Minimize**: Collapse to floating bubble icon

### 2.5 Global Features

- **Command Palette (Cmd+K)**: Global search across feeds, knowledge, analytics
- **Notification Center**: Alert bell with dropdown (source failures, trending alerts, task completions)
- **Theme**: Light / Dark mode toggle
- **i18n**: CN/EN toggle persistent in user settings
- **Responsive**: Tablet support (collapse right panel), mobile read-only feed view

---

## Part 3: Execution Phases

### Phase 1: Foundation (Week 1-2)
- [ ] Initialize Next.js 14 project with App Router
- [ ] Setup TailwindCSS 4 + shadcn/ui + Lucide icons
- [ ] Configure next-intl for CN/EN bilingual
- [ ] Setup project structure: `app/[locale]/(marketing)/` + `app/[locale]/(dashboard)/`
- [ ] Design tokens: colors, typography, spacing in tailwind.config
- [ ] Create shared UI components: Button, Card, Badge, Input, Dialog, Sheet

### Phase 2: Landing Page (Week 2-3)
- [ ] Navbar with language toggle and mobile hamburger
- [ ] Hero section with animated background (Framer Motion + Canvas/Three.js globe)
- [ ] Trust bar with auto-scrolling logos
- [ ] Feature cards section (5 core values)
- [ ] How-it-works flow diagram
- [ ] Product screenshot tabs
- [ ] ROI numbers section with scroll-triggered counters
- [ ] Use cases / testimonial cards
- [ ] Pricing table
- [ ] Footer
- [ ] SEO: meta tags, OpenGraph, structured data
- [ ] Performance: Image optimization, lazy loading, Core Web Vitals

### Phase 3: Dashboard Shell (Week 3-4)
- [ ] App layout: TopBar + LeftNav + CenterContent + RightPanel
- [ ] Left navigation with collapse/expand
- [ ] Command palette (Cmd+K) with search
- [ ] Notification center
- [ ] Theme toggle (light/dark)
- [ ] Route structure for all dashboard views
- [ ] AI Assistant panel shell with chat UI

### Phase 4: Core Views (Week 4-6)
- [ ] Dashboard Overview with widget grid
- [ ] Feed View with filter bar, card list, virtual scroll
- [ ] Feed detail slide-in panel
- [ ] Analytics sub-views (Topic/Event/Sentiment/Prediction)
- [ ] Chart components (Recharts/ECharts integration)

### Phase 5: Advanced Views (Week 6-8)
- [ ] Knowledge Graph visualization
- [ ] Structured notes with markdown editor
- [ ] Distribution platform management
- [ ] Content editor with AI rewrite
- [ ] Source management table with health indicators
- [ ] Task scheduler UI

### Phase 6: Integration (Week 8-10)
- [ ] Connect to Go backend API endpoints
- [ ] Real-time data via WebSocket/SSE for feed updates
- [ ] AI chat streaming integration (Vercel AI SDK)
- [ ] Authentication (NextAuth.js or Clerk)
- [ ] User settings persistence

### Phase 7: Polish (Week 10-12)
- [ ] Responsive design testing (tablet/mobile)
- [ ] Accessibility audit (WCAG 2.1)
- [ ] Performance optimization (bundle size, lazy routes)
- [ ] E2E tests (Playwright)
- [ ] Documentation

---

## Part 4: Directory Structure

```
frontend/
  src/
    app/
      [locale]/
        (marketing)/          # Landing page group
          page.tsx            # Landing page
          pricing/page.tsx
          layout.tsx
        (dashboard)/              # Dashboard group (auth required)
          layout.tsx              # Shell: TopBar+LeftNav+RightPanel
          dashboard/page.tsx      # 看板 — Overview
          monitor/                # 监控
            sources/page.tsx      # 数据源管理
            categories/page.tsx   # 领域分类
            health/page.tsx       # 源健康状态
            tasks/page.tsx        # 采集任务
            feed/page.tsx         # 信息流
            feed/[id]/page.tsx    # Feed item detail
            analysis/page.tsx     # 智能分析
            sentiment/page.tsx    # 舆情监控
            prediction/page.tsx   # 预测分析
            events/page.tsx       # 事件追踪
          content/                # 内容
            studio/page.tsx       # 内容中心（AI 创作工作台）
            knowledge/page.tsx    # 知识沉淀
            distribution/page.tsx # 内容分发
          observe/                # 观测
            channels/page.tsx     # 渠道运营（流量/粉丝/互动/ROI）
            agents/page.tsx       # 智能体观测（Token/工具调用/任务/模型健康）
          settings/page.tsx       # 用户设置
        layout.tsx            # Root layout (providers, i18n)
    components/
      ui/                     # shadcn/ui base components
      landing/                # Landing page sections
        Hero.tsx
        FeatureCards.tsx
        HowItWorks.tsx
        Screenshots.tsx
        ROINumbers.tsx
        UseCases.tsx
        Pricing.tsx
      dashboard/              # Dashboard components
        TopBar.tsx
        LeftNav.tsx
        AiAssistantPanel.tsx
        CommandPalette.tsx
        NotificationCenter.tsx
      feed/                   # Feed-specific components
        FeedCard.tsx
        FeedFilter.tsx
        FeedDetail.tsx
      analytics/              # Chart components
        TopicBubble.tsx
        SentimentTrend.tsx
        EventTimeline.tsx
      knowledge/              # Knowledge components
        GraphViewer.tsx
        NoteEditor.tsx
      distribution/           # Distribution components
        PlatformCard.tsx
        ContentCalendar.tsx
        ContentEditor.tsx
    lib/
      api.ts                  # Backend API client
      i18n.ts                 # i18n configuration
      store.ts                # Zustand stores
      utils.ts                # Utilities
    messages/
      en.json                 # English translations
      zh.json                 # Chinese translations
    styles/
      globals.css             # Tailwind imports + custom styles
  public/
    images/
    icons/
  next.config.ts
  tailwind.config.ts
  package.json
  tsconfig.json
```

---

## Part 5: Key Design Decisions

| Decision | Choice | Rationale |
|----------|--------|-----------|
| SSR vs CSR | SSR for landing, CSR for dashboard | SEO for marketing, performance for app |
| Chart lib | Recharts (primary) + ECharts (complex) | Recharts for simple charts, ECharts for geo/graph |
| Virtual list | @tanstack/react-virtual | Handle 1000+ feed items without DOM bloat |
| Rich editor | TipTap or Milkdown | Modern, extensible, markdown-native |
| AI chat | Vercel AI SDK + useChat | Built-in streaming, Next.js native |
| Auth | NextAuth.js v5 | Flexible, supports multiple providers |
| Globe/3D | Three.js (react-three-fiber) or Cobe | Lightweight animated globe for hero section |
