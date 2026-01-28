import { useCallback, useRef, useEffect } from 'react';
import { useTerminal } from '../context';
import { TerminalTabs } from './TerminalTabs';
import { TerminalInstance, type TerminalInstanceHandle } from './TerminalInstance';
import { ResizeHandle } from './ResizeHandle';

export function TerminalPanel() {
  const {
    isOpen,
    tabs,
    activeTabId,
    panelHeight,
    isFullscreen,
    closePanel,
    addTab,
    closeTab,
    setActiveTab,
    setPanelHeight,
    toggleFullscreen,
    activeTerminalRef,
    minHeight,
    maxHeightPercent,
  } = useTerminal();

  // Map of terminal instance refs
  const terminalRefs = useRef<Map<string, TerminalInstanceHandle>>(new Map());

  // Update activeTerminalRef when active tab changes
  useEffect(() => {
    if (activeTabId) {
      const activeRef = terminalRefs.current.get(activeTabId);
      if (activeRef) {
        activeTerminalRef.current = activeRef;
      }
    }
  }, [activeTabId, activeTerminalRef]);

  const handleResize = useCallback((deltaY: number) => {
    setPanelHeight(panelHeight + deltaY);
  }, [panelHeight, setPanelHeight]);

  const setTerminalRef = useCallback((id: string, ref: TerminalInstanceHandle | null) => {
    if (ref) {
      terminalRefs.current.set(id, ref);
    } else {
      terminalRefs.current.delete(id);
    }
  }, []);

  if (!isOpen) {
    return null;
  }

  // Calculate panel styles
  const panelStyle: React.CSSProperties = isFullscreen
    ? {
        position: 'fixed',
        top: '64px', // Header height
        left: 0,
        right: 0,
        bottom: 0,
        zIndex: 40,
      }
    : {
        position: 'fixed',
        bottom: 0,
        left: '224px', // Sidebar width (w-56 = 14rem = 224px)
        right: 0,
        height: `${panelHeight}px`,
        zIndex: 40,
      };

  return (
    <div
      className="bg-[#1e1e1e] border-t border-[#3c3c3c] flex flex-col overflow-hidden"
      style={panelStyle}
    >
      {/* Resize handle (only when not fullscreen) */}
      {!isFullscreen && (
        <ResizeHandle
          onResize={handleResize}
          minHeight={minHeight}
          maxHeightPercent={maxHeightPercent}
        />
      )}

      {/* Tab bar */}
      <TerminalTabs
        tabs={tabs}
        activeTabId={activeTabId}
        isFullscreen={isFullscreen}
        onTabSelect={setActiveTab}
        onTabClose={closeTab}
        onAddTab={addTab}
        onToggleFullscreen={toggleFullscreen}
        onClose={closePanel}
      />

      {/* Terminal instances */}
      <div className="flex-1 overflow-hidden">
        {tabs.map((tab) => (
          <TerminalInstance
            key={tab.id}
            ref={(ref) => setTerminalRef(tab.id, ref)}
            sessionId={tab.id}
            isActive={tab.id === activeTabId}
          />
        ))}
      </div>
    </div>
  );
}
