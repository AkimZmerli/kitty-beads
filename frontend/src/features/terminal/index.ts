// Context & Hook
export { TerminalProvider, useTerminal } from "./context";

// Types
export type { TerminalTab, TerminalState, TerminalContextValue } from "./types";

// Constants
export {
  TERMINAL_STORAGE_KEY,
  DEFAULT_HEIGHT,
  MIN_HEIGHT,
  MAX_HEIGHT_PERCENT,
} from "./constants";

// Components
export { TerminalPanel } from "./components/TerminalPanel";
export { TerminalTabs } from "./components/TerminalTabs";
export { TerminalInstance } from "./components/TerminalInstance";
export { ResizeHandle } from "./components/ResizeHandle";
