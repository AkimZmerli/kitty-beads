import { Search, X } from "lucide-react";
import { useRoadmapContext } from "../context";
import { FILTER_LABELS, GROUP_BY_OPTIONS } from "../constants";
import type { FilterType, GroupByOption } from "../types";

const FILTERS: FilterType[] = ["all", "epics", "open", "p1"];

export function RoadmapFilterBar() {
  const { activeFilter, setFilter, searchQuery, setSearchQuery, groupBy, setGroupBy } =
    useRoadmapContext();

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
      <div className="relative">
        <select
          value={groupBy}
          onChange={(e) => setGroupBy(e.target.value as GroupByOption)}
          className="appearance-none bg-night-surface border border-night-border rounded-lg px-4 py-2 pr-8 text-sm text-text-normal focus:outline-none focus:border-neon-cyan focus:ring-1 focus:ring-neon-cyan/30 transition-all cursor-pointer"
        >
          {GROUP_BY_OPTIONS.map((option) => (
            <option key={option.value} value={option.value}>
              {option.label}
            </option>
          ))}
        </select>
        <div className="absolute right-3 top-1/2 -translate-y-1/2 pointer-events-none text-text-muted">
          <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
          </svg>
        </div>
      </div>
    </div>
  );
}
