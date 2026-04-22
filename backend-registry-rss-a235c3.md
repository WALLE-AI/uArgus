# Backend Registry + 中间件 + News & Research 两个 Domain（修正版）

在空的 `backend/` 目录从零构建 Go 服务。分层原则：**Go 后端 = SERVER 层**，数据定义从 `src/config/` + `src/services/` 提取（feeds/classifier/geo-hub），处理流程从 `server/` handler 提取（regex 解析/Redis 交互/story tracking），中间件从 `server/gateway.ts` 完整复制。`src/services/{domain}/index.ts` 是浏览器 RPC 客户端封装（CircuitBreaker/IndexedDB/DOMParser），**不直接移植**。

---

## 最终目录结构

```
backend/
  cmd/server/
    main.go                  ← HTTP server 启动，注册所有 domain
  internal/
    registry/
      registry.go            ← DomainRegistration + Registry 核心类型
      router.go              ← 路由匹配（static O(1) + dynamic 线性）★ server/router.ts
      cors.go                ← CORS 规则 ★ server/cors.ts
      rate_limit.go          ← 全局 + 端点级滑动窗口限流 ★ server/_shared/rate-limit.ts
      premium_paths.go       ← PREMIUM_RPC_PATHS ★ src/shared/premium-paths.ts
      cache_tiers.go         ← RPC_CACHE_TIER 全量映射 + TIER_HEADERS ★ server/gateway.ts
      entitlement.go         ← 双体系鉴权 ★ server/gateway.ts + server/_shared/entitlement-check.ts
      middleware.go          ← 10 步管道组装（调用以上各模块）★ server/gateway.ts
    cache/
      client.go              ← go-redis v9 客户端封装
      fetch.go               ← CachedFetchJSON + singleflight + neg sentinel
      keys.go                ← 缓存 Key 常量
      pipeline.go            ← RunRedisPipeline（HMGET/HSET/ZADD 等）
    news/
      feeds.go               ← ★ 源自 src/config/feeds.ts（87KB 完整版）
      classifier.go          ← ★ 源自 src/services/threat-classifier.ts（完整版）
      geo_hub.go             ← ★ 源自 src/services/geo-hub-index.ts（新增）
      source_tiers.go        ← 源自 server/_shared/source-tiers.ts（两层共用）
      parser.go              ← RSS/Atom regex 解析 ★ server/news/v1/list-feed-digest.ts
      scoring.go             ← importanceScore 4权重公式
      story_tracking.go      ← Redis story:track 读写
      digest.go              ← buildDigest 主流程
      handler.go             ← newsHandler（实现 Handler 接口）
    research/
      arxiv.go               ← ★ 源自 server/research/v1/list-arxiv-papers.ts + src/generated 类型
      hackernews.go          ← ★ 源自 server/research/v1/list-hackernews-items.ts + src/generated 类型
      trending_repos.go      ← ★ 源自 server/research/v1/list-trending-repos.ts + src/generated 类型
      tech_events.go         ← 源自 server/research/v1/list-tech-events.ts（ICS+RSS+地理编码）
      city_coords.go         ← 静态地理坐标数据（api/data/city-coords）
      handler.go             ← researchHandler（实现 Handler 接口）
  go.mod
  go.sum
  .env.example
```

> ★ = 权威来源：数据定义从 `src/` 提取，处理流程从 `server/` 提取

---

## Phase 1：项目骨架（go.mod + 目录）

**文件：`backend/go.mod`**
```
module github.com/wall-ai/uargus/backend

go 1.22

require (
    github.com/redis/go-redis/v9  v9.5.1
    golang.org/x/sync             v0.7.0
)
```

**文件：`backend/.env.example`**
```
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0
SERVER_PORT=8080
ENV=development
```

---

## Phase 2：Registry 核心类型

**文件：`internal/registry/registry.go`**

```go
// Handler 是所有数据源处理函数的统一签名
type Handler func(ctx context.Context, req *http.Request) (any, error)

// DomainRegistration 对应原项目每个 api/{domain}/v1/[rpc].ts
type DomainRegistration struct {
    Domain  string             // "news", "aviation", "trade"
    Version string             // "v1"
    Routes  map[string]Handler // "list-feed-digest" → handler func
}

type Registry struct {
    router *Router
}

func NewRegistry() *Registry
func (r *Registry) Register(d DomainRegistration)
func (r *Registry) ServeHTTP(w http.ResponseWriter, req *http.Request)
```

> 注：CacheTier 不再存在 Registration 中 — 统一由 `cache_tiers.go` 的 `RPC_CACHE_TIER` 查表决定（完全镜像 gateway.ts 中心化维护原则）

---

## Phase 3：中间件与注册机制（完全以源码为准）

### 3.1 `internal/registry/cors.go`
**源**：`server/cors.ts`

```go
// 生产环境允许的 Origin 正则列表（完整复制）
var productionPatterns = []*regexp.Regexp{
    regexp.MustCompile(`^https://(.*\.)?worldmonitor\.app$`),
    regexp.MustCompile(`^https://worldmonitor-[a-z0-9-]+-elie-[a-z0-9]+\.vercel\.app$`),
    regexp.MustCompile(`^https?://tauri\.localhost(:\d+)?$`),
    regexp.MustCompile(`^https?://[a-z0-9-]+\.tauri\.localhost(:\d+)?$`),
    regexp.MustCompile(`^tauri://localhost$`),
    regexp.MustCompile(`^asset://localhost$`),
}
// 开发环境追加 localhost / 127.0.0.1

func IsAllowedOrigin(origin string) bool
func GetCORSHeaders(req *http.Request) map[string]string
// Allow-Headers: Content-Type, Authorization, X-WorldMonitor-Key, X-Api-Key,
//               X-Widget-Key, X-Pro-Key
func IsDisallowedOrigin(req *http.Request) bool
```

### 3.2 `internal/registry/cache_tiers.go`
**源**：`server/gateway.ts`（`TIER_HEADERS` + `TIER_CDN_CACHE` + `RPC_CACHE_TIER` 完整映射）

```go
type CacheTier string
const (
    TierFast        CacheTier = "fast"         // max-age=60, s-maxage=300, swr=60, sie=600
    TierMedium      CacheTier = "medium"        // max-age=120, s-maxage=600, swr=120, sie=900
    TierSlow        CacheTier = "slow"          // max-age=300, s-maxage=1800, swr=300, sie=3600
    TierSlowBrowser CacheTier = "slow-browser"  // max-age=300, no s-maxage
    TierStatic      CacheTier = "static"        // max-age=600, s-maxage=3600, swr=600, sie=14400
    TierDaily       CacheTier = "daily"         // max-age=3600, s-maxage=14400, swr=7200, sie=172800
    TierNoStore     CacheTier = "no-store"
)

// TIER_HEADERS — 完整 stale-while-revalidate + stale-if-error 值
var tierHeaders = map[CacheTier]string{
    TierFast:   "public, max-age=60, s-maxage=300, stale-while-revalidate=60, stale-if-error=600",
    TierMedium: "public, max-age=120, s-maxage=600, stale-while-revalidate=120, stale-if-error=900",
    TierSlow:   "public, max-age=300, s-maxage=1800, stale-while-revalidate=300, stale-if-error=3600",
    TierSlowBrowser: "max-age=300, stale-while-revalidate=60, stale-if-error=1800",
    TierStatic: "public, max-age=600, s-maxage=3600, stale-while-revalidate=600, stale-if-error=14400",
    TierDaily:  "public, max-age=3600, s-maxage=14400, stale-while-revalidate=7200, stale-if-error=172800",
    TierNoStore: "no-store",
}

// TIER_CDN_CACHE — Vercel CDN 专用 (CDN-Cache-Control header)
var tierCDNCache = map[CacheTier]string{ /* 完整复制 */ }  

// RPC_CACHE_TIER — 全量 80+ 端点映射（完整复制 gateway.ts）
// 例："/api/news/v1/list-feed-digest" → TierSlow
//     "/api/aviation/v1/track-aircraft" → TierNoStore
//     "/api/research/v1/list-arxiv-papers" → TierStatic
var rpcCacheTier = map[string]CacheTier{ /* 完整 80+ 条 */ }

func GetCacheTier(pathname string) CacheTier  // 查表，默认 TierMedium
```

### 3.3 `internal/registry/premium_paths.go`
**源**：`src/shared/premium-paths.ts`（单一来源，server 与 src 共用）

```go
// 完整复制 36 条 premium path
var premiumRPCPaths = map[string]bool{
    "/api/market/v1/analyze-stock":                    true,
    "/api/intelligence/v1/get-regional-snapshot":      true,
    "/api/supply-chain/v1/get-country-chokepoint-index": true,
    "/api/scenario/v1/run-scenario":                   true,
    // ... 完整 36 条
}

func IsPremiumPath(pathname string) bool
```

### 3.4 `internal/registry/rate_limit.go`
**源**：`server/_shared/rate-limit.ts`

原项目用 Upstash Redis REST API 的 `Ratelimit.slidingWindow`，Go 改用标准 Redis 实现等效的滑动窗口：

```go
// IP 提取优先级（完整复制）：cf-connecting-ip > x-real-ip > x-forwarded-for > 0.0.0.0
func getClientIP(req *http.Request) string

// 全局限流：600 req / 60s 滑动窗口
// Redis 实现：ZADD + ZREMRANGEBYSCORE + ZCOUNT + EXPIRE（原子 pipeline）
func CheckRateLimit(req *http.Request, corsHeaders map[string]string) *http.Response

// 端点级限流（完整复制 ENDPOINT_RATE_POLICIES）
var endpointRatePolicies = map[string]struct{ limit int; windowSec int }{
    "/api/news/v1/summarize-article-cache":    {3000, 60},
    "/api/intelligence/v1/classify-event":     {600, 60},
    "/api/sanctions/v1/lookup-sanction-entity": {30, 60},
    "/api/leads/v1/submit-contact":            {3, 3600},
    "/api/leads/v1/register-interest":         {5, 3600},
    "/api/scenario/v1/run-scenario":           {10, 60},
}

func HasEndpointRatePolicy(pathname string) bool
func CheckEndpointRateLimit(req *http.Request, pathname string, corsHeaders map[string]string) *http.Response
```

### 3.5 `internal/registry/entitlement.go`
**源**：`server/gateway.ts` + `server/_shared/entitlement-check.ts` + `server/_shared/auth-session.ts`

**双体系鉴权**（核心修正项）：

Gateway 同时包含两套独立的访问控制：

| 体系 | 来源 | 用途 |
|------|------|------|
| `premiumRPCPaths` | `src/shared/premium-paths.ts` | 旧式 Pro/Bearer 门控（36 条路径） |
| `endpointEntitlements` | `server/_shared/entitlement-check.ts` | 新式 Tier 门控（当前 4 条，tier ≥ 2） |

Gateway 管道中两者互补：
```go
isTierGated := GetRequiredTier(pathname) != nil
needsLegacyProGate := IsPremiumPath(pathname) && !isTierGated
// isTierGated 优先于 PREMIUM_RPC_PATHS
```

```go
// API Key 验证（X-WorldMonitor-Key 或 X-Api-Key header）
func ValidateAPIKey(req *http.Request) (valid bool, required bool, isUserKey bool, err string)

// Clerk JWT 验证（JWKS 端点 RS256）
func ResolveSessionUserID(req *http.Request) (string, error)

// ENDPOINT_ENTITLEMENTS 查表（tier-based，当前 4 条）
var endpointEntitlements = map[string]int{
    "/api/market/v1/analyze-stock":              2,
    "/api/market/v1/get-stock-analysis-history": 2,
    "/api/market/v1/backtest-stock":             2,
    "/api/market/v1/list-stored-stock-backtests": 2,
}
func GetRequiredTier(pathname string) *int

// Entitlement 检查：Redis 缓存 → Convex 降级
// Fail-closed：查询失败 → 拒绝请求
// Redis key: "entitlements:{env_prefix}:{userId}"，TTL=15min
func CheckEntitlement(req *http.Request, pathname string, corsHeaders map[string]string) *http.Response

// Entitlements 获取（请求合并 — 多个并发请求同一 userId 共享一个 in-flight promise）
func GetEntitlements(userId string) (*CachedEntitlements, error)
```

> **部署说明**：Convex API Key 通过环境变量 `CONVEX_URL` + `CONVEX_DEPLOY_KEY` 配置；初期可设为 stub（返回 nil）跳过鉴权，仅用静态 Key 保护。

### 3.6 `internal/registry/middleware.go`
**源**：`server/gateway.ts` 完整 10 步管道

```
Request
  1. isDisallowedOrigin → 403
  2. getCORSHeaders
  3. OPTIONS preflight → 204
  4. getRequiredTier (isTierGated)
  5. resolveSessionUserID (仅 tier-gated)
  6. validateAPIKey (静态 + wm_ 用户 Key)
  7. legacy Pro bearer gate (PREMIUM_RPC_PATHS && !isTierGated)
  8. checkEntitlement (tier-gated 端点)
  9. checkEndpointRateLimit (端点级)
  10. checkRateLimit (全局，仅无端点策略时)
  11. router.match (POST→GET compat：body JSON → query params)
  12. handler 执行（error boundary）
  13. 合并 CORS + side-channel headers
  14. GET 200：检测 upstreamUnavailable → no-store；否则 CACHE_TIER_OVERRIDE env → RPC_CACHE_TIER → 默认 medium
      → 写 Cache-Control + CDN-Cache-Control（仅 isAllowedOrigin） + X-Cache-Tier
  15. ETag 计算（FNV-1a，等价实现）+ If-None-Match → 304
```

**POST → GET 兼容**（完整实现，对应 gateway.ts 444-459 行）：
- Content-Length < 1MB
- body JSON → url.searchParams（scalar 值 / array append）
- 再次 router.match

---

## Phase 4：Redis 客户端封装

**文件：`internal/cache/client.go`**
```go
// 全局单例，从环境变量初始化
var Client *redis.Client

func Init(addr, password string, db int)
func Get(ctx, key) (string, error)
func Set(ctx, key, value string, ttl time.Duration) error
func MGet(ctx, keys []string) ([]any, error)
```

**文件：`internal/cache/fetch.go`**

核心：`CachedFetchJSON[T]` 泛型函数
```go
const NegSentinel = "__WM_NEG__"
const NegTTL = 120 * time.Second

var inflight singleflight.Group

func CachedFetchJSON[T any](
    ctx context.Context,
    key string,
    ttl time.Duration,
    fetcher func(ctx context.Context) (*T, error),
) (*T, error)
// 流程：Redis GET → hit 返回 → 负值 sentinel 返回 nil → 
//        singleflight.Do(key, fetcher) → 结果为 nil 写 NEG_SENTINEL →
//        结果非 nil 写 JSON + TTL → 返回
```

**文件：`internal/cache/pipeline.go`**
```go
// 批量执行 Redis 命令（对应 runRedisPipeline）
func RunPipeline(ctx context.Context, cmds [][]any) ([]redis.Cmder, error)
// 用于 story tracking 的批量写入（HINCRBY/HSET/HSETNX/ZADD/SADD/EXPIRE）
```

---

## Phase 5：News Domain 接入（全部以 src/ 为准）

### 5.1 `internal/news/feeds.go`
**源**：`src/config/feeds.ts`（87KB 完整版，非 server/_feeds.ts 子集）

```go
type Feed struct {
    Name string
    URL  string            // 单语言；多语言 feed 用 URLByLang 字段
    URLByLang map[string]string  // lang code → URL
    Lang string            // "" = universal
}

// Feed 元数据（来自 src/config/feeds.ts）
type SourceType string     // wire | gov | intel | mainstream | market | tech | other
type PropagandaRisk string // low | medium | high

type SourceRiskProfile struct {
    Risk           PropagandaRisk
    StateAffiliated string
    KnownBiases    []string
    Note           string
}

var SOURCE_TYPES       map[string]SourceType        // 100+ 源分类
var SOURCE_PROPAGANDA_RISK map[string]SourceRiskProfile // RT/TASS/CGTN 等宣传风险
var VARIANT_FEEDS      map[string]map[string][]Feed  // 5 variants，500+ feeds
var INTEL_SOURCES      []Feed

func gn(query string) string  // Google News RSS URL 构建
func GetFeedURL(f Feed, lang string) string  // 多语言 URL 选择
```

### 5.2 `internal/news/classifier.go`
**源**：`src/services/threat-classifier.ts`（完整版，非 server/_classifier.ts 精简版）

关键差异（比 server/ 多出的内容）：
```go
// 1. CRITICAL_KEYWORDS 新增：
//    declares war / all-out war / full-scale war / massive strikes /
//    military strikes / retaliatory strikes / launches strikes /
//    attack on iran / strikes iran / bombs iran / iran retaliates 等 30+ 条

// 2. HIGH_KEYWORDS 新增：
//    airstrikes / strikes / ground offensive / bombardment / shelling /
//    killed in / strike on / strikes on / attack on / launches attacks /
//    explosions / ballistic missile / cruise missile 等 20+ 条

// 3. 新增 TRAILING_BOUNDARY_KEYWORDS（特殊正则：trailing word boundary）
var trailingBoundaryKW = map[string]bool{
    "attack iran": true, "strikes iran": true, "bombs iran": true, ...
}

// 4. 新增复合升级逻辑（核心新特性）
var escalationActions = regexp.MustCompile(`\b(attack|strikes|struck|bomb|bombed|shelling|missile|retaliates|offensive|invaded)\b`)
var escalationTargets = regexp.MustCompile(`\b(iran|tehran|russia|moscow|china|beijing|taiwan|nato|us forces|american forces)\b`)

func shouldEscalateToCritical(lower string, cat EventCategory) bool {
    // HIGH 军事/冲突 + 关键地缘政治目标 → CRITICAL（confidence 0.85）
    return (cat == "conflict" || cat == "military") &&
        escalationActions.MatchString(lower) && escalationTargets.MatchString(lower)
}

// 5. EXCLUSIONS 新增：'strikes deal', 'strikes agreement', 'strikes partnership'
// 6. SHORT_KEYWORDS 新增：'strikes'

// 7. 新增 AggregateThreats（跨 item 聚合）
func AggregateThreats(items []ThreatItem) ThreatClassification {
    // level = max; category = most frequent; confidence = tier加权平均
}
```

### 5.3 `internal/news/geo_hub.go`
**源**：`src/services/geo-hub-index.ts`（全新模块）

```go
type GeoHubLocation struct {
    ID      string   // "washington", "moscow", "gaza" 等
    Name    string
    Region  string
    Country string
    Lat, Lon float64
    Type    string   // "capital" | "conflict" | "strategic" | "organization"
    Tier    string   // "critical" | "major" | "notable"
    Keywords []string
}

// 130+ 个地理实体（首都/冲突区/战略要地/国际组织/美军基地）
var GEO_HUBS = []GeoHubLocation{ ... }

// 构建倒排索引（keyword → []hubID），初始化时 sync.Once 保证线程安全
type GeoHubIndex struct {
    hubs      map[string]*GeoHubLocation
    byKeyword map[string][]string
}

type GeoHubMatch struct {
    HubID          string
    Hub            *GeoHubLocation
    Confidence     float64  // 0.5~1.0，按 keyword 长度 + conflict/critical boost
    MatchedKeyword string
}

func InferGeoHubsFromTitle(title string) []GeoHubMatch {
    // 1. 标题转小写，tokenize（按空格/标点切词）
    // 2. 遍历 byKeyword 倒排索引，matchKeyword 检查是否命中
    // 3. confidence 计算：len≥10→0.9, len≥6→0.75, len≥4→0.6, else 0.5
    //                    conflict/strategic type → +0.1, critical tier → +0.1
    // 4. 按 confidence 降序排列，返回所有命中
}
```

### 5.4 `internal/news/source_tiers.go`
**源**：`server/_shared/source-tiers.ts`（server 与 src 共用，保持不变）

### 5.5 `internal/news/parser.go`
**源**：`server/worldmonitor/news/v1/list-feed-digest.ts` regex 解析（不用 `src/services/rss.ts` 的 DOMParser）

> **架构修正**：`src/services/rss.ts` 是纯客户端代码（DOMParser/IndexedDB/mlWorker）。
> Go 后端是 SERVER 层，应复制 `server/list-feed-digest.ts` 的 regex 解析方式。

Server 已有的 regex 解析模式：
```
<title><![CDATA[...]]> / <title>...</title>  → regex
<link>...</link>                              → regex
<pubDate>...</pubDate>                        → regex
<item>...</item> / <entry>...</entry> 切块   → regex
```

客户端专属功能的服务端等价：
| 客户端 (`src/services/rss.ts`) | Go 服务端等价 |
|---|---|
| `DOMParser` | regex（复制 server 已有实现） |
| `fetchWithProxy` | `http.Client` + Railway relay 降级 |
| `IndexedDB 持久化` | Redis（`cache.Set`） |
| `feedCache` (内存 Map) | `singleflight` 合并请求 |
| `feedFailures` 冷却 | 负值 sentinel + TTL |
| `classifyWithAI()` | HTTP 调用 `/api/intelligence/v1/classify-event` |
| `mlWorker`/`ingestHeadlines` | **跳过**（客户端专属） |

```go
// ParsedItem 结构（完整版，包含 geo 字段）
type ParsedItem struct {
    Source           string
    Title            string
    Link             string
    PublishedAt      int64    // epoch ms
    IsAlert          bool
    Level            ThreatLevel
    Category         string
    Confidence       float64
    ClassSource      string   // "keyword" | "llm"
    ImportanceScore  int
    CorroborationCount int
    TitleHash        string
    Lang             string
    // ★ Geo Hub 字段（来自 src/services/rss.ts）
    Lat              float64
    Lon              float64
    LocationName     string
    // ★ Feed 元数据
    PropagandaRisk   string   // "low" | "medium" | "high"
    SourceType       string   // "wire" | "intel" | ...
}
```

### 5.6 `internal/news/scoring.go`
无变化，权重与 server/ 一致（`src/` 无独立评分逻辑）。

### 5.7 `internal/news/story_tracking.go`
无变化，Redis pipeline 结构不变。

### 5.8 `internal/news/digest.go`
**新增流程步骤**（在 classifyByKeyword 后）：
```
c. 每条 item → InferGeoHubsFromTitle(title)
   → topGeo = matches[0]（最高 confidence）
   → item.Lat = topGeo.Hub.Lat
   → item.Lon = topGeo.Hub.Lon
   → item.LocationName = topGeo.Hub.Name
```

`enrichWithAICache` 升级版：
- 批量 MGET `classify:sebuf:v3:{hash}`（原有）
- **同时**支持调用 Intelligence service HTTP API 作为降级（`classifyWithAI` 等价）

### 5.9 `internal/news/handler.go`
无变化。

---

## Phase 6：Research Domain 接入

### 6.1 `internal/research/arxiv.go`
**源**：`server/research/v1/list-arxiv-papers.ts`（服务端 handler）+ `src/generated` 类型定义

```go
// seed-only：从 Redis 读取外部脚本预填的数据
// Redis key: "research:arxiv:v1:{category}::50"
// 参数 query 接收但不影响 key（与源实现一致）
func ListArxivPapers(ctx context.Context, req *http.Request) (any, error) {
    category := req.URL.Query().Get("category")  // default: "cs.AI"
    if category == "" { category = "cs.AI" }
    key := fmt.Sprintf("research:arxiv:v1:%s::50", category)
    return cache.GetCachedJSON[ArxivPapersResponse](ctx, key)
}
```

### 6.2 `internal/research/hackernews.go`
**源**：`server/research/v1/list-hackernews-items.ts`（服务端 handler）+ `src/generated` 类型定义

```go
// ALLOWED_FEEDS: top / new / best / ask / show / job（完整复制）
// Redis key: "research:hackernews:v1:{feedType}:30"
func ListHackernewsItems(ctx context.Context, req *http.Request) (any, error)
```

### 6.3 `internal/research/trending_repos.go`
**源**：`server/research/v1/list-trending-repos.ts`（服务端 handler）+ `src/generated` 类型定义

```go
// Redis key: "research:trending:v1:{language}:{period}:50"
// language default: "python", period default: "daily"
func ListTrendingRepos(ctx context.Context, req *http.Request) (any, error)
```

### 6.4 `internal/research/tech_events.go`
**源**：`server/research/v1/list-tech-events.ts`（完整实现，无浏览器专属 API）

```go
const (
    ICS_URL        = "https://www.techmeme.com/newsy_events.ics"
    DEV_EVENTS_RSS = "https://dev.events/rss.xml"
    REDIS_CACHE_KEY = "research:tech-events:v1"
    REDIS_CACHE_TTL = 21600  // 6h
    FETCH_TIMEOUT_MS = 8000
)

// 完整复制的函数：
func fetchTextWithRelay(ctx, url string) (string, error)
// 直连 → Railway relay 降级（与 list-tech-events.ts 完全一致）

func parseICS(icsText string) []TechEvent
// BEGIN:VEVENT 正则切块，提取 SUMMARY/LOCATION/DTSTART/DTEND/URL/UID
// 解析 date YYYYMMDD → "YYYY-MM-DD"，type: earnings/ipo/conference

func parseDevEventsRSS(rssText string) []TechEvent
// <item> 正则，CDATA title，date from description "on Month Day, Year"
// 跳过过期事件；Online → virtual coords

func normalizeLocation(loc string) *TechEventCoords
// 直接查表 → 去 state/country 后缀 → 模糊包含匹配（三级查找）

// 6 个 curated 事件（GITEX Global 2026 等，完整复制）
var CURATED_EVENTS []TechEvent

func ListTechEvents(ctx context.Context, req *http.Request) (any, error) {
    // 1. CachedFetchJSON("research:tech-events:v1", 21600s, fetchTechEvents)
    // 2. fetchTechEvents：并发抓 ICS+RSS → parseICS/parseDevEventsRSS
    //    → 合并 CURATED_EVENTS → dedup(title+year) → sort by startDate
    // 3. 按 query params 过滤：type / mappable / days / limit(clamp 1-200, default 50)
    // 4. externalSourcesFailed > 0 → warn 但不 fail（curated fallback）
}
```

### 6.5 `internal/research/city_coords.go`
**源**：`api/data/city-coords`（静态地理坐标查找表）

```go
type CityCoord struct {
    Lat     float64
    Lng     float64
    Country string
    Virtual bool
}

// 500+ 城市静态 map（完整复制）
var CITY_COORDS = map[string]CityCoord{ ... }
```

### 6.6 `internal/research/handler.go`

```go
func Register(r *registry.Registry) {
    r.Register(registry.DomainRegistration{
        Domain:  "research",
        Version: "v1",
        Routes: map[string]registry.Handler{
            "list-arxiv-papers":    ListArxivPapers,
            "list-trending-repos":  ListTrendingRepos,
            "list-hackernews-items": ListHackernewsItems,
            "list-tech-events":     ListTechEvents,
        },
    })
}
```

---

## Phase 7：主入口注册

**文件：`cmd/server/main.go`**
```go
func main() {
    cache.Init(...)          // 从 env 初始化 Redis
    
    r := registry.NewRegistry()
    
    // News domain — 以 src/ 为准（完整 feeds/classifier/geo）
    news.Register(r)

    // Research domain — 以 server/ handler 为准（seed-only + tech-events live-fetch）
    research.Register(r)

    // 后续 domain 在此追加：
    // aviation.Register(r)
    
    port := os.Getenv("SERVER_PORT")
    if port == "" { port = "8080" }
    http.ListenAndServe(":"+port, r)
}
```

---

## 实现顺序

| 步骤 | 文件 | 权威来源 | 备注 |
|------|------|---------|------|
| 1 | `go.mod` + `.env.example` | — | 骨架 |
| 2 | `cache/client.go` | — | Redis 初始化 |
| 3 | `cache/fetch.go` | — | CachedFetchJSON + singleflight |
| 4 | `cache/pipeline.go` | — | 批量命令 |
| 5 | `registry/registry.go` | — | 注册类型（无 CacheTier 字段） |
| 6 | `registry/router.go` | **server/router.ts** | static+dynamic 路由 |
| 7a | `registry/cors.go` | **server/cors.ts** | 可与 7b-7e 并行 |
| 7b | `registry/cache_tiers.go` | **server/gateway.ts** | 80+ 端点映射，可与 7a,7c-7e 并行 |
| 7c | `registry/premium_paths.go` | **src/shared/premium-paths.ts** | 36 条，可与 7a,7b,7d,7e 并行 |
| 7d | `registry/rate_limit.go` | **server/_shared/rate-limit.ts** | Redis 滑动窗口，可与 7a-7c,7e 并行 |
| 7e | `registry/entitlement.go` | **server/gateway.ts + _shared/entitlement-check.ts** | 双体系鉴权，可与 7a-7d 并行 |
| 8 | `registry/middleware.go` | **server/gateway.ts** | 依赖 7a-7e，15 步管道组装 |
| 9 | `news/classifier.go` | **src/services/threat-classifier.ts** | 可与 10-13 并行 |
| 10 | `news/source_tiers.go` | server/_shared/source-tiers.ts | 可与 9,11-13 并行 |
| 11 | `news/geo_hub.go` | **src/services/geo-hub-index.ts** | 可与 9,10,12,13 并行 |
| 12 | `news/parser.go` | **server/news/v1/list-feed-digest.ts** regex | 可与 9-11,13 并行 |
| 13 | `news/feeds.go` | **src/config/feeds.ts** | 可与 9-12 并行 |
| 14 | `news/scoring.go` | server/list-feed-digest.ts | 依赖 9,10 |
| 15 | `news/story_tracking.go` | server/list-feed-digest.ts | 依赖 4 |
| 16 | `news/digest.go` | server/ + src/ 融合 | 依赖 3,9-15 |
| 17 | `news/handler.go` | — | 依赖 16 |
| 18a | `research/city_coords.go` | **api/data/city-coords** | 静态数据，可与 18b-18d 并行 |
| 18b | `research/arxiv.go` | **server/research/v1/list-arxiv-papers.ts** | 可与 18a,18c,18d 并行 |
| 18c | `research/hackernews.go` | **server/research/v1/list-hackernews-items.ts** | 可与 18a,18b,18d 并行 |
| 18d | `research/trending_repos.go` | **server/research/v1/list-trending-repos.ts** | 可与 18a-18c 并行 |
| 19 | `research/tech_events.go` | server/research/v1/list-tech-events.ts | 依赖 18a（city_coords） |
| 20 | `research/handler.go` | — | 依赖 18b-19 |
| 21 | `cmd/server/main.go` | — | 依赖 8,17,20 |

---

## 测试端点（完成后）

**News domain**
```
GET /api/news/v1/list-feed-digest?variant=full&lang=en
GET /api/news/v1/list-feed-digest?variant=tech&lang=en
GET /api/news/v1/list-feed-digest?variant=finance&lang=en
```

**Research domain**
```
GET /api/research/v1/list-arxiv-papers?category=cs.AI&page_size=20
GET /api/research/v1/list-hackernews-items?feed_type=top&page_size=30
GET /api/research/v1/list-trending-repos?language=python&period=daily&page_size=20
GET /api/research/v1/list-tech-events?type=conference&mappable=true&days=90&limit=50
```

**中间件验证**
```
OPTIONS /api/news/v1/list-feed-digest          → 204 + CORS headers
GET /api/news/v1/list-feed-digest              → 401（无 API Key）
GET /api/news/v1/list-feed-digest?_debug=1     → 200 + X-Cache-Tier: slow
GET /api/aviation/v1/track-aircraft?callsign=X → 200 + Cache-Control: no-store
```

期望响应：
```json
{
  "categories": {
    "politics": { "items": [...] },
    "tech":     { "items": [...] }
  },
  "feedStatuses": { "BBC World": "ok", "RT": "timeout" },
  "generatedAt": "2026-04-22T13:00:00Z"
}
```
