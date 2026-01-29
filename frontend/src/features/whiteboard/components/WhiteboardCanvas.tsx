import { Tldraw } from "tldraw";
import "tldraw/tldraw.css";
import { useWhiteboard } from "../hooks/useWhiteboard";
import { useClaudeExport } from "../hooks/useClaudeExport";
import { WhiteboardToolbar } from "./WhiteboardToolbar";

export function WhiteboardCanvas() {
  const { editor, isDirty, handleMount, clearCanvas } = useWhiteboard();
  const { copyToClipboard } = useClaudeExport(editor);

  return (
    <div className="h-full flex flex-col bg-card-bg rounded-xl border border-neon-magenta overflow-hidden">
      <WhiteboardToolbar
        onExportClaude={copyToClipboard}
        onClear={clearCanvas}
        isDirty={isDirty}
      />

      <div className="flex-1 relative">
        <Tldraw
          onMount={handleMount}
          persistenceKey="kitty-beads-whiteboard"
          className="absolute inset-0"
        />
      </div>
    </div>
  );
}
