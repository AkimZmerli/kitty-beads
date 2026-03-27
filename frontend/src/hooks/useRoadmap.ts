import { useQuery } from "@tanstack/react-query";
import { getAllIssues } from "../lib/api";
import type { Issue, RoadmapIssue } from "../types/api";

export interface RoadmapData {
  /** All issues organized as a tree */
  issues: RoadmapIssue[];
  /** Total issue count */
  totalCount: number;
  /** Open issue count */
  openCount: number;
  /** Completed issue count */
  completedCount: number;
}

/**
 * Build a tree structure from flat issue list
 * Groups children under their parent issues
 */
function buildIssueTree(issues: Issue[]): RoadmapIssue[] {
  const issueMap = new Map<string, RoadmapIssue>();
  const rootIssues: RoadmapIssue[] = [];

  // First pass: create map of all issues
  for (const issue of issues) {
    issueMap.set(issue.id, { ...issue, children: [] });
  }

  // Second pass: build parent-child relationships
  for (const issue of issues) {
    const roadmapIssue = issueMap.get(issue.id)!;

    // Check if this issue has a parent (ID contains a dot)
    const lastDotIndex = issue.id.lastIndexOf(".");
    if (lastDotIndex > 0) {
      const parentId = issue.id.substring(0, lastDotIndex);
      const parent = issueMap.get(parentId);

      if (parent) {
        parent.children = parent.children || [];
        parent.children.push(roadmapIssue);
      } else {
        // Parent not found, treat as root
        rootIssues.push(roadmapIssue);
      }
    } else {
      // Top-level issue (no dots in ID)
      rootIssues.push(roadmapIssue);
    }
  }

  // Sort children by priority, then by ID
  const sortIssues = (items: RoadmapIssue[]) => {
    items.sort((a, b) => {
      if (a.priority !== b.priority) {
        return a.priority - b.priority;
      }
      return a.id.localeCompare(b.id);
    });
    for (const item of items) {
      if (item.children && item.children.length > 0) {
        sortIssues(item.children);
      }
    }
  };

  sortIssues(rootIssues);

  return rootIssues;
}

/**
 * Hook to fetch and process roadmap data
 */
export function useRoadmap() {
  return useQuery<RoadmapData>({
    queryKey: ["roadmap"],
    queryFn: async () => {
      const issues = await getAllIssues();

      const tree = buildIssueTree(issues);

      const openCount = issues.filter(
        (i) => i.status === "open" || i.status === "in_progress",
      ).length;
      const completedCount = issues.filter((i) => i.status === "closed").length;

      return {
        issues: tree,
        totalCount: issues.length,
        openCount,
        completedCount,
      };
    },
    staleTime: 5000,
    refetchInterval: 30_000,
    refetchIntervalInBackground: false,
  });
}
