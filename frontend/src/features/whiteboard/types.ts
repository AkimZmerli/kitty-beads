import type { TLStoreSnapshot } from "tldraw";

export interface WhiteboardData {
  id: string;
  beadId?: string;
  name: string;
  snapshot: TLStoreSnapshot;
  createdAt?: string;
  updatedAt?: string;
}

export interface WhiteboardProps {
  beadId?: string;
}

export interface ExportOptions {
  includeShapes: boolean;
  includeText: boolean;
  format: "markdown" | "json";
}
