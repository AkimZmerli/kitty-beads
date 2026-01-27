import {
  createContext,
  useContext,
  useState,
  useCallback,
  useEffect,
  useRef,
  type ReactNode,
} from 'react';
import { useQueryClient, type QueryClient } from '@tanstack/react-query';

// Storage key for persisting panel height
const STORAGE_KEY = 'terminal-panel-height';
const DEFAULT_HEIGHT = 300;
const MIN_HEIGHT = 150;
const MAX_HEIGHT_PERCENT = 70; // 70vh max

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
  activeTerminalRef: React.MutableRefObject<{ sendCommand: (cmd: string) => void } | null>;

  // React Query client for pattern detection
  queryClient: QueryClient;

  // Height constraints
  minHeight: number;
  maxHeightPercent: number;
}

const TerminalContext = createContext<TerminalContextValue | null>(null);

// Generate unique ID for terminal tabs
function generateId(): string {
  return `term-${Date.now()}-${Math.random().toString(36).slice(2, 9)}`;
}

// Get saved height from localStorage
function getSavedHeight(): number {
  try {
    const saved = localStorage.getItem(STORAGE_KEY);
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

export function TerminalProvider({ children }: { children: ReactNode }) {
  const queryClient = useQueryClient();
  const activeTerminalRef = useRef<{ sendCommand: (cmd: string) => void } | null>(null);

  const [state, setState] = useState<TerminalState>({
    isOpen: false,
    tabs: [],
    activeTabId: null,
    panelHeight: getSavedHeight(),
    isFullscreen: false,
  });

  // Persist panel height to localStorage
  useEffect(() => {
    try {
      localStorage.setItem(STORAGE_KEY, state.panelHeight.toString());
    } catch {
      // Ignore localStorage errors
    }
  }, [state.panelHeight]);

  // Keyboard shortcut: Ctrl+` to toggle terminal
  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if (e.ctrlKey && e.key === '`') {
        e.preventDefault();
        togglePanel();
      }
    };

    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, []);

  const openPanel = useCallback(() => {
    setState(prev => {
      const newState = { ...prev, isOpen: true };
      // Auto-create first tab if no tabs exist
      if (prev.tabs.length === 0) {
        const newTab = {
          id: generateId(),
          title: 'Terminal 1',
          createdAt: Date.now(),
        };
        newState.tabs = [newTab];
        newState.activeTabId = newTab.id;
      }
      return newState;
    });
  }, []);

  const closePanel = useCallback(() => {
    setState(prev => ({ ...prev, isOpen: false, isFullscreen: false }));
  }, []);

  const togglePanel = useCallback(() => {
    setState(prev => {
      if (prev.isOpen) {
        return { ...prev, isOpen: false, isFullscreen: false };
      }
      // Opening - ensure we have a tab
      const newState = { ...prev, isOpen: true };
      if (prev.tabs.length === 0) {
        const newTab = {
          id: generateId(),
          title: 'Terminal 1',
          createdAt: Date.now(),
        };
        newState.tabs = [newTab];
        newState.activeTabId = newTab.id;
      }
      return newState;
    });
  }, []);

  const addTab = useCallback((): string => {
    const newTab = {
      id: generateId(),
      title: `Terminal ${state.tabs.length + 1}`,
      createdAt: Date.now(),
    };

    setState(prev => ({
      ...prev,
      tabs: [...prev.tabs, newTab],
      activeTabId: newTab.id,
      isOpen: true,
    }));

    return newTab.id;
  }, [state.tabs.length]);

  const closeTab = useCallback((id: string) => {
    setState(prev => {
      const newTabs = prev.tabs.filter(t => t.id !== id);
      let newActiveId = prev.activeTabId;

      // If closing active tab, switch to another
      if (prev.activeTabId === id) {
        const closingIndex = prev.tabs.findIndex(t => t.id === id);
        if (newTabs.length > 0) {
          // Prefer the tab to the left, or the first remaining tab
          newActiveId = newTabs[Math.max(0, closingIndex - 1)]?.id || newTabs[0]?.id || null;
        } else {
          newActiveId = null;
        }
      }

      return {
        ...prev,
        tabs: newTabs,
        activeTabId: newActiveId,
        // Close panel if no tabs left
        isOpen: newTabs.length > 0 ? prev.isOpen : false,
      };
    });
  }, []);

  const setActiveTab = useCallback((id: string) => {
    setState(prev => ({ ...prev, activeTabId: id }));
  }, []);

  const setPanelHeight = useCallback((height: number) => {
    const maxHeight = window.innerHeight * (MAX_HEIGHT_PERCENT / 100);
    const clampedHeight = Math.max(MIN_HEIGHT, Math.min(height, maxHeight));
    setState(prev => ({ ...prev, panelHeight: clampedHeight }));
  }, []);

  const toggleFullscreen = useCallback(() => {
    setState(prev => ({ ...prev, isFullscreen: !prev.isFullscreen }));
  }, []);

  const executeCommand = useCallback((command: string) => {
    // Ensure terminal is open
    if (!state.isOpen) {
      openPanel();
    }

    // Wait a tick for the terminal to be ready, then send command
    setTimeout(() => {
      if (activeTerminalRef.current) {
        activeTerminalRef.current.sendCommand(command);
      }
    }, 100);
  }, [state.isOpen, openPanel]);

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
    throw new Error('useTerminal must be used within a TerminalProvider');
  }
  return context;
}
