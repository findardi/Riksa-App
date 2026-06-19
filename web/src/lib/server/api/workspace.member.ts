import type { ApiResult } from '$lib/types';
import type { UpdateMemberRolePayload, WorkspaceMemberData } from '$lib/types/workspace';
import { del, get, post, put } from './client';

// Whether an email already has a Wadi account. NOTE: backend route must be POST
// (`/access/check`) — fetch/undici forbid a body on GET.
export async function checkUser(token: string, email: string): Promise<ApiResult<boolean>> {
	return post<boolean>(`/access/check`, { email }, token);
}

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
