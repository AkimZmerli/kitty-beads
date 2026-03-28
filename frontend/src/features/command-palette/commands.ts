import {
  Map,
  GitBranch,
  Target,
  PenTool,
  Code2,
  Activity,
  Settings,
  Terminal,
  Plus,
} from "lucide-react";
import type { CommandGroup } from "./types";

export interface CommandDefinition {
  id: string;
  label: string;
  icon: typeof Map;
  shortcut?: string;
  group: CommandGroup;
  path?: string;
  action?: string;
}

export const navigationCommands: CommandDefinition[] = [
  { id: "nav-roadmap", label: "Go to Roadmap", icon: Map, group: "navigation", path: "/roadmap" },
  { id: "nav-tree", label: "Go to Tree View", icon: GitBranch, group: "navigation", path: "/tree" },
  { id: "nav-kanban", label: "Go to Kanban", icon: Target, group: "navigation", path: "/kanban" },
  { id: "nav-whiteboard", label: "Go to Whiteboard", icon: PenTool, group: "navigation", path: "/whiteboard" },
  { id: "nav-editor", label: "Go to Editor", icon: Code2, group: "navigation", path: "/editor" },
  { id: "nav-activity", label: "Go to Activity", icon: Activity, group: "navigation", path: "/activity" },
  { id: "nav-diagnostics", label: "Go to Diagnostics", icon: Settings, group: "navigation", path: "/diagnostics" },
];

export const actionCommands: CommandDefinition[] = [
  { id: "action-terminal", label: "Toggle Terminal", icon: Terminal, shortcut: "Ctrl+`", group: "actions", action: "toggleTerminal" },
  { id: "action-new-bead", label: "Create New Bead", icon: Plus, group: "actions", path: "/ideation/new" },
];
