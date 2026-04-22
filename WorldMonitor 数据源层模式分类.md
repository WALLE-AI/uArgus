# WorldMonitor 数据源层模式分类

本方案把 `opensource/worldmonitor` 里约 180 条 RPC 归为 **8 种数据获取/分发模式**,每种配一个单文件级的最简例子、上游与存储画像、触发频率与性能约束。

---

## 全景归类表(按 RPC 数量近似估算)

| # | 模式 | 占比 | 典型标志 | 是否碰上游 | 缓存位置 |
|---|---|---|---|---|---|
| 1 | **纯 seed 读**(read-only) | ~55% | `getCachedJson(SEED_KEY)` | ❌ | Redis(由 cron 填) |
| 2 | **参数化 seed 读** | ~15% | `${PREFIX}:${lang}:${period}` | ❌ | Redis(变体预生成) |
| 3 | **Bootstrap 聚合读** | ~5% | 读 `*-bootstrap:v1` 多段 | ❌ | Redis |
| 4 | **在线抓取 + TTL 缓存** | ~10% | `cachedFetchJson(paramKey, ttl, fetcher)` | ✅ 经 relay 或直连 | Redis + inflight Map |
| 5 | **多源合并/派生** | ~5% | 同时读多 key + 跨 RPC 调用 | 少量 | Redis(+ 再写派生 key) |
| 6 | **LLM 按需生成** | ~3% | `callLlm(...)` + `cachedFetchJsonWithMeta` | ✅ OpenAI/Groq/etc | Redis(按内容哈希) |
| 7 | **异步作业队列** | <1% | `RPUSH scenario-queue:pending` | ❌(enqueue) | Redis List + worker |
| 8 | **空间查询 / GEO** | ~2% | `geoSearchByBox` | ❌ | Redis GEO key |
| — | **Relay 保活循环** | 基础设施 | 长连续进程,不是 RPC | ✅ | Redis |

---

## 1. 纯 seed 读(最常见)

**特征**:RPC 只从 Redis 取一把预先由 cron seeder 填好的数据,`slice(0, pageSize)` 后返回。**请求路径零外部 I/O**。

**核心代码**(`@d:\LLM\project\WALL-AI\uArgus\opensource\worldmonitor\server\worldmonitor\seismology\v1\list-earthquakes.ts:22-31`):
```ts
const seedData = await getCachedJson('seismology:earthquakes:v1', true) as EarthquakeCache | null;
return { earthquakes: (seedData?.earthquakes ?? []).slice(0, pageSize), pagination: undefined };
```

**画像**:
- 上游:USGS / ACLED / NASA / BIS / IMF / CoinGecko …(不在请求路径,在 `scripts/seed-*.mjs` cron)
- 存储:`<domain>:<resource>:v<N>`,TTL 1h~1 周,envelope `{_seed, data}` 自动剥离
- 触发频率:5 min~月(随上游变动性)
- 性能要求:Edge 1.5s Redis 超时内必须返回,通常 < 50ms

**其他例子**:`unrest/list-unrest-events`(带客户端 filter)、`research/list-aviation-news` 等。

---

## 2. 参数化 seed 读

**特征**:seeder 预生成多个变体(语言/周期/国家/地区),RPC 按参数组合 key 读。

**核心代码**(`@d:\LLM\project\WALL-AI\uArgus\opensource\worldmonitor\server\worldmonitor\research\v1\list-trending-repos.ts:18-28`):
```ts
const seedKey = `${SEED_KEY_PREFIX}:${language}:${period}:50`;
const result = await getCachedJson(seedKey, true) as ListTrendingReposResponse | null;
return { repos: (result?.repos ?? []).slice(0, pageSize) };
```

**画像**:
- 上游:GitHub trending / OpenAQ / Polymarket …
- 存储:`<domain>:<resource>:<p1>:<p2>:v<N>`,参数笛卡尔积必须可控(通常 ≤ 几十个变体),否则落入模式 4
- 性能要求:同 #1,但要求 seeder 覆盖率 ≥ 客户端常见参数组合,否则变成 miss

**其他例子**:`prediction/list-prediction-markets`(按 `category` 路由到 `bootstrap.geopolitical/tech/finance` 子段)。

---

## 3. Bootstrap 聚合读

**特征**:不是 RPC,是 `/api/bootstrap?tier=fast|slow` 一次性 Upstash pipeline GET 数十~百个 canonical key,给前端首屏灌数据。

**核心代码**(`@d:\LLM\project\WALL-AI\uArgus\opensource\worldmonitor\api\bootstrap.js:188-207`):
```js
const pipeline = keys.map((k) => ['GET', k]);
const data = await redisPipeline(pipeline, 3000);
// ... unwrapEnvelope 剥 _seed,NEG 哨兵跳过,写入 result
```

**画像**:
- fast tier(60s 浏览器 / 300s CF / 600s CDN)≈ 30 key
- slow tier(300s/1800s/3600s)≈ 80 key
- 一个 RPC 如果想加入 bootstrap,只需在对应 seeder 里同步 `aviation:delays-bootstrap:v1` 这种镜像键,再在 `BOOTSTRAP_CACHE_KEYS` 注册名字

---

## 4. 在线抓取 + TTL 缓存(关键路径上有上游)

**特征**:请求路径上调用 `cachedFetchJson(key, ttl, fetcher)`,命中走 Redis,miss 就经 Railway relay 或直连上游。自动带 **inflight 合并**(同 key 并发只发一次 upstream)和 **负向哨兵**(`__WM_NEG__` 120s)。

**核心代码**(`@d:\LLM\project\WALL-AI\uArgus\opensource\worldmonitor\server\worldmonitor\aviation\v1\get-flight-status.ts:70-95`):
```ts
const result = await cachedFetchJson(cacheKey, CACHE_TTL, async () => {
  const resp = await fetch(`${relayBase}/aviationstack?${params}`, { headers: getRelayHeaders() });
  if (!resp.ok) throw new Error(`HTTP ${resp.status}`);
  return { flights: normalize(await resp.json()), source: 'aviationstack' };
});
```

**画像**:
- 上游:参数高基数的商用/限流 API(AviationStack、SerpAPI、Finnhub、Yahoo)或公开直连(OpenSky、Open-Meteo)
- 存储:`<domain>:<resource>:<param1>:<param2>:v<N>`;TTL 60s~15min
- 性能:miss 路径受上游影响(500ms~5s);`AbortSignal.timeout(15_000)` 是默认上限
- 风险:上游 rate-limit;`ais-relay` 作为**共享上游代理**承担 key 管理 + 节流 + Decodo IP 轮换

**其他例子**:`maritime/get-vessel-snapshot`(本地内存 + relay,`@d:\LLM\project\WALL-AI\uArgus\opensource\worldmonitor\server\worldmonitor\maritime\v1\get-vessel-snapshot.ts:46-65`)、`infrastructure/reverse-geocode`、大部分 `search-*` 接口。

---

## 5. 多源合并 / 派生

**特征**:同时读多个 seed/cache key,或者跨 RPC 调其他 handler,做合并/打分/映射后返回,**有时顺便再写一个派生 key**。

**核心代码思路**(`@d:\LLM\project\WALL-AI\uArgus\opensource\worldmonitor\server\worldmonitor\supply-chain\v1\get-chokepoint-status.ts:22-31`,示意):
```ts
const [summaries, flows, warnings, vessels] = await Promise.all([
  getCachedJson(TRANSIT_SUMMARIES_KEY),
  getCachedJson(FLOWS_KEY),
  listNavigationalWarnings(ctx, ...),   // 跨 domain 直接调 handler 函数
  getVesselSnapshot(ctx, ...),
]);
// computeDisruptionScore / scoreToStatus 合并打分
await setCachedJson(CHOKEPOINT_STATUS_KEY, result, 300);  // 派生 key
```

**画像**:
- 上游:一般不直接打外部,而是聚合已有 seed + 跨 RPC 调用
- 存储:派生结果写回 Redis(再供 bootstrap / 下游 RPC 使用)
- 性能:I/O 并发是关键,全 `Promise.all`;edge 时间受**最慢的那个 Redis 取**决定

**其他例子**:`aviation/list-airport-delays`(FAA+INTL+NOTAM 三源合并)、`intelligence/get-risk-scores`、`resilience/get-resilience-ranking`。

---

## 6. LLM 按需生成

**特征**:RPC 在请求路径上 `callLlm`,把输入规约后作缓存 key(哈希),命中 Redis 就不再打 LLM。**通常 premium-gated**。

**核心代码**(`@d:\LLM\project\WALL-AI\uArgus\opensource\worldmonitor\server\worldmonitor\news\v1\summarize-article.ts:38-62`):
```ts
const isPremium = await isCallerPremium(ctx.request);
const cacheKey = getCacheKey({ headlines, provider, mode, geoContext, variant, lang });
// cachedFetchJsonWithMeta 同样带 inflight 合并 + 负向哨兵
const { data } = await cachedFetchJsonWithMeta(cacheKey, CACHE_TTL_SECONDS, async () => {
  return callLlm({ prompt, provider, ... });
});
```

**画像**:
- 上游:OpenAI / Groq / OpenRouter / Anthropic / Ollama(多提供商自动 fallback,见 `_shared/llm.ts` + `llm-health.ts`)
- 存储:内容哈希 key,TTL 小时~天;premium 用户的输出存久些
- 性能:首次 1~10s;用 `AbortSignal` 超时;同内容并发合并
- 安全:严格 prompt 注入过滤(`llm-sanitize.js`)、系统消息 append 只 premium 可用

**其他例子**:`market/analyze-stock`(AI overlay)、`intelligence/chat-analyst`、`intelligence/classify-event`。

---

## 7. 异步作业队列

**特征**:长耗时任务(> Edge 函数上限)。RPC 只负责 `RPUSH` 入队,立即返回 `{jobId, status: "pending"}`;另一个 RPC 轮询状态;Railway worker(`scripts/scenario-worker.mjs`)消费队列、写回状态。

**核心代码**(`@d:\LLM\project\WALL-AI\uArgus\opensource\worldmonitor\server\worldmonitor\scenario\v1\run-scenario.ts:48-67`):
```ts
const [depthEntry] = await runRedisPipeline([['LLEN', QUEUE_KEY]], true);   // 背压
if (depth > MAX_QUEUE_DEPTH) throw new ApiError(429, ...);
const [pushEntry] = await runRedisPipeline([['RPUSH', QUEUE_KEY, payload]], true);
return { jobId, status: 'pending', statusUrl: `/api/.../get-scenario-status?jobId=${jobId}` };
```

**画像**:
- 存储:`scenario-queue:pending`(List)+ `scenario:<jobId>`(Hash/String)+ `scenario-package:latest`(最新产物)
- worker 心跳与锁:`SET NX PX`
- 典型场景:情景推演、深度预测(`process-deep-forecast-tasks.mjs`)

---

## 8. 空间查询 / GEO

**特征**:用 Redis GEO 指令(`GEOADD` / `GEOSEARCH`)做边界盒或半径查询。适合船舶、飞机、基地等点位数据。

**核心代码**(relay 侧 `GEOADD`;RPC 侧):
```ts
const members = await geoSearchByBox('ais:positions:v1', { lon, lat, widthKm, heightKm });
// members 是 id 数组,再去 HGETALL / MGET 取属性
```

**画像**:
- 上游:AIS WebSocket(ais-relay 持续 `GEOADD`)、military flights、OpenSky
- 存储:GEO key + 伴生 Hash / String 属性
- 性能:box 查询 O(log N + M),配合 Hash pipeline 可一次返回上千点;注意 box 不要越洋(穿越经度 180° 线要拆两 box)

**其他例子**:`military/list-military-flights`、`military/list-military-bases`(自然也会退化为 pure seed read,如果基数不大)。

---

## 附 A:基础设施 — Relay 保活循环(不是 RPC)

`scripts/ais-relay.cjs` 在 Railway 上常驻,承担三类不可 cron 的持续任务:
- **AIS WebSocket 接收** → 写 GEO key(模式 8 的数据源)
- **Warm-ping 循环**(cable-health、risk-scores、chokepoint-transits 等) — 每 N 分钟算一次并写 Redis,保证 `seed-meta:<x>` 不过期
- **HTTP 上游代理**(`/aviationstack`、`/googleflights`) — 给模式 4 的 Edge RPC 使用,集中持 key 与节流

Go 迁移时:这部分是**完全独立的进程**,和 gateway 解耦,先迁或后迁都不影响 RPC 路径。

---

## 附 B:模式到 RPC 的快速查询表(样本)

| 模式 | 样例 RPC(单文件就能看懂) |
|---|---|
| 1 | `seismology/list-earthquakes`、`unrest/list-unrest-events`、`wildfire/list-wildfires` |
| 2 | `research/list-trending-repos`、`prediction/list-prediction-markets` |
| 3 | `/api/bootstrap`(非 RPC,是 aggregator) |
| 4 | `aviation/get-flight-status`、`maritime/get-vessel-snapshot`、`infrastructure/reverse-geocode` |
| 5 | `supply-chain/get-chokepoint-status`、`intelligence/get-risk-scores` |
| 6 | `news/summarize-article`、`market/analyze-stock`、`intelligence/chat-analyst` |
| 7 | `scenario/run-scenario` + `scenario/get-scenario-status` |
| 8 | `maritime/get-vessel-snapshot` 的候选层、`military/list-military-flights` |

---

## 附 C:Go 迁移时的关键差异(一句话各自)

1. **纯 seed 读**:最简,Go `rdb.Get` + `json.Unmarshal` 即可,首选先迁。
2. **参数化 seed 读**:同上,但注意 key 构造逻辑必须与 seeder 端字节对齐。
3. **Bootstrap 聚合**:Go 用 `MGet`/`Pipeline` 一次取完,注意 `__WM_NEG__` 哨兵与 envelope 剥离。
4. **在线抓取**:用 `singleflight.Group` 替代 JS inflight Map;上游 client 带 proxy/timeout/retry。
5. **多源合并**:Go `errgroup.Go` 扇出扇入,与 TS `Promise.all` 语义一致。
6. **LLM**:要么 Go 接各 LLM REST(原生 SDK 少),要么先保留 TS 端跑这类;把 prompt sanitize 一并迁。
7. **作业队列**:Go 用 `LPUSH/BRPOP`(带阻塞)或保留 Redis List polling,worker 单独二进制。
8. **GEO**:go-redis 原生支持 `GeoSearch`,无障碍。

建议迁移顺序:**1 → 2 → 8 → 4 → 5 → 3 → 7 → 6**(先把 I/O 形态简单的灭掉,复杂的 LLM/队列放最后)。

---

_方案完_。如需我继续把每种模式各选一个具体 domain 出详细 Go 蓝图(如先从 **seismology 模式 1** 打样),告诉我就做。
