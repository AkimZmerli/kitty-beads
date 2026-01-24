import { useState } from "react";
import type { RoadmapIssue } from "../../types/api";
import {
  parsePlan,
  formatCriteriaProgress,
  getPriorityColor,
  getPriorityIndicator,
} from "../../lib/planParser";

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
          ${isExpanded ? "bg-gray-100" : "bg-white hover:bg-gray-50"}
          ${depth > 0 ? "ml-6 border-l-2 border-gray-200" : ""}
        `}
        style={{ borderLeftColor: depth > 0 ? priorityColor : undefined }}
        onClick={() => setIsExpanded(!isExpanded)}
      >
        {/* Expand/collapse chevron */}
        <span
          className={`text-gray-400 transition-transform duration-150 ${
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
        <span className="text-gray-500 font-mono text-sm">{issue.id}</span>

        {/* Title */}
        <span className="font-medium text-gray-900 flex-1 truncate">
          {issue.title}
        </span>

        {/* Type badge */}
        {isEpic && (
          <span className="px-2 py-0.5 bg-purple-100 text-purple-700 text-xs rounded-full font-medium">
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
          <span className="px-2 py-0.5 bg-red-100 text-red-700 text-xs rounded-full font-medium">
            BLOCKED
          </span>
        )}
      </div>

      {/* Expanded view */}
      {isExpanded && (
        <div
          className={`
            mt-1 px-4 py-4 bg-gray-50 rounded-lg
            ${depth > 0 ? "ml-6" : ""}
          `}
        >
          {/* Summary */}
          {parsedPlan.summary && (
            <p className="text-gray-700 mb-4">{parsedPlan.summary}</p>
          )}

          {/* Stats row */}
          <div className="flex items-center gap-4 text-sm text-gray-500 mb-4">
            {/* Acceptance criteria */}
            {parsedPlan.acceptanceCriteriaCount > 0 && (
              <div className="flex items-center gap-1">
                <span>✓</span>
                <span>{formatCriteriaProgress(parsedPlan)}</span>
              </div>
            )}

            {/* Open questions */}
            {parsedPlan.openQuestions.length > 0 && (
              <div className="flex items-center gap-1">
                <span>?</span>
                <span>{parsedPlan.openQuestions.length} questions</span>
              </div>
            )}

            {/* Children count */}
            {hasChildren && (
              <div className="flex items-center gap-1">
                <span>📋</span>
                <span>{issue.children!.length} subtasks</span>
              </div>
            )}
          </div>

          {/* Open questions preview */}
          {parsedPlan.openQuestions.length > 0 && (
            <div className="mb-4">
              <div className="text-xs uppercase tracking-wider text-gray-400 mb-2">
                Open Questions
              </div>
              <ul className="text-sm text-gray-600 space-y-1">
                {parsedPlan.openQuestions.slice(0, 3).map((q, i) => (
                  <li key={i} className="flex items-start gap-2">
                    <span className="text-amber-500">?</span>
                    <span className="truncate">{q}</span>
                  </li>
                ))}
                {parsedPlan.openQuestions.length > 3 && (
                  <li className="text-gray-400 italic">
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
              className="px-3 py-1.5 bg-grassy-green text-white text-sm rounded-md
                hover:bg-green-700 transition-colors"
            >
              Edit Plan
            </button>
            <button
              onClick={(e) => {
                e.stopPropagation();
                onViewFull?.(issue.id);
              }}
              className="px-3 py-1.5 bg-gray-200 text-gray-700 text-sm rounded-md
                hover:bg-gray-300 transition-colors"
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
      return { label: "OPEN", className: "bg-blue-100 text-blue-700" };
    case "in_progress":
      return { label: "IN PROGRESS", className: "bg-yellow-100 text-yellow-700" };
    case "closed":
      return { label: "DONE", className: "bg-green-100 text-green-700" };
    case "blocked":
      return { label: "BLOCKED", className: "bg-red-100 text-red-700" };
    default:
      return { label: status.toUpperCase(), className: "bg-gray-100 text-gray-700" };
  }
}
