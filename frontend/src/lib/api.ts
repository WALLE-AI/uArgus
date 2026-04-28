import type { SemanticFeedItem, InsightsResult, TopicSuggestion, SeedEnvelope, EventAnalysis } from "./types";
import { MOCK_DIGEST, MOCK_INSIGHTS, MOCK_TOPIC_SUGGESTIONS, MOCK_EVENT_ANALYSES } from "./mock-data";

// ── Configuration ────────────────────────────────────────────────────────

const API_BASE = process.env.NEXT_PUBLIC_API_URL ?? "";

/** True when a backend API endpoint is configured */
export function isApiConfigured(): boolean {
  return API_BASE.length > 0;
}

// ── Generic fetch helper ─────────────────────────────────────────────────

async function apiFetch<T>(path: string): Promise<T> {
  if (!API_BASE) throw new Error("API not configured");

  const res = await fetch(`${API_BASE}${path}`, {
    headers: { "Content-Type": "application/json" },
    next: { revalidate: 60 },
  });

  if (!res.ok) {
    throw new Error(`API ${res.status}: ${res.statusText}`);
  }

  const json = await res.json();

  // unwrap SeedEnvelope if present
  if (json && typeof json === "object" && "_seed" in json && "data" in json) {
    return (json as SeedEnvelope<T>).data;
  }

  return json as T;
}

// ── Feed / Digest ────────────────────────────────────────────────────────

export interface FetchDigestOptions {
  variant?: string; // "full" | "tech" | "finance" | "geopolitics" | "science"
  lang?: string;    // "en" | "zh"
}

/**
 * Fetch processed digest from semantic layer.
 * Falls back to mock data when API is not available.
 */
export async function fetchDigest(
  opts: FetchDigestOptions = {}
): Promise<{ items: SemanticFeedItem[]; fromApi: boolean }> {
  const { variant = "full", lang = "en" } = opts;

  try {
    const items = await apiFetch<SemanticFeedItem[]>(
      `/api/news/v1/list-feed-digest?variant=${variant}&lang=${lang}`
    );
    return { items, fromApi: true };
  } catch {
    // fallback to mock
    return { items: MOCK_DIGEST, fromApi: false };
  }
}

// ── Insights ─────────────────────────────────────────────────────────────

/**
 * Fetch AI-generated insights/summary from semantic layer.
 * Falls back to mock data when API is not available.
 */
export async function fetchInsights(): Promise<{
  data: InsightsResult;
  fromApi: boolean;
}> {
  try {
    const data = await apiFetch<InsightsResult>("/api/news/v1/insights");
    return { data, fromApi: true };
  } catch {
    return { data: MOCK_INSIGHTS, fromApi: false };
  }
}

// ── Topic Suggestions ────────────────────────────────────────────────────

/**
 * Fetch AI topic suggestions derived from semantic analysis.
 * Falls back to mock data when API is not available.
 */
export async function fetchTopicSuggestions(): Promise<{
  data: TopicSuggestion[];
  fromApi: boolean;
}> {
  try {
    const data = await apiFetch<TopicSuggestion[]>(
      "/api/news/v1/topic-suggestions"
    );
    return { data, fromApi: true };
  } catch {
    return { data: MOCK_TOPIC_SUGGESTIONS, fromApi: false };
  }
}

// ── Event Analysis ─────────────────────────────────────────────────────────────

/**
 * Fetch per-item event analysis.
 * Falls back to mock data when API is not available.
 * Returns null if no analysis exists for the given hash.
 */
export async function fetchEventAnalysis(
  hash: string
): Promise<{ data: EventAnalysis | null; fromApi: boolean }> {
  try {
    const data = await apiFetch<EventAnalysis>(
      `/api/news/v1/event-analysis?hash=${hash}`
    );
    return { data, fromApi: true };
  } catch {
    const mock = MOCK_EVENT_ANALYSES[hash] ?? null;
    return { data: mock, fromApi: false };
  }
}
