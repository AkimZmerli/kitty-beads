import { X, Edit3, Eye } from "lucide-react";
import type { BeadPopupProps } from "../types";

export function BeadPopup({
  issue,
  position,
  onClose,
  onEditPlan,
  onViewFull,
}: BeadPopupProps) {
  const isEpic = issue.issue_type === "epic" || (issue.children?.length ?? 0) > 0;
  const isDone =
    issue.status === "done" ||
    issue.status === "closed" ||
    issue.status === "completed";

  return (
    <>
      {/* Backdrop */}
      <div
        className="fixed inset-0 z-40"
        onClick={onClose}
      />

      {/* Popup */}
      <div
        className="absolute z-50 w-80 bg-night-surface border border-night-border rounded-lg shadow-xl"
        style={{
          left: position.x,
          top: position.y,
          transform: "translate(-50%, 10px)",
        }}
      >
        {/* Header */}
        <div className="flex items-center justify-between px-4 py-3 border-b border-night-border">
          <div className="flex items-center gap-2">
            <span className="text-text-muted font-mono text-sm">{issue.id}</span>
            {isEpic && (
              <span className="px-2 py-0.5 bg-purple-900/30 text-neon-magenta text-xs rounded-full font-medium">
                EPIC
              </span>
            )}
          </div>
          <button
            onClick={onClose}
            className="text-text-muted hover:text-text-primary transition-colors"
          >
            <X className="w-4 h-4" />
          </button>
        </div>

        {/* Content */}
        <div className="p-4">
          <h3 className="text-text-primary font-medium mb-2">{issue.title}</h3>

          {/* Status & Priority */}
          <div className="flex items-center gap-2 mb-3">
            <span
              className={`px-2 py-0.5 text-xs rounded-full font-medium ${
                isDone
                  ? "bg-green-900/30 text-neon-green"
                  : issue.status === "in_progress"
                  ? "bg-yellow-900/30 text-neon-yellow"
                  : "bg-blue-900/30 text-neon-blue"
              }`}
            >
              {issue.status.toUpperCase().replace("_", " ")}
            </span>
            <span className="text-text-muted text-sm">P{issue.priority}</span>
          </div>

          {/* Description preview */}
          {issue.description && (
            <p className="text-text-secondary text-sm mb-4 line-clamp-2">
              {issue.description}
            </p>
          )}

          {/* Children count */}
          {issue.children && issue.children.length > 0 && (
            <p className="text-text-muted text-sm mb-4">
              {issue.children.length} subtask{issue.children.length > 1 ? "s" : ""}
            </p>
          )}

          {/* Actions */}
          <div className="flex gap-2">
            <button
              onClick={() => onEditPlan(issue.id)}
              className="flex items-center gap-2 px-3 py-2 bg-neon-cyan text-night-bg text-sm rounded-md
                hover:shadow-[0_0_12px_rgba(125,207,255,0.4)] transition-all font-medium flex-1"
            >
              <Edit3 className="w-4 h-4" />
              Edit Plan
            </button>
            <button
              onClick={() => onViewFull(issue.id)}
              className="flex items-center gap-2 px-3 py-2 bg-night-surface-bright text-text-normal text-sm rounded-md
                hover:bg-night-bg-highlight transition-colors"
            >
              <Eye className="w-4 h-4" />
              View
            </button>
          </div>
        </div>
      </div>
    </>
  );
}
