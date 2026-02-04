import { useCallback, useMemo } from "react";
import { Command } from "cmdk";
import { useNavigate } from "react-router-dom";
import {
  Search,
  FileText,
  AlertCircle,
  CheckCircle2,
  Circle,
} from "lucide-react";
import { useCommandPalette } from "../context";
import { navigationCommands, actionCommands } from "../commands";
import { useTerminal } from "../../terminal";
import { useRoadmap } from "../../../hooks/useRoadmap";
import type { RoadmapIssue } from "../../../types/api";

function flattenIssues(issues: RoadmapIssue[]): RoadmapIssue[] {
  const result: RoadmapIssue[] = [];
  for (const issue of issues) {
    result.push(issue);
    if (issue.children && issue.children.length > 0) {
      result.push(...flattenIssues(issue.children));
    }
  }
  return result;
}

function getStatusIcon(status: string) {
  switch (status) {
    case "closed":
      return <CheckCircle2 className="w-4 h-4 text-neon-green" />;
    case "in_progress":
      return <AlertCircle className="w-4 h-4 text-neon-cyan" />;
    default:
      return <Circle className="w-4 h-4 text-text-muted" />;
  }
}

function getPriorityBadge(priority: number) {
  const colors: Record<number, string> = {
    1: "bg-neon-pink/20 text-neon-pink",
    2: "bg-orange-500/20 text-orange-400",
    3: "bg-yellow-500/20 text-yellow-400",
    4: "bg-text-muted/20 text-text-muted",
    5: "bg-text-muted/10 text-text-muted",
  };
  return (
    <span className={`px-1.5 py-0.5 rounded text-xs font-medium ${colors[priority] || colors[4]}`}>
      P{priority}
    </span>
  );
}

export function CommandPalette() {
  const { isOpen, close } = useCommandPalette();
  const { togglePanel } = useTerminal();
  const navigate = useNavigate();
  const { data: roadmapData } = useRoadmap();

  const allIssues = useMemo(() => {
    if (!roadmapData?.issues) return [];
    return flattenIssues(roadmapData.issues);
  }, [roadmapData?.issues]);

  const handleSelect = useCallback(
    (value: string) => {
      close();

      // Handle navigation commands
      const navCommand = navigationCommands.find((c) => c.id === value);
      if (navCommand?.path) {
        navigate(navCommand.path);
        return;
      }

      // Handle action commands
      const actionCommand = actionCommands.find((c) => c.id === value);
      if (actionCommand) {
        if (actionCommand.action === "toggleTerminal") {
          togglePanel();
        } else if (actionCommand.path) {
          navigate(actionCommand.path);
        }
        return;
      }

      // Handle bead selection (navigate to ideation)
      if (value.startsWith("bead-")) {
        const issueId = value.replace("bead-", "");
        navigate(`/ideation/${issueId}`);
      }
    },
    [close, navigate, togglePanel]
  );

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 z-[100]">
      {/* Backdrop */}
      <div
        className="absolute inset-0 bg-black/60 backdrop-blur-sm"
        onClick={close}
      />

      {/* Palette */}
      <div className="absolute top-[20%] left-1/2 -translate-x-1/2 w-full max-w-xl">
        <Command
          className="bg-card-bg border border-night-border rounded-xl shadow-2xl overflow-hidden"
          onKeyDown={(e) => {
            if (e.key === "Escape") {
              close();
            }
          }}
        >
          {/* Search input */}
          <div className="flex items-center gap-3 px-4 border-b border-night-border">
            <Search className="w-5 h-5 text-text-muted" />
            <Command.Input
              placeholder="Search commands or beads..."
              className="w-full py-4 bg-transparent text-text-primary placeholder:text-text-muted outline-none"
              autoFocus
            />
          </div>

          {/* Results */}
          <Command.List className="max-h-80 overflow-y-auto p-2">
            <Command.Empty className="py-6 text-center text-text-muted">
              No results found.
            </Command.Empty>

            {/* Navigation group */}
            <Command.Group
              heading="Navigation"
              className="text-xs text-text-muted uppercase tracking-wider px-2 py-1.5"
            >
              {navigationCommands.map((cmd) => (
                <Command.Item
                  key={cmd.id}
                  value={cmd.id}
                  onSelect={handleSelect}
                  className="flex items-center gap-3 px-3 py-2.5 rounded-lg cursor-pointer text-text-secondary data-[selected=true]:bg-night-bg data-[selected=true]:text-neon-cyan transition-colors"
                >
                  <cmd.icon className="w-4 h-4" />
                  <span>{cmd.label}</span>
                </Command.Item>
              ))}
            </Command.Group>

            {/* Beads group */}
            {allIssues.length > 0 && (
              <Command.Group
                heading="Beads"
                className="text-xs text-text-muted uppercase tracking-wider px-2 py-1.5 mt-2"
              >
                {allIssues.slice(0, 20).map((issue) => (
                  <Command.Item
                    key={issue.id}
                    value={`bead-${issue.id}`}
                    keywords={[issue.title, issue.description || "", issue.id]}
                    onSelect={handleSelect}
                    className="flex items-center gap-3 px-3 py-2.5 rounded-lg cursor-pointer text-text-secondary data-[selected=true]:bg-night-bg data-[selected=true]:text-neon-cyan transition-colors"
                  >
                    {getStatusIcon(issue.status)}
                    <div className="flex-1 min-w-0">
                      <div className="flex items-center gap-2">
                        <FileText className="w-3.5 h-3.5 text-text-muted shrink-0" />
                        <span className="truncate">{issue.title}</span>
                      </div>
                      <span className="text-xs text-text-muted">{issue.id}</span>
                    </div>
                    {getPriorityBadge(issue.priority)}
                  </Command.Item>
                ))}
              </Command.Group>
            )}

            {/* Actions group */}
            <Command.Group
              heading="Actions"
              className="text-xs text-text-muted uppercase tracking-wider px-2 py-1.5 mt-2"
            >
              {actionCommands.map((cmd) => (
                <Command.Item
                  key={cmd.id}
                  value={cmd.id}
                  onSelect={handleSelect}
                  className="flex items-center gap-3 px-3 py-2.5 rounded-lg cursor-pointer text-text-secondary data-[selected=true]:bg-night-bg data-[selected=true]:text-neon-cyan transition-colors"
                >
                  <cmd.icon className="w-4 h-4" />
                  <span className="flex-1">{cmd.label}</span>
                  {cmd.shortcut && (
                    <kbd className="text-xs text-text-muted bg-night-bg px-1.5 py-0.5 rounded">
                      {cmd.shortcut}
                    </kbd>
                  )}
                </Command.Item>
              ))}
            </Command.Group>
          </Command.List>

          {/* Footer hint */}
          <div className="flex items-center justify-between px-4 py-2 border-t border-night-border text-xs text-text-muted">
            <div className="flex items-center gap-4">
              <span>
                <kbd className="bg-night-bg px-1 py-0.5 rounded mr-1">↑↓</kbd>
                navigate
              </span>
              <span>
                <kbd className="bg-night-bg px-1 py-0.5 rounded mr-1">↵</kbd>
                select
              </span>
              <span>
                <kbd className="bg-night-bg px-1 py-0.5 rounded mr-1">esc</kbd>
                close
              </span>
            </div>
          </div>
        </Command>
      </div>
    </div>
  );
}
