import type { TerminalTab } from '../../context/TerminalContext';

interface TerminalTabsProps {
  tabs: TerminalTab[];
  activeTabId: string | null;
  isFullscreen: boolean;
  onTabSelect: (id: string) => void;
  onTabClose: (id: string) => void;
  onAddTab: () => void;
  onToggleFullscreen: () => void;
  onClose: () => void;
}

export function TerminalTabs({
  tabs,
  activeTabId,
  isFullscreen,
  onTabSelect,
  onTabClose,
  onAddTab,
  onToggleFullscreen,
  onClose,
}: TerminalTabsProps) {
  return (
    <div className="flex items-center h-9 bg-[#252526] border-b border-[#3c3c3c] px-2">
      {/* Terminal icon */}
      <div className="flex items-center gap-1 text-[#cccccc] mr-2">
        <span className="text-sm font-medium">TERMINAL</span>
      </div>

      {/* Tab list */}
      <div className="flex items-center gap-1 flex-1 overflow-x-auto">
        {tabs.map((tab) => (
          <button
            key={tab.id}
            onClick={() => onTabSelect(tab.id)}
            className={`
              group flex items-center gap-2 px-3 py-1 text-sm rounded-t
              transition-colors duration-150
              ${activeTabId === tab.id
                ? 'bg-[#1e1e1e] text-white border-t-2 border-[#7dcfff]'
                : 'text-[#969696] hover:text-white hover:bg-[#2d2d2d]'
              }
            `}
          >
            {/* Terminal icon for tab */}
            <span className={`text-xs ${activeTabId === tab.id ? 'text-[#7dcfff]' : 'opacity-70'}`}>&#62;_</span>
            <span>{tab.title}</span>
            {/* Close button */}
            <button
              onClick={(e) => {
                e.stopPropagation();
                onTabClose(tab.id);
              }}
              className={`
                w-4 h-4 flex items-center justify-center rounded
                opacity-0 group-hover:opacity-100
                hover:bg-[#505050] transition-opacity
                ${activeTabId === tab.id ? 'opacity-100' : ''}
              `}
              title="Close terminal"
            >
              <span className="text-xs">×</span>
            </button>
          </button>
        ))}

        {/* Add tab button */}
        <button
          onClick={onAddTab}
          className="w-6 h-6 flex items-center justify-center text-[#969696] hover:text-[#7dcfff] hover:bg-[#2d2d2d] rounded transition-colors"
          title="New Terminal"
        >
          <span className="text-lg">+</span>
        </button>
      </div>

      {/* Right-side controls */}
      <div className="flex items-center gap-1 ml-2">
        {/* Fullscreen toggle */}
        <button
          onClick={onToggleFullscreen}
          className="w-7 h-7 flex items-center justify-center text-[#969696] hover:text-[#7dcfff] hover:bg-[#2d2d2d] rounded transition-colors"
          title={isFullscreen ? 'Exit fullscreen' : 'Fullscreen'}
        >
          {isFullscreen ? (
            // Exit fullscreen icon
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
              <path d="M8 3v3a2 2 0 0 1-2 2H3m18 0h-3a2 2 0 0 1-2-2V3m0 18v-3a2 2 0 0 1 2-2h3M3 16h3a2 2 0 0 1 2 2v3" />
            </svg>
          ) : (
            // Fullscreen icon
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
              <path d="M8 3H5a2 2 0 0 0-2 2v3m18 0V5a2 2 0 0 0-2-2h-3m0 18h3a2 2 0 0 0 2-2v-3M3 16v3a2 2 0 0 0 2 2h3" />
            </svg>
          )}
        </button>

        {/* Close panel button */}
        <button
          onClick={onClose}
          className="w-7 h-7 flex items-center justify-center text-[#969696] hover:text-white hover:bg-[#2d2d2d] rounded transition-colors"
          title="Close panel"
        >
          <span className="text-lg">×</span>
        </button>
      </div>
    </div>
  );
}
