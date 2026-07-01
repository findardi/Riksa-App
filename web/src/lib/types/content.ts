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
    created_by: string;
    created_at: string;
    updated_at: string;
}

export interface FolderTreeNode {
    id: string;
    name: string;
    number: string;
    position: number;
    children: FolderTreeNode[];
}