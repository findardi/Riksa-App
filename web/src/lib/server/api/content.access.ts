import type { ApiResult } from '$lib/types';
import type { FolderAccessData, SetFolderAccessPayload } from '$lib/types/content';
import { del, get, put } from './client';

const base = (workspaceId: string) => `/content/workspaces/${workspaceId}`;

export function getFolderAccess(
	token: string,
	workspaceId: string,
	folderId: string
): Promise<ApiResult<FolderAccessData[]>> {
	return get<FolderAccessData[]>(`${base(workspaceId)}/folders/${folderId}/access`, token);
}

export function setFolderAccess(
	token: string,
	workspaceId: string,
	folderId: string,
	p: SetFolderAccessPayload
): Promise<ApiResult<null>> {
	return put<null>(`${base(workspaceId)}/folders/${folderId}/access`, p, token);
}

export function removeFolderAccess(
	token: string,
	workspaceId: string,
	folderId: string,
	groupId: string
): Promise<ApiResult<null>> {
	return del<null>(`${base(workspaceId)}/folders/${folderId}/access/${groupId}`, token);
}
