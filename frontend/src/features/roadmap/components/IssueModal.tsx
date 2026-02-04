import { useState, useEffect } from "react";
import { MarkdownViewer } from "../../../components/ui";
import { getIssue } from "../../../lib/api";
import type { Issue } from "../../../types/api";

interface IssueModalProps {
  issueId: string;
  onClose: () => void;
}

export function IssueModal({ issueId, onClose }: IssueModalProps) {
  const [issue, setIssue] = useState<Issue | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    let cancelled = false;
    setLoading(true);

    getIssue(issueId)
      .then((data) => {
        if (!cancelled) {
          setIssue(data);
        }
      })
      .catch((err) => {
        console.error("Failed to load issue:", err);
      })
      .finally(() => {
        if (!cancelled) {
          setLoading(false);
        }
      });

    return () => {
      cancelled = true;
    };
  }, [issueId]);

  // Close on Escape
  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if (e.key === "Escape") {
        onClose();
      }
    };
    window.addEventListener("keydown", handleKeyDown);
    return () => window.removeEventListener("keydown", handleKeyDown);
  }, [onClose]);

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center">
      {/* Backdrop */}
      <div
        className="absolute inset-0 bg-black/70 animate-in fade-in duration-150"
        onClick={onClose}
      />

      {/* Modal */}
      <div className="relative bg-card-bg rounded-xl w-[90%] max-w-[800px] max-h-[85vh] overflow-auto shadow-modal border border-night-border animate-in fade-in zoom-in-95 duration-200">
        {/* Header */}
        <div className="sticky top-0 bg-card-bg flex justify-between items-center p-5 border-b border-night-border z-10">
          <h3 className="text-lg font-semibold text-text-primary">
            {loading ? "Loading..." : `${issue?.id}: ${issue?.title}`}
          </h3>
          <button
            onClick={onClose}
            className="text-2xl text-text-muted hover:text-text-primary transition-colors"
          >
            &times;
          </button>
        </div>

        {/* Content */}
        <div className="p-5">
          {loading ? (
            <div className="flex items-center justify-center py-8">
              <div className="w-6 h-6 border-2 border-neon-cyan border-t-transparent rounded-full animate-spin" />
            </div>
          ) : issue ? (
            <>
              {/* Metadata badges */}
              <div className="flex gap-2 flex-wrap mb-5">
                <span className="bg-red-900/30 text-neon-pink px-3 py-1 rounded-full text-sm font-medium">
                  P{issue.priority}
                </span>
                <span className="bg-blue-900/30 text-neon-blue px-3 py-1 rounded-full text-sm font-medium">
                  {issue.status || "open"}
                </span>
                {issue.issue_type && (
                  <span className="bg-purple-900/30 text-neon-magenta px-3 py-1 rounded-full text-sm font-medium">
                    {issue.issue_type}
                  </span>
                )}
                {issue.assignee && (
                  <span className="bg-night-surface-bright text-text-secondary px-3 py-1 rounded-full text-sm">
                    {issue.assignee}
                  </span>
                )}
              </div>

              {/* Description */}
              {issue.description && (
                <div className="mb-6">
                  <h4 className="text-neon-cyan font-semibold mb-3">
                    Description
                  </h4>
                  <MarkdownViewer content={issue.description} />
                </div>
              )}

              {/* Design/Plan */}
              {issue.design && (
                <div className="mb-6">
                  <h4 className="text-neon-cyan font-semibold mb-3">
                    Design / Plan
                  </h4>
                  <MarkdownViewer content={issue.design} />
                </div>
              )}

              {/* Acceptance Criteria */}
              {issue.acceptance_criteria && (
                <div className="mb-6">
                  <h4 className="text-neon-cyan font-semibold mb-3">
                    Acceptance Criteria
                  </h4>
                  <MarkdownViewer content={issue.acceptance_criteria} />
                </div>
              )}

              {/* No content fallback */}
              {!issue.description &&
                !issue.design &&
                !issue.acceptance_criteria && (
                  <p className="text-text-muted italic">
                    No detailed content available for this bead.
                  </p>
                )}
            </>
          ) : (
            <p className="text-neon-pink">Failed to load issue</p>
          )}
        </div>
      </div>
    </div>
  );
}
