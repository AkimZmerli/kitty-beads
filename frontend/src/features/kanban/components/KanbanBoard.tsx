import { useState, useCallback } from "react";
import { useKanban } from "../hooks/useKanban";
import { useDragDrop } from "../hooks/useDragDrop";
import { useKanbanMutation } from "../hooks/useKanbanMutation";
import { KanbanLane } from "./KanbanLane";
import { IssueModal } from "./IssueModal";
import { LANES } from "../constants";
import type { KanbanProps } from "../types";
import type { IssueCard, KanbanLanes } from "../../../types/api";

export function KanbanBoard({ featureId }: KanbanProps) {
  const { data, isLoading, error } = useKanban(featureId);
  const [selectedIssue, setSelectedIssue] = useState<string | null>(null);

  const mutation = useKanbanMutation(featureId);

  const handleMoveCard = useCallback(
    (card: IssueCard, fromColumn: keyof KanbanLanes, toColumn: keyof KanbanLanes) => {
      mutation.mutate({ card, fromColumn, toColumn });
    },
    [mutation]
  );

  const [dragState, dragHandlers] = useDragDrop(handleMoveCard);

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
      <div className="flex justify-between items-center mb-6">
        <div className="text-text-secondary text-lg">
          <span className="text-neon-green">{doneCount}</span>/{totalIssues} tasks
          completed
        </div>
        {mutation.isPending && (
          <span className="text-text-muted text-sm animate-pulse">Saving...</span>
        )}
      </div>
      <div className="grid grid-cols-4 gap-5">
        {LANES.map((config) => (
          <KanbanLane
            key={config.key}
            config={config}
            issues={lanes[config.key]}
            onCardClick={(issue) => setSelectedIssue(issue.id)}
            onDragStart={(card) => dragHandlers.onDragStart(card, config.key)}
            onDragEnd={dragHandlers.onDragEnd}
            onDragOver={() => dragHandlers.onDragOver(config.key)}
            onDragLeave={dragHandlers.onDragLeave}
            onDrop={(e) => {
              e.preventDefault();
              dragHandlers.onDrop(config.key);
            }}
            isDropTarget={
              dragState.isDragging &&
              dragState.targetColumn === config.key &&
              dragState.sourceColumn !== config.key
            }
            draggedCardId={dragState.draggedCard?.id}
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
