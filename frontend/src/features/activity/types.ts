export type ActivityType =
  | "issue_created"
  | "issue_updated"
  | "status_changed"
  | "plan_edited"
  | "comment_added"
  | "assignee_changed";

export interface ActivityEvent {
  id: string;
  type: ActivityType;
  issueId: string;
  issueTitle: string;
  actor?: string;
  timestamp: string;
  details?: {
    field?: string;
    oldValue?: string;
    newValue?: string;
    comment?: string;
  };
}

export interface ActivityFilters {
  types: ActivityType[];
  issueId?: string;
  actor?: string;
}
