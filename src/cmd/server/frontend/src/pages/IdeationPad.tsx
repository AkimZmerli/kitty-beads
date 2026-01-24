import { useState, useEffect, useCallback } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { useIdeation } from "../hooks/useIdeation";
import { MarkdownViewer } from "../components/MarkdownViewer";
import { extractSummary } from "../lib/planParser";

type ViewMode = "split" | "edit" | "preview";

export function IdeationPad() {
  const { issueId } = useParams<{ issueId: string }>();
  const navigate = useNavigate();

  const [content, setContent] = useState("");
  const [hasChanges, setHasChanges] = useState(false);
  const [viewMode, setViewMode] = useState<ViewMode>("split");

  const { issue, isLoading, error, save, isSaving } = useIdeation({
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
    },
    [issue]
  );

  const handleSave = useCallback(() => {
    save({ design: content });
  }, [save, content]);

  const handleCancel = useCallback(() => {
    if (hasChanges) {
      const confirmed = window.confirm(
        "You have unsaved changes. Are you sure you want to discard them?"
      );
      if (!confirmed) return;
    }
    navigate("/roadmap");
  }, [navigate, hasChanges]);

  if (!issueId) {
    return (
      <div className="bg-white rounded-xl p-8 border-2 border-sunny-yellow-border">
        <p className="text-red-500">No issue ID provided</p>
      </div>
    );
  }

  if (isLoading) {
    return (
      <div className="bg-white rounded-xl p-8 border-2 border-sunny-yellow-border">
        <div className="flex items-center justify-center py-12">
          <div className="text-text-muted">Loading issue...</div>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="bg-white rounded-xl p-8 border-2 border-sunny-yellow-border">
        <p className="text-red-500">Error loading issue: {error.message}</p>
        <button
          onClick={() => navigate("/roadmap")}
          className="mt-4 px-4 py-2 bg-gray-200 text-gray-700 rounded-md hover:bg-gray-300"
        >
          Back to Roadmap
        </button>
      </div>
    );
  }

  const summary = extractSummary(content);

  return (
    <div className="h-full flex flex-col bg-white rounded-xl border-2 border-sunny-yellow-border overflow-hidden">
      {/* Header */}
      <div className="flex items-center justify-between px-6 py-4 border-b border-gray-200 bg-gray-50">
        <div className="flex items-center gap-4">
          <button
            onClick={handleCancel}
            className="text-gray-500 hover:text-gray-700 transition-colors"
            title="Back to Roadmap"
          >
            <svg
              className="w-5 h-5"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M15 19l-7-7 7-7"
              />
            </svg>
          </button>
          <div>
            <h1 className="text-lg font-semibold text-gray-900">
              {issue?.title || "Ideation Pad"}
            </h1>
            <p className="text-sm text-gray-500 font-mono">{issueId}</p>
          </div>
          {hasChanges && (
            <span className="px-2 py-0.5 bg-amber-100 text-amber-700 text-xs rounded-full font-medium">
              Unsaved
            </span>
          )}
        </div>

        <div className="flex items-center gap-3">
          {/* View mode toggle */}
          <div className="flex bg-gray-200 rounded-md p-0.5">
            <button
              onClick={() => setViewMode("edit")}
              className={`px-3 py-1 text-sm rounded ${
                viewMode === "edit"
                  ? "bg-white shadow-sm text-gray-900"
                  : "text-gray-600 hover:text-gray-900"
              }`}
            >
              Edit
            </button>
            <button
              onClick={() => setViewMode("split")}
              className={`px-3 py-1 text-sm rounded ${
                viewMode === "split"
                  ? "bg-white shadow-sm text-gray-900"
                  : "text-gray-600 hover:text-gray-900"
              }`}
            >
              Split
            </button>
            <button
              onClick={() => setViewMode("preview")}
              className={`px-3 py-1 text-sm rounded ${
                viewMode === "preview"
                  ? "bg-white shadow-sm text-gray-900"
                  : "text-gray-600 hover:text-gray-900"
              }`}
            >
              Preview
            </button>
          </div>

          {/* Action buttons */}
          <button
            onClick={handleCancel}
            className="px-4 py-1.5 text-sm text-gray-600 hover:text-gray-900 transition-colors"
          >
            Cancel
          </button>
          <button
            onClick={handleSave}
            disabled={!hasChanges || isSaving}
            className={`px-4 py-1.5 text-sm rounded-md transition-colors ${
              hasChanges && !isSaving
                ? "bg-grassy-green text-white hover:bg-green-700"
                : "bg-gray-200 text-gray-400 cursor-not-allowed"
            }`}
          >
            {isSaving ? "Saving..." : "Save"}
          </button>
        </div>
      </div>

      {/* Summary preview */}
      {summary && (
        <div className="px-6 py-3 bg-blue-50 border-b border-blue-100">
          <p className="text-sm text-blue-700">
            <span className="font-medium">Roadmap summary:</span> {summary}
          </p>
        </div>
      )}

      {/* Main content area */}
      <div className="flex-1 flex overflow-hidden">
        {/* Editor */}
        {(viewMode === "edit" || viewMode === "split") && (
          <div
            className={`flex flex-col ${viewMode === "split" ? "w-1/2 border-r border-gray-200" : "flex-1"}`}
          >
            <div className="px-4 py-2 bg-gray-50 border-b border-gray-200">
              <span className="text-xs text-gray-500 uppercase tracking-wider font-medium">
                Editor
              </span>
            </div>
            <textarea
              value={content}
              onChange={(e) => handleContentChange(e.target.value)}
              className="flex-1 p-4 font-mono text-sm resize-none focus:outline-none bg-gray-900 text-gray-100"
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
            <div className="px-4 py-2 bg-gray-50 border-b border-gray-200">
              <span className="text-xs text-gray-500 uppercase tracking-wider font-medium">
                Preview
              </span>
            </div>
            <div className="flex-1 p-4 overflow-auto">
              {content ? (
                <MarkdownViewer content={content} />
              ) : (
                <p className="text-gray-400 italic">
                  Start writing to see preview...
                </p>
              )}
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
