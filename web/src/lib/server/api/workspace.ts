import type { ApiResult } from '$lib/types';
import type { CreateWorkspacePayload, WorkspaceData } from '$lib/types/workspace';
import { get, post } from './client';

// Both endpoints are JWT-protected (RequireAuth + RequireActive) — pass the token.

export async function getWorkspaces(token: string): Promise<ApiResult<WorkspaceData[]>> {
	return get<WorkspaceData[]>('/workspaces/', token);
}

export async function createWorkspace(
	token: string,
	p: CreateWorkspacePayload
): Promise<ApiResult<WorkspaceData>> {
	return post<WorkspaceData>('/workspaces/', p, token);
}
