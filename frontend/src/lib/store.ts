import { create } from "zustand";

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
}));
