package dto

type CreateFolderRequest struct {
	WorkspaceID string `json:"-"`
	CreatedBy   string `json:"-"`
	ParentID    string `json:"parent_id"`
	Name        string `json:"name" validate:"required"`
}

type MoveFolderRequest struct {
	FolderID string `json:"-"`
	ParentID string `json:"parent_id"`
}

type RenameFolderRequest struct {
	FolderID string `json:"-"`
	Name     string `json:"name" validate:"required"`
}

type SetFolderAccessRequest struct {
	WorkspaceID         string `json:"-"`
	FolderID            string `json:"-"`
	GroupID             string `json:"group_id" validate:"required,uuid"`
	CanView             bool   `json:"can_view"`
	CanDownload         bool   `json:"can_download"`
	CanWatermark        bool   `json:"can_watermark"`
	CanDownloadOriginal bool   `json:"can_download_original"`
}
