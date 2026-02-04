import { useState } from "react";
import { Save, Check, Trash2, Loader2 } from "lucide-react";
import type { PngSaveResult } from "../hooks/useClaudeExport";

type SaveState = "idle" | "saving" | "saved" | "error";

interface WhiteboardToolbarProps {
  onSaveForClaude: () => Promise<PngSaveResult>;
  onClear: () => void;
  isDirty: boolean;
}

export function WhiteboardToolbar({
  onSaveForClaude,
  onClear,
  isDirty,
}: WhiteboardToolbarProps) {
  const [saveState, setSaveState] = useState<SaveState>("idle");

  const handleSave = async () => {
    setSaveState("saving");
    const result = await onSaveForClaude();
    if (result.success) {
      setSaveState("saved");
      setTimeout(() => setSaveState("idle"), 2000);
    } else {
      setSaveState("error");
      setTimeout(() => setSaveState("idle"), 3000);
    }
  };

  const handleClear = () => {
    if (window.confirm("Are you sure you want to clear the whiteboard?")) {
      onClear();
    }
  };

  const renderButtonContent = () => {
    switch (saveState) {
      case "saving":
        return (
          <>
            <Loader2 className="w-4 h-4 animate-spin" />
            Saving...
          </>
        );
      case "saved":
        return (
          <>
            <Check className="w-4 h-4" />
            Saved!
          </>
        );
      case "error":
        return (
          <>
            <Save className="w-4 h-4" />
            Failed
          </>
        );
      default:
        return (
          <>
            <Save className="w-4 h-4" />
            Save for Claude
          </>
        );
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
        onClick={handleSave}
        disabled={saveState === "saving"}
        className={`flex items-center gap-2 px-4 py-1.5 text-sm rounded-md transition-all font-medium ${
          saveState === "error"
            ? "bg-red-600 text-white"
            : saveState === "saved"
              ? "bg-green-600 text-white"
              : "bg-neon-cyan text-night-bg hover:shadow-[0_0_12px_rgba(125,207,255,0.4)]"
        } disabled:opacity-50 disabled:cursor-not-allowed`}
      >
        {renderButtonContent()}
      </button>
    </div>
  );
}
