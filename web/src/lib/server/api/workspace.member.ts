import type { ApiResult } from '$lib/types';
import type { UpdateMemberRolePayload, WorkspaceMemberData } from '$lib/types/workspace';
import { del, get, put } from './client';

export async function getMembers(
	token: string,
	workspaceId: string
): Promise<ApiResult<WorkspaceMemberData[]>> {
	return get<WorkspaceMemberData[]>(`/access/member/${workspaceId}`, token);
}

export async function updateMemberRole(
	token: string,
	workspaceId: string,
	memberId: string,
	p: UpdateMemberRolePayload
): Promise<ApiResult<WorkspaceMemberData>> {
	return put<WorkspaceMemberData>(`/access/member/${workspaceId}/${memberId}`, p, token);
}

export async function deleteMember(
	token: string,
	workspaceId: string,
	memberId: string
): Promise<ApiResult<null>> {
	return del<null>(`/access/member/${workspaceId}/${memberId}`, token);
}
