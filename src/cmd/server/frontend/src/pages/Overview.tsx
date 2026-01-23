import type { Feature } from "../types/api";

interface OverviewProps {
  feature: Feature | null;
}

function StatCard({
  label,
  value,
  sublabel,
  isGreen,
}: {
  label: string;
  value: number;
  sublabel?: string;
  isGreen?: boolean;
}) {
  return (
    <div className="bg-white p-5 rounded-lg border border-border">
      <div className="text-xs uppercase tracking-wider text-text-muted mb-2">
        {label}
      </div>
      <div
        className={`text-4xl font-bold ${isGreen ? "text-grassy-green" : "text-text-primary"}`}
      >
        {value}
      </div>
      {sublabel && (
        <div className="text-sm text-text-light mt-1">{sublabel}</div>
      )}
    </div>
  );
}

function ProgressBar({ percentage }: { percentage: number }) {
  return (
    <div className="mt-3 bg-gray-200 h-1.5 rounded-full overflow-hidden">
      <div
        className="h-full bg-grassy-green transition-all duration-300"
        style={{ width: `${percentage}%` }}
      />
    </div>
  );
}

function ArtifactItem({
  icon,
  name,
  available,
}: {
  icon: string;
  name: string;
  available: boolean;
}) {
  return (
    <div className="flex items-center gap-3 px-4 py-3 bg-white border border-border rounded-md mb-2">
      <span>{icon}</span>
      <span className="flex-1">{name}:</span>
      <span className={available ? "text-grassy-green" : "text-text-light"}>
        {available ? "✅ Available" : "⬜ Pending"}
      </span>
    </div>
  );
}

export function Overview({ feature }: OverviewProps) {
  if (!feature) {
    return (
      <div className="bg-white rounded-xl p-12 border-2 border-sunny-yellow-border text-center">
        <h2 className="text-text-muted text-xl">No Feature Selected</h2>
        <p className="text-text-light mt-2">
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
    { key: "spec", name: "Specification", icon: "📄" },
    { key: "plan", name: "Plan", icon: "📋" },
    { key: "tasks", name: "Tasks", icon: "📝" },
    { key: "research", name: "Research", icon: "🔬" },
    { key: "quickstart", name: "Quickstart", icon: "🚀" },
    { key: "data_model", name: "Data Model", icon: "💾" },
  ];

  return (
    <div className="bg-white rounded-xl p-8 border-2 border-sunny-yellow-border">
      <h2 className="text-grassy-green text-2xl font-bold mb-6">
        Feature Overview
      </h2>

      <div className="mb-6">
        <h3 className="text-xl text-text-primary">
          Feature:{" "}
          <span className="font-normal">
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
        <div className="bg-white p-5 rounded-lg border border-border">
          <div className="text-xs uppercase tracking-wider text-text-muted mb-2">
            Completed
          </div>
          <div className="text-4xl font-bold text-grassy-green">
            {stats.done}
          </div>
          <div className="text-sm text-text-light mt-1">{progress}% done</div>
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
            icon={item.icon}
            name={item.name}
            available={!!artifacts[item.key as keyof typeof artifacts]}
          />
        ))}
      </div>
    </div>
  );
}
