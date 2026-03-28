import { Tldraw } from "tldraw";
import "tldraw/tldraw.css";
import { useWhiteboard } from "../hooks/useWhiteboard";
import { useClaudeExport } from "../hooks/useClaudeExport";
import { WhiteboardToolbar } from "./WhiteboardToolbar";

// Layout constants (must match Layout.tsx)
const HEADER_HEIGHT = 73;
const LAYOUT_PADDING = 32; // p-8 = 2rem = 32px
const TOOLBAR_HEIGHT = 44; // Approximate toolbar height

export function WhiteboardCanvas() {
  const { editor, isDirty, handleMount, clearCanvas } = useWhiteboard();
  const { exportPngToUserStorage } = useClaudeExport(editor);

  // Calculate available height: viewport - header - layout padding (top + bottom) - toolbar
  const canvasHeight = `calc(100vh - ${HEADER_HEIGHT}px - ${LAYOUT_PADDING * 2}px - ${TOOLBAR_HEIGHT}px)`;

  return (
    <div className="flex flex-col bg-card-bg rounded-xl border border-neon-magenta overflow-hidden">
      <WhiteboardToolbar
        onSaveForClaude={exportPngToUserStorage}
        onClear={clearCanvas}
        isDirty={isDirty}
      />

      <div className="relative isolate" style={{ height: canvasHeight }}>
        <Tldraw
          onMount={handleMount}
          persistenceKey="kitty-beads-whiteboard"
          className="absolute inset-0 z-0"
        />
      </div>
    </div>
  );
}
