import { createPortal } from "react-dom";
import { useEffect, useRef, useState } from "react";
import { Check } from "lucide-react";
import type { RoadmapIssue } from "../../../types/api";
import { useIssueMutation } from "../hooks/useIssueMutation";
import { PRIORITY_OPTIONS } from "../constants";
import { getPriorityColor } from "../../../lib/planParser";

interface PriorityDropdownProps {
  issue: RoadmapIssue;
  onClose: () => void;
}

export function PriorityDropdown({ issue, onClose }: PriorityDropdownProps) {
  const { mutate } = useIssueMutation();
  const [position, setPosition] = useState({ top: 0, left: 0 });
  const menuRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    // Position near the priority indicator
    const handlePosition = () => {
      const priorityButtons = document.querySelectorAll("button");
      const priorityButton = Array.from(priorityButtons).find(
        (el) =>
          el.textContent?.includes(`P${issue.priority}`) &&
          el.closest('[class*="cursor-pointer"]')
      );

      if (priorityButton) {
        const rect = priorityButton.getBoundingClientRect();
        setPosition({
          top: rect.bottom + 4,
          left: Math.min(rect.left, window.innerWidth - 180),
        });
      }
    };

    handlePosition();
  }, [issue.priority]);

  useEffect(() => {
    const handleClickOutside = (e: MouseEvent) => {
      if (menuRef.current && !menuRef.current.contains(e.target as Node)) {
        onClose();
      }
    };

    const handleEscape = (e: KeyboardEvent) => {
      if (e.key === "Escape") {
        onClose();
      }
    };

    document.addEventListener("mousedown", handleClickOutside);
    document.addEventListener("keydown", handleEscape);
    return () => {
      document.removeEventListener("mousedown", handleClickOutside);
      document.removeEventListener("keydown", handleEscape);
    };
  }, [onClose]);

  const handleSelect = (priority: number) => {
    if (priority !== issue.priority) {
      mutate({ issueId: issue.id, updates: { priority } });
    }
    onClose();
  };

  return createPortal(
    <div
      ref={menuRef}
      className="fixed z-50 bg-night-surface border border-night-border rounded-lg shadow-xl py-1 min-w-[160px] animate-in fade-in zoom-in-95 duration-100"
      style={{ top: position.top, left: position.left }}
    >
      <div className="px-3 py-1.5 text-xs text-text-muted uppercase tracking-wider border-b border-night-border mb-1">
        Set Priority
      </div>
      {PRIORITY_OPTIONS.map((option) => (
        <button
          key={option.value}
          onClick={() => handleSelect(option.value)}
          className={`w-full px-3 py-2 text-left flex items-center gap-3 hover:bg-night-bg-highlight transition-colors ${
            issue.priority === option.value
              ? "text-neon-cyan"
              : "text-text-normal"
          }`}
        >
          <span
            className="font-bold"
            style={{ color: getPriorityColor(option.value) }}
          >
            P{option.value}
          </span>
          <span className="flex-1 text-sm">{option.label.split(" - ")[1]}</span>
          {issue.priority === option.value && (
            <Check className="w-4 h-4 text-neon-cyan" />
          )}
        </button>
      ))}
    </div>,
    document.body
  );
}
