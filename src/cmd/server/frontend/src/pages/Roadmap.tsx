import { useState } from "react";
import { useRoadmap } from "../hooks/useRoadmap";
import { BeadCard } from "../components/roadmap/BeadCard";
import { MarkdownViewer } from "../components/MarkdownViewer";
import { getIssue } from "../lib/api";
import type { Issue } from "../types/api";

export function Roadmap() {
  const { data, isLoading, error } = useRoadmap();
  const [viewingIssue, setViewingIssue] = useState<Issue | null>(null);
  const [loadingIssue, setLoadingIssue] = useState(false);

  const handleEditPlan = (issueId: string) => {
    // TODO: Navigate to Ideation Pad (Phase 2)
    console.log("Edit plan for:", issueId);
    alert(`Edit Plan: ${issueId}\n\nIdeation Pad coming in Phase 2!`);
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
      <div className="bg-white rounded-xl p-8 border-2 border-sunny-yellow-border">
        <h2 className="text-grassy-green text-2xl font-bold mb-4">Roadmap</h2>
        <p className="text-text-muted text-center py-8">Loading roadmap...</p>
      </div>
    );
  }

  if (error) {
    return (
      <div className="bg-white rounded-xl p-8 border-2 border-sunny-yellow-border">
        <h2 className="text-grassy-green text-2xl font-bold mb-4">Roadmap</h2>
        <p className="text-red-500 text-center py-8">
          Error loading roadmap data
        </p>
      </div>
    );
  }

  const issues = data?.issues || [];

  return (
    <div className="bg-white rounded-xl p-8 border-2 border-sunny-yellow-border">
      {/* Header */}
      <div className="flex justify-between items-center mb-6">
        <div>
          <h2 className="text-grassy-green text-2xl font-bold">Roadmap</h2>
          <p className="text-text-secondary">
            {data?.openCount || 0} open · {data?.completedCount || 0} completed
            · {data?.totalCount || 0} total
          </p>
        </div>
      </div>

      {/* Filter tabs (future enhancement) */}
      <div className="flex gap-2 mb-6 border-b border-gray-200 pb-4">
        <button className="px-4 py-2 bg-grassy-green text-white rounded-md text-sm font-medium">
          All
        </button>
        <button className="px-4 py-2 bg-gray-100 text-gray-600 rounded-md text-sm font-medium hover:bg-gray-200">
          Epics
        </button>
        <button className="px-4 py-2 bg-gray-100 text-gray-600 rounded-md text-sm font-medium hover:bg-gray-200">
          Open
        </button>
        <button className="px-4 py-2 bg-gray-100 text-gray-600 rounded-md text-sm font-medium hover:bg-gray-200">
          P1 Only
        </button>
      </div>

      {/* Beads list */}
      {issues.length === 0 ? (
        <div className="text-center py-12 text-text-muted">
          <p className="text-lg mb-2">No beads yet</p>
          <p className="text-sm">
            Create your first bead with{" "}
            <code className="bg-gray-100 px-2 py-0.5 rounded">
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
      <div className="absolute inset-0 bg-black/50" onClick={onClose} />
      <div className="relative bg-white rounded-xl w-[90%] max-w-[800px] max-h-[85vh] overflow-auto shadow-modal">
        {/* Header */}
        <div className="sticky top-0 bg-white flex justify-between items-center p-5 border-b border-border z-10">
          <h3 className="text-lg font-semibold">
            {loading ? "Loading..." : `${issue?.id}: ${issue?.title}`}
          </h3>
          <button
            onClick={onClose}
            className="text-2xl text-text-muted hover:text-text-primary"
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
                <span className="bg-red-100 text-red-800 px-3 py-1 rounded-full text-sm font-medium">
                  P{issue.priority}
                </span>
                <span className="bg-blue-100 text-blue-800 px-3 py-1 rounded-full text-sm font-medium">
                  {issue.status || "open"}
                </span>
                {issue.issue_type && (
                  <span className="bg-purple-100 text-purple-800 px-3 py-1 rounded-full text-sm font-medium">
                    {issue.issue_type}
                  </span>
                )}
                {issue.assignee && (
                  <span className="bg-gray-100 text-text-secondary px-3 py-1 rounded-full text-sm">
                    {issue.assignee}
                  </span>
                )}
              </div>

              {/* Description */}
              {issue.description && (
                <div className="mb-6">
                  <h4 className="text-grassy-green font-semibold mb-3">
                    Description
                  </h4>
                  <MarkdownViewer content={issue.description} />
                </div>
              )}

              {/* Design/Plan */}
              {issue.design && (
                <div className="mb-6">
                  <h4 className="text-grassy-green font-semibold mb-3">
                    Design / Plan
                  </h4>
                  <MarkdownViewer content={issue.design} />
                </div>
              )}

              {/* Acceptance Criteria */}
              {issue.acceptance_criteria && (
                <div className="mb-6">
                  <h4 className="text-grassy-green font-semibold mb-3">
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
            <p className="text-red-500">Failed to load issue</p>
          )}
        </div>
      </div>
    </div>
  );
}
