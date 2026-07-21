import type { ApiResult } from '$lib/types';
import type {
	AbortMultipartPayload,
	CompleteMultipartPayload,
	CompleteUploadPayload,
	DocumentData,
	DownloadUrlData,
	InitMultipartData,
	InitMultipartPayload,
	MoveDocumentPayload,
	MultipartPartsData,
	MultipartPartUrlsData,
	MultipartPartUrlsPayload,
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
	folderId: string,
	storageKey?: string
): Promise<ApiResult<UploadUrlData>> {
	return post<UploadUrlData>(
		`${foldersBase(workspaceId)}/${folderId}/documents/upload-url`,
		storageKey ? { storage_key: storageKey } : undefined,
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

const multipartBase = (workspaceId: string, folderId: string) =>
	`${foldersBase(workspaceId)}/${folderId}/documents/multipart`;

export function initMultipart(
	token: string,
	workspaceId: string,
	folderId: string,
	p: InitMultipartPayload
): Promise<ApiResult<InitMultipartData>> {
	return post<InitMultipartData>(`${multipartBase(workspaceId, folderId)}/init`, p, token);
}

// Upstream caps a batch at 100 part numbers and the presigned URLs expire in
// 15 minutes, so callers request them in waves rather than all up front.
export function multipartPartUrls(
	token: string,
	workspaceId: string,
	folderId: string,
	p: MultipartPartUrlsPayload
): Promise<ApiResult<MultipartPartUrlsData>> {
	return post<MultipartPartUrlsData>(`${multipartBase(workspaceId, folderId)}/part-urls`, p, token);
}

// The resume read: which parts object storage already holds. Query string here,
// unlike abort, which takes the same pair as a JSON body.
export function multipartParts(
	token: string,
	workspaceId: string,
	folderId: string,
	uploadId: string,
	storageKey: string
): Promise<ApiResult<MultipartPartsData>> {
	const q = new URLSearchParams({ upload_id: uploadId, storage_key: storageKey });
	return get<MultipartPartsData>(`${multipartBase(workspaceId, folderId)}/parts?${q}`, token);
}

export function completeMultipart(
	token: string,
	workspaceId: string,
	folderId: string,
	p: CompleteMultipartPayload
): Promise<ApiResult<DocumentData>> {
	return post<DocumentData>(`${multipartBase(workspaceId, folderId)}/complete`, p, token);
}

export function abortMultipart(
	token: string,
	workspaceId: string,
	folderId: string,
	p: AbortMultipartPayload
): Promise<ApiResult<null>> {
	return del<null>(multipartBase(workspaceId, folderId), token, p);
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
