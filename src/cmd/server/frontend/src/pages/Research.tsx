import { EditableArtifact } from "../components/EditableArtifact";

interface ResearchProps {
  featureId: string | null;
}

export function Research({ featureId }: ResearchProps) {
  return (
    <EditableArtifact
      featureId={featureId}
      artifactType="research"
      title="Research"
      editable={true}
      emptyMessage="No research content yet."
    />
  );
}
