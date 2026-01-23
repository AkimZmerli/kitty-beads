import { EditableArtifact } from "../components/EditableArtifact";

interface SpecifyProps {
  featureId: string | null;
}

export function Specify({ featureId }: SpecifyProps) {
  return (
    <EditableArtifact
      featureId={featureId}
      artifactType="spec"
      title="Specification"
      emptyMessage="No specification content yet."
    />
  );
}
