import { TERMINAL_STORAGE_KEY, DEFAULT_HEIGHT, MIN_HEIGHT } from "./constants";

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
