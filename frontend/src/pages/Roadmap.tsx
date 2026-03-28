import { useEffect, useMemo } from "react";
import { useRoadmap } from "../hooks/useRoadmap";
import {
  RoadmapProvider,
  useRoadmapContext,
  BeadCard,
  RoadmapFilterBar,
  GroupedView,
  ContextMenu,
  IssueModal,
  useKeyboardNavigation,
} from "../features/roadmap";
import { flattenIssueTree } from "../features/roadmap/utils/grouping";
import type { RoadmapIssue } from "../types/api";
import type { FilterType } from "../features/roadmap/types";

function RoadmapContent() {
  // Register keyboard handlers
  useKeyboardNavigation();

  const { data, isLoading, error } = useRoadmap();
  const {
    activeFilter,
    searchQuery,
    groupBy,
    expandedIds,
    modalIssueId,
    closeModal,
    setFlattenedIssues,
  } = useRoadmapContext();

  // Filter and search issues
  const filteredIssues = useMemo(() => {
    if (!data?.issues) return [];

    let issues = data.issues;

    // Apply filter
    issues = applyFilter(issues, activeFilter);

    // Apply search
    if (searchQuery.trim()) {
      issues = applySearch(issues, searchQuery);
    }

    return issues;
  }, [data?.issues, activeFilter, searchQuery]);

  // Update flattened issues for keyboard navigation
  useEffect(() => {
    const flattened = flattenIssueTree(filteredIssues, expandedIds);
    setFlattenedIssues(flattened);
  }, [filteredIssues, expandedIds, setFlattenedIssues]);

  if (isLoading) {
    return (
      <div className="bg-card-bg rounded-xl p-8 border border-neon-magenta">
        <h2 className="text-neon-cyan text-2xl font-bold mb-4 drop-shadow-[0_0_8px_rgba(125,207,255,0.3)]">
          Roadmap
        </h2>
        <div className="flex items-center justify-center py-12">
          <div className="w-8 h-8 border-2 border-neon-cyan border-t-transparent rounded-full animate-spin" />
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="bg-card-bg rounded-xl p-8 border border-neon-magenta">
        <h2 className="text-neon-cyan text-2xl font-bold mb-4 drop-shadow-[0_0_8px_rgba(125,207,255,0.3)]">
          Roadmap
        </h2>
        <p className="text-neon-pink text-center py-8">
          Error loading roadmap data
        </p>
      </div>
    );
  }

  return (
    <div className="bg-card-bg rounded-xl p-8 border border-neon-magenta">
      {/* Header */}
      <div className="flex justify-between items-center mb-6">
        <div>
          <h2 className="text-neon-cyan text-2xl font-bold drop-shadow-[0_0_8px_rgba(125,207,255,0.3)]">
            Roadmap
          </h2>
          <p className="text-text-secondary">
            {data?.openCount || 0} open · {data?.completedCount || 0} completed
            · {data?.totalCount || 0} total
          </p>
        </div>
        {/* Keyboard shortcuts hint */}
        <div className="text-xs text-text-muted hidden md:flex items-center gap-4">
          <span>
            <kbd className="bg-night-surface px-1.5 py-0.5 rounded">j</kbd>/
            <kbd className="bg-night-surface px-1.5 py-0.5 rounded">k</kbd>{" "}
            navigate
          </span>
          <span>
            <kbd className="bg-night-surface px-1.5 py-0.5 rounded">Enter</kbd>{" "}
            expand
          </span>
          <span>
            <kbd className="bg-night-surface px-1.5 py-0.5 rounded">o</kbd> open
          </span>
          <span>
            <kbd className="bg-night-surface px-1.5 py-0.5 rounded">Space</kbd>{" "}
            menu
          </span>
        </div>
      </div>

      {/* Filter bar */}
      <RoadmapFilterBar />

      {/* Beads list */}
      {filteredIssues.length === 0 ? (
        <div className="text-center py-12 text-text-muted">
          {searchQuery || activeFilter !== "all" ? (
            <>
              <p className="text-lg mb-2">No matching beads</p>
              <p className="text-sm">
                Try adjusting your filters or search query
              </p>
            </>
          ) : (
            <>
              <p className="text-lg mb-2">No beads yet</p>
              <p className="text-sm">
                Create your first bead with{" "}
                <code className="bg-night-surface px-2 py-0.5 rounded text-neon-cyan">
                  bd create "Your task"
                </code>
              </p>
            </>
          )}
        </div>
      ) : groupBy === "none" ? (
        <div className="space-y-1">
          {filteredIssues.map((issue, index) => (
            <BeadCard
              key={issue.id}
              issue={issue}
              isLastChild={index === filteredIssues.length - 1}
            />
          ))}
        </div>
      ) : (
        <GroupedView issues={filteredIssues} groupBy={groupBy} />
      )}

      {/* Context menu (portal) */}
      <ContextMenu />

      {/* Issue modal */}
      {modalIssueId && (
        <IssueModal issueId={modalIssueId} onClose={closeModal} />
      )}
    </div>
  );
}

export function Roadmap() {
  return (
    <RoadmapProvider>
      <RoadmapContent />
    </RoadmapProvider>
  );
}

/**
 * Apply filter to issue tree
 */
function applyFilter(
  issues: RoadmapIssue[],
  filter: FilterType,
): RoadmapIssue[] {
  if (filter === "all") return issues;

  const filterIssue = (issue: RoadmapIssue): RoadmapIssue | null => {
    // Check if this issue matches the filter
    let matches = false;
    switch (filter) {
      case "epics":
        matches =
          issue.issue_type === "epic" || (issue.children?.length ?? 0) > 0;
        break;
      case "open":
        matches = issue.status === "open" || issue.status === "in_progress";
        break;
      case "p1":
        matches = issue.priority === 1;
        break;
    }

    // Filter children recursively
    const filteredChildren = issue.children
      ?.map(filterIssue)
      .filter((c): c is RoadmapIssue => c !== null);

    // Include issue if it matches or has matching children
    if (matches || (filteredChildren && filteredChildren.length > 0)) {
      return {
        ...issue,
        children: filteredChildren,
      };
    }

    return null;
  };

  return issues.map(filterIssue).filter((i): i is RoadmapIssue => i !== null);
}

/**
 * Apply search to issue tree
 */
function applySearch(issues: RoadmapIssue[], query: string): RoadmapIssue[] {
  const lowerQuery = query.toLowerCase();

  const searchIssue = (issue: RoadmapIssue): RoadmapIssue | null => {
    // Check if this issue matches the search
    const matches =
      issue.title.toLowerCase().includes(lowerQuery) ||
      issue.id.toLowerCase().includes(lowerQuery) ||
      issue.description?.toLowerCase().includes(lowerQuery);

    // Search children recursively
    const filteredChildren = issue.children
      ?.map(searchIssue)
      .filter((c): c is RoadmapIssue => c !== null);

    // Include issue if it matches or has matching children
    if (matches || (filteredChildren && filteredChildren.length > 0)) {
      return {
        ...issue,
        children: filteredChildren,
      };
    }

    return null;
  };

  return issues.map(searchIssue).filter((i): i is RoadmapIssue => i !== null);
}
