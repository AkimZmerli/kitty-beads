import { useState } from "react";
import { MarkdownViewer } from "../components/ui";
import { Code2, Eye, Columns2 } from "lucide-react";

type ViewMode = "split" | "edit" | "preview";

const DEFAULT_CONTENT = `# Scratchpad

Start typing here to draft ideas, notes, or plans.

## Tips

- Use **Markdown** for formatting
- Switch between Edit, Split, and Preview modes
- This is a local scratchpad - content is not saved yet
`;

export function Editor() {
  const [viewMode, setViewMode] = useState<ViewMode>("split");
  const [content, setContent] = useState(DEFAULT_CONTENT);

  return (
    <div className="bg-card-bg rounded-xl border border-neon-magenta h-[calc(100vh-140px)] flex flex-col">
      {/* Header */}
      <div className="flex items-center justify-between px-6 py-4 border-b border-night-border">
        <div className="flex items-center gap-3">
          <Code2 className="w-5 h-5 text-neon-cyan" />
          <h2 className="text-neon-cyan text-xl font-bold">Editor</h2>
        </div>

        {/* View mode toggle */}
        <div className="flex bg-night-surface-bright rounded-md p-0.5">
          <button
            onClick={() => setViewMode("edit")}
            className={`flex items-center gap-1.5 px-3 py-1.5 rounded text-sm transition-all ${
              viewMode === "edit"
                ? "bg-neon-cyan text-night-bg font-medium"
                : "text-text-secondary hover:text-text-primary"
            }`}
          >
            <Code2 className="w-4 h-4" />
            Edit
          </button>
          <button
            onClick={() => setViewMode("split")}
            className={`flex items-center gap-1.5 px-3 py-1.5 rounded text-sm transition-all ${
              viewMode === "split"
                ? "bg-neon-cyan text-night-bg font-medium"
                : "text-text-secondary hover:text-text-primary"
            }`}
          >
            <Columns2 className="w-4 h-4" />
            Split
          </button>
          <button
            onClick={() => setViewMode("preview")}
            className={`flex items-center gap-1.5 px-3 py-1.5 rounded text-sm transition-all ${
              viewMode === "preview"
                ? "bg-neon-cyan text-night-bg font-medium"
                : "text-text-secondary hover:text-text-primary"
            }`}
          >
            <Eye className="w-4 h-4" />
            Preview
          </button>
        </div>
      </div>

      {/* Editor/Preview area */}
      <div className="flex flex-1 min-h-0">
        {/* Editor pane */}
        {(viewMode === "edit" || viewMode === "split") && (
          <div
            className={`flex flex-col ${viewMode === "split" ? "w-1/2 border-r border-night-border" : "flex-1"}`}
          >
            <div className="px-4 py-2 text-xs text-text-muted border-b border-night-border bg-night-surface/50">
              MARKDOWN
            </div>
            <textarea
              value={content}
              onChange={(e) => setContent(e.target.value)}
              className="flex-1 w-full p-4 bg-transparent text-text-primary font-mono text-sm resize-none focus:outline-none"
              placeholder="Start writing..."
            />
          </div>
        )}

        {/* Preview pane */}
        {(viewMode === "preview" || viewMode === "split") && (
          <div
            className={`flex flex-col ${viewMode === "split" ? "w-1/2" : "flex-1"}`}
          >
            <div className="px-4 py-2 text-xs text-text-muted border-b border-night-border bg-night-surface/50">
              PREVIEW
            </div>
            <div className="flex-1 p-4 overflow-auto">
              <MarkdownViewer content={content} />
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
