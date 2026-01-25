import { useState } from "react";
import { MarkdownViewer } from "./MarkdownViewer";
import { useArtifact, useSaveArtifact } from "../hooks/useArtifact";
import { Pencil } from "lucide-react";

interface EditableArtifactProps {
  featureId: string | null;
  artifactType: string;
  title: string;
  editable?: boolean;
  emptyMessage?: string;
}

export function EditableArtifact({
  featureId,
  artifactType,
  title,
  editable = false,
  emptyMessage,
}: EditableArtifactProps) {
  const { data, isLoading, error } = useArtifact(featureId, artifactType);
  const saveMutation = useSaveArtifact(featureId, artifactType);

  const [isEditing, setIsEditing] = useState(false);
  const [editContent, setEditContent] = useState("");

  const handleEdit = () => {
    setEditContent(data?.content || "");
    setIsEditing(true);
  };

  const handleCancel = () => {
    setIsEditing(false);
    setEditContent("");
  };

  const handleSave = async () => {
    await saveMutation.mutateAsync(editContent);
    setIsEditing(false);
    setEditContent("");
  };

  if (!featureId) {
    return (
      <div className="bg-card-bg rounded-xl p-8 border border-neon-magenta">
        <h2 className="text-neon-cyan text-2xl font-bold mb-4 drop-shadow-[0_0_8px_rgba(125,207,255,0.3)]">{title}</h2>
        <p className="text-text-muted text-center py-8">
          Select a feature to view its {title.toLowerCase()}.
        </p>
      </div>
    );
  }

  if (isLoading) {
    return (
      <div className="bg-card-bg rounded-xl p-8 border border-neon-magenta">
        <h2 className="text-neon-cyan text-2xl font-bold mb-4 drop-shadow-[0_0_8px_rgba(125,207,255,0.3)]">{title}</h2>
        <p className="text-text-muted text-center py-8">Loading...</p>
      </div>
    );
  }

  if (error) {
    return (
      <div className="bg-card-bg rounded-xl p-8 border border-neon-magenta">
        <h2 className="text-neon-cyan text-2xl font-bold mb-4 drop-shadow-[0_0_8px_rgba(125,207,255,0.3)]">{title}</h2>
        <p className="text-neon-pink text-center py-8">Error loading content</p>
      </div>
    );
  }

  const hasContent = data?.exists && data?.content;
  const defaultEmpty = `No ${artifactType.replace("_", " ")} content yet.`;

  return (
    <div className="bg-card-bg rounded-xl p-8 border border-neon-magenta">
      <div className="flex justify-between items-center mb-6">
        <h2 className="text-neon-cyan text-2xl font-bold drop-shadow-[0_0_8px_rgba(125,207,255,0.3)]">{title}</h2>
        {editable && !isEditing && (
          <button
            onClick={handleEdit}
            className="flex items-center gap-2 px-4 py-2 bg-night-surface hover:bg-night-surface-bright border border-night-border rounded-md text-text-secondary transition-colors"
          >
            <Pencil className="w-4 h-4" /> Edit
          </button>
        )}
      </div>

      {isEditing ? (
        <div>
          <textarea
            value={editContent}
            onChange={(e) => setEditContent(e.target.value)}
            className="w-full min-h-[400px] p-4 border border-night-border rounded-lg font-mono text-sm leading-relaxed resize-y bg-night-surface text-text-primary focus:outline-none focus:border-neon-cyan placeholder:text-text-muted"
          />
          <div className="flex justify-end gap-3 mt-4">
            <button
              onClick={handleCancel}
              className="px-5 py-2.5 bg-night-surface hover:bg-night-surface-bright border border-night-border rounded-md text-text-secondary transition-colors"
            >
              Cancel
            </button>
            <button
              onClick={handleSave}
              disabled={saveMutation.isPending}
              className="px-5 py-2.5 bg-neon-cyan hover:shadow-[0_0_12px_rgba(125,207,255,0.4)] text-night-bg rounded-md font-semibold transition-all disabled:opacity-50"
            >
              {saveMutation.isPending ? "Saving..." : "Save Changes"}
            </button>
          </div>
        </div>
      ) : hasContent ? (
        <MarkdownViewer content={data.content!} />
      ) : (
        <p className="text-text-muted text-center py-10">
          {emptyMessage || defaultEmpty}
        </p>
      )}
    </div>
  );
}
