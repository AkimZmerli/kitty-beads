import { useQuery } from "@tanstack/react-query";
import { getAllIssues } from "../../../lib/api";
import type { ActivityEvent, ActivityType } from "../types";
import type { Issue } from "../../../types/api";

// Generate activity events from issues
// In a real implementation, this would come from an activity API endpoint
function generateActivityFromIssues(issues: Issue[]): ActivityEvent[] {
  const events: ActivityEvent[] = [];

  for (const issue of issues) {
    // Issue creation event
    if (issue.created_at) {
      events.push({
        id: `${issue.id}-created`,
        type: "issue_created",
        issueId: issue.id,
        issueTitle: issue.title,
        actor: issue.owner || "system",
        timestamp: issue.created_at,
      });
    }

    // Issue update event (if updated after creation)
    if (issue.updated_at && issue.updated_at !== issue.created_at) {
      events.push({
        id: `${issue.id}-updated-${issue.updated_at}`,
        type: "issue_updated",
        issueId: issue.id,
        issueTitle: issue.title,
        actor: issue.assignee || issue.owner || "system",
        timestamp: issue.updated_at,
      });
    }

    // Plan edited event (if has design content)
    if (issue.design && issue.updated_at) {
      events.push({
        id: `${issue.id}-plan-${issue.updated_at}`,
        type: "plan_edited",
        issueId: issue.id,
        issueTitle: issue.title,
        actor: issue.assignee || issue.owner || "system",
        timestamp: issue.updated_at,
      });
    }

    // Status change events based on current status
    if (issue.status && issue.status !== "open") {
      events.push({
        id: `${issue.id}-status-${issue.status}`,
        type: "status_changed",
        issueId: issue.id,
        issueTitle: issue.title,
        actor: issue.assignee || "system",
        timestamp: issue.updated_at || issue.created_at || new Date().toISOString(),
        details: {
          field: "status",
          newValue: issue.status,
        },
      });
    }

    // Assignee change event
    if (issue.assignee) {
      events.push({
        id: `${issue.id}-assignee`,
        type: "assignee_changed",
        issueId: issue.id,
        issueTitle: issue.title,
        timestamp: issue.updated_at || issue.created_at || new Date().toISOString(),
        details: {
          field: "assignee",
          newValue: issue.assignee,
        },
      });
    }
  }

  // Sort by timestamp descending (most recent first)
  events.sort((a, b) => {
    const dateA = new Date(a.timestamp).getTime();
    const dateB = new Date(b.timestamp).getTime();
    return dateB - dateA;
  });

  // Deduplicate by keeping only unique events
  const seen = new Set<string>();
  return events.filter((event) => {
    if (seen.has(event.id)) return false;
    seen.add(event.id);
    return true;
  });
}

interface UseActivityFeedOptions {
  types?: ActivityType[];
  issueId?: string;
  limit?: number;
}

export function useActivityFeed({
  types,
  issueId,
  limit = 50,
}: UseActivityFeedOptions = {}) {
  return useQuery({
    queryKey: ["activity", types, issueId, limit],
    queryFn: async () => {
      const issues = await getAllIssues();
      let events = generateActivityFromIssues(issues);

      // Apply filters
      if (types && types.length > 0) {
        events = events.filter((e) => types.includes(e.type));
      }
      if (issueId) {
        events = events.filter((e) => e.issueId === issueId);
      }

      // Apply limit
      return events.slice(0, limit);
    },
    refetchInterval: 10000, // Refresh every 10 seconds
  });
}
