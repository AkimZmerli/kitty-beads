import { useState } from "react";
import { Copy, Check, Trash2 } from "lucide-react";

interface WhiteboardToolbarProps {
  onExportClaude: () => Promise<boolean>;
  onClear: () => void;
  isDirty: boolean;
}

export function WhiteboardToolbar({
  onExportClaude,
  onClear,
  isDirty,
}: WhiteboardToolbarProps) {
  const [copied, setCopied] = useState(false);

  const handleExport = async () => {
    const success = await onExportClaude();
    if (success) {
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    }
  };

  const handleClear = () => {
    if (window.confirm("Are you sure you want to clear the whiteboard?")) {
      onClear();
    }
  };

  return (
    <div className="flex items-center gap-2 px-4 py-2 bg-night-bg-highlight border-b border-night-border">
      <span className="text-xs text-text-muted uppercase tracking-wider font-medium mr-4">
        Whiteboard
      </span>

      {isDirty && (
        <span className="px-2 py-0.5 bg-orange-900/30 text-neon-orange text-xs rounded-full font-medium">
          Unsaved
        </span>
      )}

      <div className="flex-1" />

      <button
        onClick={handleClear}
        className="flex items-center gap-2 px-3 py-1.5 text-sm text-text-muted hover:text-neon-pink hover:bg-night-surface-bright rounded-md transition-colors"
        title="Clear whiteboard"
      >
        <Trash2 className="w-4 h-4" />
        Clear
      </button>

      <button
        onClick={handleExport}
        className="flex items-center gap-2 px-4 py-1.5 bg-neon-cyan text-night-bg text-sm rounded-md hover:shadow-[0_0_12px_rgba(125,207,255,0.4)] transition-all font-medium"
      >
        {copied ? (
          <>
            <Check className="w-4 h-4" />
            Copied!
          </>
        ) : (
          <>
            <Copy className="w-4 h-4" />
            Export to Claude
          </>
        )}
      </button>
    </div>
  );
}
