import { createPortal } from "react-dom";
import { useEffect, useRef } from "react";
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
} from "lucide-react";
import { useRoadmapContext } from "../context";
import { useIssueMutation } from "../hooks/useIssueMutation";

export function ContextMenu() {
  const navigate = useNavigate();
  const menuRef = useRef<HTMLDivElement>(null);
  const { contextMenu, hideContextMenu, openModal, startEditingTitle } =
    useRoadmapContext();
  const { mutate } = useIssueMutation();

  useEffect(() => {
    if (!contextMenu) return;

    const handleClickOutside = (e: MouseEvent) => {
      if (menuRef.current && !menuRef.current.contains(e.target as Node)) {
        hideContextMenu();
      }
    };

    const handleEscape = (e: KeyboardEvent) => {
      if (e.key === "Escape") {
        hideContextMenu();
      }
    };

    const handleContextMenu = (e: MouseEvent) => {
      e.preventDefault();
      hideContextMenu();
    };

    document.addEventListener("mousedown", handleClickOutside);
    document.addEventListener("keydown", handleEscape);
    document.addEventListener("contextmenu", handleContextMenu);

    return () => {
      document.removeEventListener("mousedown", handleClickOutside);
      document.removeEventListener("keydown", handleEscape);
      document.removeEventListener("contextmenu", handleContextMenu);
    };
  }, [contextMenu, hideContextMenu]);

  if (!contextMenu) return null;

  const { x, y, issueId } = contextMenu;

  // Adjust position to stay within viewport
  const adjustedX = Math.min(x, window.innerWidth - 220);
  const adjustedY = Math.min(y, window.innerHeight - 350);

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
        onClick={() =>
          handleAction(() =>
            mutate({ issueId, updates: { status: "open" } })
          )
        }
      />
      <MenuItem
        icon={Clock}
        label="In Progress"
        onClick={() =>
          handleAction(() =>
            mutate({ issueId, updates: { status: "in_progress" } })
          )
        }
      />
      <MenuItem
        icon={AlertCircle}
        label="Blocked"
        onClick={() =>
          handleAction(() =>
            mutate({ issueId, updates: { status: "blocked" } })
          )
        }
      />
      <MenuItem
        icon={CheckCircle2}
        label="Done"
        onClick={() =>
          handleAction(() =>
            mutate({ issueId, updates: { status: "closed" } })
          )
        }
      />

      <Divider />

      {/* Copy section */}
      <MenuItem
        icon={Copy}
        label="Copy ID"
        onClick={() =>
          handleAction(() => {
            navigator.clipboard.writeText(issueId);
          })
        }
      />
    </div>,
    document.body
  );
}

interface MenuItemProps {
  icon: React.ComponentType<{ className?: string }>;
  label: string;
  shortcut?: string;
  onClick: () => void;
}

function MenuItem({ icon: Icon, label, shortcut, onClick }: MenuItemProps) {
  return (
    <button
      onClick={onClick}
      className="w-full px-3 py-2 text-left flex items-center gap-2.5 text-text-normal hover:bg-night-bg-highlight hover:text-neon-cyan transition-colors"
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
