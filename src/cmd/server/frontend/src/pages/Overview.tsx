import type { Feature } from "../types/api";
import {
  FileText,
  ClipboardList,
  ListTodo,
  FlaskConical,
  Rocket,
  Database,
  CheckCircle2,
  Circle,
} from "lucide-react";

interface OverviewProps {
  feature: Feature | null;
}

function StatCard({
  label,
  value,
  sublabel,
  isHighlight,
}: {
  label: string;
  value: number;
  sublabel?: string;
  isHighlight?: boolean;
}) {
  return (
    <div className="bg-card-bg p-5 rounded-lg border border-night-border">
      <div className="text-xs uppercase tracking-wider text-text-muted mb-2">
        {label}
      </div>
      <div
        className={`text-4xl font-bold ${isHighlight ? "text-neon-green" : "text-text-primary"}`}
      >
        {value}
      </div>
      {sublabel && (
        <div className="text-sm text-text-muted mt-1">{sublabel}</div>
      )}
    </div>
  );
}

function ProgressBar({ percentage }: { percentage: number }) {
  return (
    <div className="mt-3 bg-night-surface-bright h-1.5 rounded-full overflow-hidden">
      <div
        className="h-full bg-neon-green transition-all duration-300 shadow-[0_0_8px_rgba(158,206,106,0.5)]"
        style={{ width: `${percentage}%` }}
      />
    </div>
  );
}

const ARTIFACT_ICONS = {
  spec: FileText,
  plan: ClipboardList,
  tasks: ListTodo,
  research: FlaskConical,
  quickstart: Rocket,
  data_model: Database,
};

function ArtifactItem({
  iconKey,
  name,
  available,
}: {
  iconKey: string;
  name: string;
  available: boolean;
}) {
  const Icon =
    ARTIFACT_ICONS[iconKey as keyof typeof ARTIFACT_ICONS] || FileText;

  return (
    <div className="flex items-center gap-3 px-4 py-3 bg-card-bg border border-night-border rounded-md mb-2">
      <Icon className="w-5 h-5 text-text-muted" />
      <span className="flex-1 text-text-normal">{name}:</span>
      <span
        className={`flex items-center gap-2 ${available ? "text-neon-green" : "text-text-muted"}`}
      >
        {available ? (
          <>
            <CheckCircle2 className="w-4 h-4" />
            Available
          </>
        ) : (
          <>
            <Circle className="w-4 h-4" />
            Pending
          </>
        )}
      </span>
    </div>
  );
}

export function Overview({ feature }: OverviewProps) {
  if (!feature) {
    return (
      <div className="bg-card-bg rounded-xl p-12 border border-neon-magenta text-center">
        <h2 className="text-text-muted text-xl">No Feature Selected</h2>
        <p className="text-text-muted mt-2">
          Select a feature from the dropdown to view its overview.
        </p>
      </div>
    );
  }

  const stats = feature.kanban_stats || {
    planned: 0,
    doing: 0,
    for_review: 0,
    done: 0,
  };
  const total = stats.planned + stats.doing + stats.for_review + stats.done;
  const progress = total > 0 ? Math.round((stats.done / total) * 100) : 0;

  const artifacts = feature.artifacts || {};
  const artifactItems = [
    { key: "spec", name: "Specification" },
    { key: "plan", name: "Plan" },
    { key: "tasks", name: "Tasks" },
    { key: "research", name: "Research" },
    { key: "quickstart", name: "Quickstart" },
    { key: "data_model", name: "Data Model" },
  ];

  return (
    <div className="bg-card-bg rounded-xl p-8 border border-neon-magenta">
      <h2 className="text-neon-cyan text-2xl font-bold mb-6 drop-shadow-[0_0_8px_rgba(125,207,255,0.3)]">
        Feature Overview
      </h2>

      <div className="mb-6">
        <h3 className="text-xl text-text-primary">
          Feature:{" "}
          <span className="font-normal text-text-secondary">
            {feature.id}-{feature.name.toLowerCase().replace(/\s+/g, "-")}
          </span>
        </h3>
        <p className="text-text-muted">
          View and track all artifacts for this feature
        </p>
      </div>

      {/* Stats Grid */}
      <div className="grid grid-cols-4 gap-5 mb-8">
        <StatCard
          label="Total Tasks"
          value={total}
          sublabel={`${stats.planned} planned`}
        />
        <StatCard label="In Progress" value={stats.doing} />
        <StatCard label="Review" value={stats.for_review} />
        <div className="bg-card-bg p-5 rounded-lg border border-night-border">
          <div className="text-xs uppercase tracking-wider text-text-muted mb-2">
            Completed
          </div>
          <div className="text-4xl font-bold text-neon-green">{stats.done}</div>
          <div className="text-sm text-text-muted mt-1">{progress}% done</div>
          <ProgressBar percentage={progress} />
        </div>
      </div>

      {/* Artifacts Section */}
      <div>
        <h3 className="text-lg font-semibold text-text-primary mb-4">
          Available Artifacts
        </h3>
        {artifactItems.map((item) => (
          <ArtifactItem
            key={item.key}
            iconKey={item.key}
            name={item.name}
            available={!!artifacts[item.key as keyof typeof artifacts]}
          />
        ))}
      </div>
    </div>
  );
}
