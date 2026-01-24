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

async function fetchText(url: string): Promise<string> {
  const response = await fetch(url);
  if (!response.ok) {
    throw new Error(`API error: ${response.status} ${response.statusText}`);
  }
  return response.text();
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

// Constitution API
export async function getConstitution(): Promise<string> {
  return fetchText(`${API_BASE}/constitution`);
}

// Diagnostics API
export async function getDiagnostics(): Promise<DiagnosticsResponse> {
  return fetchJson<DiagnosticsResponse>(`${API_BASE}/diagnostics`);
}
