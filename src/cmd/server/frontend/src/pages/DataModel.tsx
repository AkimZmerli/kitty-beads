import { EditableArtifact } from "../components/EditableArtifact";

interface DataModelProps {
  featureId: string | null;
}

export function DataModel({ featureId }: DataModelProps) {
  return (
    <EditableArtifact
      featureId={featureId}
      artifactType="data_model"
      title="Data Model"
      emptyMessage="No data model defined yet."
    />
  );
}
