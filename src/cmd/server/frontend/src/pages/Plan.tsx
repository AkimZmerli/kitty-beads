import { EditableArtifact } from "../components/EditableArtifact";

interface PlanProps {
  featureId: string | null;
}

export function Plan({ featureId }: PlanProps) {
  return (
    <EditableArtifact
      featureId={featureId}
      artifactType="plan"
      title="Implementation Plan"
      editable={true}
      emptyMessage="No implementation plan yet."
    />
  );
}
