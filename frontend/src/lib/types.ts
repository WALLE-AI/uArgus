// ── Types aligned with backend/monitor semantic layer output ─────────────

/**
 * Aligned with backend/monitor/internal/news/parser.go ParsedItem.
 * This is the primary data structure for items that have passed through
 * the semantic pipeline (classify → score → track → enrich).
 */
export interface SemanticFeedItem {
  title: string;
  link: string;
  description?: string;
  pubDate: string;
  source: string;
  hash: string;

  // ── semantic/classify ──
  severity: number;       // 0-7, keyword classification severity
  categories: string[];   // e.g. ["conflict", "technology"]

  // ── semantic/scoring ──
  importanceScore: number; // weighted score (0-1), severity*0.55 + tier*0.2 + corr*0.15 + recency*0.1

  // ── semantic/tracking ──
  stage?: EventStage;          // story lifecycle stage
  corroboration: number;       // cross-source mention count

  // ── semantic/tiers ──
  tier: number;           // 1-4 source credibility (1=highest)
}

/** Story lifecycle stages from semantic/tracking */
export type EventStage = "BREAKING" | "DEVELOPING" | "SUSTAINED" | "FADING";

/** Severity label mapping (0-7 → display) */
export type SeverityLevel = "critical" | "high" | "medium" | "low" | "info";

export function severityToLevel(severity: number): SeverityLevel {
  if (severity >= 6) return "critical";
  if (severity >= 4) return "high";
  if (severity >= 2) return "medium";
  if (severity >= 1) return "low";
  return "info";
}

export function severityToLabel(severity: number): string {
  const map: Record<SeverityLevel, string> = {
    critical: "紧急",
    high: "高",
    medium: "中",
    low: "低",
    info: "信息",
  };
  return map[severityToLevel(severity)];
}

export function stageToLabel(stage?: EventStage): string {
  if (!stage) return "";
  const map: Record<EventStage, string> = {
    BREAKING: "突发",
    DEVELOPING: "发展中",
    SUSTAINED: "持续",
    FADING: "消退",
  };
  return map[stage];
}

export function tierToLabel(tier: number): string {
  const map: Record<number, string> = { 1: "核心", 2: "主流", 3: "一般", 4: "边缘" };
  return map[tier] ?? "未知";
}

/**
 * Aligned with backend/monitor/internal/news/insights_source.go output.
 */
export interface InsightsResult {
  summary: string;
  headlines: string[];
  generatedAt: string;
}

/**
 * AI-generated topic suggestion (derived from semantic analysis).
 */
export interface TopicSuggestion {
  title: string;
  reason: string;
  score: number;
  categories: string[];
  relatedCount: number;
}

/**
 * Aligned with backend/monitor/internal/seed/envelope.go SeedMeta.
 */
export interface SeedMeta {
  fetchedAt: number;
  recordCount: number;
  sourceVersion: string;
  schemaVersion: number;
  state: "OK" | "OK_ZERO" | "ERROR";
  failedDatasets?: string[];
  errorReason?: string;
  groupId?: string;
}

/**
 * Aligned with backend/monitor/internal/seed/envelope.go SeedEnvelope.
 */
export interface SeedEnvelope<T> {
  _seed?: SeedMeta;
  data: T;
}

/**
 * Digest metadata returned alongside feed items.
 */
export interface DigestMeta {
  variant: string;
  lang: string;
  generatedAt: string;
  totalItems: number;
}

// ── Per-item Event Analysis ──────────────────────────────────────────────

/** Named entity extracted from an article (aligned with agents.BuildNERPrompt output). */
export interface Entity {
  text: string;
  type: "PER" | "ORG" | "LOC" | "MISC";
  confidence: number;
}

/** Sentiment classification (aligned with agents.BuildSentimentPrompt output). */
export interface SentimentResult {
  label: "positive" | "negative" | "neutral";
  score: number; // 0-1
}

/** Impact assessment scope. */
export type ImpactScope = "global" | "regional" | "national" | "local";

/** Structured impact assessment for an event. */
export interface ImpactAssessment {
  scope: ImpactScope;
  sectors: string[];     // affected sectors/industries
  disruption: number;    // 0-1 disruption score (aligned with scoring/disruption.go)
  summary: string;       // 1-2 sentence impact summary
}

/** Timeline milestone for an event. */
export interface TimelineMilestone {
  date: string;  // ISO date or descriptive label
  label: string; // event description
}

/**
 * Per-item event analysis result.
 * Frontend mock-first; backend will generate via LLM agents + cache.
 */
export interface EventAnalysis {
  entities: Entity[];
  sentiment: SentimentResult;
  impact: ImpactAssessment;
  timeline: TimelineMilestone[];
  actions: string[];
  generatedAt: string;
}

export function sentimentToLabel(label: SentimentResult["label"]): string {
  const map: Record<SentimentResult["label"], string> = {
    positive: "积极",
    negative: "消极",
    neutral: "中性",
  };
  return map[label];
}

export function scopeToLabel(scope: ImpactScope): string {
  const map: Record<ImpactScope, string> = {
    global: "全球",
    regional: "区域",
    national: "国家",
    local: "地方",
  };
  return map[scope];
}
