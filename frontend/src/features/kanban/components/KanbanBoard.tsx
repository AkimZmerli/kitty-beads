import { useState } from "react";
import { useKanban } from "../hooks/useKanban";
import { KanbanLane } from "./KanbanLane";
import { IssueModal } from "./IssueModal";
import { LANES } from "../constants";
import type { KanbanProps } from "../types";

export function KanbanBoard({ featureId }: KanbanProps) {
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
