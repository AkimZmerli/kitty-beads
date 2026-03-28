import { createPortal } from "react-dom";
import { useEffect, useRef, useState } from "react";
import { Check } from "lucide-react";
import type { RoadmapIssue } from "../../../types/api";
import { useIssueMutation } from "../hooks/useIssueMutation";
import { STATUS_OPTIONS } from "../constants";

interface StatusDropdownProps {
  issue: RoadmapIssue;
  onClose: () => void;
}

export function StatusDropdown({ issue, onClose }: StatusDropdownProps) {
  const { mutate } = useIssueMutation();
  const [position, setPosition] = useState({ top: 0, left: 0 });
  const menuRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    // Position near the status badge (approximate)
    const handlePosition = () => {
      // Find the status badge button that triggered this
      const badges = document.querySelectorAll('[class*="rounded-full"]');
      const statusBadge = Array.from(badges).find(
        (el) =>
          el.textContent?.includes(issue.status?.toUpperCase() || "OPEN") &&
          el.closest('[class*="cursor-pointer"]')
      );

      if (statusBadge) {
        const rect = statusBadge.getBoundingClientRect();
        setPosition({
          top: rect.bottom + 4,
          left: Math.min(rect.left, window.innerWidth - 180),
        });
      }
    };

    handlePosition();
  }, [issue.status]);

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

  const handleSelect = (status: string) => {
    if (status !== issue.status) {
      mutate({ issueId: issue.id, updates: { status } });
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
        Set Status
      </div>
      {STATUS_OPTIONS.map((option) => (
        <button
          key={option.value}
          onClick={() => handleSelect(option.value)}
          className={`w-full px-3 py-2 text-left flex items-center gap-3 hover:bg-night-bg-highlight transition-colors ${
            issue.status === option.value ? "text-neon-cyan" : "text-text-normal"
          }`}
        >
          <span
            className={`w-2 h-2 rounded-full bg-${option.color}`}
            style={{
              backgroundColor:
                option.color === "neon-blue"
                  ? "#7aa2f7"
                  : option.color === "neon-yellow"
                    ? "#e0af68"
                    : option.color === "neon-pink"
                      ? "#f7768e"
                      : option.color === "neon-green"
                        ? "#9ece6a"
                        : undefined,
            }}
          />
          <span className="flex-1">{option.label}</span>
          {issue.status === option.value && (
            <Check className="w-4 h-4 text-neon-cyan" />
          )}
        </button>
      ))}
    </div>,
    document.body
  );
}
