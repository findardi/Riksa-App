import type { ApiResult } from '$lib/types';
import type {
	AssignMembersPayload,
	GroupMemberData,
	GroupWorkspaceData,
	UpsertGroupWorkspacePayload
} from '$lib/types/workspace';
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

// Group detail = the members assigned to it (the backend returns the member
// list, not the group's name/description).
export async function getGroupDetail(
	token: string,
	workspaceId: string,
	groupId: string
): Promise<ApiResult<GroupMemberData[]>> {
	return get<GroupMemberData[]>(`/access/workspaces/${workspaceId}/groups/${groupId}`, token);
}

// Assign workspace members to a group. Idempotent — already-assigned ids are
// skipped server-side. Returns the group's updated member list.
export async function assignMembers(
	token: string,
	workspaceId: string,
	groupId: string,
	p: AssignMembersPayload
): Promise<ApiResult<GroupMemberData[]>> {
	return post<GroupMemberData[]>(
		`/access/workspaces/${workspaceId}/groups/${groupId}/assign`,
		p,
		token
	);
}

export async function unassignMember(
	token: string,
	workspaceId: string,
	groupId: string,
	memberId: string
): Promise<ApiResult<null>> {
	return del<null>(
		`/access/workspaces/${workspaceId}/groups/${groupId}/unassign/${memberId}`,
		token
	);
}
