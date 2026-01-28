import type { KanbanLaneProps } from "../types";
import { KanbanCard } from "./KanbanCard";

export function KanbanLane({ config, issues, onCardClick }: KanbanLaneProps) {
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
