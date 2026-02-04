import type { LucideIcon } from "lucide-react";

export interface Command {
  id: string;
  label: string;
  icon?: LucideIcon;
  shortcut?: string;
  group: CommandGroup;
  action: () => void;
}

export type CommandGroup = "navigation" | "beads" | "actions";

export interface CommandPaletteState {
  isOpen: boolean;
}

export interface CommandPaletteContextValue extends CommandPaletteState {
  open: () => void;
  close: () => void;
  toggle: () => void;
}
