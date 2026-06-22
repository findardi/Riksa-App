import type { ApiResult } from '$lib/types';
import type {
	AddMemberResult,
	AddMembersPayload,
	InvitationData,
	UpdateMemberRolePayload,
	WorkspaceMemberData
} from '$lib/types/workspace';
import { del, get, post, put } from './client';

// Bulk invite. The backend resolves account existence per email internally and
// returns a per-email outcome — it never tells the caller who was registered.
export async function addMembers(
	token: string,
	workspaceId: string,
	p: AddMembersPayload
): Promise<ApiResult<AddMemberResult[]>> {
	return post<AddMemberResult[]>(`/access/member/${workspaceId}/invite`, p, token);
}

// Workspace invitations. `status` filters by an exact status; omit for all.
export async function getInvitations(
	token: string,
	workspaceId: string,
	status?: string
): Promise<ApiResult<InvitationData[]>> {
	const q = status ? `?status=${encodeURIComponent(status)}` : '';
	return get<InvitationData[]>(`/access/member/${workspaceId}/invite${q}`, token);
}

// Re-issue an invitation's token and resend its email (no body).
export async function resendInvitation(
	token: string,
	workspaceId: string,
	invitationId: string
): Promise<ApiResult<null>> {
	return post<null>(`/access/member/${workspaceId}/${invitationId}/resend`, undefined, token);
}

// Revoke a pending invitation, invalidating its link (no body).
export async function revokeInvitation(
	token: string,
	workspaceId: string,
	invitationId: string
): Promise<ApiResult<null>> {
	return post<null>(`/access/member/${workspaceId}/${invitationId}/revoke`, undefined, token);
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
