import { GripVertical } from "lucide-react";
import type { KanbanCardProps } from "../types";
import { PRIORITY_COLORS } from "../constants";

export function KanbanCard({
  issue,
  onClick,
  onDragStart,
  onDragEnd,
  isDragging,
}: KanbanCardProps) {
  const priorityColor = PRIORITY_COLORS[Math.min(issue.priority, 4)];

  return (
    <div
      draggable
      onDragStart={onDragStart}
      onDragEnd={onDragEnd}
      onClick={onClick}
      className={`
        bg-night-surface p-3 rounded-lg cursor-grab transition-all
        hover:-translate-y-0.5 hover:shadow-[0_4px_12px_rgba(0,0,0,0.3)]
        active:cursor-grabbing
        ${isDragging ? "opacity-50 scale-95" : ""}
      `}
      style={{
        borderLeft: `4px solid ${priorityColor}`,
      }}
    >
      <div className="flex justify-between items-center text-xs text-text-muted mb-1">
        <div className="flex items-center gap-1">
          <GripVertical className="w-3 h-3 opacity-50" />
          <span>{issue.id}</span>
        </div>
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
