import type { ApiResult } from '$lib/types';
import type {
	CreateWorkspaceRolePayload,
	UpdateWorkspaceRolePayload,
	WorkspaceRoleData
} from '$lib/types/workspace';
import { del, get, post, put } from './client';

export async function getRoles(
	token: string,
	workspaceId: string
): Promise<ApiResult<WorkspaceRoleData[]>> {
	return get<WorkspaceRoleData[]>(`/access/workspaces/${workspaceId}/roles`, token);
}

export async function getRole(
	token: string,
	workspaceId: string,
	roleId: string
): Promise<ApiResult<WorkspaceRoleData>> {
	return get<WorkspaceRoleData>(`/access/workspaces/${workspaceId}/roles/${roleId}`, token);
}

export async function createRole(
	token: string,
	workspaceId: string,
	p: CreateWorkspaceRolePayload
): Promise<ApiResult<WorkspaceRoleData>> {
	return post<WorkspaceRoleData>(`/access/workspaces/${workspaceId}/roles`, p, token);
}

export async function updateRole(
	token: string,
	workspaceId: string,
	roleId: string,
	p: UpdateWorkspaceRolePayload
): Promise<ApiResult<WorkspaceRoleData>> {
	return put<WorkspaceRoleData>(`/access/workspaces/${workspaceId}/roles/${roleId}`, p, token);
}

export async function deleteRole(
	token: string,
	workspaceId: string,
	roleId: string
): Promise<ApiResult<null>> {
	return del<null>(`/access/workspaces/${workspaceId}/roles/${roleId}`, token);
}

// Permission catalog — single source of truth lives in the Go `permission.All`.
export async function getPermissions(token: string): Promise<ApiResult<string[]>> {
	return get<string[]>(`/access/permissions`, token);
}
