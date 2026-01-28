import {
  createContext,
  useContext,
  useState,
  useCallback,
  useEffect,
  useRef,
  type ReactNode,
} from "react";
import { useQueryClient } from "@tanstack/react-query";
import type { TerminalState, TerminalContextValue, TerminalTab } from "./types";
import { MIN_HEIGHT, MAX_HEIGHT_PERCENT } from "./constants";
import { generateTerminalId, getSavedHeight, saveHeight } from "./utils";

const TerminalContext = createContext<TerminalContextValue | null>(null);

function createNewTab(tabNumber: number): TerminalTab {
  return {
    id: generateTerminalId(),
    title: `Terminal ${tabNumber}`,
    createdAt: Date.now(),
  };
}

export function TerminalProvider({ children }: { children: ReactNode }) {
  const queryClient = useQueryClient();
  const activeTerminalRef = useRef<{
    sendCommand: (cmd: string) => void;
  } | null>(null);

  const [state, setState] = useState<TerminalState>({
    isOpen: false,
    tabs: [],
    activeTabId: null,
    panelHeight: getSavedHeight(),
    isFullscreen: false,
  });

  // Persist panel height to localStorage
  useEffect(() => {
    saveHeight(state.panelHeight);
  }, [state.panelHeight]);

  const openPanel = useCallback(() => {
    setState((prev) => {
      const newState = { ...prev, isOpen: true };
      if (prev.tabs.length === 0) {
        const newTab = createNewTab(1);
        newState.tabs = [newTab];
        newState.activeTabId = newTab.id;
      }
      return newState;
    });
  }, []);

  const closePanel = useCallback(() => {
    setState((prev) => ({ ...prev, isOpen: false, isFullscreen: false }));
  }, []);

  const togglePanel = useCallback(() => {
    setState((prev) => {
      if (prev.isOpen) {
        return { ...prev, isOpen: false, isFullscreen: false };
      }
      const newState = { ...prev, isOpen: true };
      if (prev.tabs.length === 0) {
        const newTab = createNewTab(1);
        newState.tabs = [newTab];
        newState.activeTabId = newTab.id;
      }
      return newState;
    });
  }, []);

  // Keyboard shortcut: Ctrl+` to toggle terminal
  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if (e.ctrlKey && e.key === "`") {
        e.preventDefault();
        togglePanel();
      }
    };

    window.addEventListener("keydown", handleKeyDown);
    return () => window.removeEventListener("keydown", handleKeyDown);
  }, [togglePanel]);

  const addTab = useCallback((): string => {
    const newTab = createNewTab(state.tabs.length + 1);

    setState((prev) => ({
      ...prev,
      tabs: [...prev.tabs, newTab],
      activeTabId: newTab.id,
      isOpen: true,
    }));

    return newTab.id;
  }, [state.tabs.length]);

  const closeTab = useCallback((id: string) => {
    setState((prev) => {
      const newTabs = prev.tabs.filter((t) => t.id !== id);
      let newActiveId = prev.activeTabId;

      if (prev.activeTabId === id) {
        const closingIndex = prev.tabs.findIndex((t) => t.id === id);
        if (newTabs.length > 0) {
          newActiveId =
            newTabs[Math.max(0, closingIndex - 1)]?.id ||
            newTabs[0]?.id ||
            null;
        } else {
          newActiveId = null;
        }
      }

      return {
        ...prev,
        tabs: newTabs,
        activeTabId: newActiveId,
        isOpen: newTabs.length > 0 ? prev.isOpen : false,
      };
    });
  }, []);

  const setActiveTab = useCallback((id: string) => {
    setState((prev) => ({ ...prev, activeTabId: id }));
  }, []);

  const setPanelHeight = useCallback((height: number) => {
    const maxHeight = window.innerHeight * (MAX_HEIGHT_PERCENT / 100);
    const clampedHeight = Math.max(MIN_HEIGHT, Math.min(height, maxHeight));
    setState((prev) => ({ ...prev, panelHeight: clampedHeight }));
  }, []);

  const toggleFullscreen = useCallback(() => {
    setState((prev) => ({ ...prev, isFullscreen: !prev.isFullscreen }));
  }, []);

  const executeCommand = useCallback(
    (command: string) => {
      if (!state.isOpen) {
        openPanel();
      }

      setTimeout(() => {
        if (activeTerminalRef.current) {
          activeTerminalRef.current.sendCommand(command);
        }
      }, 100);
    },
    [state.isOpen, openPanel],
  );

  const value: TerminalContextValue = {
    ...state,
    openPanel,
    closePanel,
    togglePanel,
    addTab,
    closeTab,
    setActiveTab,
    setPanelHeight,
    toggleFullscreen,
    executeCommand,
    activeTerminalRef,
    queryClient,
    minHeight: MIN_HEIGHT,
    maxHeightPercent: MAX_HEIGHT_PERCENT,
  };

  return (
    <TerminalContext.Provider value={value}>
      {children}
    </TerminalContext.Provider>
  );
}

export function useTerminal(): TerminalContextValue {
  const context = useContext(TerminalContext);
  if (!context) {
    throw new Error("useTerminal must be used within a TerminalProvider");
  }
  return context;
}
