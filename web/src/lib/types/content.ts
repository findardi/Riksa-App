export interface CreateFolderPayload {
	name: string;
	parent_id: string;
}

// `position` is "insert before whatever currently sits at index N"; omitting it
// appends. The server clamps out-of-range values instead of erroring.
export interface MoveFolderPayload {
	parent_id: string;
	position?: number;
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

export interface ViewMetaData {
	document_id: string;
	name: string;
	mime: string;
	page_count: number;
	can_download_original: boolean;
}

export interface CompleteUploadPayload {
	name: string;
	storage_key: string;
}

export interface MoveDocumentPayload {
	folder_id: string;
	position?: number;
}

export interface FolderAccessData {
	folder_id: string;
	group_id: string;
	group_name: string;
	can_view: boolean;
	can_download: boolean;
	can_watermark: boolean;
	can_download_original: boolean;
}

export interface SetFolderAccessPayload {
	group_id: string;
	can_view: boolean;
	can_download: boolean;
	can_watermark: boolean;
	can_download_original: boolean;
}

export interface InheritedFolderAccess extends FolderAccessData {
	source_folder_id: string;
	source_folder_name: string;
}

export interface DirectFolderAccess extends FolderAccessData {
	shadows: InheritedFolderAccess | null;
}

export interface FolderAccessPanel {
	folder_id: string;
	direct: DirectFolderAccess[];
	inherited: InheritedFolderAccess[];
}
