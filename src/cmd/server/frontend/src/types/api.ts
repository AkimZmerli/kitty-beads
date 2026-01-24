// API Types for Kitty-Beads Dashboard

export interface Feature {
  id: string;
  name: string;
  description?: string;
  artifacts: {
    spec?: boolean;
    plan?: boolean;
    tasks?: boolean;
    research?: boolean;
    quickstart?: boolean;
    data_model?: boolean;
  };
  kanban_stats?: KanbanStats;
}

export interface KanbanStats {
  planned: number;
  doing: number;
  for_review: number;
  done: number;
}

export interface FeaturesResponse {
  features: Feature[];
  project_path: string;
  worktrees_root?: string;
  active_worktree?: string;
  active_mission?: Record<string, string>;
}

export interface IssueCard {
  id: string;
  title: string;
  priority: number;
  status: string;
  assignee?: string;
  labels?: string[];
  is_blocked: boolean;
  blockers?: string[];
}

export interface KanbanLanes {
  planned: IssueCard[];
  doing: IssueCard[];
  for_review: IssueCard[];
  done: IssueCard[];
}

export interface KanbanResponse {
  lanes: KanbanLanes;
}

export interface ArtifactResponse {
  exists: boolean;
  content?: string;
}

export interface Issue {
  id: string;
  title: string;
  description?: string;
  priority: number;
  status: string;
  assignee?: string;
  labels?: string[];
  design?: string;
  acceptance_criteria?: string;
  created_at?: string;
  updated_at?: string;
  issue_type?: "epic" | "task" | "bug" | "story";
  owner?: string;
}

// Roadmap types
export interface RoadmapIssue extends Issue {
  children?: RoadmapIssue[];
  isBlocked?: boolean;
  blockers?: string[];
}

export interface DiagnosticsResponse {
  status: {
    total_issues: number;
    open_issues: number;
    project_path: string;
  };
}

// Terminal WebSocket message types
export interface TerminalMessage {
  type: "input" | "resize" | "output" | "exit" | "error";
  data: string | ResizeData | ExitData | ErrorData;
}

export interface ResizeData {
  cols: number;
  rows: number;
}

export interface ExitData {
  code: number;
}

export interface ErrorData {
  message: string;
}
