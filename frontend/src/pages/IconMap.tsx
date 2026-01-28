import { useRoadmap } from "../hooks/useRoadmap";
import { Map as MapIcon, ChevronRight, Circle } from "lucide-react";
import type { RoadmapIssue } from "../types/api";

function NodeIcon({ type, status }: { type?: string; status: string }) {
  const isComplete = status === "done" || status === "completed";
  const baseClasses = "w-4 h-4";

  if (type === "epic") {
    return (
      <div
        className={`${baseClasses} rounded bg-neon-magenta/30 border border-neon-magenta flex items-center justify-center ${isComplete ? "opacity-50" : ""}`}
      >
        <span className="text-[8px] text-neon-magenta font-bold">E</span>
      </div>
    );
  }

  return (
    <Circle
      className={`${baseClasses} ${isComplete ? "text-neon-green fill-neon-green/30" : "text-neon-cyan"}`}
      strokeWidth={2}
    />
  );
}

function RoadmapNode({
  node,
  depth = 0,
}: {
  node: RoadmapIssue;
  depth?: number;
}) {
  const hasChildren = (node.children?.length ?? 0) > 0;
  const isComplete = node.status === "done" || node.status === "completed";

  return (
    <div className="relative">
      {/* Connection line from parent */}
      {depth > 0 && (
        <div
          className="absolute left-0 top-0 w-6 border-l-2 border-b-2 border-night-border rounded-bl-lg"
          style={{ height: "16px", marginLeft: "-12px" }}
        />
      )}

      {/* Node row */}
      <div
        className={`flex items-center gap-2 py-1.5 px-2 rounded hover:bg-night-surface-bright transition-colors cursor-pointer group ${isComplete ? "opacity-60" : ""}`}
        style={{ marginLeft: depth * 24 }}
      >
        {hasChildren && (
          <ChevronRight className="w-3 h-3 text-text-muted group-hover:text-neon-cyan transition-colors" />
        )}
        {!hasChildren && <span className="w-3" />}

        <NodeIcon type={node.issue_type} status={node.status} />

        <span
          className={`text-sm ${isComplete ? "text-text-muted line-through" : "text-text-primary"}`}
        >
          {node.title}
        </span>

        <span className="text-xs text-text-muted ml-auto opacity-0 group-hover:opacity-100 transition-opacity">
          {node.id}
        </span>

        {node.priority === 1 && (
          <span className="text-[10px] px-1.5 py-0.5 bg-red-900/40 text-neon-pink rounded font-medium">
            P1
          </span>
        )}
      </div>

      {/* Children */}
      {hasChildren && (
        <div className="relative">
          {/* Vertical connector line */}
          <div
            className="absolute border-l-2 border-night-border"
            style={{
              left: depth * 24 + 12,
              top: 0,
              height: "calc(100% - 16px)",
            }}
          />
          {node.children?.map((child) => (
            <RoadmapNode key={child.id} node={child} depth={depth + 1} />
          ))}
        </div>
      )}
    </div>
  );
}

export function IconMap() {
  const { data, isLoading, error } = useRoadmap();

  if (isLoading) {
    return (
      <div className="bg-card-bg rounded-xl p-8 border border-neon-magenta">
        <div className="flex items-center gap-3 mb-4">
          <MapIcon className="w-6 h-6 text-neon-cyan" />
          <h2 className="text-neon-cyan text-2xl font-bold drop-shadow-[0_0_8px_rgba(125,207,255,0.3)]">
            Roadmap
          </h2>
        </div>
        <p className="text-text-muted text-center py-8">Loading roadmap...</p>
      </div>
    );
  }

  if (error) {
    return (
      <div className="bg-card-bg rounded-xl p-8 border border-neon-magenta">
        <div className="flex items-center gap-3 mb-4">
          <MapIcon className="w-6 h-6 text-neon-cyan" />
          <h2 className="text-neon-cyan text-2xl font-bold drop-shadow-[0_0_8px_rgba(125,207,255,0.3)]">
            Roadmap
          </h2>
        </div>
        <p className="text-neon-pink text-center py-8">
          Error loading roadmap data
        </p>
      </div>
    );
  }

  const issues = data?.issues || [];

  return (
    <div className="bg-card-bg rounded-xl p-8 border border-neon-magenta">
      {/* Header */}
      <div className="flex justify-between items-center mb-6">
        <div className="flex items-center gap-3">
          <MapIcon className="w-6 h-6 text-neon-cyan drop-shadow-[0_0_6px_rgba(125,207,255,0.5)]" />
          <div>
            <h2 className="text-neon-cyan text-2xl font-bold drop-shadow-[0_0_8px_rgba(125,207,255,0.3)]">
              Roadmap
            </h2>
            <p className="text-text-secondary text-sm">
              {data?.openCount || 0} open · {data?.completedCount || 0}{" "}
              completed
            </p>
          </div>
        </div>

        {/* Legend */}
        <div className="flex items-center gap-4 text-xs text-text-muted">
          <div className="flex items-center gap-1.5">
            <div className="w-3 h-3 rounded bg-neon-magenta/30 border border-neon-magenta" />
            <span>Epic</span>
          </div>
          <div className="flex items-center gap-1.5">
            <Circle className="w-3 h-3 text-neon-cyan" strokeWidth={2} />
            <span>Task</span>
          </div>
          <div className="flex items-center gap-1.5">
            <Circle
              className="w-3 h-3 text-neon-green fill-neon-green/30"
              strokeWidth={2}
            />
            <span>Done</span>
          </div>
        </div>
      </div>

      {/* Tree */}
      {issues.length === 0 ? (
        <div className="text-center py-12 text-text-muted">
          <p className="text-lg mb-2">No beads yet</p>
          <p className="text-sm">
            Create your first bead with{" "}
            <code className="bg-night-surface px-2 py-0.5 rounded text-neon-cyan">
              bd create "Your task"
            </code>
          </p>
        </div>
      ) : (
        <div className="space-y-0.5">
          {issues.map((issue) => (
            <RoadmapNode key={issue.id} node={issue} />
          ))}
        </div>
      )}
    </div>
  );
}
