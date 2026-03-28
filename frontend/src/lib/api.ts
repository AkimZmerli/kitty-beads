// API Client for Kitty-Beads Dashboard

import type {
  FeaturesResponse,
  KanbanResponse,
  ArtifactResponse,
  Issue,
  DiagnosticsResponse,
} from "../types/api";

const API_BASE = "/api";

async function fetchJson<T>(url: string, options?: RequestInit): Promise<T> {
  const response = await fetch(url, options);
  if (!response.ok) {
    throw new Error(`API error: ${response.status} ${response.statusText}`);
  }
  return response.json();
}

// Features API
export async function getFeatures(): Promise<FeaturesResponse> {
  return fetchJson<FeaturesResponse>(`${API_BASE}/features`);
}

// Kanban API
export async function getKanban(featureId: string): Promise<KanbanResponse> {
  return fetchJson<KanbanResponse>(`${API_BASE}/kanban/${featureId}`);
}

// Artifact API
export async function getArtifact(
  featureId: string,
  artifactType: string,
): Promise<ArtifactResponse> {
  return fetchJson<ArtifactResponse>(
    `${API_BASE}/artifact/${featureId}/${artifactType}`,
  );
}

export async function saveArtifact(
  featureId: string,
  artifactType: string,
  content: string,
): Promise<void> {
  await fetch(`${API_BASE}/artifact/${featureId}/${artifactType}`, {
    method: "PUT",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ content }),
  });
}

// Issue API
export async function getIssue(issueId: string): Promise<Issue> {
  return fetchJson<Issue>(`${API_BASE}/issues/${issueId}`);
}

export async function getAllIssues(): Promise<Issue[]> {
  return fetchJson<Issue[]>(`${API_BASE}/issues`);
}

export async function updateIssue(
  issueId: string,
  updates: Partial<Issue>,
): Promise<Issue> {
  return fetchJson<Issue>(`${API_BASE}/issues/${issueId}`, {
    method: "PUT",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(updates),
  });
}

export async function deleteIssue(
  issueId: string,
  force = false,
): Promise<void> {
  const url = force
    ? `${API_BASE}/issues/${issueId}?force=true`
    : `${API_BASE}/issues/${issueId}`;
  const response = await fetch(url, { method: "DELETE" });
  if (!response.ok) {
    throw new Error(`API error: ${response.status} ${response.statusText}`);
  }
}

// Diagnostics API
export async function getDiagnostics(): Promise<DiagnosticsResponse> {
  return fetchJson<DiagnosticsResponse>(`${API_BASE}/diagnostics`);
}

// Whiteboard API
export interface WhiteboardSaveResponse {
  success: boolean;
  path: string;
  latest: string;
  message?: string;
  file_size: number;
}

export async function saveWhiteboardToUser(
  imageData: string,
): Promise<WhiteboardSaveResponse> {
  return fetchJson<WhiteboardSaveResponse>(`${API_BASE}/whiteboard/user/save`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ image_data: imageData }),
  });
}
