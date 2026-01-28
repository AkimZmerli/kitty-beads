import type { QueryClient } from "@tanstack/react-query";

export interface TerminalTab {
  id: string;
  title: string;
  createdAt: number;
}

export interface TerminalState {
  isOpen: boolean;
  tabs: TerminalTab[];
  activeTabId: string | null;
  panelHeight: number;
  isFullscreen: boolean;
}

export interface TerminalContextValue extends TerminalState {
  // Panel actions
  openPanel: () => void;
  closePanel: () => void;
  togglePanel: () => void;

  // Tab actions
  addTab: () => string;
  closeTab: (id: string) => void;
  setActiveTab: (id: string) => void;

  // Panel sizing
  setPanelHeight: (height: number) => void;
  toggleFullscreen: () => void;

  // Command execution (auto-opens terminal)
  executeCommand: (command: string) => void;

  // For sending commands to active terminal
  activeTerminalRef: React.MutableRefObject<{
    sendCommand: (cmd: string) => void;
  } | null>;

  // React Query client for pattern detection
  queryClient: QueryClient;

  // Height constraints
  minHeight: number;
  maxHeightPercent: number;
}
