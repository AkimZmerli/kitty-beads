import { TERMINAL_STORAGE_KEY, DEFAULT_HEIGHT, MIN_HEIGHT } from "./constants";
import type { TerminalTab } from "./types";

const TERMINAL_TABS_KEY = "terminal-tabs";
const TERMINAL_ACTIVE_TAB_KEY = "terminal-active-tab";

export function generateTerminalId(): string {
  return `term-${Date.now()}-${Math.random().toString(36).slice(2, 9)}`;
}

export function getSavedHeight(): number {
  try {
    const saved = localStorage.getItem(TERMINAL_STORAGE_KEY);
    if (saved) {
      const height = parseInt(saved, 10);
      if (!isNaN(height) && height >= MIN_HEIGHT) {
        return height;
      }
    }
  } catch {
    // Ignore localStorage errors
  }
  return DEFAULT_HEIGHT;
}

export function saveHeight(height: number): void {
  try {
    localStorage.setItem(TERMINAL_STORAGE_KEY, height.toString());
  } catch {
    // Ignore localStorage errors
  }
}

/** Persisted shape — only stable fields, not runtime content/history */
interface PersistedTab {
  id: string;
  title: string;
}

export function getSavedTabs(): TerminalTab[] | null {
  try {
    const saved = localStorage.getItem(TERMINAL_TABS_KEY);
    if (saved) {
      const parsed: PersistedTab[] = JSON.parse(saved);
      if (Array.isArray(parsed) && parsed.length > 0) {
        return parsed.map((t) => ({
          id: t.id,
          title: t.title,
          createdAt: Date.now(),
        }));
      }
    }
  } catch {
    // Ignore localStorage errors
  }
  return null;
}

export function saveTabs(tabs: TerminalTab[]): void {
  try {
    const persisted: PersistedTab[] = tabs.map((t) => ({
      id: t.id,
      title: t.title,
    }));
    localStorage.setItem(TERMINAL_TABS_KEY, JSON.stringify(persisted));
  } catch {
    // Ignore localStorage errors
  }
}

export function getSavedActiveTabId(): string | null {
  try {
    return localStorage.getItem(TERMINAL_ACTIVE_TAB_KEY);
  } catch {
    return null;
  }
}

export function saveActiveTabId(id: string | null): void {
  try {
    if (id === null) {
      localStorage.removeItem(TERMINAL_ACTIVE_TAB_KEY);
    } else {
      localStorage.setItem(TERMINAL_ACTIVE_TAB_KEY, id);
    }
  } catch {
    // Ignore localStorage errors
  }
}
