import { useState, useMemo } from "react";
import { GitBranch, ChevronDown } from "lucide-react";
import { useRoadmap } from "../hooks/useRoadmap";
import { TreeGraph } from "../features/tree-graph";
import type { RoadmapIssue } from "../types/api";

export function TreeView() {
  const { data, isLoading, error } = useRoadmap();
  const [selectedEpicId, setSelectedEpicId] = useState<string | null>(null);

  // Get list of epics (top-level issues or those with children)
  const epics = useMemo(() => {
    if (!data?.issues) return [];

    const allEpics: RoadmapIssue[] = [];

    function findEpics(issues: RoadmapIssue[]) {
      for (const issue of issues) {
        const isEpic =
          issue.issue_type === "epic" ||
          (issue.children && issue.children.length > 0);
        if (isEpic) {
          allEpics.push(issue);
        }
        // Also check nested children for epics
        if (issue.children) {
          findEpics(issue.children);
        }
      }
    }

    findEpics(data.issues);
    return allEpics;
  }, [data?.issues]);

  // Auto-select first epic if none selected
  const selectedEpic = useMemo(() => {
    if (selectedEpicId) {
      return epics.find((e) => e.id === selectedEpicId) || null;
    }
    return epics[0] || null;
  }, [epics, selectedEpicId]);

  if (isLoading) {
    return (
      <div className="bg-card-bg rounded-xl p-8 border border-neon-magenta">
        <div className="flex items-center gap-3 mb-4">
          <GitBranch className="w-6 h-6 text-neon-cyan" />
          <h2 className="text-neon-cyan text-2xl font-bold drop-shadow-[0_0_8px_rgba(125,207,255,0.3)]">
            Tree View
          </h2>
        </div>
        <p className="text-text-muted text-center py-8">Loading tree data...</p>
      </div>
    );
  }

  if (error) {
    return (
      <div className="bg-card-bg rounded-xl p-8 border border-neon-magenta">
        <div className="flex items-center gap-3 mb-4">
          <GitBranch className="w-6 h-6 text-neon-cyan" />
          <h2 className="text-neon-cyan text-2xl font-bold drop-shadow-[0_0_8px_rgba(125,207,255,0.3)]">
            Tree View
          </h2>
        </div>
        <p className="text-neon-pink text-center py-8">
          Error loading tree data
        </p>
      </div>
    );
  }

  return (
    <div className="bg-card-bg rounded-xl p-8 border border-neon-magenta">
      {/* Header */}
      <div className="flex justify-between items-center mb-6">
        <div className="flex items-center gap-3">
          <GitBranch className="w-6 h-6 text-neon-cyan drop-shadow-[0_0_6px_rgba(125,207,255,0.5)]" />
          <div>
            <h2 className="text-neon-cyan text-2xl font-bold drop-shadow-[0_0_8px_rgba(125,207,255,0.3)]">
              Tree View
            </h2>
            <p className="text-text-secondary text-sm">
              Visualize epic hierarchy
            </p>
          </div>
        </div>

        {/* Epic Selector */}
        <div className="relative">
          <select
            value={selectedEpic?.id || ""}
            onChange={(e) => setSelectedEpicId(e.target.value)}
            className="appearance-none bg-night-surface border border-night-border rounded-lg px-4 py-2 pr-10 text-text-primary cursor-pointer focus:outline-none focus:border-neon-cyan hover:bg-night-surface-bright transition-colors"
          >
            {epics.length === 0 && <option value="">No epics available</option>}
            {epics.map((epic) => (
              <option key={epic.id} value={epic.id}>
                {epic.id}: {epic.title}
              </option>
            ))}
          </select>
          <ChevronDown className="absolute right-3 top-1/2 -translate-y-1/2 w-4 h-4 text-text-muted pointer-events-none" />
        </div>
      </div>

      {/* Legend */}
      <div className="flex items-center gap-6 mb-6 text-xs text-text-muted">
        <div className="flex items-center gap-2">
          <div className="w-4 h-4 rounded bg-purple-900/30 border border-neon-magenta" />
          <span>Epic</span>
        </div>
        <div className="flex items-center gap-2">
          <div className="w-4 h-4 rounded bg-cyan-900/30 border border-neon-cyan" />
          <span>Task</span>
        </div>
        <div className="flex items-center gap-2">
          <div className="w-4 h-4 rounded bg-green-900/30 border border-neon-green" />
          <span>Done</span>
        </div>
        <div className="flex items-center gap-2">
          <div className="w-4 h-4 rounded bg-yellow-900/30 border border-neon-yellow" />
          <span>In Progress</span>
        </div>
      </div>

      {/* Tree Graph */}
      {epics.length === 0 ? (
        <div className="text-center py-12 text-text-muted">
          <p className="text-lg mb-2">No epics yet</p>
          <p className="text-sm">
            Create your first epic with{" "}
            <code className="bg-night-surface px-2 py-0.5 rounded text-neon-cyan">
              bd create --type epic "Your epic"
            </code>
          </p>
        </div>
      ) : (
        <TreeGraph rootIssue={selectedEpic} />
      )}
    </div>
  );
}
