import { useState } from "react";
import { Activity as ActivityIcon } from "lucide-react";
import { useActivityFeed } from "../hooks/useActivityFeed";
import { ActivityItem } from "./ActivityItem";
import { ActivityFilters } from "./ActivityFilters";
import type { ActivityType } from "../types";

export function ActivityFeed() {
  const [selectedTypes, setSelectedTypes] = useState<ActivityType[]>([]);

  const { data: events, isLoading, error } = useActivityFeed({
    types: selectedTypes.length > 0 ? selectedTypes : undefined,
  });

  return (
    <div className="h-full flex flex-col bg-card-bg rounded-xl border border-neon-magenta overflow-hidden">
      {/* Header */}
      <div className="flex items-center gap-3 px-6 py-4 border-b border-night-border bg-night-bg-highlight">
        <ActivityIcon className="w-6 h-6 text-neon-cyan drop-shadow-[0_0_6px_rgba(125,207,255,0.5)]" />
        <div>
          <h2 className="text-neon-cyan text-xl font-bold drop-shadow-[0_0_8px_rgba(125,207,255,0.3)]">
            Activity Feed
          </h2>
          <p className="text-text-muted text-sm">
            Recent changes and updates
          </p>
        </div>
      </div>

      {/* Filters */}
      <ActivityFilters
        selectedTypes={selectedTypes}
        onTypesChange={setSelectedTypes}
      />

      {/* Content */}
      <div className="flex-1 overflow-auto">
        {isLoading ? (
          <div className="flex items-center justify-center py-12">
            <p className="text-text-muted">Loading activity...</p>
          </div>
        ) : error ? (
          <div className="flex items-center justify-center py-12">
            <p className="text-neon-pink">Error loading activity</p>
          </div>
        ) : events && events.length > 0 ? (
          <div className="divide-y divide-night-border">
            {events.map((event) => (
              <ActivityItem key={event.id} event={event} />
            ))}
          </div>
        ) : (
          <div className="flex flex-col items-center justify-center py-12 text-center">
            <ActivityIcon className="w-12 h-12 text-text-muted mb-4 opacity-50" />
            <p className="text-text-muted text-lg mb-2">No activity yet</p>
            <p className="text-text-secondary text-sm">
              Activity will appear here as issues are created and updated
            </p>
          </div>
        )}
      </div>
    </div>
  );
}
