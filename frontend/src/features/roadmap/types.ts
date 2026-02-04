import type { LucideIcon } from "lucide-react";

export type FilterType = "all" | "epics" | "open" | "p1";
export type GroupByOption = "none" | "status" | "priority" | "assignee";

export interface ContextMenuAction {
  id: string;
  label: string;
  icon: LucideIcon;
  shortcut?: string;
  action: (issueId: string) => void;
  divider?: never;
}

export interface ContextMenuDivider {
  divider: true;
  id?: never;
  label?: never;
  icon?: never;
  shortcut?: never;
  action?: never;
}

export type ContextMenuItem = ContextMenuAction | ContextMenuDivider;

export interface StatusOption {
  value: string;
  label: string;
  color: string;
}

export interface PriorityOption {
  value: number;
  label: string;
  color: string;
}

export interface IssueGroup {
  key: string;
  label: string;
  color: string;
  count: number;
}
