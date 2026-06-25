import type { ApiResult } from '$lib/types';
import type { GroupWorkspaceData, UpsertGroupWorkspacePayload } from '$lib/types/workspace';
import { del, get, post, put } from './client';

export async function createGroup(
	token: string,
	workspaceId: string,
	p: UpsertGroupWorkspacePayload
): Promise<ApiResult<GroupWorkspaceData>> {
	return post<GroupWorkspaceData>(`/access/workspaces/${workspaceId}/groups`, p, token);
}

export async function updateGroup(
	token: string,
	workspaceId: string,
	groupId: string,
	p: UpsertGroupWorkspacePayload
): Promise<ApiResult<GroupWorkspaceData>> {
	return put<GroupWorkspaceData>(`/access/workspaces/${workspaceId}/groups/${groupId}`, p, token);
}

export async function getGroups(
	token: string,
	workspaceId: string
): Promise<ApiResult<GroupWorkspaceData[]>> {
	return get<GroupWorkspaceData[]>(`/access/workspaces/${workspaceId}/groups`, token);
}

export async function deleteGroup(
	token: string,
	workspaceId: string,
	groupId: string
): Promise<ApiResult<null>> {
	return del<null>(`/access/workspaces/${workspaceId}/groups/${groupId}`, token);
}
