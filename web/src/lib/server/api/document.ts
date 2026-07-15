import type { ApiResult } from '$lib/types';
import type {
	CompleteUploadPayload,
	DocumentData,
	DownloadUrlData,
	MoveDocumentPayload,
	UploadUrlData,
	ViewMetaData
} from '$lib/types/content';
import { API_URL, del, get, patch, post } from './client';

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

export function getViewMeta(
	token: string,
	workspaceId: string,
	documentId: string
): Promise<ApiResult<ViewMetaData>> {
	return get<ViewMetaData>(`${documentsBase(workspaceId)}/${documentId}/view`, token);
}

// Raw upstream response for the page-image proxy. This endpoint streams a
// watermarked PNG (Content-Type image/png), not a JSON envelope, so it bypasses
// the typed client entirely — the proxy route forwards the body and status.
export function fetchViewPage(
	token: string,
	workspaceId: string,
	documentId: string,
	page: number | string
): Promise<Response> {
	return fetch(`${API_URL}${documentsBase(workspaceId)}/${documentId}/pages/${page}`, {
		headers: { authorization: `Bearer ${token}` }
	});
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
