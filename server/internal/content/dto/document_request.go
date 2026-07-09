package dto

type CompleteUploadRequest struct {
	WorkspaceID string `json:"-"`
	FolderID    string `json:"-"`
	UploadedBy  string `json:"-"`
	Name        string `json:"name" validate:"required"`
	StorageKey  string `json:"storage_key" validate:"required"`
}

type CompleteVersionRequest struct {
	WorkspaceID string `json:"-"`
	DocumentID  string `json:"-"`
	UploadedBy  string `json:"-"`
	StorageKey  string `json:"storage_key" validate:"required"`
}

type MoveDocumentRequest struct {
	WorkspaceID string `json:"-"`
	DocumentID  string `json:"-"`
	FolderID    string `json:"folder_id" validate:"required,uuid"`
}
