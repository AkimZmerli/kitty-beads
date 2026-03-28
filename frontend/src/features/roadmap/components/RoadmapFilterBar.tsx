import { createPortal } from "react-dom";
import { useEffect, useRef, useState } from "react";
import { Search, X, Check, ChevronDown } from "lucide-react";
import { useRoadmapContext } from "../context";
import { FILTER_LABELS, GROUP_BY_OPTIONS } from "../constants";
import type { FilterType, GroupByOption } from "../types";

const FILTERS: FilterType[] = ["all", "epics", "open", "p1"];

export function RoadmapFilterBar() {
  const { activeFilter, setFilter, searchQuery, setSearchQuery, groupBy, setGroupBy } =
    useRoadmapContext();
  const [dropdownOpen, setDropdownOpen] = useState(false);
  const [dropdownPos, setDropdownPos] = useState({ top: 0, left: 0 });
  const triggerRef = useRef<HTMLButtonElement>(null);
  const menuRef = useRef<HTMLDivElement>(null);

  const activeLabel = GROUP_BY_OPTIONS.find((o) => o.value === groupBy)?.label ?? "No grouping";

  const openDropdown = () => {
    if (triggerRef.current) {
      const rect = triggerRef.current.getBoundingClientRect();
      setDropdownPos({
        top: rect.bottom + 4,
        left: Math.min(rect.left, window.innerWidth - 200),
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

  return (
    <div className="flex items-center gap-4 mb-6">
      {/* Search */}
      <div className="relative flex-1 max-w-md">
        <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-text-muted" />
        <input
          type="text"
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          placeholder="Search by title, ID, or description..."
          className="w-full pl-10 pr-10 py-2 bg-night-surface border border-night-border rounded-lg text-text-primary placeholder:text-text-muted focus:outline-none focus:border-neon-cyan focus:ring-1 focus:ring-neon-cyan/30 transition-all"
        />
        {searchQuery && (
          <button
            onClick={() => setSearchQuery("")}
            className="absolute right-3 top-1/2 -translate-y-1/2 text-text-muted hover:text-text-primary transition-colors"
          >
            <X className="w-4 h-4" />
          </button>
        )}
      </div>

      {/* Filter tabs */}
      <div className="flex gap-1 bg-night-surface rounded-lg p-1">
        {FILTERS.map((filter) => (
          <button
            key={filter}
            onClick={() => setFilter(filter)}
            className={`px-3 py-1.5 rounded-md text-sm font-medium transition-all ${
              activeFilter === filter
                ? "bg-neon-cyan text-night-bg shadow-[0_0_10px_rgba(125,207,255,0.3)]"
                : "text-text-normal hover:bg-night-bg-highlight"
            }`}
          >
            {FILTER_LABELS[filter]}
          </button>
        ))}
      </div>

      {/* Group by selector */}
      <button
        ref={triggerRef}
        onClick={dropdownOpen ? () => setDropdownOpen(false) : openDropdown}
        className={`flex items-center gap-2 bg-night-surface border rounded-lg px-4 py-2 text-sm transition-all cursor-pointer ${
          dropdownOpen
            ? "border-neon-cyan text-text-primary ring-1 ring-neon-cyan/30"
            : "border-night-border text-text-normal hover:border-night-border-hover"
        }`}
      >
        <span>{activeLabel}</span>
        <ChevronDown
          className={`w-4 h-4 text-text-muted transition-transform ${dropdownOpen ? "rotate-180" : ""}`}
        />
      </button>

      {dropdownOpen &&
        createPortal(
          <div
            ref={menuRef}
            className="fixed z-50 bg-night-surface border border-night-border rounded-lg shadow-xl py-1 min-w-[190px] animate-in fade-in zoom-in-95 duration-100"
            style={{ top: dropdownPos.top, left: dropdownPos.left }}
          >
            <div className="px-3 py-1.5 text-xs text-text-muted uppercase tracking-wider border-b border-night-border mb-1">
              Group by
            </div>
            {GROUP_BY_OPTIONS.map((option) => (
              <button
                key={option.value}
                onClick={() => {
                  setGroupBy(option.value as GroupByOption);
                  setDropdownOpen(false);
                }}
                className={`w-full px-3 py-2 text-left flex items-center gap-3 hover:bg-night-bg-highlight transition-colors ${
                  groupBy === option.value ? "text-neon-cyan" : "text-text-normal"
                }`}
              >
                <span className="flex-1">{option.label}</span>
                {groupBy === option.value && <Check className="w-4 h-4 text-neon-cyan" />}
              </button>
            ))}
          </div>,
          document.body
        )}
    </div>
  );
}
