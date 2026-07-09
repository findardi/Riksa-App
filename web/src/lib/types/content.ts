export interface CreateFolderPayload {
	name: string;
	parent_id: string;
}

export interface MoveFolderPayload {
	parent_id: string;
}

export interface RenameFolderPayload {
	name: string;
}

export interface FolderData {
	id: string;
	workspace_id: string;
	parent_id: string;
	name: string;
	position: number;
	is_default: boolean;
	created_by: string;
	created_at: string;
	updated_at: string;
}

export interface FolderTreeNode {
	id: string;
	name: string;
	number: string;
	position: number;
	is_default: boolean;
	children: FolderTreeNode[];
}

export interface DocumentData {
	id: string;
	folder_id: string;
	name: string;
	version_no: number;
	mime: string;
	size: number;
	created_at: string;
	updated_at: string;
}

export interface UploadUrlData {
	upload_url: string;
	storage_key: string;
}

export interface DownloadUrlData {
	download_url: string;
}

export interface CompleteUploadPayload {
	name: string;
	storage_key: string;
}

export interface MoveDocumentPayload {
	folder_id: string;
}
