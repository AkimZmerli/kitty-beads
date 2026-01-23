import { EditableArtifact } from "../components/EditableArtifact";

interface QuickstartProps {
  featureId: string | null;
}

export function Quickstart({ featureId }: QuickstartProps) {
  return (
    <EditableArtifact
      featureId={featureId}
      artifactType="quickstart"
      title="Quickstart Guide"
      emptyMessage="No quickstart guide yet."
    />
  );
}
