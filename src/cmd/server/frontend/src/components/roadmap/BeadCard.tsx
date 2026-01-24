import { useState } from "react";
import type { RoadmapIssue } from "../../types/api";
import {
  parsePlan,
  formatCriteriaProgress,
  getPriorityColor,
  getPriorityIndicator,
} from "../../lib/planParser";
import { CheckCircle2, HelpCircle, ListTodo } from "lucide-react";

interface BeadCardProps {
  issue: RoadmapIssue;
  depth?: number;
  onEditPlan?: (issueId: string) => void;
  onViewFull?: (issueId: string) => void;
}

export function BeadCard({
  issue,
  depth = 0,
  onEditPlan,
  onViewFull,
}: BeadCardProps) {
  const [isExpanded, setIsExpanded] = useState(false);

  const parsedPlan = parsePlan(issue.design || issue.description);
  const priorityColor = getPriorityColor(issue.priority);
  const priorityIndicator = getPriorityIndicator(issue.priority);

  const hasChildren = issue.children && issue.children.length > 0;
  const isEpic = issue.issue_type === "epic" || hasChildren;

  const statusBadge = getStatusBadge(issue.status);

  return (
    <div className="mb-2">
      {/* Collapsed view - one-liner */}
      <div
        className={`
          flex items-center gap-3 px-4 py-3 rounded-lg cursor-pointer
          transition-all duration-150 ease-in-out
          ${isExpanded ? "bg-night-bg-highlight" : "bg-night-surface hover:bg-night-bg-highlight"}
          ${depth > 0 ? "ml-6 border-l-2 border-night-border" : ""}
        `}
        style={{ borderLeftColor: depth > 0 ? priorityColor : undefined }}
        onClick={() => setIsExpanded(!isExpanded)}
      >
        {/* Expand/collapse chevron */}
        <span
          className={`text-text-muted transition-transform duration-150 ${
            isExpanded ? "rotate-90" : ""
          }`}
        >
          {hasChildren || parsedPlan.hasContent ? "▶" : "•"}
        </span>

        {/* Priority indicator */}
        <span style={{ color: priorityColor }} className="font-bold">
          {priorityIndicator} P{issue.priority}
        </span>

        {/* Issue ID */}
        <span className="text-text-muted font-mono text-sm">{issue.id}</span>

        {/* Title */}
        <span className="font-medium text-text-bright flex-1 truncate">
          {issue.title}
        </span>

        {/* Type badge */}
        {isEpic && (
          <span className="px-2 py-0.5 bg-purple-900/30 text-neon-magenta text-xs rounded-full font-medium">
            EPIC
          </span>
        )}

        {/* Status badge */}
        <span
          className={`px-2 py-0.5 text-xs rounded-full font-medium ${statusBadge.className}`}
        >
          {statusBadge.label}
        </span>

        {/* Blocked indicator */}
        {issue.isBlocked && (
          <span className="px-2 py-0.5 bg-red-900/30 text-neon-pink text-xs rounded-full font-medium">
            BLOCKED
          </span>
        )}
      </div>

      {/* Expanded view */}
      {isExpanded && (
        <div
          className={`
            mt-1 px-4 py-4 bg-night-surface-bright rounded-lg
            ${depth > 0 ? "ml-6" : ""}
          `}
        >
          {/* Summary */}
          {parsedPlan.summary && (
            <p className="text-text-normal mb-4">{parsedPlan.summary}</p>
          )}

          {/* Stats row */}
          <div className="flex items-center gap-4 text-sm text-text-muted mb-4">
            {/* Acceptance criteria */}
            {parsedPlan.acceptanceCriteriaCount > 0 && (
              <div className="flex items-center gap-1">
                <CheckCircle2 className="w-4 h-4 text-neon-green" />
                <span>{formatCriteriaProgress(parsedPlan)}</span>
              </div>
            )}

            {/* Open questions */}
            {parsedPlan.openQuestions.length > 0 && (
              <div className="flex items-center gap-1">
                <HelpCircle className="w-4 h-4 text-neon-orange" />
                <span>{parsedPlan.openQuestions.length} questions</span>
              </div>
            )}

            {/* Children count */}
            {hasChildren && (
              <div className="flex items-center gap-1">
                <ListTodo className="w-4 h-4 text-neon-blue" />
                <span>{issue.children!.length} subtasks</span>
              </div>
            )}
          </div>

          {/* Open questions preview */}
          {parsedPlan.openQuestions.length > 0 && (
            <div className="mb-4">
              <div className="text-xs uppercase tracking-wider text-text-muted mb-2">
                Open Questions
              </div>
              <ul className="text-sm text-text-normal space-y-1">
                {parsedPlan.openQuestions.slice(0, 3).map((q, i) => (
                  <li key={i} className="flex items-start gap-2">
                    <span className="text-neon-orange">?</span>
                    <span className="truncate">{q}</span>
                  </li>
                ))}
                {parsedPlan.openQuestions.length > 3 && (
                  <li className="text-text-muted italic">
                    +{parsedPlan.openQuestions.length - 3} more...
                  </li>
                )}
              </ul>
            </div>
          )}

          {/* Action buttons */}
          <div className="flex gap-2">
            <button
              onClick={(e) => {
                e.stopPropagation();
                onEditPlan?.(issue.id);
              }}
              className="px-3 py-1.5 bg-neon-cyan text-night-bg text-sm rounded-md
                hover:shadow-[0_0_12px_rgba(125,207,255,0.4)] transition-all font-medium"
            >
              Edit Plan
            </button>
            <button
              onClick={(e) => {
                e.stopPropagation();
                onViewFull?.(issue.id);
              }}
              className="px-3 py-1.5 bg-night-surface text-text-normal text-sm rounded-md
                hover:bg-night-bg-highlight transition-colors"
            >
              View Full
            </button>
          </div>
        </div>
      )}

      {/* Render children when parent is expanded */}
      {isExpanded && hasChildren && (
        <div className="mt-2">
          {issue.children!.map((child) => (
            <BeadCard
              key={child.id}
              issue={child}
              depth={depth + 1}
              onEditPlan={onEditPlan}
              onViewFull={onViewFull}
            />
          ))}
        </div>
      )}
    </div>
  );
}

function getStatusBadge(status: string): { label: string; className: string } {
  switch (status) {
    case "open":
      return { label: "OPEN", className: "bg-blue-900/30 text-neon-blue" };
    case "in_progress":
      return { label: "IN PROGRESS", className: "bg-yellow-900/30 text-neon-yellow" };
    case "closed":
      return { label: "DONE", className: "bg-green-900/30 text-neon-green" };
    case "blocked":
      return { label: "BLOCKED", className: "bg-red-900/30 text-neon-pink" };
    default:
      return { label: status.toUpperCase(), className: "bg-night-surface-bright text-text-normal" };
  }
}
