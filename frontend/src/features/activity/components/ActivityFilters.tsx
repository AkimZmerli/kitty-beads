import type { ActivityType } from "../types";

interface ActivityFiltersProps {
  selectedTypes: ActivityType[];
  onTypesChange: (types: ActivityType[]) => void;
}

const FILTER_OPTIONS: { type: ActivityType; label: string }[] = [
  { type: "issue_created", label: "Created" },
  { type: "issue_updated", label: "Updated" },
  { type: "status_changed", label: "Status" },
  { type: "plan_edited", label: "Plans" },
  { type: "assignee_changed", label: "Assigned" },
];

export function ActivityFilters({
  selectedTypes,
  onTypesChange,
}: ActivityFiltersProps) {
  const toggleType = (type: ActivityType) => {
    if (selectedTypes.includes(type)) {
      onTypesChange(selectedTypes.filter((t) => t !== type));
    } else {
      onTypesChange([...selectedTypes, type]);
    }
  };

  const allSelected = selectedTypes.length === 0;

  return (
    <div className="flex items-center gap-2 px-4 py-3 border-b border-night-border">
      <span className="text-xs text-text-muted uppercase tracking-wider font-medium mr-2">
        Filter:
      </span>
      <button
        onClick={() => onTypesChange([])}
        className={`px-3 py-1 text-sm rounded-md transition-colors ${
          allSelected
            ? "bg-neon-cyan text-night-bg"
            : "bg-night-surface text-text-muted hover:text-text-primary"
        }`}
      >
        All
      </button>
      {FILTER_OPTIONS.map(({ type, label }) => (
        <button
          key={type}
          onClick={() => toggleType(type)}
          className={`px-3 py-1 text-sm rounded-md transition-colors ${
            selectedTypes.includes(type)
              ? "bg-neon-cyan text-night-bg"
              : "bg-night-surface text-text-muted hover:text-text-primary"
          }`}
        >
          {label}
        </button>
      ))}
    </div>
  );
}
