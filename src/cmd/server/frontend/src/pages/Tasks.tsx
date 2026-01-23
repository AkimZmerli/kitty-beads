import { EditableArtifact } from "../components/EditableArtifact";

interface TasksProps {
  featureId: string | null;
}

export function Tasks({ featureId }: TasksProps) {
  return (
    <EditableArtifact
      featureId={featureId}
      artifactType="tasks"
      title="Task List"
      emptyMessage="No tasks defined yet."
    />
  );
}
