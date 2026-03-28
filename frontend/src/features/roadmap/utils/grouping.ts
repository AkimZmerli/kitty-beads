import type { RoadmapIssue } from "../../../types/api";
import type { GroupByOption, IssueGroup } from "../types";
import { STATUS_OPTIONS, PRIORITY_OPTIONS } from "../constants";

/**
 * Flatten a tree of issues into a flat list
 * Respects expanded state - only includes children of expanded issues
 */
export function flattenIssueTree(
  issues: RoadmapIssue[],
  expandedIds: Set<string>
): RoadmapIssue[] {
  const result: RoadmapIssue[] = [];

  function traverse(items: RoadmapIssue[]) {
    for (const issue of items) {
      result.push(issue);
      if (issue.children?.length && expandedIds.has(issue.id)) {
        traverse(issue.children);
      }
    }
  }

  traverse(issues);
  return result;
}

/**
 * Group issues by a specified field
 */
export function groupIssues(
  issues: RoadmapIssue[],
  groupBy: GroupByOption
): Map<string, RoadmapIssue[]> {
  if (groupBy === "none") {
    return new Map([["all", issues]]);
  }

  const groups = new Map<string, RoadmapIssue[]>();

  // Flatten the tree first for grouping
  const flatIssues = flattenAllIssues(issues);

  for (const issue of flatIssues) {
    let key: string;
    switch (groupBy) {
      case "status":
        key = issue.status || "open";
        break;
      case "priority":
        key = `p${issue.priority}`;
        break;
      case "assignee":
        key = issue.assignee || "Unassigned";
        break;
      default:
        key = "all";
    }

    if (!groups.has(key)) {
      groups.set(key, []);
    }
    groups.get(key)!.push(issue);
  }

  return groups;
}

/**
 * Flatten all issues in tree (ignoring expanded state)
 */
function flattenAllIssues(issues: RoadmapIssue[]): RoadmapIssue[] {
  const result: RoadmapIssue[] = [];

  function traverse(items: RoadmapIssue[]) {
    for (const issue of items) {
      result.push(issue);
      if (issue.children?.length) {
        traverse(issue.children);
      }
    }
  }

  traverse(issues);
  return result;
}

/**
 * Get group metadata (label, color, order)
 */
export function getGroupInfo(
  key: string,
  groupBy: GroupByOption
): IssueGroup {
  switch (groupBy) {
    case "status": {
      const option = STATUS_OPTIONS.find((o) => o.value === key);
      return {
        key,
        label: option?.label || key,
        color: option?.color || "text-muted",
        count: 0,
      };
    }
    case "priority": {
      const priority = parseInt(key.replace("p", ""));
      const option = PRIORITY_OPTIONS.find((o) => o.value === priority);
      return {
        key,
        label: option?.label || `P${priority}`,
        color: option?.color || "text-muted",
        count: 0,
      };
    }
    case "assignee":
      return {
        key,
        label: key,
        color: "neon-cyan",
        count: 0,
      };
    default:
      return {
        key,
        label: "All",
        color: "text-normal",
        count: 0,
      };
  }
}

/**
 * Get sorted group keys based on groupBy type
 */
export function getSortedGroupKeys(
  groups: Map<string, RoadmapIssue[]>,
  groupBy: GroupByOption
): string[] {
  const keys = Array.from(groups.keys());

  switch (groupBy) {
    case "status":
      // Sort by status order
      const statusOrder = STATUS_OPTIONS.map((o) => o.value);
      return keys.sort(
        (a, b) => statusOrder.indexOf(a) - statusOrder.indexOf(b)
      );

    case "priority":
      // Sort by priority number
      return keys.sort((a, b) => {
        const aNum = parseInt(a.replace("p", ""));
        const bNum = parseInt(b.replace("p", ""));
        return aNum - bNum;
      });

    case "assignee":
      // Sort alphabetically, with Unassigned at end
      return keys.sort((a, b) => {
        if (a === "Unassigned") return 1;
        if (b === "Unassigned") return -1;
        return a.localeCompare(b);
      });

    default:
      return keys;
  }
}
