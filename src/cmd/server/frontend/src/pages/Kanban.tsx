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

// Tokyo Night themed lane colors with subtle tints
const LANES: LaneConfig[] = [
  {
    key: "planned",
    title: "Planned",
    bgColor: "rgba(125, 207, 255, 0.08)",
    borderColor: "#7dcfff",
  },
  {
    key: "doing",
    title: "In Progress",
    bgColor: "rgba(224, 175, 104, 0.08)",
    borderColor: "#e0af68",
  },
  {
    key: "for_review",
    title: "Review",
    bgColor: "rgba(122, 162, 247, 0.08)",
    borderColor: "#7aa2f7",
  },
  {
    key: "done",
    title: "Done",
    bgColor: "rgba(158, 206, 106, 0.08)",
    borderColor: "#9ece6a",
  },
];

const PRIORITY_COLORS = ["#f7768e", "#ff9e64", "#e0af68", "#9ece6a", "#565f89"];

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
      className="bg-night-surface p-3 rounded-lg cursor-pointer transition-all hover:-translate-y-0.5 hover:shadow-[0_4px_12px_rgba(0,0,0,0.3)]"
      style={{
        borderLeft: `4px solid ${priorityColor}`,
      }}
    >
      <div className="flex justify-between items-center text-xs text-text-muted mb-1">
        <span>{issue.id}</span>
        {issue.is_blocked && (
          <span className="text-neon-pink font-semibold text-[0.7rem]">
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

export function Kanban({ featureId }: KanbanProps) {
  const { data, isLoading, error } = useKanban(featureId);
  const [selectedIssue, setSelectedIssue] = useState<string | null>(null);

  if (!featureId) {
    return (
      <div className="bg-card-bg rounded-xl p-8 border border-neon-magenta">
        <h2 className="text-neon-cyan text-2xl font-bold mb-4 drop-shadow-[0_0_8px_rgba(125,207,255,0.3)]">
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
      <div className="bg-card-bg rounded-xl p-8 border border-neon-magenta">
        <h2 className="text-neon-cyan text-2xl font-bold mb-4 drop-shadow-[0_0_8px_rgba(125,207,255,0.3)]">
          Kanban Board
        </h2>
        <p className="text-text-muted text-center py-8">Loading...</p>
      </div>
    );
  }

  if (error) {
    return (
      <div className="bg-card-bg rounded-xl p-8 border border-neon-magenta">
        <h2 className="text-neon-cyan text-2xl font-bold mb-4 drop-shadow-[0_0_8px_rgba(125,207,255,0.3)]">
          Kanban Board
        </h2>
        <p className="text-neon-pink text-center py-8">
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
    <div className="bg-card-bg rounded-xl p-8 border border-neon-magenta">
      <h2 className="text-neon-cyan text-2xl font-bold mb-4 drop-shadow-[0_0_8px_rgba(125,207,255,0.3)]">
        Kanban Board
      </h2>
      <div className="text-text-secondary text-lg mb-6">
        <span className="text-neon-green">{doneCount}</span>/{totalIssues} tasks
        completed
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
