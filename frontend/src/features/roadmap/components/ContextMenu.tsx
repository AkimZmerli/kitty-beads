import { createPortal } from "react-dom";
import { useEffect, useRef, useState } from "react";
import { useNavigate } from "react-router-dom";
import {
  Eye,
  Edit,
  Circle,
  Clock,
  CheckCircle2,
  AlertCircle,
  Copy,
  Pencil,
  Trash2,
} from "lucide-react";
import { useRoadmapContext } from "../context";
import { useIssueMutation } from "../hooks/useIssueMutation";
import { useDeleteIssueMutation } from "../hooks/useDeleteIssueMutation";

export function ContextMenu() {
  const navigate = useNavigate();
  const menuRef = useRef<HTMLDivElement>(null);
  const { contextMenu, hideContextMenu, openModal, startEditingTitle, flattenedIssues } =
    useRoadmapContext();
  const { mutate } = useIssueMutation();
  const { mutate: deleteIssue } = useDeleteIssueMutation();
  const [confirmingDelete, setConfirmingDelete] = useState(false);

  useEffect(() => {
    if (!contextMenu) return;

    const { issueId } = contextMenu;
    const issue = flattenedIssues.find((i) => i.id === issueId);
    const isEpic = issue?.issue_type === "epic";

    const handleClickOutside = (e: MouseEvent) => {
      if (menuRef.current && !menuRef.current.contains(e.target as Node)) {
        hideContextMenu();
      }
    };

    const handleContextMenu = (e: MouseEvent) => {
      e.preventDefault();
      hideContextMenu();
    };

    const handleKeyDown = (e: KeyboardEvent) => {
      // Don't intercept if user is typing in an input
      if (e.target instanceof HTMLInputElement || e.target instanceof HTMLTextAreaElement) return;

      switch (e.key) {
        case "Escape":
          e.preventDefault();
          if (confirmingDelete) {
            setConfirmingDelete(false);
          } else {
            hideContextMenu();
          }
          break;
        case "o":
          e.preventDefault();
          hideContextMenu();
          openModal(issueId);
          break;
        case "e":
          e.preventDefault();
          hideContextMenu();
          navigate(`/ideation/${issueId}`);
          break;
        case "r":
          e.preventDefault();
          hideContextMenu();
          startEditingTitle(issueId);
          break;
        case "1":
          e.preventDefault();
          mutate({ issueId, updates: { status: "open" } });
          hideContextMenu();
          break;
        case "2":
          e.preventDefault();
          mutate({ issueId, updates: { status: "in_progress" } });
          hideContextMenu();
          break;
        case "3":
          e.preventDefault();
          mutate({ issueId, updates: { status: "blocked" } });
          hideContextMenu();
          break;
        case "d":
          e.preventDefault();
          mutate({ issueId, updates: { status: "closed" } });
          hideContextMenu();
          break;
        case "c":
          e.preventDefault();
          navigator.clipboard.writeText(issueId);
          hideContextMenu();
          break;
        case "x":
          e.preventDefault();
          if (confirmingDelete) {
            deleteIssue({ issueId, force: isEpic });
            hideContextMenu();
          } else {
            setConfirmingDelete(true);
          }
          break;
      }
    };

    document.addEventListener("mousedown", handleClickOutside);
    document.addEventListener("keydown", handleKeyDown);
    document.addEventListener("contextmenu", handleContextMenu);

    return () => {
      document.removeEventListener("mousedown", handleClickOutside);
      document.removeEventListener("keydown", handleKeyDown);
      document.removeEventListener("contextmenu", handleContextMenu);
    };
  }, [contextMenu, confirmingDelete, hideContextMenu, openModal, startEditingTitle, navigate, mutate, deleteIssue, flattenedIssues]);

  if (!contextMenu) {
    if (confirmingDelete) setConfirmingDelete(false);
    return null;
  }

  const { x, y, issueId } = contextMenu;
  const issue = flattenedIssues.find((i) => i.id === issueId);
  const isEpic = issue?.issue_type === "epic";

  // Adjust position to stay within viewport
  const adjustedX = Math.min(x, window.innerWidth - 220);
  const adjustedY = Math.min(y, window.innerHeight - 380);

  const handleAction = (action: () => void) => {
    action();
    hideContextMenu();
  };

  return createPortal(
    <div
      ref={menuRef}
      className="fixed z-50 bg-night-surface border border-night-border rounded-md shadow-xl py-0.5 min-w-[150px] animate-in fade-in zoom-in-95 duration-100 text-sm"
      style={{ top: adjustedY, left: adjustedX }}
    >
      {/* View/Edit section */}
      <MenuItem
        icon={Eye}
        label="View Full"
        shortcut="o"
        onClick={() => handleAction(() => openModal(issueId))}
      />
      <MenuItem
        icon={Edit}
        label="Edit Plan"
        shortcut="e"
        onClick={() => handleAction(() => navigate(`/ideation/${issueId}`))}
      />
      <MenuItem
        icon={Pencil}
        label="Edit Title"
        shortcut="r"
        onClick={() => handleAction(() => startEditingTitle(issueId))}
      />

      <Divider />

      {/* Status section */}
      <div className="px-2 py-0.5 text-[10px] text-text-muted uppercase tracking-wider">
        Set Status
      </div>
      <MenuItem
        icon={Circle}
        label="Open"
        shortcut="1"
        onClick={() =>
          handleAction(() => mutate({ issueId, updates: { status: "open" } }))
        }
      />
      <MenuItem
        icon={Clock}
        label="In Progress"
        shortcut="2"
        onClick={() =>
          handleAction(() => mutate({ issueId, updates: { status: "in_progress" } }))
        }
      />
      <MenuItem
        icon={AlertCircle}
        label="Blocked"
        shortcut="3"
        onClick={() =>
          handleAction(() => mutate({ issueId, updates: { status: "blocked" } }))
        }
      />
      <MenuItem
        icon={CheckCircle2}
        label="Done"
        shortcut="d"
        onClick={() =>
          handleAction(() => mutate({ issueId, updates: { status: "closed" } }))
        }
      />

      <Divider />

      {/* Copy section */}
      <MenuItem
        icon={Copy}
        label="Copy ID"
        shortcut="c"
        onClick={() =>
          handleAction(() => navigator.clipboard.writeText(issueId))
        }
      />

      <Divider />

      {/* Danger section */}
      {confirmingDelete ? (
        <div className="px-3 py-2">
          <div className="text-xs text-red-400 mb-2">
            {isEpic ? "Delete epic and all tasks?" : "Delete this issue?"}
          </div>
          <div className="flex gap-1.5">
            <button
              autoFocus
              onClick={() => handleAction(() => deleteIssue({ issueId, force: isEpic }))}
              className="flex-1 px-2 py-1 text-xs bg-red-500/20 text-red-300 hover:bg-red-500/30 rounded transition-colors"
            >
              Delete <kbd className="opacity-60">x</kbd>
            </button>
            <button
              onClick={() => setConfirmingDelete(false)}
              className="flex-1 px-2 py-1 text-xs bg-night-bg text-text-muted hover:text-text-normal rounded transition-colors"
            >
              Cancel <kbd className="opacity-60">esc</kbd>
            </button>
          </div>
        </div>
      ) : (
        <MenuItem
          icon={Trash2}
          label={isEpic ? "Delete Epic" : "Delete"}
          shortcut="x"
          danger
          onClick={() => setConfirmingDelete(true)}
        />
      )}
    </div>,
    document.body
  );
}

interface MenuItemProps {
  icon: React.ComponentType<{ className?: string }>;
  label: string;
  shortcut?: string;
  danger?: boolean;
  onClick: () => void;
}

function MenuItem({ icon: Icon, label, shortcut, danger, onClick }: MenuItemProps) {
  return (
    <button
      onClick={onClick}
      className={`w-full px-3 py-2 text-left flex items-center gap-2.5 transition-colors ${
        danger
          ? "text-red-400 hover:bg-red-500/10 hover:text-red-300"
          : "text-text-normal hover:bg-night-bg-highlight hover:text-neon-cyan"
      }`}
    >
      <Icon className="w-4 h-4" />
      <span className="flex-1">{label}</span>
      {shortcut && (
        <kbd className="text-xs text-text-muted bg-night-bg px-1.5 py-0.5 rounded">
          {shortcut}
        </kbd>
      )}
    </button>
  );
}

function Divider() {
  return <div className="my-1 border-t border-night-border" />;
}
