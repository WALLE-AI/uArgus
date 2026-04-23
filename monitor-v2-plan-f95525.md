# Go Monitor v2: News + Research 执行方案（注册描述符 + 预处理 Pipeline + 语义工具箱 + 采集工具箱 + 统一缓存/存储层）

`backend/monitor` 迁移 News + Research 数据采集引擎。六大核心抽象层：`internal/registry/`（SourceSpec 声明式注册+多模式调度+健康）、`internal/preprocess/`（预处理 Pipeline，可组合阶段）、`internal/semantic/`（语义能力工具箱，6+3 子包）、`internal/fetcher/`（采集能力工具箱，6 子包）、`internal/cache/`（统一缓存层：Client+FetchThrough+KeyRegistry+Metrics）、`internal/seed/`（数据发布层：Runner+Lock+Envelope+Publish+Meta+TTL）。各 Source 按需组合预处理、语义和采集能力，共享缓存/存储基础设施。

---

## 一、技术架构图

```
┌───────────────────────────────────────────────────────────────────────┐
│                       backend/monitor (Go)                            │
│                                                                       │
│  cmd/monitor/main.go                                                  │
│    Config → Registry.New() → news.RegisterAll() → research.RegisterAll()
│    → Registry.Boot() → healthz → graceful shutdown                    │
│                                                                       │
│  ┌─────────────────────────────────────────────────────────────────┐  │
│  │ internal/registry/  ★ 声明式注册模块                             │  │
│  │  SourceSpec │ Schedule(Cron/Interval/OnDemand)                   │  │
│  │  Source 接口 │ Registry 容器 │ HealthTracker │ Scheduler         │  │
│  └──────┬──────────────────────────────────────────────────────────┘  │
│         │ Source.Run()                                                 │
│         ▼                                                             │
│  ┌─────────────────────────────────────────────────────────────────┐  │
│  │ internal/preprocess/  ★ 预处理 Pipeline（可组合阶段）            │  │
│  │  Stage 接口 │ FormatMapper │ MultiSourceMerger                  │  │
│  └──────┬──────────────────────────────────────────────────────────┘  │
│         │                                                             │
│         ▼                                                             │
│  ┌─────────────────────────────────────────────────────────────────┐  │
│  │ internal/semantic/  ★ 能力工具箱（各子包独立，Source 按需组合）     │  │
│  │  classify/ │ scoring/ │ tracking/ │ enrichment/ │ tiers/ │ agents/│  │
│  │  ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ │  │
│  │  geo/ (⏳) │ anomaly/ (⏳) │ fusion/ (P2 实装)                  │  │
│  └─────────────────────────────────────────────────────────────────┘  │
│                                                                       │
│  ┌──────────────┐ ┌──────────────┐ ┌──────────────┐                  │
│  │ news/        │ │ research/    │ │ (future)     │                  │
│  │ DigestSource │ │ ArxivSource  │ │ cyber/       │                  │
│  │  ← classify  │ │ HNSource     │ │  ← preprocess│                  │
│  │  ← scoring   │ │ EventsSource │ │  ← scoring   │                  │
│  │  ← tracking  │ │ TrendSource  │ │ resilience/  │                  │
│  │  ← enrichment│ │  (无 semantic)│ │  ← scoring   │                  │
│  │  ← tiers     │ └──────────────┘ │  ← tracking  │                  │
│  │ InsightsSource│                  │ correlation/ │                  │
│  │  ← agents    │                   │  ← fusion    │                  │
│  └──────┬───────┘                   └──────────────┘                  │
│         │                                                             │
│  ┌──────▼──────────────────────────────────────────────────────────┐  │
│  │ internal/fetcher/ ★ 采集能力工具箱                             │  │
│  │  proxy/ │ ratelimit/ │ pool/ │ fallback/ │ pagination/ │ parser/│  │
│  └──────┬──────────────────────────────────────────────────────────┘  │
│  ┌──────▼──────────────────────────────────────────────────────────┐  │
│  │ internal/cache/ ★ 统一缓存层                                    │  │
│  │  Client(Upstash) │ FetchThrough(singleflight+neg+SWR)         │  │
│  │  KeyRegistry(SourceSpec→key+TTL+maxStale) │ Pipeline(batch)   │  │
│  │  Metrics(Prometheus hit/miss/latency) │ Sidecar(LRU+TTL)     │  │
│  └─────────────────────────────────────────────────────────────────┘  │
│  ┌─────────────────────────────────────────────────────────────────┐  │
│  │ internal/seed/ ★ 数据发布层                                     │  │
│  │  Runner(lock→fetch→publish→meta→unlock)                        │  │
│  │  Lock(NX+PX+Lua) │ Envelope(struct 单源) │ Publish(MULTI/EXEC)│  │
│  │  Meta(seed-meta:*) │ TTL(SourceSpec.DataTTL)                  │  │
│  └─────────────────────────────────────────────────────────────────┘  │
└──────────┬──────────────────────┬─────────────────────────────────────┘
           ▼                      ▼
    ┌──────────┐          ┌──────────────┐       ┌────────────────┐
    │ Redis    │          │ Upstream APIs│       │ backend/agents │
    │ (Upstash)│          │ RSS/arXiv/…  │       │ LLMProvider    │
    └──────────┘          └──────────────┘       └────────────────┘
```

---

## 二、核心新增模块

### 2.1 声明式注册模块 `internal/registry/`

**SourceSpec 声明符**（统一 4 种调度模式，消除元数据碎片化）：
```go
// 每个 Source 在构造时填写 SourceSpec，调度、健康、TTL 集中声明
type SourceSpec struct {
    // 身份
    Domain   string   // "news", "cyber", "resilience"
    Resource string   // "digest", "threats", "scores"
    Version  int      // 1, 2

    // 调度（统一原 Railway Cron / Bundle / Relay setInterval / RPC warm-ping）
    Schedule Schedule // Cron("*/5 * * * *") | Interval(5*time.Minute) | OnDemand
    LockTTL  time.Duration

    // 数据
    CanonicalKey string   // "news:digest:v1:full:en"
    ExtraKeys    []string // bootstrap 镜像 key
    DataTTL      time.Duration

    // 健康（替代 health.js 手工 SEED_META 表）
    MaxStaleDuration time.Duration
    MinRecordCount   int
}

// Schedule 类型：统一原来 4 种调度模式
type Schedule interface{ scheduleTag() }
type CronSchedule    struct { Expr string }             // 原 Railway Cron / Bundle
type IntervalSchedule struct { Every time.Duration }    // 原 Relay setInterval
type OnDemandSchedule struct {}                          // 原 RPC warm-ping
```

**Source 接口**（融合 `go-engine-review` 建议）：
```go
type Source interface {
    Name() string             // "news:digest:full:en"
    Spec() SourceSpec         // 声明式元数据
    Dependencies() []string   // 依赖的其他 Source Name
    Run(ctx context.Context) (*FetchResult, error)
}

type FetchResult struct {
    Data        any
    ExtraKeys   []ExtraKey
    Metrics     FetchMetrics
    NeedsEnrich bool
}

type FetchMetrics struct {
    Duration       time.Duration
    RecordCount    int
    UpstreamStatus int
    BytesRead      int64
}
```

**Registry 容器**（支持 Cron/Interval/OnDemand 三种调度）：
```go
type Registry struct {
    sources map[string]Source
    health  *HealthTracker
    sched   *Scheduler        // 统一调度器（内部区分 cron/ticker/ondemand）
    seed    *seed.Runner
}

func (r *Registry) Register(s Source)
func (r *Registry) Boot(ctx context.Context) error      // 依赖排序 → 按 Schedule 类型注册 → 启动
func (r *Registry) Shutdown(ctx context.Context) error
func (r *Registry) HealthSnapshot() map[string]HealthStatus
func (r *Registry) HealthHandler() http.Handler          // /healthz endpoint
func (r *Registry) TriggerOnDemand(name string) error    // RPC warm-ping 触发
```

**HealthTracker**（从 SourceSpec 读取阈值，无需手工 SEED_META 表）：
```go
type HealthStatus struct {
    LastSuccessAt    time.Time
    LastAttemptAt    time.Time
    ConsecutiveFails int
    LastRecordCount  int
    AvgDuration      time.Duration
    State            string  // "healthy" | "degraded" | "failing"
    // 新增：从 SourceSpec 自动读取
    MaxStaleDuration time.Duration  // 来自 Spec().MaxStaleDuration
    MinRecordCount   int            // 来自 Spec().MinRecordCount
}
```

### 2.2 语义层能力工具箱 `internal/semantic/`

**设计原则**：不同数据源的语义层差异极大（news=关键词分类+评分+跟踪，seismology=地理空间proximity，economic=加权指数，climate=Z-score异常检测）。因此不设统一 Pipeline，而是拆成**独立子包**，各 Source 按需组合。

#### classify/ — 关键词分类
```go
type Classifier struct { levels []KeywordLevel }

func NewNewsClassifier() *Classifier        // 预设 7 层 geo+tech 关键词
func NewCustomClassifier(levels []KeywordLevel) *Classifier

func (c *Classifier) Classify(title string, opts ...ClassifyOption) ClassificationResult
```

#### scoring/ — 通用加权评分（通用 WeightedScorer + 多 Profile 预设）
```go
type ScoringProfile struct {
    Name    string
    Factors []Factor  // {Name, Weight}
    Clamp   [2]float64
}

type WeightedScorer struct { profile ScoringProfile }

// 预设 Profile：
func NewImportanceScorer() *WeightedScorer    // news: severity*0.55+tier*0.2+corr*0.15+recency*0.1
func NewDisruptionScorer() *WeightedScorer    // chokepoint: threatLevel*0.4+warningCount*0.3+severity*0.2+anomaly*0.1
func NewGoalpostScorer() *WeightedScorer      // resilience: direction+goalposts+weight → normalized 0-100
func NewCustomScorer(p ScoringProfile) *WeightedScorer

func (s *WeightedScorer) Score(input ScoringInput) float64
```

#### tracking/ — 通用有状态追踪（StoryTracking + IntervalTracking）
```go
// 通用 StatefulTracker 接口
type StatefulTracker interface {
    Read(ctx context.Context, keys []string) (map[string]TrackInfo, error)
    Write(ctx context.Context, items []Trackable) error
}

// 新闻故事生命周期追踪（BREAKING→DEVELOPING→SUSTAINED→FADING）
type StoryTracker struct { rdb cache.Client }
func (t *StoryTracker) Read(ctx, hashes) (map[string]TrackInfo, error)
func (t *StoryTracker) Write(ctx, items) error
func (t *StoryTracker) ComputeCorroboration(items []Hashable) map[string]int

// 韧性分数区间追踪（stable/improving/declining）
type IntervalTracker struct { rdb cache.Client }
func (t *IntervalTracker) Read(ctx, countryCode) (*ScoreInterval, error)
func (t *IntervalTracker) Write(ctx, countryCode, score) error

type HashDedup struct{}  // 通用 hash 去重（news/unrest/aviation 复用）
```

#### enrichment/ — AI 缓存增强
```go
func EnrichBatch(ctx, rdb, items []ClassifiedItem) error  // MGET classify:sebuf:v3:{hash}
```

#### tiers/ — 源可信度分级
```go
func GetTier(sourceName string) int  // 1-4
```

#### agents/ — LLM/Embedding 客户端
```go
type AgentsClient interface {
    Summarize(ctx context.Context, texts []string, opts SummarizeOpts) (string, error)
    Classify(ctx context.Context, text string) (*AiClassification, error)   // 未来
    Embed(ctx context.Context, text string) ([]float32, error)              // 未来
}
```

#### fusion/ — 跨域派生（P2 实装）
```go
// 相关性融合器（读 digest+market+unrest → 交叉信号卡片）
type CorrelationFuser struct { rdb cache.Client }
func (f *CorrelationFuser) Fuse(ctx context.Context, sources []FuseInput) ([]CorrelationCard, error)

// 跨源融合器（读多域 seed key → 融合信号）
type CrossSourceFuser struct { rdb cache.Client }
func (f *CrossSourceFuser) Fuse(ctx context.Context, keys []string) ([]CrossSignal, error)
```

#### 未来扩展子包（当前只建骨架）
- **geo/** — Haversine / ProximityChecker / CityCoords / Geocoder
- **anomaly/** — ZScore / SpikeDetect / RollingBaseline

### 2.3 预处理 Pipeline `internal/preprocess/`

**设计原则**：数据源的预处理复杂度差异极大（P0 无预处理 55% → P4 有状态追踪 5%）。拆成可组合的 Stage 接口，各 Source 声明所需阶段。

```go
// Pipeline 阶段接口
type Stage interface {
    Process(ctx context.Context, raw any) (any, error)
}

// 预置 Stage
func FormatMapper(mapping FieldMapping) Stage      // P1: 字段映射/单位换算/时间标准化
func MultiSourceMerger(keys []string) Stage        // P2: 多源并行读取+合并

// Pipeline 组合器
type Pipeline struct { stages []Stage }
func NewPipeline(stages ...Stage) *Pipeline
func (p *Pipeline) Run(ctx context.Context, raw any) (any, error)
```

**各 Source 的 Pipeline 组合**：

| Source | Pipeline |
|--------|----------|
| ArxivSource | `[]` (P0, 无预处理) |
| IMFMacroSource | `[FormatMapper(weo_pivot)]` |
| DigestSource | `[FormatMapper(rss_parse)]` + semantic scoring/tracking |
| ChokepointSource | `[MultiSourceMerger(4_keys)]` + semantic scoring |
| ResilienceSource | `[MultiSourceMerger(19_dims)]` + semantic scoring/tracking |
| CyberSource | `[FormatMapper(normalize), MultiSourceMerger(3_feeds)]` + semantic scoring |

### 2.4 采集能力工具箱 `internal/fetcher/`

**设计原则**：不同数据源的采集策略差异极大（代理策略 4 种、限速策略 4 种、分页策略 3 种、并发控制 3 种、降级链、响应解析 7 种格式）。因此不设统一 Fetch()，而是拆成独立子包，各 Source 按需组合。

#### proxy/ — 代理策略
```go
type ProxyStrategy interface {
    Fetch(ctx context.Context, url string, opts ...FetchOption) (*http.Response, error)
}
func DirectFirst(resolver ProxyResolver) ProxyStrategy   // 直连优先，代理备用
func ProxyFirst(resolver ProxyResolver) ProxyStrategy    // 代理优先（FRED）
func CurlOnly(resolver ProxyResolver) ProxyStrategy      // 仅 curl（Yahoo）
func TwoLegCascade(connect, curl ProxyResolver) ProxyStrategy // CONNECT→curl（Open-Meteo）
```

#### ratelimit/ — 限速
```go
type RateLimiter interface { Wait(ctx context.Context) error }
func NewFixedInterval(d time.Duration) RateLimiter   // arXiv 3s, GDELT 20s
func NewRetryAfterHandler() *RetryAfterHandler       // 429 Retry-After 解析
type KeyRotator struct { keys []string }             // API Key 轮换
```

#### pool/ — 并发控制
```go
func BoundedPool(ctx, n int, items, fn) []Result   // 固定并发+allSettled
func FanOut(ctx, fns ...func() error) error         // errgroup 扇出
func Sequential(ctx, items, fn, cooldown) []T       // 顺序+冷却
```

#### fallback/ — 降级链
```go
type FallbackChain[T any] struct { Providers []Provider[T] }
func (c *FallbackChain[T]) Execute(ctx) (T, error)  // 首个成功即返回
```

#### pagination/ — 分页
```go
type IDBatchFetcher struct { BatchSize, Concurrency int }
func (f *IDBatchFetcher) FetchByIDs(ctx, ids, fn) ([]T, error)  // HN 模式
// 未来: OffsetPaginator (ArcGIS), PagePaginator
```

#### parser/ — 响应解析
```go
func ParseAtomXML(body []byte) ([]AtomEntry, error)  // arXiv
func ParseRssXML(body []byte) ([]RssItem, error)     // news RSS
func ParseICS(body []byte) ([]ICSEvent, error)        // Techmeme
// 未来: ParseCSV, ParseSDMX
```

### 2.5 统一缓存层 `internal/cache/`

**设计原则**：TS v1 的读侧 (`redis.ts`) 和写侧 (`_seed-utils.mjs`) 是完全割裂的两套 Redis 客户端，超时/前缀/错误策略各不相同。Go v2 统一为单一 `cache.Client` 接口，消除双轨维护。同时将 4 处散落的 key 定义（`cache-keys.ts`、`health.js` BOOTSTRAP/STANDALONE/SEED_META、seed 脚本、RPC handler）统一为 `KeyRegistry`，从 `SourceSpec` 自动推导。

#### Client — 统一 Redis 客户端（解决 v1 P-1 读写割裂 + P-5 前缀不一致）
```go
// Client 是唯一的 Redis 交互入口，读侧和写侧共用
type Client interface {
    Get(ctx context.Context, key string) ([]byte, error)
    Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
    Del(ctx context.Context, keys ...string) error
    Pipeline(ctx context.Context, cmds []Cmd) ([]Result, error)
    Expire(ctx context.Context, key string, ttl time.Duration) error
    Eval(ctx context.Context, script string, keys []string, args ...any) (any, error)
    // Geo/Hash 等特殊操作
    GeoSearchByBox(ctx context.Context, key string, lng, lat, w, h float64) ([]GeoMember, error)
    HMGet(ctx context.Context, key string, fields ...string) (map[string]string, error)
}

// ClientConfig 统一超时/重试配置（取代 v1 的硬编码 1.5s/5s/15s）
type ClientConfig struct {
    ReadTimeout     time.Duration  // default 1.5s (edge read)
    WriteTimeout    time.Duration  // default 5s
    PipelineTimeout time.Duration  // default 5s
    SeedTimeout     time.Duration  // default 15s (seed write)
    MaxRetries      int            // default 1
}

func NewUpstashClient(url, token string, cfg ClientConfig) Client
```

#### FetchThrough — 增强版 read-through 缓存（解决 v1 P-6 仅进程内 coalescing + P-8 NEG_SENTINEL 不安全）
```go
// FetchThrough 替代 v1 cachedFetchJson，增加 stale-while-revalidate 和类型安全
type FetchThrough[T any] struct {
    client    Client
    group     singleflight.Group  // 进程级 request coalescing
    metrics   *Metrics            // Prometheus 埋点
}

type FetchOpts[T any] struct {
    Key                  string
    TTL                  time.Duration
    Fetcher              func(ctx context.Context) (*T, error)
    NegativeTTL          time.Duration        // default 120s; 类型安全 sentinel（非裸字符串）
    StaleWhileRevalidate time.Duration        // >0 则返回过期数据同时后台刷新
    Validator            func(raw []byte) bool // 可选校验钩子
}

func (ft *FetchThrough[T]) Fetch(ctx context.Context, opts FetchOpts[T]) (*T, error)
```

#### KeyRegistry — 统一 key 注册表（解决 v1 P-3 碎片化 + P-4 TTL 散乱 + P-9 health 手工维护）
```go
// KeyEntry 集中声明一个 Redis key 的全部属性
// 从 SourceSpec 自动推导：key、seedMetaKey、TTL、maxStale、tier
type KeyEntry struct {
    Key          string        // Redis key, e.g. "seismology:earthquakes:v1"
    SeedMetaKey  string        // 自动推导: "seed-meta:{domain}:{resource}"
    DataTTL      time.Duration // 来自 SourceSpec.DataTTL
    MaxStale     time.Duration // 来自 SourceSpec.MaxStaleDuration
    Tier         string        // "fast" | "slow"
    Category     string        // "bootstrap" | "standalone" | "on-demand" | "derived"
}

type KeyRegistry struct {
    entries map[string]KeyEntry
}

func NewKeyRegistryFromSpecs(specs []SourceSpec) *KeyRegistry
func (r *KeyRegistry) Get(name string) (KeyEntry, bool)
func (r *KeyRegistry) BootstrapKeys() map[string]string    // 替代 v1 BOOTSTRAP_CACHE_KEYS
func (r *KeyRegistry) SeedMetaEntries() map[string]KeyEntry // 替代 v1 SEED_META
func (r *KeyRegistry) Validate() error                      // CI lint: TTL 一致性校验
```

#### Metrics — 缓存可观测性（解决 v1 P-10 缺少可观测性）
```go
type Metrics struct {
    hits     prometheus.Counter     // cache_hit_total
    misses   prometheus.Counter     // cache_miss_total
    errors   prometheus.Counter     // cache_error_total
    timeouts prometheus.Counter     // cache_timeout_total
    latency  prometheus.Histogram   // cache_latency_seconds
    writeBytes prometheus.Histogram // cache_write_bytes
}

func NewMetrics(reg prometheus.Registerer) *Metrics
func (m *Metrics) RecordHit(key string, d time.Duration)
func (m *Metrics) RecordMiss(key string, d time.Duration)
func (m *Metrics) RecordError(key string, errType string)
func (m *Metrics) RecordWrite(key string, bytes int, ttl time.Duration)
```

#### Sidecar — 本地 LRU+TTL 内存缓存（保留 v1 良好设计）
```go
// 仅在 Tauri sidecar 模式下启用的 L1 内存缓存
type SidecarCache struct {
    maxEntries int           // 500
    maxBytes   int64         // 50MB
    sweepEvery time.Duration // 60s
}
func NewSidecarCache(opts ...SidecarOption) *SidecarCache
func (c *SidecarCache) Get(key string) ([]byte, bool)
func (c *SidecarCache) Set(key string, value []byte, ttl time.Duration)
func (c *SidecarCache) Stats() SidecarStats
```

### 2.6 数据发布层 `internal/seed/`

**设计原则**：保留 v1 `runSeed` 的韧性模型（失败→延长 TTL→不丢数据）和 Lua CAS 锁释放。改进：(1) 信封协议用 Go struct 单源定义，消除 v1 三副本同步问题；(2) 原子发布改用 MULTI/EXEC 事务，消除 staging key 残留风险；(3) TTL 从 SourceSpec 中央读取，不再硬编码。

#### Runner — runSeed 编排器
```go
// Runner 编排完整的 seed 生命周期（保留 v1 韧性模型）
type Runner struct {
    client  cache.Client
    metrics *cache.Metrics
    lock    *Lock
    pub     *Publisher
    meta    *MetaWriter
    ttl     *TTLManager
}

// Run 执行: lock → fetch → publish → meta → verify → unlock
// 失败时: extendTTL → 不丢数据（保留 v1 设计）
func (r *Runner) Run(ctx context.Context, src registry.Source) error
```

#### Lock — 分布式锁（保留 v1 最佳实践）
```go
type Lock struct { client cache.Client }

// Acquire: SET seed-lock:{domain}:{resource} NX PX ttl
func (l *Lock) Acquire(ctx context.Context, domain, resource, runID string, ttl time.Duration) (bool, error)

// Release: EVAL Lua CAS — 仅释放自己持有的锁
func (l *Lock) Release(ctx context.Context, domain, resource, runID string) error
```

#### Envelope — 信封协议（解决 v1 P-2 三副本）
```go
// SeedMeta 元数据 — Go struct 单源定义，彻底消除 v1 三文件手工同步
type SeedMeta struct {
    FetchedAt     int64    `json:"fetchedAt"`
    RecordCount   int      `json:"recordCount"`
    SourceVersion string   `json:"sourceVersion"`
    SchemaVersion int      `json:"schemaVersion"`
    State         string   `json:"state"`          // "OK" | "OK_ZERO" | "ERROR"
    FailedDatasets []string `json:"failedDatasets,omitempty"`
    ErrorReason   string   `json:"errorReason,omitempty"`
    GroupID       string   `json:"groupId,omitempty"`
}

type SeedEnvelope struct {
    Seed *SeedMeta `json:"_seed"`
    Data any       `json:"data"`
}

func Unwrap(raw []byte) (*SeedMeta, any, error)          // 兼容 legacy + contract 格式
func Build(meta SeedMeta, data any) (*SeedEnvelope, error)
```

#### Publish — 原子发布（解决 v1 P-7 staging key 残留）
```go
type Publisher struct { client cache.Client }

// AtomicPublish 使用 MULTI/EXEC 事务替代 v1 的 staging→canonical→del 三步操作
// 消除 staging key 残留风险
func (p *Publisher) AtomicPublish(ctx context.Context, key string, data []byte, ttl time.Duration) error

// PublishWithEnvelope: build envelope → validate → MULTI/EXEC
func (p *Publisher) PublishWithEnvelope(ctx context.Context, key string, meta SeedMeta, data any, ttl time.Duration) error
```

#### Meta — 新鲜度元数据
```go
type MetaWriter struct { client cache.Client }

// WriteFreshness: SET seed-meta:{domain}:{resource} {fetchedAt, recordCount, status} EX metaTTL
func (w *MetaWriter) WriteFreshness(ctx context.Context, domain, resource string, count int, source string, ttl time.Duration) error
```

#### TTL — TTL 管理
```go
type TTLManager struct { client cache.Client }

// ExtendExisting: EXPIRE pipeline — 失败时延长现有数据 TTL（保留 v1 韧性设计）
func (t *TTLManager) ExtendExisting(ctx context.Context, keys []string, ttl time.Duration) error

// ResolveTTL: 从 SourceSpec.DataTTL 读取，不再硬编码
func ResolveTTL(spec registry.SourceSpec) time.Duration
func ResolveMetaTTL(spec registry.SourceSpec) time.Duration // max(7d, DataTTL)
```

---

## 三、目录结构

```
backend/monitor/
  cmd/monitor/
    main.go                      ← config → registry.Boot() → shutdown
  internal/
    config/
      config.go                  ← 统一配置
    registry/                    ← ★ 声明式注册模块
      source.go                  ← Source 接口 + FetchResult + FetchMetrics
      source_spec.go             ← SourceSpec 声明符 + Schedule 类型 (Cron/Interval/OnDemand)
      registry.go                ← Registry 容器（注册/Boot/Shutdown/TriggerOnDemand）
      scheduler.go               ← Scheduler（统一 cron/ticker/ondemand 调度）
      health.go                  ← HealthTracker（从 SourceSpec 读取阈值）
      registry_test.go
    preprocess/                  ← ★ 预处理 Pipeline
      pipeline.go                ← Stage 接口 + Pipeline 组合器
      format.go                  ← FormatMapper (P1: 字段映射/单位换算)
      merger.go                  ← MultiSourceMerger (P2: 多源合并)
      preprocess_test.go
    semantic/                    ← ★ 语义层能力工具箱
      classify/                  ← 关键词分类
        classifier.go            ← Classifier + Classify()
        keywords.go              ← 7 层关键词表数据
        exclusions.go            ← 排除词 + 短词 Set
        classifier_test.go
      scoring/                   ← 通用加权评分 + 多 Profile
        weighted.go              ← WeightedScorer + ScoringProfile
        importance.go            ← news ImportanceProfile 预设
        disruption.go            ← chokepoint DisruptionProfile 预设
        goalpost.go              ← resilience GoalpostProfile 预设
        scoring_test.go
      tracking/                  ← 通用有状态追踪
        tracker.go               ← StatefulTracker 接口
        story_tracker.go         ← StoryTracker (新闻生命周期)
        interval_tracker.go      ← IntervalTracker (韧性分数趋势)
        deduplicator.go          ← HashDedup (通用去重)
        accumulator.go           ← Accumulator (Redis ZADD)
        tracking_test.go
      enrichment/                ← AI 缓存增强
        ai_cache.go              ← EnrichBatch (MGET classify:sebuf:v3)
        enrichment_test.go
      tiers/                     ← 源可信度分级
        source_tiers.go          ← GetTier(name) → 1-4
        source_tiers_test.go
      agents/                    ← LLM/Embedding 客户端
        client.go                ← AgentsClient 接口
        http_client.go           ← HTTP 实现
        mock_client.go           ← 测试 mock
      fusion/                    ← P2 实装：跨域派生
        correlation.go           ← CorrelationFuser
        cross_source.go          ← CrossSourceFuser
        fusion_test.go
      geo/                       ← ⏳ 预留骨架
        haversine.go             ← Distance()
      anomaly/                   ← ⏳ 预留骨架
        zscore.go                ← ZScore()
    fetcher/                     ← ★ 采集能力工具箱
      client.go                  ← 基础 HTTP 客户端（timeout/UA/abort）
      proxy/                     ← 代理策略
        strategy.go              ← ProxyStrategy 接口 + 4 种实现
        connect.go               ← HTTPS CONNECT 隧道
        resolver.go              ← 代理配置解析
      ratelimit/                 ← 限速
        limiter.go               ← FixedInterval + TokenBucket
        retry_after.go           ← 429 Retry-After 解析
        backoff.go               ← 指数/线性退避
      pool/                      ← 并发控制
        bounded.go               ← BoundedPool (N 并发 + allSettled)
        fanout.go                ← FanOut (errgroup)
      fallback/                  ← 降级链
        chain.go                 ← FallbackChain[T]
      pagination/                ← 分页
        id_batch.go              ← IDBatchFetcher (HN 模式)
      parser/                    ← 响应解析
        xml.go                   ← RSS/Atom XML 解析
        ics.go                   ← ICS 日历解析
    cache/                       ← ★ 统一缓存层（解决 v1 读写割裂+key 碎片化+TTL 散乱+可观测性缺失）
      client.go                  ← Client 接口 + NewUpstashClient + ClientConfig (~120)
      fetch.go                   ← FetchThrough[T] + singleflight + neg sentinel + SWR (~100)
      pipeline.go                ← Pipeline 批量操作 (~60)
      key_registry.go            ← KeyEntry + KeyRegistry + NewKeyRegistryFromSpecs + Validate (~150)
      metrics.go                 ← Metrics(Prometheus hit/miss/error/latency) (~80)
      sidecar.go                 ← SidecarCache LRU+TTL 内存缓存 (Tauri 模式) (~60)
      cache_test.go
    seed/                        ← ★ 数据发布层（保留韧性模型+改进原子发布+信封单源）
      runner.go                  ← Runner 编排: lock→fetch→publish→meta→verify→unlock (~200)
      lock.go                    ← Lock NX+PX + EVAL Lua CAS 释放 (~60)
      envelope.go                ← SeedMeta struct + SeedEnvelope + Unwrap/Build (~80)
      publish.go                 ← Publisher AtomicPublish MULTI/EXEC (~80)
      meta.go                    ← MetaWriter WriteFreshness (~40)
      ttl.go                     ← TTLManager ExtendExisting + ResolveTTL (~40)
      seed_test.go
    news/
      feeds.go                   ← 5 variants 100+ feeds
      parser.go                  ← RSS/Atom regex 解析
      digest_source.go           ← DigestSource（组合 classify+scoring+tracking+enrichment+tiers）
      insights_source.go         ← InsightsSource（组合 agents）
      register.go                ← RegisterAll(reg)
    research/
      arxiv_source.go            ← ArxivSource（无 semantic 依赖）
      hackernews_source.go       ← HackerNewsSource
      tech_events_source.go      ← TechEventsSource
      trending_source.go         ← TrendingSource
      city_coords.go
      register.go                ← RegisterAll(reg)
  go.mod / go.sum / .env.example
```

---

## 四、数据流

### 4.1 News DigestSource（组合 5 个语义能力）

```
news.RegisterAll(reg):
  DigestSource{variant:"full",lang:"en",
    classifier: classify.NewNewsClassifier(),      ← semantic/classify
    scorer:     scoring.NewImportanceScorer(),      ← semantic/scoring
    tracker:    tracking.NewStoryTracker(rdb),      ← semantic/tracking
    enricher:   enrichment.New(rdb),                ← semantic/enrichment
    tiers:      tiers.Load(),                       ← semantic/tiers
  }

Cron tick → Registry.dispatch(DigestSource)
  1. seed.AcquireLock() + registry.health.RecordAttempt()
  2. pool.BoundedPool(20) × proxy.DirectFirst(relay)   ← fetcher/pool + fetcher/proxy
  3. parser.ParseRssXML() → []ParsedItem               ← fetcher/parser
  ┌── DigestSource 内部调用各 semantic 子包 ────────────────┐
  │ 4. s.classifier.Classify(title, WithVariant(variant))  │
  │ 5. s.enricher.EnrichBatch(ctx, items)                  │
  │ 6. s.tracker.ComputeCorroboration(items)               │
  │ 7. s.scorer.Score(ScoringInput{components...})         │
  │ 8. s.tracker.ReadAndWriteTracks(ctx, items)            │
  └─────────────────────────────────────────────────────────┘
  9. 排序+截断 → toProtoItem
  ┌── seed.Runner.Run() 统一编排（cache/seed 调用链）──────────────┐
  │ 10a. seed.Lock.Acquire(NX+PX)                  ← seed/lock    │
  │ 10b. seed.Publisher.PublishWithEnvelope(MULTI/EXEC)             │
  │      → cache.Client.Pipeline([SET canonical, SET meta])         │
  │      → cache.Metrics.RecordWrite(key, bytes, ttl)  ← O-7      │
  │ 10c. seed.MetaWriter.WriteFreshness()           ← seed/meta    │
  │ 10d. seed.TTLManager.ExtendExisting() (on fail) ← seed/ttl    │
  │ 10e. seed.Lock.Release(EVAL Lua CAS)            ← seed/lock   │
  └─────────────────────────────────────────────────────────────────┘
  11. registry.health.RecordSuccess(metrics)
```

### 4.2 News InsightsSource（只用 agents）

```
Cron tick → Registry.dispatch(InsightsSource)
  1. cache.FetchThrough.Fetch("news:digest:v1:full:en")    ← cache/fetch (read-through)
  2. extractTopHeadlines(10)
  3. s.agents.Summarize(ctx, headlines, {Mode: "brief"})   ← semantic/agents
  4. seed.Publisher.AtomicPublish("news:insights:v1", MULTI/EXEC)  ← seed/publish
```

### 4.3 Research ArxivSource（用 ratelimit + parser）

```
Cron tick → Registry.dispatch(ArxivSource)
  1. ratelimit.NewFixedInterval(3s)                      ← fetcher/ratelimit
  2. fetcher.client.Get(arXiv API) × 3 类别
  3. parser.ParseAtomXML() → []Paper                     ← fetcher/parser
  4. seed.Runner.Run() → Lock + Publisher.AtomicPublish(MULTI/EXEC) + Meta + Unlock
```

### 4.4 未来：CyberSource（preprocess + scoring）

```
cyber.RegisterAll(reg):
  CyberSource{
    pipeline: preprocess.NewPipeline(
      preprocess.FormatMapper(cyberNormalize),        ← preprocess/format
      preprocess.MultiSourceMerger(3_feeds),          ← preprocess/merger
    ),
    scorer: scoring.NewDisruptionScorer(),             ← semantic/scoring
  }

Interval(2h) → Registry.dispatch(CyberSource)
  1. pipeline.Run(ctx, raw)
     → FormatMapper: parseFeodoRecord + parseUrlhausRecord + normalize
     → MultiSourceMerger: Feodo + URLhaus + OTX → 合并去重
  2. s.scorer.Score(ScoringInput{severity...})
  3. seed.AtomicPublish("cyber:threats:v2")
```

### 4.5 未来：CorrelationSource（fusion 实装）

```
correlation.RegisterAll(reg):
  CorrelationSource{
    fuser: fusion.NewCorrelationFuser(rdb),            ← semantic/fusion
  }

Interval(15m) → Registry.dispatch(CorrelationSource)
  1. s.fuser.Fuse(ctx, [digest, market, unrest, resilience])
     → 读多域 seed key → 交叉相关信号计算
  2. seed.AtomicPublish("correlation:cards-bootstrap:v1")
```

### 4.6 未来：SeismologySource（只用 geo）

```
Cron tick → Registry.dispatch(EarthquakeSource)
  1. fetchUSGS()
  2. s.proximity.NearestMatch(coords, nuclearTestSites, 100km)  ← semantic/geo
  3. computeConcernScore()
  4. seed.AtomicPublish()
```

---

## 五、实施步骤

### Phase 0：骨架 + 配置（1 天）
- `go.mod` / `.env.example` / `config/config.go` / `cmd/monitor/main.go` 骨架

### Phase 1A：声明式注册模块 `internal/registry/`（1.5 天）
- `source.go` Source 接口 + `source_spec.go` SourceSpec 声明符 + Schedule 类型
- `registry.go` 容器 + `scheduler.go`（统一 Cron/Interval/OnDemand 调度）
- `health.go` HealthTracker（从 SourceSpec 读取 MaxStaleDuration/MinRecordCount）
- 含 MockSource 单元测试

### Phase 1B：语义层能力工具箱 `internal/semantic/`（2.5 天）
| 子包 | 文件 | LOC | 源对标 |
|------|------|-----|--------|
| `classify/` | classifier.go + keywords.go + exclusions.go | ~250 | `_classifier.ts` (253行) |
| `scoring/` | weighted.go + importance.go + disruption.go + goalpost.go | ~180 | news importanceScore + chokepoint disruptionScore + resilience goalposts |
| `tiers/` | source_tiers.go | ~120 | `shared/source-tiers.json` |
| `tracking/` | tracker.go + story_tracker.go + interval_tracker.go + deduplicator.go + accumulator.go | ~300 | news storyTracking + resilience scoreInterval |
| `enrichment/` | ai_cache.go | ~60 | enrichWithAiCache |
| `agents/` | client.go + http_client.go + mock_client.go | ~100 | agents HTTP client |
| `fusion/` | correlation.go + cross_source.go | ~150 | seed-correlation + seed-cross-source-signals |
| `geo/` ⏳ | haversine.go (骨架) | ~20 | 预留 |
| `anomaly/` ⏳ | zscore.go (骨架) | ~20 | 预留 |

### Phase 1C：采集能力工具箱 `internal/fetcher/`（2 天）
| 子包 | 文件 | LOC | 用于 |
|------|------|-----|------|
| `client.go` | 基础 HTTP | ~80 | 所有 Source |
| `proxy/` | strategy + connect + resolver | ~150 | news (relay), 未来 FRED/Yahoo |
| `ratelimit/` | limiter + backoff + retry_after | ~120 | arXiv (3s), 未来 GDELT/CoinGecko |
| `pool/` | bounded + fanout | ~120 | news (20并发), research (4源并行) |
| `fallback/` | chain | ~80 | trending (OSSInsight→GitHub) |
| `pagination/` | id_batch | ~60 | HN (top500 batch) |
| `parser/` | xml + ics | ~150 | arXiv/RSS/Techmeme |

### Phase 1E：预处理 Pipeline `internal/preprocess/`（0.5 天）
| 文件 | LOC | 功能 |
|------|-----|------|
| `pipeline.go` | ~60 | Stage 接口 + Pipeline 组合器 |
| `format.go` | ~80 | FormatMapper (P1: 字段映射/单位换算) |
| `merger.go` | ~60 | MultiSourceMerger (P2: 多源并行读取+合并) |

### Phase 1F-cache：统一缓存层 `internal/cache/`（1.5 天）
| 文件 | LOC | 功能 | 对标 v1 优化 |
|------|-----|------|-------------|
| `client.go` | ~120 | Client 接口 + Upstash 实现 + ClientConfig | O-1: 统一客户端 |
| `fetch.go` | ~100 | FetchThrough[T] + singleflight + neg sentinel + SWR | O-5: 增强 read-through |
| `pipeline.go` | ~60 | 批量操作 | 保留 |
| `key_registry.go` | ~150 | KeyEntry + KeyRegistry + SourceSpec 推导 + Validate | O-3+O-4: 统一注册表+TTL |
| `metrics.go` | ~80 | Prometheus counter/histogram | O-7: 可观测性 |
| `sidecar.go` | ~60 | LRU+TTL 内存缓存 (Tauri) | 保留良好设计 |

### Phase 1F-seed：数据发布层 `internal/seed/`（1 天）
| 文件 | LOC | 功能 | 对标 v1 优化 |
|------|-----|------|-------------|
| `runner.go` | ~200 | Runner 编排 (lock→fetch→publish→meta→unlock) | 保留韧性模型 |
| `lock.go` | ~60 | NX+PX + EVAL Lua CAS 释放 | 保留最佳实践 |
| `envelope.go` | ~80 | SeedMeta struct + Unwrap/Build | O-2: 单源定义 |
| `publish.go` | ~80 | AtomicPublish MULTI/EXEC 事务 | O-6: 原子发布 |
| `meta.go` | ~40 | WriteFreshness | 保留 |
| `ttl.go` | ~40 | ExtendExisting + ResolveTTL(SourceSpec) | O-4: 中央 TTL |

### Phase 2：News 数据层（2-3 天）
| 文件 | 内容 |
|------|------|
| `news/feeds.go` | 5 variants 100+ feeds (~400 LOC) |
| `news/parser.go` | RSS regex 解析 (~180 LOC) |
| `news/digest_source.go` | DigestSource 实现 Source，编排 buildDigest (~200 LOC) |
| `news/insights_source.go` | InsightsSource 调 semantic.Summarize (~120 LOC) |
| `news/register.go` | RegisterAll → 6 Sources |

### Phase 3：Research 数据层 + fusion 实装（1.5-2 天）
- `arxiv_source.go` / `hackernews_source.go` / `tech_events_source.go` / `trending_source.go`
- 各实现 Source 接口，独立调度
- `register.go` → RegisterAll → 4 Sources
- fusion/ 实装：CorrelationFuser + CrossSourceFuser 集成测试

### Phase 4：入口 + 可观测性（1 天）
- `main.go` 完整版：registry.Boot() → healthz → shutdown
- HealthTracker → `/healthz` JSON 输出各源健康状态
- **cache.Metrics 接入**（对应 v1 分析 O-7）：
  - `cache_hit_total` / `cache_miss_total` / `cache_error_total` / `cache_timeout_total` (Prometheus Counter)
  - `cache_latency_seconds` / `cache_write_bytes` (Prometheus Histogram)
  - `seed_run_duration_seconds` / `seed_publish_bytes` (Prometheus Histogram)
  - `/healthz` 补充 Redis INFO stats（key count, memory usage, connected clients）
  - 所有 `cache.Client` 操作自动埋点，`FetchThrough` 的 hit/miss 自动上报

### Phase 5：测试（贯穿）
- semantic 纯函数单元测试
- Registry MockSource 测试
- Redis 集成（miniredis）
- HTTP 录制回放（go-vcr）

---

## 六、源码对标映射表（更新版）

| Go 文件 | 源 TS/MJS | 关键逻辑 |
|---------|-----------|----------|
| **registry/source.go** | `go-engine-review` 建议 | Source 接口 + FetchResult |
| **registry/source_spec.go** | `api/health.js` SEED_META + Railway cron + relay setInterval | SourceSpec 声明符 + Schedule(Cron/Interval/OnDemand) |
| **registry/scheduler.go** | `_bundle-runner.mjs` + `ais-relay.cjs` setInterval | 统一调度器 |
| **registry/registry.go** | `backend-registry-rss` DomainRegistration | 容器+调度+依赖排序 |
| **registry/health.go** | `api/health.js` SEED_META (~110条) | HealthTracker（从 SourceSpec 自动读阈值） |
| **preprocess/pipeline.go** | 各 seed 脚本内散的预处理 | Stage 接口 + Pipeline 组合器 |
| **preprocess/format.go** | cyber `_shared.ts` parseRecord + BIS SDMX | FormatMapper (P1) |
| **preprocess/merger.go** | `get-chokepoint-status` 4源合并 | MultiSourceMerger (P2) |
| **semantic/classify/** | `_classifier.ts` (253行) | 7 层关键词 + 排除词 + 正则缓存 |
| **semantic/scoring/** | news importanceScore + chokepoint disruptionScore + resilience goalposts | 通用 WeightedScorer + 3 Profile 预设 |
| **semantic/tiers/** | `source-tiers.json` | 100+ 源 → tier 1-4 |
| **semantic/tracking/** | news storyTracking + resilience scoreInterval | StatefulTracker + StoryTracker + IntervalTracker |
| **semantic/enrichment/** | `list-feed-digest.ts` enrichWithAiCache | MGET classify:sebuf:v3:{hash} |
| **semantic/agents/** | `seed-insights.mjs` agents 调用 | AgentsClient 接口 + HTTP 实现 |
| **semantic/fusion/** | `seed-correlation.mjs` + `seed-cross-source-signals.mjs` | CorrelationFuser + CrossSourceFuser |
| **fetcher/proxy/** | `_shared/relay.ts` + `_seed-utils.mjs` proxy helpers | DirectFirst (relay) + ProxyFirst + CurlOnly + TwoLeg |
| **fetcher/ratelimit/** | arXiv sleep(3s), Yahoo parseRetryAfterMs | FixedInterval + RetryAfter + Backoff |
| **fetcher/pool/** | `list-feed-digest.ts` batch=20, `seed-research.mjs` Promise.all | BoundedPool + FanOut |
| **fetcher/fallback/** | `seed-token-panels.mjs` CoinGecko→CoinPaprika | FallbackChain[T] |
| **fetcher/pagination/** | `seed-research.mjs` HN batch fetch | IDBatchFetcher |
| **fetcher/parser/** | `list-feed-digest.ts` RSS regex, arXiv Atom | ParseRssXML + ParseAtomXML + ParseICS |
| **cache/client.go** | `_shared/redis.ts` getCachedJson/setCachedJson + `_seed-utils.mjs` redisCommand | O-1: 统一 Client 接口（消除读写割裂） |
| **cache/fetch.go** | `_shared/redis.ts` cachedFetchJson + inflight Map | O-5: FetchThrough + singleflight + SWR |
| **cache/key_registry.go** | `_shared/cache-keys.ts` BOOTSTRAP_CACHE_KEYS + `api/health.js` SEED_META 130+ 条 | O-3+O-4: 统一注册表 + TTL 自动推导 |
| **cache/metrics.go** | 新增（v1 仅 console.warn 日志） | O-7: Prometheus 可观测性 |
| **cache/sidecar.go** | `_shared/sidecar-cache.ts` (114行) | 保留 LRU+TTL 设计 |
| **seed/runner.go** | `_seed-utils.mjs` runSeed (777-994行) | 保留韧性模型 (失败→extendTTL) |
| **seed/lock.go** | `_seed-utils.mjs` acquireLock/releaseLock | 保留 Lua CAS 最佳实践 |
| **seed/envelope.go** | `_seed-envelope-source.mjs` + `seed-envelope.ts` + `api/_seed-envelope.js` (三副本) | O-2: Go struct 单源定义 |
| **seed/publish.go** | `_seed-utils.mjs` atomicPublish (staging→canonical→del) | O-6: MULTI/EXEC 事务替代 |
| **seed/meta.go** | `_seed-utils.mjs` writeFreshnessMetadata | 保留 |
| **seed/ttl.go** | `_seed-utils.mjs` extendExistingTtl + 各 seed 硬编码 TTL | O-4: 从 SourceSpec.DataTTL 读取 |
| news/feeds.go | `_feeds.ts` (437行) | 5 variants × 多 category |
| news/parser.go | `list-feed-digest.ts` 185-268行 | regex + CDATA + entity |
| news/digest_source.go | `list-feed-digest.ts` 477-629行 | buildDigest（组合 classify+scoring+tracking+enrichment+tiers） |
| news/insights_source.go | `seed-insights.mjs` | 读 digest → agents.Summarize → 写 |
| research/*_source.go | `seed-research.mjs` | arXiv/HN/events/trending（无 semantic 依赖） |

---

## 七、工作量估算

| Phase | 内容 | 天数 | LOC |
|-------|------|------|-----|
| P0 骨架+配置 | go.mod/config/main | 1 | ~200 |
| P1A 声明式注册 | registry/ (SourceSpec+Scheduler+Health) | 1.5 | ~420 |
| P1B 语义层工具箱 | semantic/ (6 子包 + 2 骨架 + fusion 实装) | 2.5 | ~1200 |
| P1C 采集工具箱 | fetcher/ (6 子包) | 2 | ~760 |
| P1E 预处理 Pipeline | preprocess/ (pipeline+format+merger) | 0.5 | ~200 |
| P1F-cache 统一缓存层 | cache/ (Client+FetchThrough+KeyRegistry+Metrics+Sidecar) | 1.5 | ~570 |
| P1F-seed 数据发布层 | seed/ (Runner+Lock+Envelope+Publish+Meta+TTL) | 1 | ~500 |
| P2 News | feeds/digest_source/insights_source | 2-3 | ~740 |
| P3 Research+fusion | 4 source + fusion 实装 + coords + register | 1.5-2 | ~650 |
| P4 入口+观测 | main/healthz | 1 | ~200 |
| P5 测试 | 贯穿 | +2 | ~700 |
| **合计** | | **16.5-19 天** | **~6140** |

---

## 八、与 v1 方案的主要变化

| 维度 | v1 | v2（本方案） |
|------|-----|-------------|
| 数据源注册 | `register.go` 直接注册 cron | SourceSpec 声明符 → Registry 容器 → 统一 Cron/Interval/OnDemand 调度 |
| 预处理 | 分散在各 seed 脚本/RPC handler | `preprocess/` Pipeline（可组合 Stage） |
| 语义处理 | 分散在 `news/classifier.go` 等 | `semantic/` 能力工具箱（6 子包 + fusion 实装），Source 按需组合 |
| 数据采集 | `fetcher/http.go + relay.go` | `fetcher/` 采集工具箱（6 子包），Source 按需组合 |
| 健康检查 | 手动 healthz + health.js SEED_META | HealthTracker 从 SourceSpec 自动读取阈值 |
| Redis 客户端 | 读 (`redis.ts`) / 写 (`_seed-utils.mjs`) 完全割裂，2 套超时/前缀/错误策略 | `cache.Client` 统一接口，单点维护 |
| Key 管理 | 4 处散落 (`cache-keys.ts` + `health.js` + seed 脚本 + handler 内联) | `cache.KeyRegistry` 从 SourceSpec 自动推导 |
| 信封协议 | 3 文件手工同步 (`_seed-envelope-source.mjs` + `.ts` + `.js`) | `seed.Envelope` Go struct 单源定义 |
| 原子发布 | staging key → canonical → del（非原子，staging 残留风险） | `seed.Publisher` MULTI/EXEC 事务 |
| TTL 配置 | 硬编码在 118 个 seed 脚本 | `SourceSpec.DataTTL` 中央声明 + TTLManager |
| 缓存可观测性 | console.warn/error 日志（无命中率、无延迟分布） | `cache.Metrics` Prometheus counter/histogram |
| 扩展性 | 新 domain 需改 main.go | 只需实现 Source + 引入所需 preprocess/semantic/fetcher 子包 |
| main.go | 直接操作 cron | registry.Boot() 一行启动 |

---

## 九、语义层各子包复用矩阵

| 子包 | news Digest | news Insights | research | 未来 cyber | 未来 resilience | 未来 chokepoint | 未来 correlation |
|------|:-:|:-:|:-:|:-:|:-:|:-:|:-:|
| `classify/` | ✅ | — | — | — | — | — | — |
| `scoring/` Importance | ✅ | — | — | — | — | — | — |
| `scoring/` Disruption | — | — | — | ✅ severity | — | ✅ disruptionScore | — |
| `scoring/` Goalpost | — | — | — | — | ✅ 19维度 | — | — |
| `tracking/` Story | ✅ | — | — | — | — | — | — |
| `tracking/` Interval | — | — | — | — | ✅ scoreInterval | — | — |
| `enrichment/` | ✅ | — | — | — | — | — | — |
| `tiers/` | ✅ | — | — | — | — | — | — |
| `agents/` | — | ✅ | — | — | — | — | — |
| `fusion/` Correlation | — | — | — | — | — | — | ✅ |
| `fusion/` CrossSource | — | — | — | — | — | — | ✅ |
| `geo/` | — | — | ⚠️ city_coords | — | — | — | — |
| `anomaly/` | — | — | — | — | — | — | — |

---

## 十、采集层各子包复用矩阵

| 子包 | news Digest | news Insights | research arXiv | research HN | research Events | research Trending | 未来 PortWatch | 未来 Market |
|------|:-:|:-:|:-:|:-:|:-:|:-:|:-:|:-:|
| `proxy/` | ✅ DirectFirst | — | — | — | ✅ relay | — | — | ✅ CurlOnly |
| `ratelimit/` | — | — | ✅ 3s | — | — | — | — | ✅ 429 |
| `pool/bounded` | ✅ 20并发 | — | — | — | — | — | ✅ 5并发 | — |
| `pool/fanout` | — | — | — | — | — | — | — | — |
| `fallback/` | — | — | — | — | — | ✅ OSS→GH | — | ✅ AV→FH→YH |
| `pagination/` | — | — | — | ✅ IDBatch | — | — | ✅ Offset | — |
| `parser/xml` | ✅ RSS | — | ✅ Atom | — | — | — | — | — |
| `parser/ics` | — | — | — | — | ✅ | — | — | — |

---

## 十一、预处理层复用矩阵

| Stage | news Digest | 未来 cyber | 未来 resilience | 未来 chokepoint | 未来 economic |
|-------|:-:|:-:|:-:|:-:|:-:|
| `FormatMapper` | ✅ RSS parse | ✅ normalize | — | — | ✅ WEO pivot |
| `MultiSourceMerger` | — | ✅ 3 feeds | ✅ 19 dims | ✅ 4 sources | — |
| (P0 无预处理) | — | — | — | — | — |

---

## 十二、缓存/存储层复用矩阵

| 子模块 | news Digest | news Insights | research | 未来 cyber | 未来 resilience | 未来 chokepoint | 未来 correlation | 未来 market |
|--------|:-:|:-:|:-:|:-:|:-:|:-:|:-:|:-:|
| `cache.Client` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| `cache.FetchThrough` | ✅ read-through | ✅ 读 digest | — | ✅ | ✅ | ✅ | ✅ 读多域 | ✅ |
| `cache.KeyRegistry` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| `cache.Pipeline` | ✅ batch write | — | — | — | ✅ 19 dims | ✅ 4 sources | ✅ 读多域 | — |
| `cache.Metrics` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| `cache.Sidecar` | ✅ (Tauri) | ✅ (Tauri) | ✅ (Tauri) | — | — | — | — | — |
| `seed.Runner` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| `seed.Lock` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| `seed.Envelope` | ✅ contract | — legacy | ✅ contract | ✅ | ✅ | ✅ | ✅ | ✅ |
| `seed.Publisher` | ✅ MULTI/EXEC | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| `seed.Meta` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| `seed.TTL` | ✅ 86400s | ✅ 86400s | ✅ 86400s | ✅ 172800s | ✅ 604800s | ✅ 172800s | ✅ 86400s | ✅ 86400s |

**说明**：cache/seed 层是全 Source 共用的基础设施（复用率 100%），与 semantic/fetcher 的"按需组合"模式不同。表中标注的差异点主要体现在：
- **FetchThrough**：CorrelationSource 需读多域 seed key (digest+market+unrest+resilience)
- **Pipeline**：多维度合并 Source 需批量读取多个 key
- **Sidecar**：仅 Tauri 桌面模式下启用
- **Envelope**：部分 legacy Source 尚未迁移到 contract 模式
- **TTL**：各 Source 的 DataTTL 不同，但均从 SourceSpec 中央声明
