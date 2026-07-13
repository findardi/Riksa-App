import type { ApiResult } from '$lib/types';
import type { WorkspaceRoleData } from '$lib/types/workspace';
import { get } from './client';

// Roles are fixed system roles (owner/admin/guest), seeded at workspace creation.
// They are read-only via the API — there is no create/update/delete surface.
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

// Permission catalog — single source of truth lives in the Go `permission.All`.
// Used read-only to render each role's granted permissions grouped by resource.
export async function getPermissions(token: string): Promise<ApiResult<string[]>> {
	return get<string[]>(`/access/permissions`, token);
}
