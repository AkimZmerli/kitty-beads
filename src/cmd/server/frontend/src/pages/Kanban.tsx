import { useState, useEffect } from "react";
import { useKanban } from "../hooks/useKanban";
import type { IssueCard } from "../types/api";
import { getIssue } from "../lib/api";
import { MarkdownViewer } from "../components/MarkdownViewer";

interface KanbanProps {
  featureId: string | null;
}

interface LaneConfig {
  key: keyof import("../types/api").KanbanLanes;
  title: string;
  bgColor: string;
  borderColor: string;
}

const LANES: LaneConfig[] = [
  {
    key: "planned",
    title: "Planned",
    bgColor: "#e0f2fe",
    borderColor: "#0284c7",
  },
  {
    key: "doing",
    title: "In Progress",
    bgColor: "#fef3c7",
    borderColor: "#d97706",
  },
  {
    key: "for_review",
    title: "Review",
    bgColor: "#e0e7ff",
    borderColor: "#4f46e5",
  },
  { key: "done", title: "Done", bgColor: "#dcfce7", borderColor: "#16a34a" },
];

const PRIORITY_COLORS = ["#ef4444", "#f97316", "#eab308", "#22c55e", "#6b7280"];

function KanbanCard({
  issue,
  onClick,
}: {
  issue: IssueCard;
  onClick: () => void;
}) {
  const priorityColor = PRIORITY_COLORS[Math.min(issue.priority, 4)];

  return (
    <div
      onClick={onClick}
      className="bg-white p-3 rounded-lg cursor-pointer transition-transform hover:-translate-y-0.5"
      style={{
        borderLeft: `4px solid ${priorityColor}`,
        boxShadow: "0 1px 3px rgba(0,0,0,0.1)",
      }}
    >
      <div className="flex justify-between items-center text-xs text-text-muted mb-1">
        <span>{issue.id}</span>
        {issue.is_blocked && (
          <span className="text-red-500 font-semibold text-[0.7rem]">
            BLOCKED
          </span>
        )}
      </div>
      <div className="font-medium text-text-primary text-sm">{issue.title}</div>
      {issue.assignee && (
        <div className="text-xs text-text-muted mt-2">{issue.assignee}</div>
      )}
    </div>
  );
}

function KanbanLane({
  config,
  issues,
  onCardClick,
}: {
  config: LaneConfig;
  issues: IssueCard[];
  onCardClick: (issue: IssueCard) => void;
}) {
  return (
    <div
      className="rounded-xl p-4 min-h-[300px]"
      style={{
        background: config.bgColor,
        borderTop: `4px solid ${config.borderColor}`,
      }}
    >
      <h3 className="text-sm uppercase tracking-wider text-text-secondary font-semibold mb-4">
        {config.title} <span className="opacity-60">({issues.length})</span>
      </h3>
      <div className="flex flex-col gap-2.5">
        {issues.map((issue) => (
          <KanbanCard
            key={issue.id}
            issue={issue}
            onClick={() => onCardClick(issue)}
          />
        ))}
      </div>
    </div>
  );
}

interface IssueModalProps {
  issueId: string;
  onClose: () => void;
}

function IssueModal({ issueId, onClose }: IssueModalProps) {
  const [issue, setIssue] = useState<import("../types/api").Issue | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    setLoading(true);
    getIssue(issueId)
      .then(setIssue)
      .finally(() => setLoading(false));
  }, [issueId]);

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center">
      <div className="absolute inset-0 bg-black/50" onClick={onClose} />
      <div className="relative bg-white rounded-xl w-[90%] max-w-[700px] max-h-[80vh] overflow-auto shadow-modal">
        <div className="flex justify-between items-center p-5 border-b border-border">
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
        <div className="p-5">
          {loading ? (
            <p className="text-text-muted">Loading issue details...</p>
          ) : issue ? (
            <>
              <div className="flex gap-2 flex-wrap mb-5">
                <span className="bg-red-100 text-red-800 px-3 py-1 rounded-full text-sm font-medium">
                  P{issue.priority}
                </span>
                <span className="bg-blue-100 text-blue-800 px-3 py-1 rounded-full text-sm font-medium">
                  {issue.status || "open"}
                </span>
                {issue.assignee && (
                  <span className="bg-gray-100 text-text-secondary px-3 py-1 rounded-full text-sm">
                    {issue.assignee}
                  </span>
                )}
              </div>
              <MarkdownViewer
                content={issue.description || "*No description provided*"}
              />
              {issue.design && (
                <>
                  <hr className="my-5 border-border" />
                  <h4 className="text-grassy-green font-semibold mb-3">
                    Design
                  </h4>
                  <MarkdownViewer content={issue.design} />
                </>
              )}
              {issue.acceptance_criteria && (
                <>
                  <hr className="my-5 border-border" />
                  <h4 className="text-grassy-green font-semibold mb-3">
                    Acceptance Criteria
                  </h4>
                  <MarkdownViewer content={issue.acceptance_criteria} />
                </>
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

export function Kanban({ featureId }: KanbanProps) {
  const { data, isLoading, error } = useKanban(featureId);
  const [selectedIssue, setSelectedIssue] = useState<string | null>(null);

  if (!featureId) {
    return (
      <div className="bg-white rounded-xl p-8 border-2 border-sunny-yellow-border">
        <h2 className="text-grassy-green text-2xl font-bold mb-4">
          Kanban Board
        </h2>
        <p className="text-text-muted text-center py-8">
          Select a feature to view its kanban board.
        </p>
      </div>
    );
  }

  if (isLoading) {
    return (
      <div className="bg-white rounded-xl p-8 border-2 border-sunny-yellow-border">
        <h2 className="text-grassy-green text-2xl font-bold mb-4">
          Kanban Board
        </h2>
        <p className="text-text-muted text-center py-8">Loading...</p>
      </div>
    );
  }

  if (error) {
    return (
      <div className="bg-white rounded-xl p-8 border-2 border-sunny-yellow-border">
        <h2 className="text-grassy-green text-2xl font-bold mb-4">
          Kanban Board
        </h2>
        <p className="text-red-500 text-center py-8">
          Error loading kanban data
        </p>
      </div>
    );
  }

  const lanes = data?.lanes || {
    planned: [],
    doing: [],
    for_review: [],
    done: [],
  };
  const totalIssues = Object.values(lanes).flat().length;
  const doneCount = lanes.done.length;

  return (
    <div className="bg-white rounded-xl p-8 border-2 border-sunny-yellow-border">
      <h2 className="text-grassy-green text-2xl font-bold mb-4">
        Kanban Board
      </h2>
      <div className="text-text-secondary text-lg mb-6">
        {doneCount}/{totalIssues} tasks completed
      </div>
      <div className="grid grid-cols-4 gap-5">
        {LANES.map((config) => (
          <KanbanLane
            key={config.key}
            config={config}
            issues={lanes[config.key]}
            onCardClick={(issue) => setSelectedIssue(issue.id)}
          />
        ))}
      </div>
      {selectedIssue && (
        <IssueModal
          issueId={selectedIssue}
          onClose={() => setSelectedIssue(null)}
        />
      )}
    </div>
  );
}
