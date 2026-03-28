import { useState, useMemo, useEffect, useRef } from "react";
import { createPortal } from "react-dom";
import { GitBranch, ChevronDown, Check } from "lucide-react";
import { useRoadmap } from "../hooks/useRoadmap";
import { TreeGraph } from "../features/tree-graph";
import type { RoadmapIssue } from "../types/api";

export function TreeView() {
  const { data, isLoading, error } = useRoadmap();
  const [selectedEpicId, setSelectedEpicId] = useState<string | null>(null);
  const [dropdownOpen, setDropdownOpen] = useState(false);
  const [dropdownPos, setDropdownPos] = useState({ top: 0, left: 0 });
  const triggerRef = useRef<HTMLButtonElement>(null);
  const menuRef = useRef<HTMLDivElement>(null);

  // Get list of epics (top-level issues or those with children)
  const epics = useMemo(() => {
    if (!data?.issues) return [];

    const allEpics: RoadmapIssue[] = [];

    function findEpics(issues: RoadmapIssue[]) {
      for (const issue of issues) {
        const isEpic =
          issue.issue_type === "epic" ||
          (issue.children && issue.children.length > 0);
        if (isEpic) {
          allEpics.push(issue);
        }
        // Also check nested children for epics
        if (issue.children) {
          findEpics(issue.children);
        }
      }
    }

    findEpics(data.issues);
    return allEpics;
  }, [data?.issues]);

  // Auto-select first epic if none selected
  const selectedEpic = useMemo(() => {
    if (selectedEpicId) {
      return epics.find((e) => e.id === selectedEpicId) || null;
    }
    return epics[0] || null;
  }, [epics, selectedEpicId]);

  const openDropdown = () => {
    if (triggerRef.current) {
      const rect = triggerRef.current.getBoundingClientRect();
      setDropdownPos({
        top: rect.bottom + 4,
        left: Math.min(rect.right - 280, window.innerWidth - 290),
      });
    }
    setDropdownOpen(true);
  };

  useEffect(() => {
    if (!dropdownOpen) return;

    const handleClickOutside = (e: MouseEvent) => {
      if (
        menuRef.current &&
        !menuRef.current.contains(e.target as Node) &&
        !triggerRef.current?.contains(e.target as Node)
      ) {
        setDropdownOpen(false);
      }
    };

    const handleEscape = (e: KeyboardEvent) => {
      if (e.key === "Escape") setDropdownOpen(false);
    };

    document.addEventListener("mousedown", handleClickOutside);
    document.addEventListener("keydown", handleEscape);
    return () => {
      document.removeEventListener("mousedown", handleClickOutside);
      document.removeEventListener("keydown", handleEscape);
    };
  }, [dropdownOpen]);

  if (isLoading) {
    return (
      <div className="bg-card-bg rounded-xl p-8 border border-neon-magenta">
        <div className="flex items-center gap-3 mb-4">
          <GitBranch className="w-6 h-6 text-neon-cyan" />
          <h2 className="text-neon-cyan text-2xl font-bold drop-shadow-[0_0_8px_rgba(125,207,255,0.3)]">
            Tree View
          </h2>
        </div>
        <p className="text-text-muted text-center py-8">Loading tree data...</p>
      </div>
    );
  }

  if (error) {
    return (
      <div className="bg-card-bg rounded-xl p-8 border border-neon-magenta">
        <div className="flex items-center gap-3 mb-4">
          <GitBranch className="w-6 h-6 text-neon-cyan" />
          <h2 className="text-neon-cyan text-2xl font-bold drop-shadow-[0_0_8px_rgba(125,207,255,0.3)]">
            Tree View
          </h2>
        </div>
        <p className="text-neon-pink text-center py-8">
          Error loading tree data
        </p>
      </div>
    );
  }

  return (
    <div className="bg-card-bg rounded-xl p-8 border border-neon-magenta">
      {/* Header */}
      <div className="flex justify-between items-center mb-6">
        <div className="flex items-center gap-3">
          <GitBranch className="w-6 h-6 text-neon-cyan drop-shadow-[0_0_6px_rgba(125,207,255,0.5)]" />
          <div>
            <h2 className="text-neon-cyan text-2xl font-bold drop-shadow-[0_0_8px_rgba(125,207,255,0.3)]">
              Tree View
            </h2>
            <p className="text-text-secondary text-sm">
              Visualize epic hierarchy
            </p>
          </div>
        </div>

        {/* Epic Selector */}
        <button
          ref={triggerRef}
          onClick={dropdownOpen ? () => setDropdownOpen(false) : openDropdown}
          className={`flex items-center gap-2 bg-night-surface border rounded-lg px-4 py-2 text-sm transition-all cursor-pointer max-w-[280px] ${
            dropdownOpen
              ? "border-neon-cyan text-text-primary ring-1 ring-neon-cyan/30"
              : "border-night-border text-text-normal hover:border-night-border-hover"
          }`}
        >
          <span className="truncate flex-1 text-left">
            {selectedEpic
              ? `${selectedEpic.id}: ${selectedEpic.title}`
              : "No epics available"}
          </span>
          <ChevronDown
            className={`w-4 h-4 text-text-muted flex-shrink-0 transition-transform ${dropdownOpen ? "rotate-180" : ""}`}
          />
        </button>

        {dropdownOpen &&
          createPortal(
            <div
              ref={menuRef}
              className="fixed z-50 bg-night-surface border border-night-border rounded-lg shadow-xl py-1 w-[280px] animate-in fade-in zoom-in-95 duration-100 max-h-80 overflow-y-auto"
              style={{ top: dropdownPos.top, left: dropdownPos.left }}
            >
              <div className="px-3 py-1.5 text-xs text-text-muted uppercase tracking-wider border-b border-night-border mb-1">
                Select Epic
              </div>
              {epics.length === 0 ? (
                <div className="px-3 py-2 text-sm text-text-muted">
                  No epics available
                </div>
              ) : (
                epics.map((epic) => (
                  <button
                    key={epic.id}
                    onClick={() => {
                      setSelectedEpicId(epic.id);
                      setDropdownOpen(false);
                    }}
                    className={`w-full px-3 py-2 text-left flex items-center gap-3 hover:bg-night-bg-highlight transition-colors ${
                      selectedEpic?.id === epic.id
                        ? "text-neon-cyan"
                        : "text-text-normal"
                    }`}
                  >
                    <span className="flex-1 truncate text-sm">
                      <span className="text-text-muted mr-1.5">{epic.id}:</span>
                      {epic.title}
                    </span>
                    {selectedEpic?.id === epic.id && (
                      <Check className="w-4 h-4 text-neon-cyan flex-shrink-0" />
                    )}
                  </button>
                ))
              )}
            </div>,
            document.body,
          )}
      </div>
    </div>

      {/* Legend */}
      <div className="flex items-center gap-6 mb-6 text-xs text-text-muted">
        <div className="flex items-center gap-2">
          <div className="w-4 h-4 rounded bg-purple-900/30 border border-neon-magenta" />
          <span>Epic</span>
        </div>
        <div className="flex items-center gap-2">
          <div className="w-4 h-4 rounded bg-cyan-900/30 border border-neon-cyan" />
          <span>Task</span>
        </div>
        <div className="flex items-center gap-2">
          <div className="w-4 h-4 rounded bg-green-900/30 border border-neon-green" />
          <span>Done</span>
        </div>
        <div className="flex items-center gap-2">
          <div className="w-4 h-4 rounded bg-yellow-900/30 border border-neon-yellow" />
          <span>In Progress</span>
        </div>
      </div>

      {/* Tree Graph */}
      {epics.length === 0 ? (
        <div className="text-center py-12 text-text-muted">
          <p className="text-lg mb-2">No epics yet</p>
          <p className="text-sm">
            Create your first epic with{" "}
            <code className="bg-night-surface px-2 py-0.5 rounded text-neon-cyan">
              bd create --type epic "Your epic"
            </code>
          </p>
        </div>
      ) : (
        <TreeGraph rootIssue={selectedEpic} />
      )}
    </div>
  );
}
