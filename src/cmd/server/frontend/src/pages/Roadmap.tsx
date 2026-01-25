import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { useRoadmap } from "../hooks/useRoadmap";
import { BeadCard } from "../components/roadmap/BeadCard";
import { MarkdownViewer } from "../components/MarkdownViewer";
import { getIssue } from "../lib/api";
import type { Issue } from "../types/api";

export function Roadmap() {
  const navigate = useNavigate();
  const { data, isLoading, error } = useRoadmap();
  const [viewingIssue, setViewingIssue] = useState<Issue | null>(null);
  const [loadingIssue, setLoadingIssue] = useState(false);

  const handleEditPlan = (issueId: string) => {
    navigate(`/ideation/${issueId}`);
  };

  const handleViewFull = async (issueId: string) => {
    setLoadingIssue(true);
    try {
      const issue = await getIssue(issueId);
      setViewingIssue(issue);
    } catch (err) {
      console.error("Failed to load issue:", err);
    } finally {
      setLoadingIssue(false);
    }
  };

  if (isLoading) {
    return (
      <div className="bg-card-bg rounded-xl p-8 border border-neon-magenta">
        <h2 className="text-neon-cyan text-2xl font-bold mb-4 drop-shadow-[0_0_8px_rgba(125,207,255,0.3)]">
          Roadmap
        </h2>
        <p className="text-text-muted text-center py-8">Loading roadmap...</p>
      </div>
    );
  }

  if (error) {
    return (
      <div className="bg-card-bg rounded-xl p-8 border border-neon-magenta">
        <h2 className="text-neon-cyan text-2xl font-bold mb-4 drop-shadow-[0_0_8px_rgba(125,207,255,0.3)]">
          Roadmap
        </h2>
        <p className="text-neon-pink text-center py-8">
          Error loading roadmap data
        </p>
      </div>
    );
  }

  const issues = data?.issues || [];

  return (
    <div className="bg-card-bg rounded-xl p-8 border border-neon-magenta">
      {/* Header */}
      <div className="flex justify-between items-center mb-6">
        <div>
          <h2 className="text-neon-cyan text-2xl font-bold drop-shadow-[0_0_8px_rgba(125,207,255,0.3)]">
            Roadmap
          </h2>
          <p className="text-text-secondary">
            {data?.openCount || 0} open · {data?.completedCount || 0} completed
            · {data?.totalCount || 0} total
          </p>
        </div>
      </div>

      {/* Filter tabs */}
      <div className="flex gap-2 mb-6 border-b border-night-border pb-4">
        <button className="px-4 py-2 bg-neon-cyan text-night-bg rounded-md text-sm font-medium shadow-[0_0_10px_rgba(125,207,255,0.3)]">
          All
        </button>
        <button className="px-4 py-2 bg-night-surface text-text-normal rounded-md text-sm font-medium hover:bg-night-surface-bright transition-colors">
          Epics
        </button>
        <button className="px-4 py-2 bg-night-surface text-text-normal rounded-md text-sm font-medium hover:bg-night-surface-bright transition-colors">
          Open
        </button>
        <button className="px-4 py-2 bg-night-surface text-text-normal rounded-md text-sm font-medium hover:bg-night-surface-bright transition-colors">
          P1 Only
        </button>
      </div>

      {/* Beads list */}
      {issues.length === 0 ? (
        <div className="text-center py-12 text-text-muted">
          <p className="text-lg mb-2">No beads yet</p>
          <p className="text-sm">
            Create your first bead with{" "}
            <code className="bg-night-surface px-2 py-0.5 rounded text-neon-cyan">
              bd create "Your task"
            </code>
          </p>
        </div>
      ) : (
        <div className="space-y-1">
          {issues.map((issue) => (
            <BeadCard
              key={issue.id}
              issue={issue}
              onEditPlan={handleEditPlan}
              onViewFull={handleViewFull}
            />
          ))}
        </div>
      )}

      {/* Full view modal */}
      {(viewingIssue || loadingIssue) && (
        <IssueModal
          issue={viewingIssue}
          loading={loadingIssue}
          onClose={() => setViewingIssue(null)}
        />
      )}
    </div>
  );
}

interface IssueModalProps {
  issue: Issue | null;
  loading: boolean;
  onClose: () => void;
}

function IssueModal({ issue, loading, onClose }: IssueModalProps) {
  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center">
      <div className="absolute inset-0 bg-black/70" onClick={onClose} />
      <div className="relative bg-card-bg rounded-xl w-[90%] max-w-[800px] max-h-[85vh] overflow-auto shadow-modal border border-night-border">
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
            <p className="text-text-muted">Loading issue details...</p>
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
