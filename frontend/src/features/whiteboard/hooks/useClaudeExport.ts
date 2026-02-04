import { useCallback, useState } from "react";
import type { Editor, TLShape } from "tldraw";
import { saveWhiteboardToUser } from "../../../lib/api";

interface ExportResult {
  markdown: string;
  shapes: TLShape[];
}

export type ExportStatus = "idle" | "exporting" | "success" | "error";

export interface PngSaveResult {
  success: boolean;
  path?: string;
  latest?: string;
  error?: string;
}

// Helper to safely get text from shape props
function getShapeText(shape: TLShape): string | null {
  const props = shape.props as Record<string, unknown>;
  if ("text" in props && typeof props.text === "string") {
    return props.text;
  }
  if ("richText" in props && props.richText) {
    // Handle rich text format
    const richText = props.richText as { content?: Array<{ text?: string }> };
    if (richText.content) {
      return richText.content.map((c) => c.text || "").join("");
    }
  }
  return null;
}

// Helper to get geo type
function getGeoType(shape: TLShape): string | null {
  const props = shape.props as Record<string, unknown>;
  if ("geo" in props && typeof props.geo === "string") {
    return props.geo;
  }
  return null;
}

export function useClaudeExport(editor: Editor | null) {
  const [status, setStatus] = useState<ExportStatus>("idle");
  const [lastError, setLastError] = useState<string | null>(null);

  const exportToMarkdown = useCallback((): ExportResult | null => {
    if (!editor) return null;

    const shapes = editor.getCurrentPageShapes();
    const sections: string[] = [];

    // Group shapes by type
    const notes: TLShape[] = [];
    const texts: TLShape[] = [];
    const geos: TLShape[] = [];
    const arrows: TLShape[] = [];
    const others: TLShape[] = [];

    for (const shape of shapes) {
      switch (shape.type) {
        case "note":
          notes.push(shape);
          break;
        case "text":
          texts.push(shape);
          break;
        case "geo":
          geos.push(shape);
          break;
        case "arrow":
          arrows.push(shape);
          break;
        default:
          others.push(shape);
      }
    }

    // Build markdown
    sections.push("# Whiteboard Export\n");

    if (notes.length > 0) {
      sections.push("## Notes\n");
      for (const note of notes) {
        const text = getShapeText(note) || "(empty note)";
        sections.push(`- ${text}`);
      }
      sections.push("");
    }

    if (texts.length > 0) {
      sections.push("## Text Elements\n");
      for (const textShape of texts) {
        const text = getShapeText(textShape) || "(empty text)";
        sections.push(`- ${text}`);
      }
      sections.push("");
    }

    if (geos.length > 0) {
      sections.push("## Shapes\n");
      for (const geo of geos) {
        const geoType = getGeoType(geo) || "shape";
        const label = getShapeText(geo) || "(no label)";
        sections.push(`- **${geoType}**: ${label}`);
      }
      sections.push("");
    }

    if (arrows.length > 0) {
      sections.push("## Connections\n");
      sections.push(`- ${arrows.length} connection(s) drawn`);
      sections.push("");
    }

    if (others.length > 0) {
      sections.push("## Other Elements\n");
      sections.push(`- ${others.length} other shape(s)`);
      sections.push("");
    }

    // Summary
    sections.push("---\n");
    sections.push(
      `*Exported from Kitty-Beads Whiteboard: ${shapes.length} total elements*`
    );

    return {
      markdown: sections.join("\n"),
      shapes,
    };
  }, [editor]);

  const copyToClipboard = useCallback(async () => {
    const result = exportToMarkdown();
    if (result) {
      await navigator.clipboard.writeText(result.markdown);
      return true;
    }
    return false;
  }, [exportToMarkdown]);

  const exportPngToUserStorage = useCallback(async (): Promise<PngSaveResult> => {
    if (!editor) {
      return { success: false, error: "No editor available" };
    }

    setStatus("exporting");
    setLastError(null);

    try {
      // Get all shape IDs for export
      const shapeIds = [...editor.getCurrentPageShapeIds()];
      if (shapeIds.length === 0) {
        setStatus("error");
        setLastError("No shapes to export");
        return { success: false, error: "No shapes to export" };
      }

      // Export to blob using tldraw v4 API
      const { blob } = await editor.toImage(shapeIds, {
        format: "png",
        background: true,
        padding: 32,
      });

      // Convert blob to base64
      const reader = new FileReader();
      const base64Promise = new Promise<string>((resolve, reject) => {
        reader.onloadend = () => {
          const result = reader.result as string;
          resolve(result);
        };
        reader.onerror = reject;
      });
      reader.readAsDataURL(blob);
      const base64Data = await base64Promise;

      // Save to user storage via API
      const response = await saveWhiteboardToUser(base64Data);

      if (response.success) {
        setStatus("success");
        return {
          success: true,
          path: response.path,
          latest: response.latest,
        };
      } else {
        setStatus("error");
        setLastError(response.message || "Save failed");
        return { success: false, error: response.message };
      }
    } catch (err) {
      const errorMsg = err instanceof Error ? err.message : "Export failed";
      setStatus("error");
      setLastError(errorMsg);
      return { success: false, error: errorMsg };
    }
  }, [editor]);

  return {
    exportToMarkdown,
    copyToClipboard,
    exportPngToUserStorage,
    status,
    lastError,
  };
}
