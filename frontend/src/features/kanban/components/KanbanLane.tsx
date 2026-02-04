import type { KanbanLaneProps } from "../types";
import { KanbanCard } from "./KanbanCard";

export function KanbanLane({
  config,
  issues,
  onCardClick,
  onDragStart,
  onDragEnd,
  onDragOver,
  onDragLeave,
  onDrop,
  isDropTarget,
  draggedCardId,
}: KanbanLaneProps) {
  const handleDragOver = (e: React.DragEvent) => {
    e.preventDefault();
    onDragOver?.(e);
  };

  return (
    <div
      onDragOver={handleDragOver}
      onDragLeave={onDragLeave}
      onDrop={onDrop}
      className={`
        rounded-xl p-4 min-h-[300px] transition-all duration-200
        ${isDropTarget ? "ring-2 ring-neon-cyan ring-opacity-50 scale-[1.01]" : ""}
      `}
      style={{
        background: isDropTarget
          ? `linear-gradient(180deg, rgba(125, 207, 255, 0.15) 0%, ${config.bgColor} 100%)`
          : config.bgColor,
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
            onDragStart={(e) => {
              e.dataTransfer.effectAllowed = "move";
              onDragStart?.(issue);
            }}
            onDragEnd={() => onDragEnd?.()}
            isDragging={draggedCardId === issue.id}
          />
        ))}
      </div>

      {/* Drop indicator when empty */}
      {issues.length === 0 && isDropTarget && (
        <div className="border-2 border-dashed border-neon-cyan/30 rounded-lg p-4 text-center text-text-muted text-sm">
          Drop here
        </div>
      )}
    </div>
  );
}
