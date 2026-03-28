import { useState, useEffect, useCallback } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { useIdeation } from "../hooks/useIdeation";
import { MarkdownViewer } from "../components/ui";
import { CommentThread } from "../features/comments";
import { extractSummary } from "../lib/planParser";
import {
  ChevronLeft,
  AlertTriangle,
  RefreshCw,
  X,
  MessageCircle,
  PanelRightClose,
  PanelRightOpen,
} from "lucide-react";

type ViewMode = "split" | "edit" | "preview";

export function IdeationPad() {
  const { issueId } = useParams<{ issueId: string }>();
  const navigate = useNavigate();

  const [content, setContent] = useState("");
  const [hasChanges, setHasChanges] = useState(false);
  const [viewMode, setViewMode] = useState<ViewMode>("split");
  const [showComments, setShowComments] = useState(false);

  const {
    issue,
    isLoading,
    error,
    save,
    isSaving,
    hasExternalChanges,
    startEditing,
    refreshContent,
    dismissConflict,
  } = useIdeation({
    issueId: issueId || "",
    onSaveSuccess: () => {
      setHasChanges(false);
    },
  });

  // Initialize content when issue loads
  useEffect(() => {
    if (issue) {
      const initialContent = issue.design || issue.description || "";
      setContent(initialContent);
      setHasChanges(false);
    }
  }, [issue]);

  const handleContentChange = useCallback(
    (newContent: string) => {
      setContent(newContent);
      const original = issue?.design || issue?.description || "";
      setHasChanges(newContent !== original);
      // Track that user started editing
      startEditing(original);
    },
    [issue, startEditing],
  );

  const handleRefresh = useCallback(() => {
    // Refresh will reload the issue, triggering the useEffect to update content
    refreshContent();
    setHasChanges(false);
  }, [refreshContent]);

  const handleSave = useCallback(() => {
    save({ design: content });
  }, [save, content]);

  const handleCancel = useCallback(() => {
    if (hasChanges) {
      const confirmed = window.confirm(
        "You have unsaved changes. Are you sure you want to discard them?",
      );
      if (!confirmed) return;
    }
    navigate("/roadmap");
  }, [navigate, hasChanges]);

  if (!issueId) {
    return (
      <div className="bg-card-bg rounded-xl p-8 border border-neon-magenta">
        <p className="text-neon-pink">No issue ID provided</p>
      </div>
    );
  }

  if (isLoading) {
    return (
      <div className="bg-card-bg rounded-xl p-8 border border-neon-magenta">
        <div className="flex items-center justify-center py-12">
          <div className="text-text-muted">Loading issue...</div>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="bg-card-bg rounded-xl p-8 border border-neon-magenta">
        <p className="text-neon-pink">Error loading issue: {error.message}</p>
        <button
          onClick={() => navigate("/roadmap")}
          className="mt-4 px-4 py-2 bg-night-surface text-text-normal rounded-md hover:bg-night-surface-bright transition-colors"
        >
          Back to Roadmap
        </button>
      </div>
    );
  }

  const summary = extractSummary(content);

  return (
    <div className="h-full flex flex-col bg-card-bg rounded-xl border border-neon-magenta overflow-hidden">
      {/* Header */}
      <div className="flex items-center justify-between px-6 py-4 border-b border-night-border bg-night-bg-highlight">
        <div className="flex items-center gap-4">
          <button
            onClick={handleCancel}
            className="text-text-muted hover:text-text-primary transition-colors"
            title="Back to Roadmap"
          >
            <ChevronLeft className="w-5 h-5" />
          </button>
          <div>
            <h1 className="text-lg font-semibold text-text-bright">
              {issue?.title || "Ideation Pad"}
            </h1>
            <p className="text-sm text-text-muted font-mono">{issueId}</p>
          </div>
          {hasChanges && (
            <span className="px-2 py-0.5 bg-orange-900/30 text-neon-orange text-xs rounded-full font-medium">
              Unsaved
            </span>
          )}
        </div>

        <div className="flex items-center gap-3">
          {/* View mode toggle */}
          <div className="flex bg-night-surface-bright rounded-md p-0.5">
            <button
              onClick={() => setViewMode("edit")}
              className={`px-3 py-1 text-sm rounded transition-colors ${
                viewMode === "edit"
                  ? "bg-night-bg shadow-sm text-text-bright"
                  : "text-text-muted hover:text-text-primary"
              }`}
            >
              Edit
            </button>
            <button
              onClick={() => setViewMode("split")}
              className={`px-3 py-1 text-sm rounded transition-colors ${
                viewMode === "split"
                  ? "bg-night-bg shadow-sm text-text-bright"
                  : "text-text-muted hover:text-text-primary"
              }`}
            >
              Split
            </button>
            <button
              onClick={() => setViewMode("preview")}
              className={`px-3 py-1 text-sm rounded transition-colors ${
                viewMode === "preview"
                  ? "bg-night-bg shadow-sm text-text-bright"
                  : "text-text-muted hover:text-text-primary"
              }`}
            >
              Preview
            </button>
          </div>

          {/* Comments toggle */}
          <button
            onClick={() => setShowComments(!showComments)}
            className={`flex items-center gap-2 px-3 py-1.5 text-sm rounded-md transition-colors ${
              showComments
                ? "bg-neon-cyan text-night-bg"
                : "bg-night-surface text-text-muted hover:text-text-primary"
            }`}
            title={showComments ? "Hide comments" : "Show comments"}
          >
            <MessageCircle className="w-4 h-4" />
            {showComments ? (
              <PanelRightClose className="w-4 h-4" />
            ) : (
              <PanelRightOpen className="w-4 h-4" />
            )}
          </button>

          {/* Action buttons */}
          <button
            onClick={handleCancel}
            className="px-4 py-1.5 text-sm text-text-muted hover:text-text-primary transition-colors"
          >
            Cancel
          </button>
          <button
            onClick={handleSave}
            disabled={!hasChanges || isSaving}
            className={`px-4 py-1.5 text-sm rounded-md transition-all ${
              hasChanges && !isSaving
                ? "bg-neon-cyan text-night-bg hover:shadow-[0_0_12px_rgba(125,207,255,0.4)]"
                : "bg-night-surface text-text-muted cursor-not-allowed"
            }`}
          >
            {isSaving ? "Saving..." : "Save"}
          </button>
        </div>
      </div>

      {/* External changes warning banner */}
      {hasExternalChanges && (
        <div className="px-6 py-3 bg-orange-900/20 border-b border-orange-900/30 flex items-center justify-between">
          <div className="flex items-center gap-3">
            <AlertTriangle className="w-5 h-5 text-neon-orange" />
            <p className="text-sm text-neon-orange">
              <span className="font-medium">External changes detected!</span>{" "}
              This issue was modified by another session.
            </p>
          </div>
          <div className="flex items-center gap-2">
            <button
              onClick={handleRefresh}
              className="flex items-center gap-1.5 px-3 py-1.5 bg-orange-900/30 text-neon-orange text-sm rounded-md hover:bg-orange-900/40 transition-colors"
            >
              <RefreshCw className="w-4 h-4" />
              Refresh
            </button>
            <button
              onClick={dismissConflict}
              className="p-1.5 text-text-muted hover:text-text-primary transition-colors"
              title="Dismiss"
            >
              <X className="w-4 h-4" />
            </button>
          </div>
        </div>
      )}

      {/* Summary preview */}
      {summary && (
        <div className="px-6 py-3 bg-blue-900/20 border-b border-blue-900/30">
          <p className="text-sm text-neon-blue">
            <span className="font-medium">Roadmap summary:</span> {summary}
          </p>
        </div>
      )}

      {/* Main content area */}
      <div className="flex-1 flex overflow-hidden">
        {/* Editor/Preview section */}
        <div
          className={`flex-1 flex overflow-hidden ${showComments ? "border-r border-night-border" : ""}`}
        >
          {/* Editor */}
          {(viewMode === "edit" || viewMode === "split") && (
            <div
              className={`flex flex-col ${viewMode === "split" ? "w-1/2 border-r border-night-border" : "flex-1"}`}
            >
              <div className="px-4 py-2 bg-night-bg-highlight border-b border-night-border">
                <span className="text-xs text-text-muted uppercase tracking-wider font-medium">
                  Editor
                </span>
              </div>
              <textarea
                value={content}
                onChange={(e) => handleContentChange(e.target.value)}
                className="flex-1 p-4 font-mono text-sm resize-none focus:outline-none bg-night-bg text-text-bright placeholder:text-text-muted"
                placeholder="Write your plan in markdown..."
                spellCheck={false}
              />
            </div>
          )}

          {/* Preview */}
          {(viewMode === "preview" || viewMode === "split") && (
            <div
              className={`flex flex-col ${viewMode === "split" ? "w-1/2" : "flex-1"}`}
            >
              <div className="px-4 py-2 bg-night-bg-highlight border-b border-night-border">
                <span className="text-xs text-text-muted uppercase tracking-wider font-medium">
                  Preview
                </span>
              </div>
              <div className="flex-1 p-4 overflow-auto bg-night-surface">
                {content ? (
                  <MarkdownViewer content={content} />
                ) : (
                  <p className="text-text-muted italic">
                    Start writing to see preview...
                  </p>
                )}
              </div>
            </div>
          )}
        </div>

        {/* Comments panel */}
        {showComments && issueId && (
          <div className="w-80 flex-shrink-0 bg-night-surface">
            <CommentThread issueId={issueId} className="h-full" />
          </div>
        )}
      </div>
    </div>
  );
}
