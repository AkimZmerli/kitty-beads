import type { StatusOption, PriorityOption } from "./types";

export const STATUS_OPTIONS: StatusOption[] = [
  { value: "open", label: "Open", color: "neon-blue" },
  { value: "in_progress", label: "In Progress", color: "neon-yellow" },
  { value: "blocked", label: "Blocked", color: "neon-pink" },
  { value: "closed", label: "Done", color: "neon-green" },
];

export const PRIORITY_OPTIONS: PriorityOption[] = [
  { value: 1, label: "P1 - Urgent", color: "neon-pink" },
  { value: 2, label: "P2 - High", color: "neon-orange" },
  { value: 3, label: "P3 - Medium", color: "neon-yellow" },
  { value: 4, label: "P4 - Low", color: "text-muted" },
];

export const FILTER_LABELS: Record<string, string> = {
  all: "All",
  epics: "Epics",
  open: "Open",
  p1: "P1 Only",
};

export const GROUP_BY_OPTIONS = [
  { value: "none", label: "No grouping" },
  { value: "status", label: "Group by Status" },
  { value: "priority", label: "Group by Priority" },
  { value: "assignee", label: "Group by Assignee" },
];

export const KEYBOARD_SHORTCUTS = {
  MOVE_DOWN: ["j", "ArrowDown"] as string[],
  MOVE_UP: ["k", "ArrowUp"] as string[],
  EXPAND_COLLAPSE: ["Enter"] as string[],
  OPEN_MODAL: ["o"] as string[],
  CLOSE: ["Escape"] as string[],
  EDIT: ["e"] as string[],
  CONTEXT_MENU: [" "] as string[],
};
