import { useState, useEffect } from "react";
import type { Issue } from "../../../types/api";
import { getIssue } from "../../../lib/api";
import { MarkdownViewer } from "../../../components/ui";
import type { IssueModalProps } from "../types";

export function IssueModal({ issueId, onClose }: IssueModalProps) {
  const [issue, setIssue] = useState<Issue | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    setLoading(true);
    getIssue(issueId)
      .then(setIssue)
      .finally(() => setLoading(false));
  }, [issueId]);

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center">
      <div className="absolute inset-0 bg-black/70" onClick={onClose} />
      <div className="relative bg-card-bg rounded-xl w-[90%] max-w-[700px] max-h-[80vh] overflow-auto shadow-modal border border-night-border">
        <div className="flex justify-between items-center p-5 border-b border-night-border">
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
        <div className="p-5">
          {loading ? (
            <p className="text-text-muted">Loading issue details...</p>
          ) : issue ? (
            <>
              <div className="flex gap-2 flex-wrap mb-5">
                <span className="bg-red-900/30 text-neon-pink px-3 py-1 rounded-full text-sm font-medium">
                  P{issue.priority}
                </span>
                <span className="bg-blue-900/30 text-neon-blue px-3 py-1 rounded-full text-sm font-medium">
                  {issue.status || "open"}
                </span>
                {issue.assignee && (
                  <span className="bg-night-surface-bright text-text-secondary px-3 py-1 rounded-full text-sm">
                    {issue.assignee}
                  </span>
                )}
              </div>
              <MarkdownViewer
                content={issue.description || "*No description provided*"}
              />
              {issue.design && (
                <>
                  <hr className="my-5 border-night-border" />
                  <h4 className="text-neon-cyan font-semibold mb-3">Design</h4>
                  <MarkdownViewer content={issue.design} />
                </>
              )}
              {issue.acceptance_criteria && (
                <>
                  <hr className="my-5 border-night-border" />
                  <h4 className="text-neon-cyan font-semibold mb-3">
                    Acceptance Criteria
                  </h4>
                  <MarkdownViewer content={issue.acceptance_criteria} />
                </>
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
