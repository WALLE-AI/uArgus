import { create } from "zustand";
import type { SemanticFeedItem, InsightsResult, TopicSuggestion } from "./types";

interface AppState {
  sidebarCollapsed: boolean;
  toggleSidebar: () => void;
  setSidebarCollapsed: (v: boolean) => void;

  assistantOpen: boolean;
  toggleAssistant: () => void;
  setAssistantOpen: (v: boolean) => void;

  commandPaletteOpen: boolean;
  setCommandPaletteOpen: (v: boolean) => void;

  theme: "light" | "dark";
  setTheme: (v: "light" | "dark") => void;
  toggleTheme: () => void;

  // ── Semantic feed state ──
  feedItems: SemanticFeedItem[];
  setFeedItems: (items: SemanticFeedItem[]) => void;
  feedLoading: boolean;
  setFeedLoading: (v: boolean) => void;
  feedError: string | null;
  setFeedError: (v: string | null) => void;
  feedFromApi: boolean;
  setFeedFromApi: (v: boolean) => void;
  feedVariant: string;
  setFeedVariant: (v: string) => void;

  // ── Insights state ──
  insights: InsightsResult | null;
  setInsights: (v: InsightsResult | null) => void;
  insightsLoading: boolean;
  setInsightsLoading: (v: boolean) => void;

  // ── Topic suggestions state ──
  topicSuggestions: TopicSuggestion[];
  setTopicSuggestions: (v: TopicSuggestion[]) => void;
}

export const useAppStore = create<AppState>((set) => ({
  sidebarCollapsed: false,
  toggleSidebar: () => set((s) => ({ sidebarCollapsed: !s.sidebarCollapsed })),
  setSidebarCollapsed: (v) => set({ sidebarCollapsed: v }),

  assistantOpen: true,
  toggleAssistant: () => set((s) => ({ assistantOpen: !s.assistantOpen })),
  setAssistantOpen: (v) => set({ assistantOpen: v }),

  commandPaletteOpen: false,
  setCommandPaletteOpen: (v) => set({ commandPaletteOpen: v }),

  theme: "light",
  setTheme: (v) => set({ theme: v }),
  toggleTheme: () => set((s) => ({ theme: s.theme === "light" ? "dark" : "light" })),

  // ── Semantic feed state ──
  feedItems: [],
  setFeedItems: (items) => set({ feedItems: items }),
  feedLoading: false,
  setFeedLoading: (v) => set({ feedLoading: v }),
  feedError: null,
  setFeedError: (v) => set({ feedError: v }),
  feedFromApi: false,
  setFeedFromApi: (v) => set({ feedFromApi: v }),
  feedVariant: "full",
  setFeedVariant: (v) => set({ feedVariant: v }),

  // ── Insights state ──
  insights: null,
  setInsights: (v) => set({ insights: v }),
  insightsLoading: false,
  setInsightsLoading: (v) => set({ insightsLoading: v }),

  // ── Topic suggestions state ──
  topicSuggestions: [],
  setTopicSuggestions: (v) => set({ topicSuggestions: v }),
}));
