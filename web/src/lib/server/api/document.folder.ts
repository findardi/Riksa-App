import type { ApiResult } from '$lib/types';
import type {
	BulkCreateFolderData,
	BulkCreateFolderPayload,
	CreateFolderPayload,
	FolderData,
	FolderTreeNode,
	MoveFolderPayload,
	RenameFolderPayload
} from '$lib/types/content';
import { del, get, patch, post, put } from './client';

const foldersBase = (workspaceId: string) => `/content/workspaces/${workspaceId}/folders`;

export function getFoldersTree(
	token: string,
	workspaceId: string
): Promise<ApiResult<FolderTreeNode[]>> {
	return get<FolderTreeNode[]>(foldersBase(workspaceId), token);
}

export function createFolder(
	token: string,
	workspaceId: string,
	p: CreateFolderPayload
): Promise<ApiResult<FolderData>> {
	return post<FolderData>(foldersBase(workspaceId), p, token);
}

export function bulkCreateFolders(
	token: string,
	workspaceId: string,
	p: BulkCreateFolderPayload
): Promise<ApiResult<BulkCreateFolderData>> {
	return post<BulkCreateFolderData>(`${foldersBase(workspaceId)}/bulk`, p, token);
}

export function renameFolder(
	token: string,
	workspaceId: string,
	folderId: string,
	p: RenameFolderPayload
): Promise<ApiResult<FolderData>> {
	return put<FolderData>(`${foldersBase(workspaceId)}/${folderId}`, p, token);
}

export function moveFolder(
	token: string,
	workspaceId: string,
	folderId: string,
	p: MoveFolderPayload
): Promise<ApiResult<null>> {
	return patch<null>(`${foldersBase(workspaceId)}/${folderId}/move`, p, token);
}

export function deleteFolder(
	token: string,
	workspaceId: string,
	folderId: string
): Promise<ApiResult<null>> {
	return del<null>(`${foldersBase(workspaceId)}/${folderId}`, token);
}
