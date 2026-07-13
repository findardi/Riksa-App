import type { ApiResult } from '$lib/types';
import type {
	CompleteUploadPayload,
	DocumentData,
	DownloadUrlData,
	MoveDocumentPayload,
	UploadUrlData
} from '$lib/types/content';
import { del, get, patch, post } from './client';

const foldersBase = (workspaceId: string) => `/content/workspaces/${workspaceId}/folders`;
const documentsBase = (workspaceId: string) => `/content/workspaces/${workspaceId}/documents`;

export function listDocuments(
	token: string,
	workspaceId: string,
	folderId: string
): Promise<ApiResult<DocumentData[]>> {
	return get<DocumentData[]>(`${foldersBase(workspaceId)}/${folderId}/documents`, token);
}

export function requestUploadUrl(
	token: string,
	workspaceId: string,
	folderId: string
): Promise<ApiResult<UploadUrlData>> {
	return post<UploadUrlData>(
		`${foldersBase(workspaceId)}/${folderId}/documents/upload-url`,
		undefined,
		token
	);
}

export function completeUpload(
	token: string,
	workspaceId: string,
	folderId: string,
	p: CompleteUploadPayload
): Promise<ApiResult<DocumentData>> {
	return post<DocumentData>(`${foldersBase(workspaceId)}/${folderId}/documents`, p, token);
}

export function getDownloadUrl(
	token: string,
	workspaceId: string,
	documentId: string
): Promise<ApiResult<DownloadUrlData>> {
	return get<DownloadUrlData>(`${documentsBase(workspaceId)}/${documentId}/download`, token);
}

export function moveDocument(
	token: string,
	workspaceId: string,
	documentId: string,
	p: MoveDocumentPayload
): Promise<ApiResult<null>> {
	return patch<null>(`${documentsBase(workspaceId)}/${documentId}/move`, p, token);
}

export function deleteDocument(
	token: string,
	workspaceId: string,
	documentId: string
): Promise<ApiResult<null>> {
	return del<null>(`${documentsBase(workspaceId)}/${documentId}`, token);
}
